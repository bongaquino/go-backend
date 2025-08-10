package start

import (
	"fmt"

	ioc "bongaquino/server/core/container"
	"bongaquino/server/core/env"
	"bongaquino/server/core/logger"

	"github.com/gin-gonic/gin"
)

// InitializeKernel sets up the Gin engine and starts the server
func InitializeKernel() {
	// Load application environment variables
	env := env.LoadEnv()

	// Set Gin mode
	if env.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize the Gin engine
	engine := gin.Default()

	// Set MaxMultipartMemory to 2GB for large file uploads
	engine.MaxMultipartMemory = 2 << 30 // 2GB

	// Initialize IoC container
	container := ioc.NewContainer()

	// Setup CORS
	SetupCORS(engine)

	// Register middleware
	RegisterMiddleware(engine, container)

	// Register routes
	RegisterRoutes(engine, container)

	// Start the server on the specified port
	address := fmt.Sprintf(":%d", env.Port)

	if err := engine.Run(address); err != nil {
		logger.Log.Fatal("failed to start server", logger.Error(err))
	}
}
