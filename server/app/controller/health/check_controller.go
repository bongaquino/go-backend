package health

import (
	"net/http"

	"bongaquino/server/app/helper"
	"bongaquino/server/config"

	"github.com/gin-gonic/gin"
)

// CheckController handles health-related endpoints
type CheckController struct{}

// NewCheckController initializes a new CheckController
func NewCheckController() *CheckController {
	return &CheckController{}
}

// Handles the health check endpoint
func (hc *CheckController) Handle(ctx *gin.Context) {
	appConfig := config.LoadAppConfig()

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"name":    appConfig.AppName,
		"version": appConfig.AppVersion,
		"healthy": true,
	}, nil)
}
