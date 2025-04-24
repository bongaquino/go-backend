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

	// User Routes
	userGroup := engine.Group("/users")
	{
		userGroup.POST("/register", container.Controllers.Users.Register.Handle)
		userGroup.POST("/forgot-password", container.Controllers.Users.ForgotPassword.Handle)
		userGroup.POST("/reset-password", container.Controllers.Users.ResetPassword.Handle)
		userGroup.POST("/verify-account", container.Controllers.Users.VerifyAccount.Handle)
		userGroup.POST("/resend-verification-token", container.Controllers.Users.ResendVerificationToken.Handle)
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
			mfaGroup.POST("/generate-otp", container.Controllers.MFA.Generate.Handle)
			mfaGroup.POST("/enable", container.Controllers.MFA.Enable.Handle)
			mfaGroup.POST("/disable", container.Controllers.MFA.Disable.Handle)
		}
	}

	// Profile Routes
	profileGroup := engine.Group("/profile")
	profileGroup.Use(container.Middleware.Authn.Handle, container.Middleware.Verified.Handle)
	{
		profileGroup.GET("/me", container.Controllers.Profile.Me.Handle)
	}

	// Network Routes
	networkGroup := engine.Group("/network")
	{
		networkGroup.GET("/get-swarm-address", container.Controllers.Network.GetSwarmAddress.Handle)
	}
}
