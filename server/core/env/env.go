package env

import (
	"koneksi/server/core/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Env holds the environment variables
type Env struct {
	AppName              string `envconfig:"APP_NAME" default:"Koneksi"`
	AppVersion           string `envconfig:"APP_VERSION" default:"1.0.0"`
	Port                 int    `envconfig:"PORT" default:"3000"`
	Mode                 string `envconfig:"MODE" default:"debug"`
	MongoHost            string `envconfig:"MONGO_HOST" default:"mongo"`
	MongoPort            int    `envconfig:"MONGO_PORT" default:"27017"`
	MongoUser            string `envconfig:"MONGO_USER" default:"koneksi_user"`
	MongoPassword        string `envconfig:"MONGO_PASSWORD" default:"koneksi_password"`
	MongoDatabase        string `envconfig:"MONGO_DATABASE" default:"koneksi"`
	RedisHost            string `envconfig:"REDIS_HOST" default:"redis"`
	RedisPort            int    `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword        string `envconfig:"REDIS_PASSWORD" default:""`
	RedisPrefix          string `envconfig:"REDIS_PREFIX" default:"koneksi"`
	JWTSecret            string `envconfig:"JWT_SECRET" default:""`
	JWTTokenExpiration   int    `envconfig:"JWT_TOKEN_EXPIRATION" default:"3600"`
	JWTRefreshExpiration int    `envconfig:"JWT_REFRESH_EXPIRATION" default:"86400"`
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
