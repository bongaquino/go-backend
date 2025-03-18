package container

import (
	"koneksi/services/account/app/controllers/health"
	"koneksi/services/account/app/repositories/user"
	"koneksi/services/account/app/services/mongo"
)

// Container holds the dependencies for the application
type Container struct {
	// Services
	MongoService *mongo.MongoService

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

	// Initialize repositories
	userRepository := user.NewUserRepository(mongoService)

	// Initialize middleware

	// Initialize controllers
	healthController := health.NewHealthController()

	return &Container{
		MongoService:     mongoService,
		UserRepository:   userRepository,
		HealthController: healthController,
	}
}
