package serviceaccounts

import (
	"net/http"

	"bongaquino/server/app/dto"
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"

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

func (gc *GenerateController) Handle(ctx *gin.Context) {
	var request dto.GenerateServiceAccountDTO

	if err := gc.validatePayload(ctx, &request); err != nil {
		return
	}

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

	// Set other request fields
	userIDStr := userID.(string)
	request.UserID = &userIDStr
	request.ClientID = &clientID
	request.ClientSecret = &clientSecret

	// Create a new service account
	_, err = gc.serviceAccountService.CreateServiceAccount(ctx, &request)
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

func (rc *GenerateController) validatePayload(ctx *gin.Context, request *dto.GenerateServiceAccountDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	return nil
}
