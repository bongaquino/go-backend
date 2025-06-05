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

func (rc *RevokeController) Handle(ctx *gin.Context) {
	var request struct {
		ClientID string `json:"client_id" binding:"required"`
	}

	// Validate the payload
	if err := rc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Revoke the service account
	err := rc.serviceAccountService.DeleteServiceAccount(ctx, userID.(string), request.ClientID)
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

// validatePayload validates the incoming request payload
func (rc *RevokeController) validatePayload(ctx *gin.Context, request any) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
