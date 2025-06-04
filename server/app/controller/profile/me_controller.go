package profile

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MeController handles health-related endpoints
type MeController struct {
	userService *service.UserService
}

// NewMeController initializes a new MeController
func NewMeController(userService *service.UserService) *MeController {
	return &MeController{
		userService: userService,
	}
}

// Handles the health check endpoint
func (hc *MeController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Fetch the user details
	user, profile, setting, role, limit, err := hc.userService.GetUserInfo(ctx.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Sanitize the user object by removing sensitive fields
	sanitizedUser := gin.H{
		"id":          user.ID,
		"email":       user.Email,
		"is_verified": user.IsVerified,
	}

	// Sanitize the profile object by removing sensitive fields
	sanitizedProfile := gin.H{
		"first_name": profile.FirstName,
		"last_name":  profile.LastName,
	}

	// Sanitize the role object by removing sensitive fields
	sanitizedRole := gin.H{
		"id":   role.ID,
		"name": role.Name,
	}

	// Sanitize the limit object by removing sensitive fields
	sanitizedLimit := gin.H{
		"limit": limit.BytesLimit,
		"used":  limit.BytesUsage,
	}

	// Sanitize the setting object by removing sensitive fields
	sanitizedSettings := gin.H{
		"backup_cycle":                  setting.BackupCycle,
		"backup_custom_day":             setting.BackupCustomDay,
		"notifications_frequency":       setting.NotificationsFrequency,
		"recovery_priority_order":       setting.RecoveryPriorityOrder,
		"recovery_custom_order":         setting.RecoveryCustomOrder,
		"is_mfa_enabled":                setting.IsMFAEnabled,
		"is_realtime_backup_enabled":    setting.IsRealtimeBackupEnabled,
		"is_email_notification_enabled": setting.IsEmailNotificationEnabled,
		"is_sms_notification_enabled":   setting.IsSMSNotificationEnabled,
		"is_version_history_enabled":    setting.IsVersionHistoryEnabled,
	}

	// Return the user info
	helper.FormatResponse(ctx, "success", http.StatusOK, "user info retrieved successfully", gin.H{
		"user":     sanitizedUser,
		"profile":  sanitizedProfile,
		"settings": sanitizedSettings,
		"role":     sanitizedRole,
		"limit":    sanitizedLimit,
	}, nil)
}
