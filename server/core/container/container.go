package container

import (
	"koneksi/server/app/controllers/health"
	"koneksi/server/app/controllers/tokens"
	"koneksi/server/app/controllers/users"
	"koneksi/server/app/middleware"
	"koneksi/server/app/providers"
	"koneksi/server/app/repositories"
	"koneksi/server/app/services"
	"koneksi/server/database"
)

// Container holds the dependencies for the application
type Container struct {
	// Providers
	mongoProvider *providers.MongoProvider
	RedisService  *providers.RedisProvider
	JwtService    *providers.JwtProvider

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

	// Services
	UserService *services.UserService

	// Middleware
	AuthnMiddleware    *middleware.AuthnMiddleware
	AuthzMiddleware    *middleware.AuthzMiddleware
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
	// Initialize providers
	mongoProvider := providers.NewMongoProvider()
	redisProvider := providers.NewRedisProvider()
	JwtProvider := providers.NewJwtProvider(redisProvider)

	// Initialize repositories
	permissionRepository := repositories.NewPermissionRepository(mongoProvider)
	policyRepository := repositories.NewPolicyRepository(mongoProvider)
	policyPermissionRepository := repositories.NewPolicyPermissionRepository(mongoProvider)
	profileRepository := repositories.NewProfileRepository(mongoProvider)
	roleRepository := repositories.NewRoleRepository(mongoProvider)
	rolePermissionRepository := repositories.NewRolePermissionRepository(mongoProvider)
	serviceAccountRepository := repositories.NewServiceAccountRepository(mongoProvider)
	userRepository := repositories.NewUserRepository(mongoProvider)
	userRoleRepository := repositories.NewUserRoleRepository(mongoProvider)

	// Initialize services
	userService := services.NewUserService(userRepository, profileRepository, roleRepository, userRoleRepository)

	// Run database migrations
	database.MigrateCollections(mongoProvider)

	// Run database seeders
	database.SeedCollections(permissionRepository, roleRepository, rolePermissionRepository)

	// Initialize middleware
	authnMiddleware := middleware.NewAuthnMiddleware(JwtProvider)
	authzMiddleware := middleware.NewAuthzMiddleware(userRoleRepository, roleRepository)
	verifiedMiddleware := middleware.NewVerifiedMiddleware(userRepository)

	// Initialize controllers
	checkHealthController := health.NewCheckHealthController()
	registerController := users.NewRegisterController(userService)
	requestTokenController := tokens.NewRequestTokenController(userRepository, JwtProvider)
	refreshTokenController := tokens.NewRefreshTokenController(userRepository, JwtProvider)
	revokeTokenController := tokens.NewRevokeTokenController(userRepository, JwtProvider)

	// Return the container
	return &Container{
		mongoProvider:              mongoProvider,
		RedisService:               redisProvider,
		JwtService:                 JwtProvider,
		PermissionRepository:       permissionRepository,
		PolicyRepository:           policyRepository,
		PolicyPermissionRepository: policyPermissionRepository,
		ProfileRepository:          profileRepository,
		RoleRepository:             roleRepository,
		RolePermissionRepository:   rolePermissionRepository,
		ServiceAccountRepository:   serviceAccountRepository,
		UserRepository:             userRepository,
		UserRoleRepository:         userRoleRepository,
		UserService:                userService,
		AuthnMiddleware:            authnMiddleware,
		AuthzMiddleware:            authzMiddleware,
		VerifiedMiddleware:         verifiedMiddleware,
		CheckHealthController:      checkHealthController,
		RegisterController:         registerController,
		RequestTokenController:     requestTokenController,
		RefreshTokenController:     refreshTokenController,
		RevokeTokenController:      revokeTokenController,
	}
}
