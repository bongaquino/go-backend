package health

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/config"

	"github.com/gin-gonic/gin"
)

// CheckHealthController handles health-related endpoints
type CheckHealthController struct{}

// NewCheckHealthController initializes a new CheckHealthController
func NewCheckHealthController() *CheckHealthController {
	return &CheckHealthController{}
}

// Handles the health check endpoint
func (hc *CheckHealthController) Handle(ctx *gin.Context) {
	appConfig := config.LoadAppConfig()

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"name":    appConfig.AppName,
		"version": appConfig.AppVersion,
		"healthy": true,
	}, nil)
}
