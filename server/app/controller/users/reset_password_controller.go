package users

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// ResetPasswordController handles resetting user passwords
type ResetPasswordController struct {
	userService *service.UserService
}

// NewResetPasswordController initializes a new ResetPasswordController
func NewResetPasswordController(userService *service.UserService) *ResetPasswordController {
	return &ResetPasswordController{
		userService: userService,
	}
}

// Handle processes the reset password request
func (rpc *ResetPasswordController) Handle(ctx *gin.Context) {
	var request struct {
		Email              string `json:"email" binding:"required,email"`
		ResetCode          string `json:"reset_code" binding:"required"`
		NewPassword        string `json:"new_password" binding:"required,min=8"`
		ConfirmNewPassword string `json:"confirm_new_password" binding:"required,eqfield=NewPassword"`
	}

	// Validate the request payload
	if err := ctx.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return
	}

	// Check if new passwords match
	if request.NewPassword != request.ConfirmNewPassword {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "new passwords do not match", nil, nil)
		return
	}

	// Check if new passwords pass validation
	isValid, validationErr := helper.ValidatePassword(request.NewPassword)
	if !isValid || validationErr != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, validationErr.Error(), nil, nil)
		return
	}

	// Reset the password using the UserService
	err := rpc.userService.ResetPassword(ctx.Request.Context(), request.Email, request.ResetCode, request.NewPassword)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, "password reset successfully", nil, nil)
}
