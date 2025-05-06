package mfa

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// EnableMFAController handles OTP verification for MFA
type EnableMFAController struct {
	mfaService *service.MFAService
}

// NewEnableMFAController initializes a new EnableMFAController
func NewEnableMFAController(mfaService *service.MFAService) *EnableMFAController {
	return &EnableMFAController{
		mfaService: mfaService,
	}
}

// Handle verifies the OTP provided by the user
func (voc *EnableMFAController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Parse the OTP from the request body
	var request struct {
		OTP string `json:"otp" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Verify the OTP
	isValid, err := voc.mfaService.VerifyOTP(ctx.Request.Context(), userID.(string), request.OTP)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}
	if !isValid {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "invalid OTP", nil, nil)
		return
	}

	// Enable MFA for the user
	err = voc.mfaService.EnableMFA(ctx.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, "MFA enabled successfully", nil, nil)
}
