package start

import (
	ioc "koneksi/services/iam/core/container"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the application's routes
func RegisterRoutes(engine *gin.Engine, container *ioc.Container) {
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
}
