package config

import "koneksi/services/iam/core/env"

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisPrefix   string
}

func LoadRedisConfig() *RedisConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &RedisConfig{
		RedisHost:     envVars.RedisHost,
		RedisPort:     envVars.RedisPort,
		RedisPassword: envVars.RedisPassword,
		RedisPrefix:   envVars.RedisPrefix,
	}
}
