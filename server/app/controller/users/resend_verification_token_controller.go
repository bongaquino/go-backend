// server/app/controller/users/resend_verification_token_controller.go
package users

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// ResendVerificationCodeController handles resending verification tokens
type ResendVerificationCodeController struct {
	userService  *service.UserService
	emailService *service.EmailService
}

// NewResendVerificationCodeController initializes a new ResendVerificationCodeController
func NewResendVerificationCodeController(userService *service.UserService, emailService *service.EmailService) *ResendVerificationCodeController {
	return &ResendVerificationCodeController{
		userService:  userService,
		emailService: emailService,
	}
}

func (rvtc *ResendVerificationCodeController) Handle(ctx *gin.Context) {
	// Extract userID from the user token
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "unauthorized", nil, nil)
		return
	}

	// Get user email from the UserService
	user, _, _, _, _, _ := rvtc.userService.GetUserInfo(ctx.Request.Context(), userID.(string))

	// Resend verification code using the UserService
	code, err := rvtc.userService.GenerateVerificationCode(ctx.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	// Send the verification email
	err = rvtc.emailService.SendVerificationCode(user.Email, code)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to send verification email", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, "verification code resent successfully", nil, nil)
}
