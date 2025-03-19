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
	PermissionRepository       *repositories.PermissionRepository
	PolicyRepository           *repositories.PolicyRepository
	PolicyPermissionRepository *repositories.PolicyPermissionRepository
	ProfileRepository          *repositories.ProfileRepository
	RoleRepository             *repositories.RoleRepository
	RolePermissionRepository   *repositories.RolePermissionRepository
	UserRepository             *repositories.UserRepository

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
	permissionRepository := repositories.NewPermissionRepository(mongoService)
	policyRepository := repositories.NewPolicyRepository(mongoService)
	policyPermissionRepository := repositories.NewPolicyPermissionRepository(mongoService)
	profileRepository := repositories.NewProfileRepository(mongoService)
	roleRepository := repositories.NewRoleRepository(mongoService)
	rolePermissionRepository := repositories.NewRolePermissionRepository(mongoService)
	userRepository := repositories.NewUserRepository(mongoService)

	// Initialize middleware

	// Initialize controllers
	healthController := health.NewHealthController()

	return &Container{
		MongoService:               mongoService,
		PermissionRepository:       permissionRepository,
		PolicyRepository:           policyRepository,
		PolicyPermissionRepository: policyPermissionRepository,
		ProfileRepository:          profileRepository,
		RoleRepository:             roleRepository,
		RolePermissionRepository:   rolePermissionRepository,
		UserRepository:             userRepository,
		HealthController:           healthController,
	}
}
