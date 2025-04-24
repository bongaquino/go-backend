// server/app/controller/users/resend_verification_token_controller.go
package users

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// ResendVerificationTokenController handles resending verification tokens
type ResendVerificationTokenController struct {
	userService *service.UserService
	emailService  *service.EmailService
}

// NewResendVerificationTokenController initializes a new ResendVerificationTokenController
func NewResendVerificationTokenController(userService *service.UserService, emailService *service.EmailService) *ResendVerificationTokenController {
	return &ResendVerificationTokenController{
		userService: userService,
		emailService: emailService,
	}
}

func (rvtc *ResendVerificationTokenController) Handle(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Validate the request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Resend verification token using the UserService
	token, err := rvtc.userService.GenerateVerificationToken(c.Request.Context(), request.Email)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	// Send the verification email
	err = rvtc.emailService.SendVerificationCode(request.Email, token)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to send verification email", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "verification token resent successfully", nil, nil)
}