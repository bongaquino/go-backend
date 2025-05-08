package start

import (
	ioc "koneksi/server/core/container"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the application's routes
func RegisterRoutes(engine *gin.Engine, container *ioc.Container) {
	// Check Health Route
	engine.GET("/", container.Controllers.Health.Check.Handle)
	engine.GET("/check-health", container.Controllers.Health.Check.Handle)

	// Fetch Constants Route
	engine.GET("/fetch-constants", container.Controllers.Constants.Fetch.Handle)

	// User Routes
	userGroup := engine.Group("/users")
	{
		userGroup.POST("/register", container.Controllers.Users.Register.Handle)
		userGroup.POST("/forgot-password", container.Controllers.Users.ForgotPassword.Handle)
		userGroup.POST("/reset-password", container.Controllers.Users.ResetPassword.Handle)
		userGroup.Use(container.Middleware.Authn.Handle).POST("/verify-account", container.Controllers.Users.VerifyAccount.Handle)
		userGroup.Use(container.Middleware.Authn.Handle).POST("/resend-verification-code", container.Controllers.Users.ResendVerificationCode.Handle)
	}

	// Token Routes
	tokenGroup := engine.Group("/tokens")
	{
		tokenGroup.POST("/request", container.Controllers.Tokens.Request.Handle)
		tokenGroup.POST("/verify-otp", container.Controllers.Tokens.Verify.Handle)
		tokenGroup.POST("/refresh", container.Controllers.Tokens.Refresh.Handle)
		tokenGroup.POST("/revoke", container.Controllers.Tokens.Revoke.Handle)
	}

	// Settings Routes
	settingsGroup := engine.Group("/settings")
	settingsGroup.Use(container.Middleware.Authn.Handle, container.Middleware.Verified.Handle)
	{
		settingsGroup.POST("/change-password", container.Controllers.Settings.ChangePassword.Handle)

		// MFA Routes
		mfaGroup := settingsGroup.Group("/mfa")
		{
			mfaGroup.POST("/generate-otp", container.Controllers.Settings.MFA.Generate.Handle)
			mfaGroup.POST("/enable", container.Controllers.Settings.MFA.Enable.Handle)
			mfaGroup.POST("/disable", container.Controllers.Settings.MFA.Disable.Handle)
		}
	}

	// Profile Routes
	profileGroup := engine.Group("/profile")
	profileGroup.Use(container.Middleware.Authn.Handle)
	{
		profileGroup.GET("/me", container.Controllers.Profile.Me.Handle)
	}

	// Network Routes
	networkGroup := engine.Group("/network")
	{
		networkGroup.GET("/get-swarm-address", container.Controllers.Network.GetSwarmAddress.Handle)
	}

	// Admin Routes
	adminGroup := engine.Group("/admin")
	adminGroup.Use(container.Middleware.Authn.Handle, container.Middleware.Authz.Handle([]string{"system_admin"}))
	{
		// User Management Routes
		adminGroup.GET("users/list", container.Controllers.Admin.Users.List.Handle)
		adminGroup.POST("users/create", container.Controllers.Admin.Users.Create.Handle)
		adminGroup.GET("users/:userID/read", container.Controllers.Admin.Users.Read.Handle)
		adminGroup.PUT("users/:userID/update", container.Controllers.Admin.Users.Update.Handle)
		adminGroup.GET("users/search", container.Controllers.Admin.Users.Search.Handle)
		// Organization Management Routes
		adminGroup.GET("organizations/list", container.Controllers.Admin.Organizations.List.Handle)
		adminGroup.POST("organizations/create", container.Controllers.Admin.Organizations.Create.Handle)
		adminGroup.PUT("organizations/:orgID/update", container.Controllers.Admin.Organizations.Update.Handle)
	}
}
