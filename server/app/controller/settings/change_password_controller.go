package settings

import (
	"fmt"
	"net/http"

	"bongaquino/server/app/dto"
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"

	"github.com/gin-gonic/gin"
)

// ChangePasswordController handles changing user passwords
type ChangePasswordController struct {
	userService *service.UserService
}

// NewChangePasswordController initializes a new ChangePasswordController
func NewChangePasswordController(userService *service.UserService) *ChangePasswordController {
	return &ChangePasswordController{
		userService: userService,
	}
}

// Handle processes the change password request
func (cpc *ChangePasswordController) Handle(ctx *gin.Context) {
	var request dto.ChangePasswordDTO

	// Validate the payload
	if err := cpc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Change the password using the UserService
	err := cpc.userService.ChangePassword(ctx.Request.Context(), userID.(string), &request)
	if err != nil {
		if err.Error() == "new password must be different from the old password" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "new password must be different from the old password", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "password change failed", nil, nil)
		return
	}

	// Return success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "password changed successfully", nil, nil)
}

// validatePayload validates the incoming request payload
func (cpc *ChangePasswordController) validatePayload(ctx *gin.Context, request *dto.ChangePasswordDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	if request.NewPassword != request.ConfirmNewPassword {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "new passwords do not match", nil, nil)
		return fmt.Errorf("new passwords do not match")
	}
	// Check if new passwords pass validation
	isValid, validationErr := helper.ValidatePassword(request.NewPassword)
	if !isValid {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, validationErr.Error(), nil, nil)
		return validationErr
	}
	return nil
}
