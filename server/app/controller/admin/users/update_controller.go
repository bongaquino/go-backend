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
func (lc *UpdateController) Handle(ctx *gin.Context) {
	// Get userID from path parameters
	userID := ctx.Param("userID")
	if userID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "userID is required", nil, nil)
		return
	}

	// Get request body
	var request dto.UpdateUserDTO
	if err := lc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Update user
	err := lc.userService.UpdateUser(ctx, userID, &request)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update user", nil, nil)
		return
	}

	// Respond with success
	// helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
	// 	"user":    user,
	// 	"profile": profile,
	// }, nil)
}

func (rc *UpdateController) validatePayload(ctx *gin.Context, request *dto.UpdateUserDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	// Check if new passwords pass validation
	isValid, validationErr := helper.ValidatePassword(request.Password)
	if !isValid || validationErr != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, validationErr.Error(), nil, nil)
		return validationErr
	}
	return nil
}
