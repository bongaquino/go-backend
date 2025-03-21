package tokens

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/provider"
	"koneksi/server/app/repository"

	"github.com/gin-gonic/gin"
)

// RequestTokenController handles user authentication and token generation
type RequestTokenController struct {
	userRepo   *repository.UserRepository
	jwtService *provider.JwtProvider
}

// NewRequestTokenController initializes a new RequestTokenController
func NewRequestTokenController(userRepo *repository.UserRepository, jwtService *provider.JwtProvider) *RequestTokenController {
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
		helper.FormatResponse(c, "error", http.StatusUnauthorized, "invalid credentials", nil, nil)
		return
	}

	// Get user rolw

	// Verify password using the helper function
	if !helper.CheckPasswordHash(request.Password, user.Password) {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, "invalid credentials", nil, nil)
		return
	}

	// Generate both access & refresh tokens
	accessToken, refreshToken, err := rc.jwtService.GenerateTokens(user.ID.Hex(), &user.Email, nil)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, "could not generate tokens", nil, nil)
		return
	}

	// Respond with tokens
	helper.FormatResponse(c, "success", http.StatusOK, "request token successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// validatePayload validates the incoming request payload
func (rc *RequestTokenController) validatePayload(c *gin.Context, request any) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
