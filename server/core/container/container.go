package container

import (
	"koneksi/server/app/controller/admin/organizations"
	adminUsers "koneksi/server/app/controller/admin/users"
	"koneksi/server/app/controller/constants"
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

type Providers struct {
	Mongo    *provider.MongoProvider
	Redis    *provider.RedisProvider
	JWT      *provider.JWTProvider
	Postmark *provider.PostmarkProvider
	IPFS     *provider.IPFSProvider
}

type Repositories struct {
	Permission           *repository.PermissionRepository
	Policy               *repository.PolicyRepository
	PolicyPermission     *repository.PolicyPermissionRepository
	Profile              *repository.ProfileRepository
	Role                 *repository.RoleRepository
	RolePermission       *repository.RolePermissionRepository
	ServiceAccount       *repository.ServiceAccountRepository
	User                 *repository.UserRepository
	UserRole             *repository.UserRoleRepository
	Organization         *repository.OrganizationRepository
	OrganizationUserRole *repository.OrganizationUserRoleRepository
}

type Services struct {
	User         *service.UserService
	Token        *service.TokenService
	MFA          *service.MFAService
	Email        *service.EmailService
	IPFS         *service.IPFSService
	Organization *service.OrganizationService
}

type Middleware struct {
	Authn    *middleware.AuthnMiddleware
	Authz    *middleware.AuthzMiddleware
	Verified *middleware.VerifiedMiddleware
	Locked   *middleware.LockedMiddleware
}

type Controllers struct {
	Health struct {
		Check *health.CheckController
	}
	Constants struct {
		Fetch *constants.FetchController
	}
	Users struct {
		Register               *users.RegisterController
		ForgotPassword         *users.ForgotPasswordController
		ResetPassword          *users.ResetPasswordController
		VerifyAccount          *users.VerifyAccountController
		ResendVerificationCode *users.ResendVerificationCodeController
	}
	Tokens struct {
		Request *tokens.RequestTokenController
		Verify  *tokens.VerifyOTPController
		Refresh *tokens.RefreshTokenController
		Revoke  *tokens.RevokeTokenController
	}
	Settings struct {
		ChangePassword *settings.ChangePasswordController
		MFA            struct {
			Generate *mfa.GenerateOTPController
			Enable   *mfa.EnableMFAController
			Disable  *mfa.DisableMFAController
		}
	}
	Profile struct {
		Me *profile.MeController
	}
	Network struct {
		GetSwarmAddress *network.GetSwarmAddressController
	}
	Admin struct {
		Users struct {
			List   *adminUsers.ListController
			Create *adminUsers.CreateController
			Read   *adminUsers.ReadController
			Update *adminUsers.UpdateController
		}
		Organizations struct {
			List *organizations.ListController
		}
	}
}

type Container struct {
	Providers    Providers
	Repositories Repositories
	Services     Services
	Middleware   Middleware
	Controllers  Controllers
}

func initProviders() Providers {
	mongo := provider.NewMongoProvider()
	redis := provider.NewRedisProvider()
	postmark := provider.NewPostmarkProvider()
	jwt := provider.NewJWTProvider(redis)
	ipfs := provider.NewIPFSProvider()
	return Providers{mongo, redis, jwt, postmark, ipfs}
}

func initRepositories(p Providers) Repositories {
	return Repositories{
		Permission:           repository.NewPermissionRepository(p.Mongo),
		Policy:               repository.NewPolicyRepository(p.Mongo),
		PolicyPermission:     repository.NewPolicyPermissionRepository(p.Mongo),
		Profile:              repository.NewProfileRepository(p.Mongo),
		Role:                 repository.NewRoleRepository(p.Mongo),
		RolePermission:       repository.NewRolePermissionRepository(p.Mongo),
		ServiceAccount:       repository.NewServiceAccountRepository(p.Mongo),
		User:                 repository.NewUserRepository(p.Mongo),
		UserRole:             repository.NewUserRoleRepository(p.Mongo),
		Organization:         repository.NewOrganizationRepository(p.Mongo),
		OrganizationUserRole: repository.NewOrganizationUserRoleRepository(p.Mongo),
	}
}

func initServices(p Providers, r Repositories) Services {
	user := service.NewUserService(r.User, r.Profile, r.Role, r.UserRole, p.Redis)
	email := service.NewEmailService(p.Postmark)
	mfa := service.NewMFAService(r.User, p.Redis)
	ipfs := service.NewIPFSService(p.IPFS)
	token := service.NewTokenService(r.User, p.JWT, mfa, p.Redis)
	organization := service.NewOrganizationService(r.Organization)
	return Services{user, token, mfa, email, ipfs, organization}
}

func initMiddleware(p Providers, r Repositories) Middleware {
	return Middleware{
		Authn:    middleware.NewAuthnMiddleware(p.JWT),
		Authz:    middleware.NewAuthzMiddleware(r.UserRole, r.Role),
		Verified: middleware.NewVerifiedMiddleware(r.User),
		Locked:   middleware.NewLockedMiddleware(r.User),
	}
}

func initControllers(s Services) Controllers {
	return Controllers{
		Health: struct {
			Check *health.CheckController
		}{
			Check: health.NewCheckController(),
		},
		Constants: struct {
			Fetch *constants.FetchController
		}{
			Fetch: constants.NewFetchController(s.User),
		},
		Users: struct {
			Register               *users.RegisterController
			ForgotPassword         *users.ForgotPasswordController
			ResetPassword          *users.ResetPasswordController
			VerifyAccount          *users.VerifyAccountController
			ResendVerificationCode *users.ResendVerificationCodeController
		}{
			Register:               users.NewRegisterController(s.User, s.Token, s.Email),
			ForgotPassword:         users.NewForgotPasswordController(s.User, s.Email),
			ResetPassword:          users.NewResetPasswordController(s.User),
			VerifyAccount:          users.NewVerifyAccountController(s.User),
			ResendVerificationCode: users.NewResendVerificationCodeController(s.User, s.Email),
		},
		Tokens: struct {
			Request *tokens.RequestTokenController
			Verify  *tokens.VerifyOTPController
			Refresh *tokens.RefreshTokenController
			Revoke  *tokens.RevokeTokenController
		}{
			Request: tokens.NewRequestTokenController(s.Token, s.User, s.MFA),
			Verify:  tokens.NewVerifyOTPController(s.Token, s.MFA),
			Refresh: tokens.NewRefreshTokenController(s.Token),
			Revoke:  tokens.NewRevokeTokenController(s.Token),
		},
		Settings: struct {
			ChangePassword *settings.ChangePasswordController
			MFA            struct {
				Generate *mfa.GenerateOTPController
				Enable   *mfa.EnableMFAController
				Disable  *mfa.DisableMFAController
			}
		}{
			ChangePassword: settings.NewChangePasswordController(s.User),
			MFA: struct {
				Generate *mfa.GenerateOTPController
				Enable   *mfa.EnableMFAController
				Disable  *mfa.DisableMFAController
			}{
				Generate: mfa.NewGenerateOTPController(s.MFA),
				Enable:   mfa.NewEnableMFAController(s.MFA),
				Disable:  mfa.NewDisableMFAController(s.MFA, s.User),
			},
		},
		Profile: struct {
			Me *profile.MeController
		}{
			Me: profile.NewMeController(s.User),
		},
		Network: struct {
			GetSwarmAddress *network.GetSwarmAddressController
		}{
			GetSwarmAddress: network.NewGetSwarmAddressController(s.IPFS),
		},
		Admin: struct {
			Users struct {
				List   *adminUsers.ListController
				Create *adminUsers.CreateController
				Read   *adminUsers.ReadController
				Update *adminUsers.UpdateController
			}
			Organizations struct {
				List *organizations.ListController
			}
		}{
			Users: struct {
				List   *adminUsers.ListController
				Create *adminUsers.CreateController
				Read   *adminUsers.ReadController
				Update *adminUsers.UpdateController
			}{
				List:   adminUsers.NewListController(s.User),
				Create: adminUsers.NewCreateController(s.User, s.Token, s.Email),
				Read:   adminUsers.NewReadController(s.User),
				Update: adminUsers.NewUpdateController(s.User),
			},
			Organizations: struct {
				List *organizations.ListController
			}{
				List: organizations.NewListController(s.Organization),
			},
		},
	}
}

func NewContainer() *Container {
	providers := initProviders()
	repositories := initRepositories(providers)
	services := initServices(providers, repositories)
	middlewares := initMiddleware(providers, repositories)
	controllers := initControllers(services)

	database.MigrateCollections(providers.Mongo)

	database.SeedCollections(
		repositories.Permission,
		repositories.Role,
		repositories.RolePermission,
	)

	return &Container{
		Providers:    providers,
		Repositories: repositories,
		Services:     services,
		Middleware:   middlewares,
		Controllers:  controllers,
	}
}
