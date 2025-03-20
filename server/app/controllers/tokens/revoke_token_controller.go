package tokens

import (
	"net/http"

	"koneksi/server/app/helpers"
	"koneksi/server/app/repositories"
	"koneksi/server/app/services"

	"github.com/gin-gonic/gin"
)

// RevokeTokenController handles revoking refresh tokens
type RevokeTokenController struct {
	userRepo   *repositories.UserRepository
	jwtService *services.JWTService
}

// NewRevokeTokenController initializes a new RevokeTokenController
func NewRevokeTokenController(userRepo *repositories.UserRepository, jwtService *services.JWTService) *RevokeTokenController {
	return &RevokeTokenController{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Handle processes the revoke token request (Logout)
func (rc *RevokeTokenController) Handle(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Validate the payload
	if err := rc.validatePayload(c, &request); err != nil {
		return
	}

	// Validate refresh token
	claims, err := rc.jwtService.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		helpers.FormatResponse(c, "error", http.StatusUnauthorized, "Invalid or expired refresh token", nil, nil)
		return
	}

	// Check if user exists
	user, err := rc.userRepo.ReadUserByEmail(c.Request.Context(), *claims.Email)
	if err != nil || user == nil {
		helpers.FormatResponse(c, "error", http.StatusUnauthorized, "User no longer exists", nil, nil)
		return
	}

	// Revoke the refresh token (remove from Redis)
	err = rc.jwtService.RevokeRefreshToken(user.ID.Hex())
	if err != nil {
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Could not revoke refresh token", nil, nil)
		return
	}

	// Return success response
	helpers.FormatResponse(c, "success", http.StatusOK, "Refresh token revoked successfully", nil, nil)
}

// validatePayload validates the incoming request payload
func (rc *RevokeTokenController) validatePayload(c *gin.Context, request any) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helpers.FormatResponse(c, "error", http.StatusBadRequest, "Invalid input", nil, nil)
		return err
	}
	return nil
}
