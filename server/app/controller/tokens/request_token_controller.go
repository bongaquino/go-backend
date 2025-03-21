package tokens

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// RequestTokenController handles user authentication and token generation
type RequestTokenController struct {
	tokenService *service.TokenService
}

// NewRequestTokenController initializes a new RequestTokenController
func NewRequestTokenController(tokenService *service.TokenService) *RequestTokenController {
	return &RequestTokenController{
		tokenService: tokenService,
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

	// Authenticate user and generate tokens
	accessToken, refreshToken, err := rc.tokenService.AuthenticateUser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, err.Error(), nil, nil)
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
