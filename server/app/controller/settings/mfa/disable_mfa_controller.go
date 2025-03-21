package mfa

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// DisableMFAController handles disabling MFA for a user
type DisableMFAController struct {
	mfaService *service.MFAService
}

// NewDisableMFAController initializes a new DisableMFAController
func NewDisableMFAController(mfaService *service.MFAService) *DisableMFAController {
	return &DisableMFAController{
		mfaService: mfaService,
	}
}

// Handle disables MFA for the user
func (dmc *DisableMFAController) Handle(c *gin.Context) {
	// Extract user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Disable MFA for the user
	err := dmc.mfaService.DisableMFA(c.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to disable MFA", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "MFA disabled successfully", nil, nil)
}
