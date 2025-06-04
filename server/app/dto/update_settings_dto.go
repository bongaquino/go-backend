package dto

type UpdateSettingsDTO struct {
	BackupCustomDay            int    `json:"backup_custom_day"`
	BackupCycle                string `json:"backup_cycle"`
	NotificationsFrequency     string `json:"notifications_frequency"`
	RecoveryCustomOrder        string `json:"recovery_custom_order"`
	RecoveryPriorityOrder      string `json:"recovery_priority_order"`
	IsEmailNotificationEnabled bool   `json:"is_email_notification_enabled"`
	IsRealtimeBackupEnabled    bool   `json:"is_realtime_backup_enabled"`
	IsSMSNotificationEnabled   bool   `json:"is_sms_notification_enabled"`
	IsVersionHistoryEnabled    bool   `json:"is_version_history_enabled"`
}
