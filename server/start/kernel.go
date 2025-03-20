package start

import (
	"fmt"

	ioc "koneksi/server/core/container"
	"koneksi/server/core/env"
	"koneksi/server/core/logger"

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
