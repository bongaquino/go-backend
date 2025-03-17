package container

import (
	"koneksi/orchestrator/app/controllers/health"
	"koneksi/orchestrator/app/repositories/user"
	"koneksi/orchestrator/app/services/mongo"
	"koneksi/orchestrator/app/services/redis"
)

// Container holds the dependencies for the application
type Container struct {
	// Services
	MongoService *mongo.MongoService
	RedisService *redis.RedisService

	// Repositories
	UserRepository *user.UserRepository

	// Middleware

	// Controllers
	HealthController *health.HealthController
}

// NewContainer initializes a new IoC container
func NewContainer() *Container {
	// Initialize services
	mongoService := mongo.NewMongoService()
	redisService := redis.NewRedisService()

	// Initialize repositories
	userRepository := user.NewUserRepository(mongoService)

	// Initialize middleware

	// Initialize controllers
	healthController := health.NewHealthController()

	return &Container{
		MongoService:     mongoService,
		RedisService:     redisService,
		UserRepository:   userRepository,
		HealthController: healthController,
	}
}
