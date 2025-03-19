package container

import (
	"koneksi/services/iam/app/controllers/health"
	"koneksi/services/iam/app/repositories"
	"koneksi/services/iam/app/services"
	"koneksi/services/iam/database"
)

// Container holds the dependencies for the application
type Container struct {
	// Services
	MongoService *services.MongoService

	// Repositories
	PermissionRepository       *repositories.PermissionRepository
	PolicyRepository           *repositories.PolicyRepository
	PolicyPermissionRepository *repositories.PolicyPermissionRepository
	ProfileRepository          *repositories.ProfileRepository
	RoleRepository             *repositories.RoleRepository
	RolePermissionRepository   *repositories.RolePermissionRepository
	ServiceAccountRepository   *repositories.ServiceAccountRepository
	UserRepository             *repositories.UserRepository
	UserRoleRepository         *repositories.UserRoleRepository

	// Middleware

	// Controllers
	CheckHealthController *health.CheckHealthController
}

// NewContainer initializes a new IoC container
func NewContainer() *Container {
	// Initialize services
	mongoService := services.NewMongoService()

	// Initialize repositories
	permissionRepository := repositories.NewPermissionRepository(mongoService)
	policyRepository := repositories.NewPolicyRepository(mongoService)
	policyPermissionRepository := repositories.NewPolicyPermissionRepository(mongoService)
	profileRepository := repositories.NewProfileRepository(mongoService)
	roleRepository := repositories.NewRoleRepository(mongoService)
	rolePermissionRepository := repositories.NewRolePermissionRepository(mongoService)
	serviceAccountRepository := repositories.NewServiceAccountRepository(mongoService)
	userRepository := repositories.NewUserRepository(mongoService)
	userRoleRepository := repositories.NewUserRoleRepository(mongoService)

	// Run database migrations
	database.MigrateCollections(mongoService)

	// Run database seeders
	database.SeedCollections(permissionRepository, roleRepository, rolePermissionRepository)

	// Initialize middleware

	// Initialize controllers
	healthController := health.NewCheckHealthController()

	return &Container{
		MongoService:               mongoService,
		PermissionRepository:       permissionRepository,
		PolicyRepository:           policyRepository,
		PolicyPermissionRepository: policyPermissionRepository,
		ProfileRepository:          profileRepository,
		RoleRepository:             roleRepository,
		RolePermissionRepository:   rolePermissionRepository,
		ServiceAccountRepository:   serviceAccountRepository,
		UserRepository:             userRepository,
		UserRoleRepository:         userRoleRepository,
		CheckHealthController:      healthController,
	}
}
