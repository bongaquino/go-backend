package config

import "bongaquino/server/core/env"

// AppConfig holds the application configuration
type AppConfig struct {
	AppName    string
	AppVersion string
	AppKey     string
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
		AppKey:     envVars.AppKey,
		Mode:       envVars.Mode,
		Port:       envVars.Port,
	}
}
