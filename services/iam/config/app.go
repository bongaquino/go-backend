package config

import "koneksi/services/iam/core/env"

// AppConfig holds the application configuration
type AppConfig struct {
	AppName    string
	AppVersion string
	Mode       string
	Port       int
}

func LoadAppConfig() *AppConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &AppConfig{
		AppName:    envVars.AppName,
		AppVersion: envVars.AppVersion,
		Mode:       envVars.Mode,
		Port:       envVars.Port,
	}
}
