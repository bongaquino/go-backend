package files

import (
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GenerateLinkController struct {
	fsService *service.FSService
}

// NewGenerateLinkController initializes a new GenerateLinkController
func NewGenerateLinkController(
	fsService *service.FSService,
) *GenerateLinkController {
	return &GenerateLinkController{
		fsService: fsService,
	}
}

func (sc *GenerateLinkController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Get the file ID from the URL parameters
	fileID := ctx.Param("fileID")

	// Check if user ID owns the file ID
	isOwned, _ := sc.fsService.CheckFileOwnership(ctx, fileID, userID.(string))
	if !isOwned {
		helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
		return
	}

	// Verify if duration is provided in the request body
	requestBody := make(map[string]any)
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "duration is required", nil, nil)
		return
	}
	if _, ok := requestBody["duration"]; !ok || requestBody["duration"] == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "duration is required", nil, nil)
		return
	}

	// Save the ile ID in Redis
	durationVal, ok := requestBody["duration"].(float64)
	if !ok {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "duration must be a number", nil, nil)
		return
	}

	// Generate file key
	fileKey, _ := helper.GenerateFileKey(fileID)

	// Set the file key in Redis with the specified duration
	duration := time.Duration(int(durationVal)) * time.Second
	if err := sc.fsService.SetTemporaryFileKey(ctx, fileKey, fileID, duration); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to cache file ID", nil, nil)
		return
	}
	responseBody := map[string]any{
		"duration": duration.String(),
		"file_key": fileKey,
	}

	// Return success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "file shared successfully", responseBody, nil)
}
