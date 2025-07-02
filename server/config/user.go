package config

// UserConfig holds the User configuration
type UserConfig struct {
	DefaultBytesLimit                 int64
	DefaultBackupCycle                string
	DefaultBackupCustomDay            int
	DefaultNotificationsFrequency     string
	DefaultRecoveryPriorityOrder      string
	DefaultRecoveryCustomOrder        string
	DefaultIsMFAEnabled               bool
	DefaultIsRealtimeBackupEnabled    bool
	DefaultIsEmailNotificationEnabled bool
	DefaultIsSMSNotificationEnabled   bool
	DefaultIsVersionHistoryEnabled    bool
	BackupCycleOptions                []string
	NotificationsFrequencyOptions     []string
	RecoveryPriorityOrderOptions      []string
	RecoveryCustomOrderOptions        []string
}

func LoadUserConfig() *UserConfig {
	// Create the configuration from environment variables
	return &UserConfig{
		// DefaultBytesLimit is set to 5GB
		DefaultBytesLimit: 5 * 1024 * 1024 * 1024,

		// DefaultBackupCycle is set to "daily"
		DefaultBackupCycle: "daily",

		// DefaultBackupCustomDay is set to 1 (Monday)
		DefaultBackupCustomDay: 1,

		// DefaultNotificationsFrequency is set to "daily"
		DefaultNotificationsFrequency: "immediately",

		// DefaultRecoveryPriorityOrder is set to "default"
		DefaultRecoveryPriorityOrder: "default",

		// DefaultRecoveryCustomOrder is set to "general"
		DefaultRecoveryCustomOrder: "general",

		// DefaultIsMFAEnabled is set to false
		DefaultIsMFAEnabled: false,

		// DefaultIsRealtimeBackupEnabled is set to false
		DefaultIsRealtimeBackupEnabled: false,

		// DefaultIsEmailNotificationEnabled is set to false
		DefaultIsEmailNotificationEnabled: false,

		// DefaultIsSMSNotificationEnabled is set to false
		DefaultIsSMSNotificationEnabled: false,

		// DefaultIsVersionHistoryEnabled is set to false
		DefaultIsVersionHistoryEnabled: false,

		BackupCycleOptions: []string{
			"daily",
			"weekly",
			"monthly",
			"custom",
		},

		NotificationsFrequencyOptions: []string{
			"immediately",
			"daily",
			"weekly",
			"monthly",
			"never",
		},

		RecoveryPriorityOrderOptions: []string{
			"default",
			"custom",
		},

		RecoveryCustomOrderOptions: []string{
			"general",
			"critical",
		},
	}
}
