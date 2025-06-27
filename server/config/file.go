package config

// FileConfig holds the File configuration
type FileConfig struct {
	DefaultAccess   string
	AccessOptions   []string
	PrivateAccess   string
	PublicAccess    string
	TemporaryAccess string
	PasswordAccess  string
	EmailAccess     string
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
		PrivateAccess:  "private",
		PublicAccess:   "public",
		PasswordAccess: "password",
		EmailAccess:    "email",
	}
}
