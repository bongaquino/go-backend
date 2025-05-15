package config

// UserConfig holds the User configuration
type UserConfig struct {
	DefaultBytesLimit int64
}

func LoadUserConfig() *UserConfig {
	// Create the configuration from environment variables
	return &UserConfig{
		// DefaultBytesLimit is set to 2GB
		DefaultBytesLimit: 2 * 1024 * 1024 * 1024,
	}
}
