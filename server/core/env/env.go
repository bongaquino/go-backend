package env

import (
	"fmt"
	"bongaquino/server/core/logger"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Env holds the environment variables
type Env struct {
	AppName               string `envconfig:"APP_NAME" default:"bongaquino"`
	AppVersion            string `envconfig:"APP_VERSION" default:"1.0.0"`
	AppKey                string `envconfig:"APP_KEY" required:"true"`
	Port                  int    `envconfig:"PORT" default:"3000"`
	Mode                  string `envconfig:"MODE" default:"debug"`
	MongoHost             string `envconfig:"MONGO_HOST" default:"mongo"`
	MongoPort             int    `envconfig:"MONGO_PORT" default:"27017"`
	MongoUser             string `envconfig:"MONGO_USER" default:"bongaquino_user"`
	MongoPassword         string `envconfig:"MONGO_PASSWORD" default:"bongaquino_password"`
	MongoDatabase         string `envconfig:"MONGO_DATABASE" default:"bongaquino"`
	MongoConnectionString string `envconfig:"MONGO_CONNECTION_STRING" default:""`
	RedisHost             string `envconfig:"REDIS_HOST" default:"redis"`
	RedisPort             int    `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword         string `envconfig:"REDIS_PASSWORD"`
	RedisPrefix           string `envconfig:"REDIS_PREFIX" required:"true"`
	JWTSecret             string `envconfig:"JWT_SECRET" required:"true"`
	JWTTokenExpiration    int    `envconfig:"JWT_TOKEN_EXPIRATION" default:"3600"`
	JWTRefreshExpiration  int    `envconfig:"JWT_REFRESH_EXPIRATION" default:"86400"`
	PostmarkAPIKey        string `envconfig:"POSTMARK_API_KEY" required:"true"`
	PostmarkFrom          string `envconfig:"POSTMARK_FROM" required:"true"`
	IPFSNodeURL           string `envconfig:"IPFS_NODE_URL" required:"true"`
	IPFSDownloadURL       string `envconfig:"IPFS_DOWNLOAD_URL" required:"true"`
}

// LoadEnv loads and validates environment variables
func LoadEnv() *Env {
	var env Env

	// Check MODE from environment first
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "debug" // fallback default
	}

	// Only load .env file if not in release mode
	if mode != "release" {
		if err := godotenv.Load(); err != nil {
			logger.Log.Fatal("no .env file found")
		}
	}

	// Load environment variables into the struct
	err := envconfig.Process("", &env)
	if err != nil {
		logger.Log.Info("failed to load environment variables: " + err.Error())
	}

	// Debug: Log JWT_SECRET value (masked for security)
	if env.JWTSecret != "" {
		logger.Log.Info("JWT_SECRET loaded successfully", logger.String("length", fmt.Sprintf("%d", len(env.JWTSecret))))
	} else {
		logger.Log.Error("JWT_SECRET is empty or not found")
	}

	return &env
}
