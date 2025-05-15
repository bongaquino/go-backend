package serviceaccounts

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/config"

	"github.com/gin-gonic/gin"
)

type RevokeController struct{}

func NewRevokeController() *RevokeController {
	return &RevokeController{}
}

func (hc *RevokeController) Handle(ctx *gin.Context) {
	appConfig := config.LoadAppConfig()

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"name":    appConfig.AppName,
		"version": appConfig.AppVersion,
		"healthy": true,
	}, nil)
}
