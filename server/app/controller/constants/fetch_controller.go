package constants

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"koneksi/server/config"

	"github.com/gin-gonic/gin"
)

// FetchController handles health-related endpoints
type FetchController struct {
	userService *service.UserService
	orgService  *service.OrganizationService
}

// NewFetchController initializes a new FetchController
func NewFetchController(userService *service.UserService, orgService *service.OrganizationService) *FetchController {
	return &FetchController{
		userService: userService,
		orgService:  orgService,
	}
}

// Handles the health check endpoint
func (fc *FetchController) Handle(ctx *gin.Context) {

	roles, err := fc.userService.ListRoles(ctx)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	policies, err := fc.orgService.ListPolicies(ctx)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	permissions, err := fc.orgService.ListPermissions(ctx)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	userConfig := config.LoadUserConfig()
	settings := gin.H{
		"backup_cycle_options":                  userConfig.BackupCycleOptions,
		"notifications_frequency_options":       userConfig.NotificationsFrequencyOptions,
		"recovery_priority_order_options":       userConfig.RecoveryPriorityOrderOptions,
		"recovery_custom_order_options":         userConfig.RecoveryCustomOrderOptions,
		"default_bytes_limit":                   userConfig.DefaultBytesLimit,
		"default_backup_cycle":                  userConfig.DefaultBackupCycle,
		"default_backup_custom_day":             userConfig.DefaultBackupCustomDay,
		"default_notifications_frequency":       userConfig.DefaultNotificationsFrequency,
		"default_recovery_priority_order":       userConfig.DefaultRecoveryPriorityOrder,
		"default_recovery_custom_order":         userConfig.DefaultRecoveryCustomOrder,
		"default_is_mfa_enabled":                userConfig.DefaultIsMFAEnabled,
		"default_is_realtime_backup_enabled":    userConfig.DefaultIsRealtimeBackupEnabled,
		"default_is_email_notification_enabled": userConfig.DefaultIsEmailNotificationEnabled,
		"default_is_sms_notification_enabled":   userConfig.DefaultIsSMSNotificationEnabled,
		"default_is_version_history_enabled":    userConfig.DefaultIsVersionHistoryEnabled,
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"roles":       roles,
		"policies":    policies,
		"permissions": permissions,
		"settings":    settings,
	}, nil)
}
