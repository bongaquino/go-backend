package settings

import (
	"fmt"
	"net/http"

	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"koneksi/server/config"

	"github.com/gin-gonic/gin"
)

// UpdateController handles changing user passwords
type UpdateController struct {
	userService *service.UserService
}

// NewUpdateController initializes a new UpdateController
func NewUpdateController(userService *service.UserService) *UpdateController {
	return &UpdateController{
		userService: userService,
	}
}

func (uc *UpdateController) Handle(ctx *gin.Context) {
	var request dto.UpdateSettingsDTO

	// Validate the payload
	if err := uc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Update the settings using the UserService
	err := uc.userService.UpdateUserSettings(ctx.Request.Context(), userID.(string), &request)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "settings update failed", nil, nil)
	}

	// Return success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "settings updated successfully", nil, nil)
}

// validatePayload validates the incoming request payload
func (uc *UpdateController) validatePayload(ctx *gin.Context, request *dto.UpdateSettingsDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	// Load user configuration
	userConfig := config.LoadUserConfig()

	if request.BackupCycle != "" && !helper.Contains(userConfig.BackupCycleOptions, request.BackupCycle) {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid backup cycle", nil, nil)
		return fmt.Errorf("invalid backup cycle: %s", request.BackupCycle)
	}
	if request.NotificationsFrequency != "" && !helper.Contains(userConfig.NotificationsFrequencyOptions, request.NotificationsFrequency) {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid notifications frequency", nil, nil)
		return fmt.Errorf("invalid notifications frequency: %s", request.NotificationsFrequency)
	}
	if request.RecoveryPriorityOrder != "" && !helper.Contains(userConfig.RecoveryPriorityOrderOptions, request.RecoveryPriorityOrder) {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid recovery priority order", nil, nil)
		return fmt.Errorf("invalid recovery priority order: %s", request.RecoveryPriorityOrder)
	}
	if request.RecoveryCustomOrder != "" && !helper.Contains(userConfig.RecoveryCustomOrderOptions, request.RecoveryCustomOrder) {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid recovery custom order", nil, nil)
		return fmt.Errorf("invalid recovery custom order: %s", request.RecoveryCustomOrder)
	}
	if request.BackupCustomDay < 1 || request.BackupCustomDay > 31 {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "backup custom day must be between 1 and 31", nil, nil)
		return fmt.Errorf("backup custom day must be between 1 and 31: %d", request.BackupCustomDay)
	}
	return nil
}
