package serviceaccounts

import (
	"net/http"

	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"

	"github.com/gin-gonic/gin"
)

type BrowseController struct {
	serviceAccountService *service.ServiceAccountService
}

func NewBrowseController(serviceAccountService *service.ServiceAccountService) *BrowseController {
	return &BrowseController{
		serviceAccountService: serviceAccountService,
	}
}

func (hc *BrowseController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Get the list of service accounts for the user
	serviceAccounts, err := hc.serviceAccountService.ListServiceAccounts(ctx, userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to retrieve service accounts", nil, nil)
		return
	}

	// Redact sensitive information
	for _, account := range serviceAccounts {
		account.ClientSecret = "REDACTED"
	}

	// Respond with the list of service accounts
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, serviceAccounts, nil)
}
