package serviceaccounts

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

type GenerateController struct {
	serviceAccountService *service.ServiceAccountService
}

func NewGenerateController(serviceAccountService *service.ServiceAccountService) *GenerateController {
	return &GenerateController{
		serviceAccountService: serviceAccountService,
	}
}

func (hc *GenerateController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Generate client credentials
	clientID, clientSecret, err := service.GenerateClientCredentials()
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to generate client credentials", nil, nil)
		return
	}

	// Create a new service account
	_, _, err = hc.serviceAccountService.CreateServiceAccount(ctx, userID.(string), clientID, clientSecret)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to create service account", nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"client_id":     clientID,
		"client_secret": clientSecret,
	}, nil)
}
