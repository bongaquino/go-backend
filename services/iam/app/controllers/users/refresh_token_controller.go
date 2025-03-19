package users

import (
	"net/http"

	"koneksi/services/iam/app/helpers"
	"koneksi/services/iam/app/repositories"
	"koneksi/services/iam/app/services"

	"github.com/gin-gonic/gin"
)

// RefreshTokenController handles the JWT refresh process
type RefreshTokenController struct {
	userRepo   *repositories.UserRepository
	jwtService *services.JWTService
}

// NewRefreshTokenController initializes a new RefreshTokenController
func NewRefreshTokenController(userRepo *repositories.UserRepository, jwtService *services.JWTService) *RefreshTokenController {
	return &RefreshTokenController{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Handle processes the refresh token request
func (rc *RefreshTokenController) Handle(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Validate the payload
	if err := rc.validatePayload(c, &request); err != nil {
		return
	}

	// Validate refresh token using ValidateRefreshToken
	claims, err := rc.jwtService.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		helpers.FormatResponse(c, "error", http.StatusUnauthorized, "Invalid or expired refresh token", nil, nil)
		return
	}

	// Check if user still exists
	user, err := rc.userRepo.ReadUserByEmail(c.Request.Context(), claims.Email)
	if err != nil || user == nil {
		helpers.FormatResponse(c, "error", http.StatusUnauthorized, "User no longer exists", nil, nil)
		return
	}

	// Generate new access & refresh tokens
	accessToken, refreshToken, err := rc.jwtService.GenerateTokens(user.ID.Hex(), user.Email)
	if err != nil {
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Could not generate new tokens", nil, nil)
		return
	}

	// Return new tokens
	helpers.FormatResponse(c, "success", http.StatusOK, "Token refreshed successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// validatePayload validates the incoming request payload
func (rc *RefreshTokenController) validatePayload(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helpers.FormatResponse(c, "error", http.StatusBadRequest, "Invalid input", nil, nil)
		return err
	}
	return nil
}
