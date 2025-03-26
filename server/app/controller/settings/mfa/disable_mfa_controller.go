package mfa

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// DisableMFAController handles disabling MFA for a user
type DisableMFAController struct {
	mfaService  *service.MFAService
	userService *service.UserService
}

// NewDisableMFAController initializes a new DisableMFAController
func NewDisableMFAController(mfaService *service.MFAService, userService *service.UserService) *DisableMFAController {
	return &DisableMFAController{
		mfaService:  mfaService,
		userService: userService,
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

	// Parse the password from the request body
	var request struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Validate the password
	isValid, err := dmc.userService.ValidatePassword(c.Request.Context(), userID.(string), request.Password)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to validate password", nil, nil)
		return
	}
	if !isValid {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, "invalid password", nil, nil)
		return
	}

	// Disable MFA for the user
	err = dmc.mfaService.DisableMFA(c.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to disable MFA", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "MFA disabled successfully", nil, nil)
}
