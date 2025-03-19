package container

import (
	"koneksi/services/account/app/controllers/health"
	"koneksi/services/account/app/repositories"
	"koneksi/services/account/app/services/mongo"
	"koneksi/services/account/database"
)

// Container holds the dependencies for the application
type Container struct {
	// Services
	MongoService *mongo.MongoService

	// Repositories
	UserRepository *repositories.UserRepository

	// Middleware

	// Controllers
	HealthController *health.HealthController
}

// NewContainer initializes a new IoC container
func NewContainer() *Container {
	// Initialize services
	mongoService := mongo.NewMongoService()

	// Run database migrations
	database.MigrateCollections(mongoService)

	// Run database seeders

	// Initialize repositories
	userRepository := repositories.NewUserRepository(mongoService)

	// Initialize middleware

	// Initialize controllers
	healthController := health.NewHealthController()

	return &Container{
		MongoService:     mongoService,
		UserRepository:   userRepository,
		HealthController: healthController,
	}
}
