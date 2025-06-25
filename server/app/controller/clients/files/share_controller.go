package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"koneksi/server/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShareController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
}

// NewShareController initializes a new ShareController
func NewShareController(fsService *service.FSService, ipfsService *service.IPFSService) *ShareController {
	return &ShareController{
		fsService:   fsService,
		ipfsService: ipfsService,
	}
}

func (sc *ShareController) Handle(ctx *gin.Context) {
	// Load file configuration
	fileConfig := config.LoadFileConfig()

	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Get the file ID from the URL parameters
	fileID := ctx.Param("fileID")

	// Get the access type from the query parameters (default: "private")
	accessType := ctx.DefaultQuery("access", fileConfig.DefaultAccess)

	// Check if access type is valid by checking against the allowed options
	if !helper.Contains(fileConfig.AccessOptions, accessType) {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid access type", nil, nil)
		return
	}

	// Validate the access type and prepare the request body if needed
	var requestBody map[string]any
	switch accessType {
	case fileConfig.PublicAccess, fileConfig.PrivateAccess:
		// No body needed
	case fileConfig.PasswordAccess:
		requestBody = make(map[string]any)
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "password is required for password access", nil, nil)
			return
		}
		if _, ok := requestBody["password"]; !ok || requestBody["password"] == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "password is required for password access", nil, nil)
			return
		}
	case fileConfig.EmailAccess:
		requestBody = make(map[string]any)
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "emails are required for email access", nil, nil)
			return
		}
		emails, ok := requestBody["emails"].([]any)
		if !ok || len(emails) == 0 {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "at least one email is required for email access", nil, nil)
			return
		}
	}

	// Check if the file ID is in valid format
	if fileID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
		return
	}
	if _, err := primitive.ObjectIDFromHex(fileID); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
		return
	}

	// Update the file's access type using the FS service
	err := sc.fsService.UpdateFileAccess(ctx, fileID, userID.(string), accessType)
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error updating file access", nil, nil)
		return
	}

	// Return success response
	if accessType == fileConfig.EmailAccess {
		helper.FormatResponse(ctx, "success", http.StatusOK, "file shared successfully", nil, requestBody)
	}
	helper.FormatResponse(ctx, "success", http.StatusOK, "file shared successfully", nil, nil)
}
