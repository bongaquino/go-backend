package config

// FileConfig holds the File configuration
type FileConfig struct {
	DefaultAccess string
	AccessOptions []string
}

func LoadFileConfig() *FileConfig {
	// Create the configuration from environment variables
	return &FileConfig{
		// DefaultAccess is set to "private"
		DefaultAccess: "private",
		AccessOptions: []string{
			"private",
			"public",
			"password",
			"email",
		},
	}
}
