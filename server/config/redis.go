package config

import (
	"koneksi/server/core/env"
	"time"
)

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	RedisHost               string
	RedisPort               int
	RedisPassword           string
	RedisPrefix             string
	PasswordResetCodeExpiry time.Duration
	VerificationCodeExpiry  time.Duration
}

func LoadRedisConfig() *RedisConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &RedisConfig{
		RedisHost:               envVars.RedisHost,
		RedisPort:               envVars.RedisPort,
		RedisPassword:           envVars.RedisPassword,
		RedisPrefix:             envVars.RedisPrefix,
		PasswordResetCodeExpiry: 30 * time.Minute, // Default to 30 minutes
		VerificationCodeExpiry:  24 * time.Hour,   // Default to 24 hours
	}
}
