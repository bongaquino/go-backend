package dashboard

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CollectMetricsController handles health-related endpoints
type CollectMetricsController struct {
	userService *service.UserService
}

// NewCollectMetricsController initializes a new CollectMetricsController
func NewCollectMetricsController(userService *service.UserService) *CollectMetricsController {
	return &CollectMetricsController{
		userService: userService,
	}
}

// Handles the health check endpoint
func (cmc *CollectMetricsController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Fetch user metrics using the user service
	limit, directoryCount, fileCount, svcAccCount, err := cmc.userService.CollectMetrics(ctx, userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to collect metrics", nil, err)
		return
	}

	// Prepare the response data
	data := gin.H{
		"byte_limit":            limit.BytesLimit,
		"byte_usage":            limit.BytesUsage,
		"directory_count":       directoryCount,
		"file_count":            fileCount,
		"service_account_count": svcAccCount,
	}

	// Send the response
	helper.FormatResponse(ctx, "success", http.StatusOK, "metrics collected successfully", data, nil)
}
