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

func (rvtc *ResendVerificationCodeController) Handle(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Validate the request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Resend verification code using the UserService
	code, err := rvtc.userService.GenerateVerificationCode(c.Request.Context(), request.Email)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	// Send the verification email
	err = rvtc.emailService.SendVerificationCode(request.Email, code)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to send verification email", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "verification code resent successfully", nil, nil)
}
