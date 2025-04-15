package config

import "koneksi/server/core/env"

// JWTConfig holds the JWT configuration
type JWTConfig struct {
	JWTSecret            string
	JWTTokenExpiration   int
	JWTRefreshExpiration int
}

func LoadJWTConfig() *JWTConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &JWTConfig{
		JWTSecret:            envVars.JWTSecret,
		JWTTokenExpiration:   envVars.JWTTokenExpiration,
		JWTRefreshExpiration: envVars.JWTRefreshExpiration,
	}
}
