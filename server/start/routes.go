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
	}

	// Token Routes
	tokenGroup := engine.Group("/tokens")
	{
		tokenGroup.POST("/request", container.RequestTokenController.Handle)
		tokenGroup.POST("/refresh", container.RefreshTokenController.Handle)
		tokenGroup.POST("/revoke", container.RevokeTokenController.Handle)
	}

	// Test route
	protected := engine.Group("/protected")
	// protected.Use(container.AuthnMiddleware.HandleAuthn, container.AuthzMiddleware.HandleAuthz([]string{"user"}), container.VerifiedMiddleware.HandleVerified)
	protected.Use(container.AuthnMiddleware.HandleAuthn, container.VerifiedMiddleware.HandleVerified)
	{
		protected.GET("/test", func(c *gin.Context) {
			userID := c.GetString("userID")
			c.JSON(200, gin.H{"message": "Access granted", "userID": userID})
		})
	}
}
