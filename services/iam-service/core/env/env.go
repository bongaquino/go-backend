package env

import (
	"argo/core/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Env holds the environment variables
type Env struct {
	AppName       string `envconfig:"APP_NAME" default:"Argo"`
	AppVersion    string `envconfig:"APP_VERSION" default:"1.0.0"`
	Port          int    `envconfig:"PORT" default:"8080"`
	Mode          string `envconfig:"MODE" default:"debug"`
	MongoHost     string `envconfig:"MONGO_HOST" default:"mongo"`
	MongoPort     int    `envconfig:"MONGO_PORT" default:"27017"`
	MongoUser     string `envconfig:"MONGO_USER" default:"argo_user"`
	MongoPassword string `envconfig:"MONGO_PASSWORD" default:"argo_password"`
	MongoDatabase string `envconfig:"MONGO_DATABASE" default:"argo"`
	RedisHost     string `envconfig:"REDIS_HOST" default:"redis"`
	RedisPort     int    `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	// Add more environment variables here
}

// LoadEnv loads and validates environment variables
func LoadEnv() *Env {
	var env Env

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal("no .env file found")
	}

	// Load environment variables into the struct
	err := envconfig.Process("", &env)
	if err != nil {
		logger.Log.Fatal("failed to load environment variables: " + err.Error())
	}

	return &env
}
