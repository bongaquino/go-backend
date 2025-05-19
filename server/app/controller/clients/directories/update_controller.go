package directories

import (
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
}

// NewUpdateController initializes a new UpdateController
func NewUpdateController(fsService *service.FSService, ipfsService *service.IPFSService) *UpdateController {
	return &UpdateController{
		fsService:   fsService,
		ipfsService: ipfsService,
	}
}

func (uc *UpdateController) Handle(ctx *gin.Context) {
	// Validate the request payload
	var request dto.UpdateDirectoryDTO
	if err := uc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Get the directory ID from the URL parameters
	directoryID := ctx.Param("directoryID")

	// Check id directoryID is same as request.DirectoryID
	if directoryID == *request.DirectoryID {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "parent directory cannot be same as current directory", nil, nil)
		return
	}

	// Update the directory using the fsService
	err := uc.fsService.UpdateDirectory(ctx, directoryID, userID.(string), &request)
	if err != nil {
		if err.Error() == "directory not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
			return
		}
		if err.Error() == "cannot update root directory" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "cannot update root directory", nil, nil)
			return
		}
		if err.Error() == "parent directory not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "parent directory not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update directory", nil, nil)
		return
	}

	// Return success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "directory updated successfully", nil, nil)
}

func (uc *UpdateController) validatePayload(ctx *gin.Context, request *dto.UpdateDirectoryDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
