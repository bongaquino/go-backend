package start

import (
	ioc "koneksi/services/iam/core/container"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up the application's routes
func RegisterRoutes(engine *gin.Engine, container *ioc.Container) {
	// Swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Check health
	engine.GET("/check-health", container.CheckHealthController.Handle)

	// Register
	engine.POST("/register", container.RegisterController.Handle)

	// Request token
	engine.POST("/request-token", container.RequestTokenController.Handle)

	// Refresh token
	engine.POST("/refresh-token", container.RefreshTokenController.Handle)

	// Revoke token
	engine.POST("/revoke-token", container.RevokeTokenController.Handle)
}
