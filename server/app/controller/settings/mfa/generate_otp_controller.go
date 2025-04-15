package mfa

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// GenerateOTPController handles generating OTP secrets for MFA
type GenerateOTPController struct {
	mfaService *service.MFAService
}

// NewGenerateOTPController initializes a new GenerateOTPController
func NewGenerateOTPController(mfaService *service.MFAService) *GenerateOTPController {
	return &GenerateOTPController{
		mfaService: mfaService,
	}
}

// Handle generates an OTP secret and QR code for the user
func (goc *GenerateOTPController) Handle(c *gin.Context) {
	// Extract user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		helper.FormatResponse(c, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Generate OTP secret and QR code
	otpSecret, qrCode, err := goc.mfaService.GenerateOTP(c.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with the OTP secret and QR code
	helper.FormatResponse(c, "success", http.StatusOK, "OTP generated successfully", gin.H{
		"otp_secret": otpSecret,
		"qr_code":    qrCode,
	}, nil)
}
