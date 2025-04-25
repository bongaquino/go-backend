package container

import (
	"koneksi/server/app/controller/health"
	"koneksi/server/app/controller/network"
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

// Providers
type Providers struct {
	Mongo    *provider.MongoProvider
	Redis    *provider.RedisProvider
	JWT      *provider.JWTProvider
	Postmark *provider.PostmarkProvider
	IPFS     *provider.IPFSProvider
}

// Repositories
type Repositories struct {
	Permission       *repository.PermissionRepository
	Policy           *repository.PolicyRepository
	PolicyPermission *repository.PolicyPermissionRepository
	Profile          *repository.ProfileRepository
	Role             *repository.RoleRepository
	RolePermission   *repository.RolePermissionRepository
	ServiceAccount   *repository.ServiceAccountRepository
	User             *repository.UserRepository
	UserRole         *repository.UserRoleRepository
}

// Services
type Services struct {
	User  *service.UserService
	Token *service.TokenService
	MFA   *service.MFAService
	Email *service.EmailService
	IPFS  *service.IPFSService
}

// Middleware
type Middleware struct {
	Authn    *middleware.AuthnMiddleware
	Authz    *middleware.AuthzMiddleware
	Verified *middleware.VerifiedMiddleware
	Locked   *middleware.LockedMiddleware
}

// Controller Groups
type HealthControllers struct {
	Check *health.CheckHealthController
}

type UserControllers struct {
	Register       *users.RegisterController
	ForgotPassword *users.ForgotPasswordController
	ResetPassword  *users.ResetPasswordController
}

type TokenControllers struct {
	Request *tokens.RequestTokenController
	Verify  *tokens.VerifyOTPController
	Refresh *tokens.RefreshTokenController
	Revoke  *tokens.RevokeTokenController
}

type SettingsControllers struct {
	ChangePassword *settings.ChangePasswordController
}

type MFAControllers struct {
	Generate *mfa.GenerateOTPController
	Enable   *mfa.EnableMFAController
	Disable  *mfa.DisableMFAController
}

type ProfileControllers struct {
	Me *profile.MeController
}

type NetworkControllers struct {
	GetSwarmAddress *network.GetSwarmAddressController
}

// Grouped Controllers
type Controllers struct {
	Health   *HealthControllers
	Users    *UserControllers
	Tokens   *TokenControllers
	Settings *SettingsControllers
	MFA      *MFAControllers
	Profile  *ProfileControllers
	Network  *NetworkControllers
}

// Container
type Container struct {
	Providers    *Providers
	Repositories *Repositories
	Services     *Services
	Middleware   *Middleware
	Controllers  *Controllers
}

// NewContainer
func NewContainer() *Container {
	// Providers
	providers := &Providers{
		Mongo:    provider.NewMongoProvider(),
		Redis:    provider.NewRedisProvider(),
		Postmark: provider.NewPostmarkProvider(),
	}
	providers.JWT = provider.NewJWTProvider(providers.Redis)
	providers.IPFS = provider.NewIPFSProvider()

	// Repositories
	repositories := &Repositories{
		Permission:       repository.NewPermissionRepository(providers.Mongo),
		Policy:           repository.NewPolicyRepository(providers.Mongo),
		PolicyPermission: repository.NewPolicyPermissionRepository(providers.Mongo),
		Profile:          repository.NewProfileRepository(providers.Mongo),
		Role:             repository.NewRoleRepository(providers.Mongo),
		RolePermission:   repository.NewRolePermissionRepository(providers.Mongo),
		ServiceAccount:   repository.NewServiceAccountRepository(providers.Mongo),
		User:             repository.NewUserRepository(providers.Mongo),
		UserRole:         repository.NewUserRoleRepository(providers.Mongo),
	}

	// Services
	services := &Services{
		User:  service.NewUserService(repositories.User, repositories.Profile, repositories.Role, repositories.UserRole, providers.Redis),
		MFA:   service.NewMFAService(repositories.User, providers.Redis),
		Email: service.NewEmailService(providers.Postmark),
		IPFS:  service.NewIPFSService(providers.IPFS),
	}
	services.Token = service.NewTokenService(repositories.User, providers.JWT, services.MFA, providers.Redis)

	// Middleware
	middlewares := &Middleware{
		Authn:    middleware.NewAuthnMiddleware(providers.JWT),
		Authz:    middleware.NewAuthzMiddleware(repositories.UserRole, repositories.Role),
		Verified: middleware.NewVerifiedMiddleware(repositories.User),
		Locked:   middleware.NewLockedMiddleware(repositories.User),
	}

	// Controllers
	controllers := &Controllers{
		Health: &HealthControllers{
			Check: health.NewCheckHealthController(),
		},
		Users: &UserControllers{
			Register:       users.NewRegisterController(services.User, services.Token),
			ForgotPassword: users.NewForgotPasswordController(services.User, services.Email),
			ResetPassword:  users.NewResetPasswordController(services.User),
		},
		Tokens: &TokenControllers{
			Request: tokens.NewRequestTokenController(services.Token, services.User, services.MFA),
			Verify:  tokens.NewVerifyOTPController(services.Token, services.MFA),
			Refresh: tokens.NewRefreshTokenController(services.Token),
			Revoke:  tokens.NewRevokeTokenController(services.Token),
		},
		Settings: &SettingsControllers{
			ChangePassword: settings.NewChangePasswordController(services.User),
		},
		MFA: &MFAControllers{
			Generate: mfa.NewGenerateOTPController(services.MFA),
			Enable:   mfa.NewEnableMFAController(services.MFA),
			Disable:  mfa.NewDisableMFAController(services.MFA, services.User),
		},
		Profile: &ProfileControllers{
			Me: profile.NewMeController(services.User),
		},
		Network: &NetworkControllers{
			GetSwarmAddress: network.NewGetSwarmAddressController(services.IPFS),
		},
	}

	// Run migrations & seeders
	database.MigrateCollections(providers.Mongo)
	database.SeedCollections(repositories.Permission, repositories.Role, repositories.RolePermission)

	return &Container{
		Providers:    providers,
		Repositories: repositories,
		Services:     services,
		Middleware:   middlewares,
		Controllers:  controllers,
	}
}
