package mfa

import (
	"net/http"

	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"

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
func (goc *GenerateOTPController) Handle(ctx *gin.Context) {
	// Extract email from the user token
	email, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "unauthorized", nil, nil)
		return
	}

	// Generate OTP secret and QR code
	otpSecret, qrCode, err := goc.mfaService.GenerateOTP(ctx.Request.Context(), email.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with the OTP secret and QR code
	helper.FormatResponse(ctx, "success", http.StatusOK, "OTP generated successfully", gin.H{
		"otp_secret": otpSecret,
		"qr_code":    qrCode,
	}, nil)
}
