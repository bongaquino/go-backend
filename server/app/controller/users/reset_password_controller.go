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
func (rpc *ResetPasswordController) Handle(c *gin.Context) {
	var request struct {
		Email              string `json:"email" binding:"required,email"`
		ResetCode          string `json:"reset_code" binding:"required"`
		NewPassword        string `json:"new_password" binding:"required,min=8"`
		ConfirmNewPassword string `json:"confirm_new_password" binding:"required,eqfield=NewPassword"`
	}

	// Validate the request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Check if new passwords match
	if request.NewPassword != request.ConfirmNewPassword {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "new passwords do not match", nil, nil)
		return
	}

	// Reset the password using the UserService
	err := rpc.userService.ResetPassword(c.Request.Context(), request.Email, request.ResetCode, request.NewPassword)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "password reset successfully", nil, nil)
}
