package health

import (
	"net/http"

	"koneksi/services/iam/app/helpers"
	"koneksi/services/iam/config"

	"github.com/gin-gonic/gin"
)

// CheckHealthController handles health-related endpoints
type CheckHealthController struct{}

// NewCheckHealthController initializes a new CheckHealthController
func NewCheckHealthController() *CheckHealthController {
	return &CheckHealthController{}
}

// Check processes a health check request
// @Summary Health check
// @Description Get the health status of the server
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func (hc *CheckHealthController) Handle(c *gin.Context) {
	appConfig := config.LoadAppConfig()

	// Respond with success
	helpers.FormatResponse(c, "success", http.StatusOK, nil, gin.H{
		"name":    appConfig.AppName,
		"version": appConfig.AppVersion,
		"healthy": true,
	}, nil)
}
