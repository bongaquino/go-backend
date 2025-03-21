package mfa

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// VerifyOTPController handles OTP verification for MFA
type VerifyOTPController struct {
	mfaService *service.MFAService
}

// NewVerifyOTPController initializes a new VerifyOTPController
func NewVerifyOTPController(mfaService *service.MFAService) *VerifyOTPController {
	return &VerifyOTPController{
		mfaService: mfaService,
	}
}

// Handle verifies the OTP provided by the user
func (voc *VerifyOTPController) Handle(c *gin.Context) {
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

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "OTP verified successfully", nil, nil)
}
