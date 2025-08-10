package tokens

import (
	"net/http"

	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"

	"github.com/gin-gonic/gin"
)

// VerifyOTPController handles OTP verification and token generation
type VerifyOTPController struct {
	tokenService *service.TokenService
	mfaService   *service.MFAService
}

// NewVerifyOTPController initializes a new VerifyOTPController
func NewVerifyOTPController(tokenService *service.TokenService, mfaService *service.MFAService) *VerifyOTPController {
	return &VerifyOTPController{
		tokenService: tokenService,
		mfaService:   mfaService,
	}
}

// Handle verifies the OTP and issues tokens
func (vc *VerifyOTPController) Handle(ctx *gin.Context) {
	var request struct {
		LoginCode string `json:"login_code" binding:"required"`
		OTP       string `json:"otp" binding:"required"`
	}

	// Validate the payload
	if err := vc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Verify the OTP
	accessToken, refreshToken, err := vc.tokenService.AuthenticateLoginCode(ctx.Request.Context(), request.LoginCode, request.OTP)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "invalid login code or OTP", nil, nil)
		return
	}

	// Respond with tokens
	helper.FormatResponse(ctx, "success", http.StatusOK, "OTP verified successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// validatePayload validates the incoming request payload
func (vc *VerifyOTPController) validatePayload(ctx *gin.Context, request any) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	return nil
}
