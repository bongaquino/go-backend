package container

import (
	"koneksi/server/app/controller/admin/organizations"
	"koneksi/server/app/controller/admin/organizations/members"
	adminUsers "koneksi/server/app/controller/admin/users"
	"koneksi/server/app/controller/clients/peers"
	"koneksi/server/app/controller/constants"
	"koneksi/server/app/controller/health"
	"koneksi/server/app/controller/network"
	"koneksi/server/app/controller/profile"
	"koneksi/server/app/controller/serviceaccounts"
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
	Limit                *repository.LimitRepository
	Directory            *repository.DirectoryRepository
	File                 *repository.FileRepository
}

type Services struct {
	User           *service.UserService
	Token          *service.TokenService
	MFA            *service.MFAService
	Email          *service.EmailService
	IPFS           *service.IPFSService
	Organization   *service.OrganizationService
	ServiceAccount *service.ServiceAccountService
	Directory      *service.DirectoryService
	File           *service.FileService
}

type Middleware struct {
	Authn    *middleware.AuthnMiddleware
	Authz    *middleware.AuthzMiddleware
	Verified *middleware.VerifiedMiddleware
	Locked   *middleware.LockedMiddleware
	API      *middleware.APIMiddleware
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
	ServiceAccounts struct {
		Browse   *serviceaccounts.BrowseController
		Generate *serviceaccounts.GenerateController
		Revoke   *serviceaccounts.RevokeController
	}
	Clients struct {
		Peers struct {
			Fetch *peers.FetchController
		}
	}
	Admin struct {
		Users struct {
			List   *adminUsers.ListController
			Create *adminUsers.CreateController
			Read   *adminUsers.ReadController
			Update *adminUsers.UpdateController
			Search *adminUsers.SearchController
		}
		Organizations struct {
			List    *organizations.ListController
			Create  *organizations.CreateController
			Read    *organizations.ReadController
			Update  *organizations.UpdateController
			Members struct {
				Add        *members.AddController
				UpdateRole *members.UpdateRoleController
				Remove     *members.RemoveController
			}
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
		Limit:                repository.NewLimitRepository(p.Mongo),
		Directory:            repository.NewDirectoryRepository(p.Mongo),
		File:                 repository.NewFileRepository(p.Mongo),
	}
}

func initServices(p Providers, r Repositories) Services {
	user := service.NewUserService(r.User, r.Profile, r.Role, r.UserRole, r.Limit, p.Redis)
	email := service.NewEmailService(p.Postmark)
	mfa := service.NewMFAService(r.User, p.Redis)
	ipfs := service.NewIPFSService(p.IPFS)
	token := service.NewTokenService(r.User, p.JWT, mfa, p.Redis)
	organization := service.NewOrganizationService(r.Organization, r.Policy, r.Permission,
		r.OrganizationUserRole, r.User, r.Role)
	serviceAccount := service.NewServiceAccountService(r.ServiceAccount, r.User, r.Limit)
	directory := service.NewDirectoryService(r.Directory)
	file := service.NewFileService(r.File)
	return Services{user, token, mfa, email, ipfs, organization, serviceAccount, directory, file}
}

func initMiddleware(p Providers, r Repositories) Middleware {
	return Middleware{
		Authn:    middleware.NewAuthnMiddleware(p.JWT),
		Authz:    middleware.NewAuthzMiddleware(r.UserRole, r.Role),
		Verified: middleware.NewVerifiedMiddleware(r.User),
		Locked:   middleware.NewLockedMiddleware(r.User),
		API:      middleware.NewAPIMiddleware(r.ServiceAccount),
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
			Fetch: constants.NewFetchController(s.User, s.Organization),
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
		ServiceAccounts: struct {
			Browse   *serviceaccounts.BrowseController
			Generate *serviceaccounts.GenerateController
			Revoke   *serviceaccounts.RevokeController
		}{
			Browse:   serviceaccounts.NewBrowseController(s.ServiceAccount),
			Generate: serviceaccounts.NewGenerateController(s.ServiceAccount),
			Revoke:   serviceaccounts.NewRevokeController(s.ServiceAccount),
		},
		Clients: struct {
			Peers struct {
				Fetch *peers.FetchController
			}
		}{
			Peers: struct {
				Fetch *peers.FetchController
			}{
				Fetch: peers.NewFetchController(s.IPFS),
			},
		},
		Admin: struct {
			Users struct {
				List   *adminUsers.ListController
				Create *adminUsers.CreateController
				Read   *adminUsers.ReadController
				Update *adminUsers.UpdateController
				Search *adminUsers.SearchController
			}
			Organizations struct {
				List    *organizations.ListController
				Create  *organizations.CreateController
				Read    *organizations.ReadController
				Update  *organizations.UpdateController
				Members struct {
					Add        *members.AddController
					UpdateRole *members.UpdateRoleController
					Remove     *members.RemoveController
				}
			}
		}{
			Users: struct {
				List   *adminUsers.ListController
				Create *adminUsers.CreateController
				Read   *adminUsers.ReadController
				Update *adminUsers.UpdateController
				Search *adminUsers.SearchController
			}{
				List:   adminUsers.NewListController(s.User),
				Create: adminUsers.NewCreateController(s.User, s.Token, s.Email),
				Read:   adminUsers.NewReadController(s.User),
				Update: adminUsers.NewUpdateController(s.User),
				Search: adminUsers.NewSearchController(s.User),
			},
			Organizations: struct {
				List    *organizations.ListController
				Create  *organizations.CreateController
				Read    *organizations.ReadController
				Update  *organizations.UpdateController
				Members struct {
					Add        *members.AddController
					UpdateRole *members.UpdateRoleController
					Remove     *members.RemoveController
				}
			}{
				List:   organizations.NewListController(s.Organization),
				Create: organizations.NewCreateController(s.Organization),
				Read:   organizations.NewReadController(s.Organization),
				Update: organizations.NewUpdateController(s.Organization),
				Members: struct {
					Add        *members.AddController
					UpdateRole *members.UpdateRoleController
					Remove     *members.RemoveController
				}{
					Add:        members.NewAddController(s.Organization),
					UpdateRole: members.NewUpdateRoleController(s.Organization),
					Remove:     members.NewRemoveController(s.Organization),
				},
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
		repositories.Policy,
		repositories.PolicyPermission,
	)

	return &Container{
		Providers:    providers,
		Repositories: repositories,
		Services:     services,
		Middleware:   middlewares,
		Controllers:  controllers,
	}
}
