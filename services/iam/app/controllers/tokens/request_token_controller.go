package tokens

import (
	"net/http"

	"koneksi/services/iam/app/helpers"
	"koneksi/services/iam/app/repositories"
	"koneksi/services/iam/app/services"

	"github.com/gin-gonic/gin"
)

// RequestTokenController handles user authentication and token generation
type RequestTokenController struct {
	userRepo   *repositories.UserRepository
	jwtService *services.JWTService
}

// NewRequestTokenController initializes a new RequestTokenController
func NewRequestTokenController(userRepo *repositories.UserRepository, jwtService *services.JWTService) *RequestTokenController {
	return &RequestTokenController{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Handle processes the login request and returns an access & refresh token
func (rc *RequestTokenController) Handle(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	// Validate the payload
	if err := rc.validatePayload(c, &request); err != nil {
		return
	}

	// Check if user exists
	user, err := rc.userRepo.ReadUserByEmail(c.Request.Context(), request.Email)
	if err != nil || user == nil {
		helpers.FormatResponse(c, "error", http.StatusUnauthorized, "Invalid credentials", nil, nil)
		return
	}

	// Verify password using the helper function
	if !helpers.CheckPasswordHash(request.Password, user.Password) {
		helpers.FormatResponse(c, "error", http.StatusUnauthorized, "Invalid credentials", nil, nil)
		return
	}

	// Generate both access & refresh tokens
	accessToken, refreshToken, err := rc.jwtService.GenerateTokens(user.ID.Hex(), user.Email)
	if err != nil {
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Could not generate tokens", nil, nil)
		return
	}

	// Respond with tokens
	helpers.FormatResponse(c, "success", http.StatusOK, "Request token successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// validatePayload validates the incoming request payload
func (rc *RequestTokenController) validatePayload(c *gin.Context, request any) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helpers.FormatResponse(c, "error", http.StatusBadRequest, "Invalid input", nil, nil)
		return err
	}
	return nil
}
