package serviceaccounts

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

type RevokeController struct {
	serviceAccountService *service.ServiceAccountService
}

func NewRevokeController(serviceAccountService *service.ServiceAccountService) *RevokeController {
	return &RevokeController{
		serviceAccountService: serviceAccountService,
	}
}

func (hc *RevokeController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Get client ID from query parameters
	clientID := ctx.Query("client_id")
	if clientID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "client_id is required", nil, nil)
		return
	}

	// Revoke the service account
	err := hc.serviceAccountService.DeleteServiceAccount(ctx, userID.(string), clientID)
	if err != nil {
		if err.Error() == "service account not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "service account not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to revoke service account", nil, nil)
		return
	}
	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, "service account revoked successfully", nil, nil)
}
