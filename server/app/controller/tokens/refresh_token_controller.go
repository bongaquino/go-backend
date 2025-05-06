package tokens

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// RefreshTokenController handles the JWT refresh process
type RefreshTokenController struct {
	tokenService *service.TokenService
}

// NewRefreshTokenController initializes a new RefreshTokenController
func NewRefreshTokenController(tokenService *service.TokenService) *RefreshTokenController {
	return &RefreshTokenController{
		tokenService: tokenService,
	}
}

// Handle processes the refresh token request
func (rc *RefreshTokenController) Handle(ctx *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Validate the payload
	if err := rc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Refresh tokens using the TokenService
	accessToken, refreshToken, err := rc.tokenService.RefreshTokens(ctx.Request.Context(), request.RefreshToken)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	// Return the new tokens
	helper.FormatResponse(ctx, "success", http.StatusOK, "token refreshed successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// validatePayload validates the incoming request payload
func (rc *RefreshTokenController) validatePayload(ctx *gin.Context, request any) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
