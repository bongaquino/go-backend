package start

import (
	ioc "koneksi/services/iam/core/container"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up the application's routes
func RegisterRoutes(engine *gin.Engine, container *ioc.Container) {
	// Swagger Route
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health Check Route
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
}
