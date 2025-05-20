package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

type DeleteController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
}

// NewDeleteController initializes a new DeleteController
func NewDeleteController(fsService *service.FSService, ipfsService *service.IPFSService) *DeleteController {
	return &DeleteController{
		fsService:   fsService,
		ipfsService: ipfsService,
	}
}

func (dc *DeleteController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Get the file ID from the URL parameters
	fileID := ctx.Param("fileID")

	// Check if the file ID is in valid format
	if fileID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
		return
	}
	if _, err := primitive.ObjectIDFromHex(fileID); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
		return
	}

	// Delete the file using the fsService
	err := dc.fsService.DeleteFile(ctx, fileID, userID.(string))
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to delete file", nil, nil)
		return
	}

	// If the file is successfully deleted, return a success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "file deleted successfully", nil, nil)
}
