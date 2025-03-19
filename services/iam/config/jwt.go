package config

import "koneksi/services/iam/core/env"

// JwtConfig holds the Jwt configuration
type JwtConfig struct {
	JwtSecret     string
	JwtExpiration int
}

func LoadJwtConfig() *JwtConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &JwtConfig{
		JwtSecret:     envVars.JwtSecret,
		JwtExpiration: envVars.JwtExpiration,
	}
}
