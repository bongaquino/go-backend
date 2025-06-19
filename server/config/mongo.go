package config

import (
	"fmt"
	"net/url"
	"koneksi/server/core/env"
)

// MongoConfig holds the MongoDB configuration
type MongoConfig struct {
	MongoHost     string
	MongoPort     int
	MongoUser     string
	MongoPassword string
	MongoDatabase string
	MongoConnectionString string
}

func LoadMongoConfig() *MongoConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &MongoConfig{
		MongoHost:     envVars.MongoHost,
		MongoPort:     envVars.MongoPort,
		MongoUser:     envVars.MongoUser,
		MongoPassword: envVars.MongoPassword,
		MongoDatabase: envVars.MongoDatabase,
		MongoConnectionString: envVars.MongoConnectionString,
	}
}

func (config *MongoConfig) GetMongoUri() string {
    if config.MongoConnectionString != "" {
        return config.MongoConnectionString
    }

    // URL-encode username and password to handle special characters
    encodedUser := url.QueryEscape(config.MongoUser)
    encodedPassword := url.QueryEscape(config.MongoPassword)

    return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
        encodedUser, encodedPassword, config.MongoHost, config.MongoPort, config.MongoDatabase, config.MongoDatabase)

}
