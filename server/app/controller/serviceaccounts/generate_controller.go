package serviceaccounts

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/config"

	"github.com/gin-gonic/gin"
)

type GenerateController struct{}

func NewGenerateController() *GenerateController {
	return &GenerateController{}
}

func (hc *GenerateController) Handle(ctx *gin.Context) {
	appConfig := config.LoadAppConfig()

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"name":    appConfig.AppName,
		"version": appConfig.AppVersion,
		"healthy": true,
	}, nil)
}
