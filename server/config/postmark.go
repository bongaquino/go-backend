package config

import "koneksi/server/core/env"

// PostmarkConfig holds the Postmark configuration
type PostmarkConfig struct {
	PostmarkAPIKey string
	PostmarkFrom   string
}

func LoadPostmarkConfig() *PostmarkConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &PostmarkConfig{
		PostmarkAPIKey: envVars.PostmarkAPIKey,
		PostmarkFrom:   envVars.PostmarkFrom,
	}
}
