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

	// Fetch the file to get its size
	file, err := dc.fsService.ReadFileByIDUserID(ctx, fileID, userID.(string))
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error reading file", nil, nil)
		return
	}

	// Delete the file using the fsService
	err = dc.fsService.DeleteFile(ctx, fileID, userID.(string))
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to delete file", nil, nil)
		return
	}

	err = dc.fsService.RecalculateDirectorySizeAndParents(ctx, file.DirectoryID.Hex(), userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to recalculate directory sizes", nil, nil)
		return
	}

	// Get user limits (usage)
	userLimit, err := dc.userService.GetUserLimits(ctx, userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to get user limits", nil, nil)
		return
	}

	// Compute the new usage
	newUsage := userLimit.BytesUsage - file.Size
	if newUsage < 0 {
		newUsage = 0
	}

	// Update user usage
	err = dc.userService.UpdateUserUsage(ctx, userID.(string), newUsage)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update user usage", nil, nil)
		return
	}

	// If the file is successfully deleted, return a success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "file deleted successfully", nil, nil)
}
