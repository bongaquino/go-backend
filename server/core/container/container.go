package container

import (
	"koneksi/server/app/controller/health"
	"koneksi/server/app/controller/settings"
	"koneksi/server/app/controller/settings/mfa"
	"koneksi/server/app/controller/tokens"
	"koneksi/server/app/controller/users"
	"koneksi/server/app/middleware"
	"koneksi/server/app/provider"
	"koneksi/server/app/repository"
	"koneksi/server/app/service"
	"koneksi/server/database"
)

// Container holds the dependencies for the application
type Container struct {
	// Providers
	mongoProvider *provider.MongoProvider
	RedisProvider *provider.RedisProvider
	JWTProvider   *provider.JWTProvider

	// Repository
	PermissionRepository       *repository.PermissionRepository
	PolicyRepository           *repository.PolicyRepository
	PolicyPermissionRepository *repository.PolicyPermissionRepository
	ProfileRepository          *repository.ProfileRepository
	RoleRepository             *repository.RoleRepository
	RolePermissionRepository   *repository.RolePermissionRepository
	ServiceAccountRepository   *repository.ServiceAccountRepository
	UserRepository             *repository.UserRepository
	UserRoleRepository         *repository.UserRoleRepository

	// Service
	UserService  *service.UserService
	TokenService *service.TokenService
	MFAService   *service.MFAService

	// Middleware
	AuthnMiddleware    *middleware.AuthnMiddleware
	AuthzMiddleware    *middleware.AuthzMiddleware
	VerifiedMiddleware *middleware.VerifiedMiddleware

	// Controllers
	CheckHealthController    *health.CheckHealthController
	RegisterController       *users.RegisterController
	RequestTokenController   *tokens.RequestTokenController
	RefreshTokenController   *tokens.RefreshTokenController
	RevokeTokenController    *tokens.RevokeTokenController
	ChangePasswordController *settings.ChangePasswordController
	GenerateOTPController    *mfa.GenerateOTPController
	VerifyOTPController      *mfa.VerifyOTPController
}

// NewContainer initializes a new IoC container
func NewContainer() *Container {
	// Initialize provider
	mongoProvider := provider.NewMongoProvider()
	redisProvider := provider.NewRedisProvider()
	JWTProvider := provider.NewJWTProvider(redisProvider)

	// Initialize repository
	permissionRepository := repository.NewPermissionRepository(mongoProvider)
	policyRepository := repository.NewPolicyRepository(mongoProvider)
	policyPermissionRepository := repository.NewPolicyPermissionRepository(mongoProvider)
	profileRepository := repository.NewProfileRepository(mongoProvider)
	roleRepository := repository.NewRoleRepository(mongoProvider)
	rolePermissionRepository := repository.NewRolePermissionRepository(mongoProvider)
	serviceAccountRepository := repository.NewServiceAccountRepository(mongoProvider)
	userRepository := repository.NewUserRepository(mongoProvider)
	userRoleRepository := repository.NewUserRoleRepository(mongoProvider)

	// Initialize service
	userService := service.NewUserService(userRepository, profileRepository, roleRepository, userRoleRepository)
	tokenService := service.NewTokenService(userRepository, JWTProvider)
	mfaService := service.NewMFAService(userRepository)

	// Run database migrations
	database.MigrateCollections(mongoProvider)

	// Run database seeders
	database.SeedCollections(permissionRepository, roleRepository, rolePermissionRepository)

	// Initialize middleware
	authnMiddleware := middleware.NewAuthnMiddleware(JWTProvider)
	authzMiddleware := middleware.NewAuthzMiddleware(userRoleRepository, roleRepository)
	verifiedMiddleware := middleware.NewVerifiedMiddleware(userRepository)

	// Initialize controllers
	checkHealthController := health.NewCheckHealthController()
	registerController := users.NewRegisterController(userService)
	requestTokenController := tokens.NewRequestTokenController(tokenService)
	refreshTokenController := tokens.NewRefreshTokenController(tokenService)
	revokeTokenController := tokens.NewRevokeTokenController(tokenService)
	changePasswordController := settings.NewChangePasswordController(userService)
	generateOTPController := mfa.NewGenerateOTPController(mfaService)
	verifyOTPController := mfa.NewVerifyOTPController(mfaService)

	// Return the container
	return &Container{
		mongoProvider:              mongoProvider,
		RedisProvider:              redisProvider,
		JWTProvider:                JWTProvider,
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
		TokenService:               tokenService,
		MFAService:                 mfaService,
		AuthnMiddleware:            authnMiddleware,
		AuthzMiddleware:            authzMiddleware,
		VerifiedMiddleware:         verifiedMiddleware,
		CheckHealthController:      checkHealthController,
		RegisterController:         registerController,
		RequestTokenController:     requestTokenController,
		RefreshTokenController:     refreshTokenController,
		RevokeTokenController:      revokeTokenController,
		ChangePasswordController:   changePasswordController,
		GenerateOTPController:      generateOTPController,
		VerifyOTPController:        verifyOTPController,
	}
}
