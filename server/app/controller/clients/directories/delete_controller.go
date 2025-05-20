package directories

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

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

	// Get the directory ID from the URL parameters
	directoryID := ctx.Param("directoryID")

	// Check if the directory ID is not empty
	if directoryID == ":directory" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "directory ID is required", nil, nil)
		return
	}

	// Delete the directory using the fsService
	err := dc.fsService.DeleteDirectory(ctx, directoryID, userID.(string))
	if err != nil {
		if err.Error() == "directory not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
			return
		}
		if err.Error() == "cannot delete root directory" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "cannot delete root directory", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to delete directory", nil, nil)
		return
	}

	// If the directory is deleted successfully, return a success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "directory deleted successfully", nil, nil)
}
