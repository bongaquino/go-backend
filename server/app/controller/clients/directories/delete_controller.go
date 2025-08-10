package directories

import (
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
	userService *service.UserService
}

// NewDeleteController initializes a new DeleteController
func NewDeleteController(fsService *service.FSService, ipfsService *service.IPFSService, userService *service.UserService) *DeleteController {
	return &DeleteController{
		fsService:   fsService,
		ipfsService: ipfsService,
		userService: userService,
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

	// Check if the directory ID is in valid format
	if directoryID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "directory ID is required", nil, nil)
		return
	}
	if _, err := primitive.ObjectIDFromHex(directoryID); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid directory ID format", nil, nil)
		return
	}

	// Read the directory to get its parent before deletion
	dir, _, _, err := dc.fsService.ReadDirectory(ctx, directoryID, userID.(string))
	if err != nil || dir == nil {
		helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
		return
	}

	// Delete the directory using the fsService
	err= dc.fsService.DeleteDirectory(ctx, directoryID, userID.(string))
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

    if dir.DirectoryID != nil {
		err := dc.fsService.RecalculateDirectorySizeAndParents(ctx, dir.DirectoryID.Hex(), userID.(string))
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to recalculate parent directory sizes", nil, nil)
			return
		}
	}
	// Get user limits
	userLimit, err := dc.userService.GetUserLimits(ctx, userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to get user limits", nil, nil)
		return
	}

	// Compute the new usage by subtracting the directory size from the current usage
	// Since the directory is deleted, we need to decrease the usage
	dirSize := dir.Size
	newUsage := userLimit.BytesUsage - dirSize
	if newUsage < 0 {
		newUsage = 0 // Ensure usage doesn't go negative
	}

	// Update user usage
	err = dc.userService.UpdateUserUsage(ctx, userID.(string), newUsage)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update user usage", nil, nil)
		return
	}

	// If the directory is deleted successfully, return a success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "directory deleted successfully", nil, nil)
}
