package start

import (
	ioc "bongaquino/server/core/container"

	"github.com/gin-gonic/gin"
)

// RegisterMiddleware sets up the middleware for the Gin engine
func RegisterMiddleware(engine *gin.Engine, container *ioc.Container) {
	// Global middleware
}
