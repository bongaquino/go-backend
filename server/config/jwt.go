package config

import "koneksi/server/core/env"

// JwtConfig holds the Jwt configuration
type JwtConfig struct {
	JwtSecret            string
	JwtTokenExpiration   int
	JwtRefreshExpiration int
}

func LoadJwtConfig() *JwtConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &JwtConfig{
		JwtSecret:            envVars.JwtSecret,
		JwtTokenExpiration:   envVars.JwtTokenExpiration,
		JwtRefreshExpiration: envVars.JwtRefreshExpiration,
	}
}
