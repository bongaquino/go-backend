package start

import (
	"koneksi/server/app/middleware"
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

	// Protected route
	protected := engine.Group("/protected")
	protected.Use(middleware.AuthMiddleware(container.JwtService))
	{
		protected.GET("/read", func(c *gin.Context) {
			userID := c.GetString("userID")
			email := c.GetString("email")
			c.JSON(200, gin.H{"message": "Access granted", "userID": userID, "email": email})
		})
	}
}
