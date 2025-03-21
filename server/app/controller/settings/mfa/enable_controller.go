package mfa

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// EnableController handles OTP verification for MFA
type EnableController struct {
	mfaService *service.MFAService
}

// NewEnableController initializes a new EnableController
func NewEnableController(mfaService *service.MFAService) *EnableController {
	return &EnableController{
		mfaService: mfaService,
	}
}

// Handle verifies the OTP provided by the user
func (voc *EnableController) Handle(c *gin.Context) {
	// Extract user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Parse the OTP from the request body
	var request struct {
		OTP string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Verify the OTP
	isValid, err := voc.mfaService.VerifyOTP(c.Request.Context(), userID.(string), request.OTP)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	if !isValid {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, "invalid OTP", nil, nil)
		return
	}

	// Enable MFA for the user
	err = voc.mfaService.EnableMFA(c.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "OTP verified successfully", nil, nil)
}
