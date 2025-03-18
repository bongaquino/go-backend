package start

import (
	ioc "koneksi/orchestrator/core/container"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up the application's routes
func RegisterRoutes(engine *gin.Engine, container *ioc.Container) {
	// Health check endpoint
	engine.GET("/", container.HealthController.Check)

	// Register Swagger handler
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
