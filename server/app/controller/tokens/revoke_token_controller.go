package tokens

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// RevokeTokenController handles revoking refresh tokens
type RevokeTokenController struct {
	tokenService *service.TokenService
}

// NewRevokeTokenController initializes a new RevokeTokenController
func NewRevokeTokenController(tokenService *service.TokenService) *RevokeTokenController {
	return &RevokeTokenController{
		tokenService: tokenService,
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

	// Revoke the token using the TokenService
	err := rc.tokenService.RevokeToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	// Return success response
	helper.FormatResponse(c, "success", http.StatusOK, "token revoked successfully", nil, nil)
}

// validatePayload validates the incoming request payload
func (rc *RevokeTokenController) validatePayload(c *gin.Context, request any) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
