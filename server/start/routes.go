package start

import (
	ioc "koneksi/server/core/container"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the application's routes
func RegisterRoutes(engine *gin.Engine, container *ioc.Container) {
	// Check Health Route
	engine.GET("/", container.CheckHealthController.Handle)
	engine.GET("/check-health", container.CheckHealthController.Handle)

	// User Routes
	userGroup := engine.Group("/users")
	{
		userGroup.POST("/register", container.RegisterController.Handle)
		userGroup.POST("/forgot-password", container.ForgotPasswordController.Handle)
		// userGroup.POST("/reset-password", container.ResetPasswordController.Handle)
	}

	// Token Routes
	tokenGroup := engine.Group("/tokens")
	{
		tokenGroup.POST("/request", container.RequestTokenController.Handle)
		tokenGroup.POST("/refresh", container.RefreshTokenController.Handle)
		tokenGroup.POST("/revoke", container.RevokeTokenController.Handle)
	}

	// Settings Routes
	settingsGroup := engine.Group("/settings")
	settingsGroup.Use(container.AuthnMiddleware.HandleAuthn, container.VerifiedMiddleware.HandleVerified)
	{
		settingsGroup.POST("/change-password", container.ChangePasswordController.Handle)

		// MFA Routes
		mfaGroup := settingsGroup.Group("/mfa")
		{
			mfaGroup.POST("/generate-otp", container.GenerateOTPController.Handle)
			mfaGroup.POST("/enable", container.EnableMFAController.Handle)
			mfaGroup.POST("/disable", container.DisableMFAController.Handle)
		}
	}
}
