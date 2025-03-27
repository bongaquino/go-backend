package config

import (
	"fmt"

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

    return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
        config.MongoUser, config.MongoPassword, config.MongoHost, config.MongoPort, config.MongoDatabase, config.MongoDatabase)

}
