package profile

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"koneksi/server/config"

	"github.com/gin-gonic/gin"
)

// MeController handles health-related endpoints
type MeController struct {
	userService *service.UserService
}

// NewMeController initializes a new MeController
func NewMeController(userService *service.UserService) *MeController {
	return &MeController{
		userService: userService,
	}
}

// Handles the health check endpoint
func (hc *MeController) Handle(c *gin.Context) {
	appConfig := config.LoadAppConfig()

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, nil, gin.H{
		"name":    appConfig.AppName,
		"version": appConfig.AppVersion,
		"healthy": true,
	}, nil)
}
