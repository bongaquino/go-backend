package tokens

import (
	"net/http"

	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"

	"github.com/gin-gonic/gin"
)

// RefreshController handles the JWT refresh process
type RefreshController struct {
	tokenService *service.TokenService
}

// NewRefreshController initializes a new RefreshController
func NewRefreshController(tokenService *service.TokenService) *RefreshController {
	return &RefreshController{
		tokenService: tokenService,
	}
}

// Handle processes the refresh token request
func (rc *RefreshController) Handle(ctx *gin.Context) {
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
func (rc *RefreshController) validatePayload(ctx *gin.Context, request any) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	return nil
}
