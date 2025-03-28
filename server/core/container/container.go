package container

import (
	"koneksi/server/app/controller/health"
	"koneksi/server/app/controller/profile"
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
	MongoProvider    *provider.MongoProvider
	RedisProvider    *provider.RedisProvider
	JWTProvider      *provider.JWTProvider
	PostmarkProvider *provider.PostmarkProvider
	IPFSProvider     *provider.IPFSProvider

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
	EmailService *service.EmailService

	// Middleware
	AuthnMiddleware    *middleware.AuthnMiddleware
	AuthzMiddleware    *middleware.AuthzMiddleware
	VerifiedMiddleware *middleware.VerifiedMiddleware

	// Controllers
	CheckHealthController    *health.CheckHealthController
	RegisterController       *users.RegisterController
	ForgotPasswordController *users.ForgotPasswordController
	ResetPasswordController  *users.ResetPasswordController
	RequestTokenController   *tokens.RequestTokenController
	VerifyOTPController      *tokens.VerifyOTPController
	RefreshTokenController   *tokens.RefreshTokenController
	RevokeTokenController    *tokens.RevokeTokenController
	ChangePasswordController *settings.ChangePasswordController
	GenerateOTPController    *mfa.GenerateOTPController
	EnableMFAController      *mfa.EnableMFAController
	DisableMFAController     *mfa.DisableMFAController
	MeController             *profile.MeController
}

// NewContainer initializes a new IoC container
func NewContainer() *Container {
	// Initialize provider
	mongoProvider := provider.NewMongoProvider()
	redisProvider := provider.NewRedisProvider()
	JWTProvider := provider.NewJWTProvider(redisProvider)
	postmarkProvider := provider.NewPostmarkProvider()
	ipfsProvider := provider.NewIPFSProvider()

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
	userService := service.NewUserService(userRepository, profileRepository,
		roleRepository, userRoleRepository, redisProvider)
	mfaService := service.NewMFAService(userRepository, redisProvider)
	tokenService := service.NewTokenService(userRepository, JWTProvider, mfaService)
	emailService := service.NewEmailService(postmarkProvider)

	// Initialize middleware
	authnMiddleware := middleware.NewAuthnMiddleware(JWTProvider)
	authzMiddleware := middleware.NewAuthzMiddleware(userRoleRepository, roleRepository)
	verifiedMiddleware := middleware.NewVerifiedMiddleware(userRepository)

	// Initialize controllers
	checkHealthController := health.NewCheckHealthController()
	registerController := users.NewRegisterController(userService, tokenService)
	forgotPasswordController := users.NewForgotPasswordController(userService, emailService)
	resetPasswordController := users.NewResetPasswordController(userService)
	requestTokenController := tokens.NewRequestTokenController(tokenService, userService, mfaService)
	verifyOTPController := tokens.NewVerifyOTPController(tokenService, mfaService)
	refreshTokenController := tokens.NewRefreshTokenController(tokenService)
	revokeTokenController := tokens.NewRevokeTokenController(tokenService)
	changePasswordController := settings.NewChangePasswordController(userService)
	generateOTPController := mfa.NewGenerateOTPController(mfaService)
	enableMFAController := mfa.NewEnableMFAController(mfaService)
	disableMFAController := mfa.NewDisableMFAController(mfaService, userService)
	meController := profile.NewMeController(userService)

	// Run database migrations
	database.MigrateCollections(mongoProvider)

	// Run database seeders
	database.SeedCollections(permissionRepository, roleRepository, rolePermissionRepository)

	// Return the container
	return &Container{
		MongoProvider:              mongoProvider,
		RedisProvider:              redisProvider,
		JWTProvider:                JWTProvider,
		PostmarkProvider:           postmarkProvider,
		IPFSProvider:               ipfsProvider,
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
		EmailService:               emailService,
		MFAService:                 mfaService,
		AuthnMiddleware:            authnMiddleware,
		AuthzMiddleware:            authzMiddleware,
		VerifiedMiddleware:         verifiedMiddleware,
		CheckHealthController:      checkHealthController,
		RegisterController:         registerController,
		ForgotPasswordController:   forgotPasswordController,
		ResetPasswordController:    resetPasswordController,
		RequestTokenController:     requestTokenController,
		VerifyOTPController:        verifyOTPController,
		RefreshTokenController:     refreshTokenController,
		RevokeTokenController:      revokeTokenController,
		ChangePasswordController:   changePasswordController,
		GenerateOTPController:      generateOTPController,
		EnableMFAController:        enableMFAController,
		DisableMFAController:       disableMFAController,
		MeController:               meController,
	}
}
