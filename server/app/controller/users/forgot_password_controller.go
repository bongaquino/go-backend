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

func (fpc *ForgotPasswordController) Handle(ctx *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Validate the request payload
	if err := ctx.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return
	}

	// Generate a password reset code
	resetCode, err := fpc.userService.GeneratePasswordResetCode(ctx.Request.Context(), request.Email)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Send the reset code via email
	err = fpc.emailService.SendPasswordResetCode(request.Email, resetCode)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to send reset code", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, "password reset code sent successfully", nil, nil)
}
