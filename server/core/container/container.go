package container

import (
	"koneksi/server/app/controllers/health"
	"koneksi/server/app/controllers/tokens"
	"koneksi/server/app/controllers/users"
	"koneksi/server/app/middleware"
	"koneksi/server/app/repositories"
	"koneksi/server/app/services"
	"koneksi/server/database"
)

// Container holds the dependencies for the application
type Container struct {
	// Services
	MongoService *services.MongoService
	RedisService *services.RedisService
	JwtService   *services.JWTService

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
	AuthnMiddleware    *middleware.AuthnMiddleware
	VerifiedMiddleware *middleware.VerifiedMiddleware

	// Controllers
	CheckHealthController  *health.CheckHealthController
	RegisterController     *users.RegisterController
	RequestTokenController *tokens.RequestTokenController
	RefreshTokenController *tokens.RefreshTokenController
	RevokeTokenController  *tokens.RevokeTokenController
}

// NewContainer initializes a new IoC container
func NewContainer() *Container {
	// Initialize services
	mongoService := services.NewMongoService()
	redisService := services.NewRedisService()
	jwtService := services.NewJWTService(redisService)

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
	authnMiddleware := middleware.NewAuthnMiddleware(jwtService)
	verifiedMiddleware := middleware.NewVerifiedMiddleware(userRepository)

	// Initialize controllers
	checkHealthController := health.NewCheckHealthController()
	registerController := users.NewRegisterController(userRepository, profileRepository, roleRepository, userRoleRepository)
	requestTokenController := tokens.NewRequestTokenController(userRepository, jwtService)
	refreshTokenController := tokens.NewRefreshTokenController(userRepository, jwtService)
	revokeTokenController := tokens.NewRevokeTokenController(userRepository, jwtService)

	// Return the container
	return &Container{
		MongoService:               mongoService,
		RedisService:               redisService,
		JwtService:                 jwtService,
		PermissionRepository:       permissionRepository,
		PolicyRepository:           policyRepository,
		PolicyPermissionRepository: policyPermissionRepository,
		ProfileRepository:          profileRepository,
		RoleRepository:             roleRepository,
		RolePermissionRepository:   rolePermissionRepository,
		ServiceAccountRepository:   serviceAccountRepository,
		UserRepository:             userRepository,
		UserRoleRepository:         userRoleRepository,
		AuthnMiddleware:            authnMiddleware,
		VerifiedMiddleware:         verifiedMiddleware,
		CheckHealthController:      checkHealthController,
		RegisterController:         registerController,
		RequestTokenController:     requestTokenController,
		RefreshTokenController:     refreshTokenController,
		RevokeTokenController:      revokeTokenController,
	}
}
