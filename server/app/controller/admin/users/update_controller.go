package users

import (
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateController struct {
	userService *service.UserService
}

// NewUpdateController initializes a new UpdateController
func NewUpdateController(userService *service.UserService) *UpdateController {
	return &UpdateController{
		userService: userService,
	}
}

// Handle handles the health check request
func (uc *UpdateController) Handle(ctx *gin.Context) {
	// Get userID from path parameters
	userID := ctx.Param("userID")
	if userID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "userID is required", nil, nil)
		return
	}

	// Get request body
	var request dto.UpdateUserDTO
	if err := uc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Update user
	user, profile, userRole, roleName, err := uc.userService.UpdateUser(ctx, userID, &request)
	if err != nil {
		if err.Error() == "role not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "role not found", nil, nil)
			return
		}
		if err.Error() == "user not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "user not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update user", nil, nil)
		return
	}

	// Check if userRole is nil
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"user": gin.H{
			"email": user.Email,
		},
		"profile": profile,
		"user_role": gin.H{
			"role_id":   userRole.RoleID,
			"role_name": roleName,
		},
	}, nil)
}

func (rc *UpdateController) validatePayload(ctx *gin.Context, request *dto.UpdateUserDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	// Only validate password if it's provided (not empty)
	if request.Password != "" {
		// Check if new password passes validation
		isValid, validationErr := helper.ValidatePassword(request.Password)
		if !isValid || validationErr != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, validationErr.Error(), nil, nil)
			return validationErr
		}
	}
	return nil
}
