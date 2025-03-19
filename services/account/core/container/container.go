package container

import (
	"koneksi/services/iam/app/controllers/health"
	"koneksi/services/iam/app/repositories"
	"koneksi/services/iam/app/services/mongo"
	"koneksi/services/iam/database"
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
	ServiceAccountRepository   *repositories.ServiceAccountRepository
	UserRepository             *repositories.UserRepository
	UserRoleRepository         *repositories.UserRoleRepository

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
	serviceAccountRepository := repositories.NewServiceAccountRepository(mongoService)
	userRepository := repositories.NewUserRepository(mongoService)
	userRoleRepository := repositories.NewUserRoleRepository(mongoService)

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
		ServiceAccountRepository:   serviceAccountRepository,
		UserRepository:             userRepository,
		UserRoleRepository:         userRoleRepository,
		HealthController:           healthController,
	}
}
