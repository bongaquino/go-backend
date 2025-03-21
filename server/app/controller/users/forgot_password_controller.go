package users

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

type ForgotPasswordController struct {
	userService  *service.UserService
	emailService *service.EmailService
}

func NewForgotPasswordController(userService *service.UserService, emailService *service.EmailService) *ForgotPasswordController {
	return &ForgotPasswordController{
		userService:  userService,
		emailService: emailService,
	}
}

func (fpc *ForgotPasswordController) Handle(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Validate the request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Generate a password reset code
	resetCode, err := fpc.userService.GeneratePasswordResetCode(c.Request.Context(), request.Email)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Send the reset code via email
	err = fpc.emailService.SendPasswordResetCode(request.Email, resetCode)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to send reset code", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "password reset code sent successfully", nil, nil)
}
