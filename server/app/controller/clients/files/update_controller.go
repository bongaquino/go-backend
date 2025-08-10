package files

import (
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var request dto.UpdateFileDTO
	if err := uc.validatePayload(ctx, &request); err != nil {
		return
	}
	// Limit file name length to 255 characters while preserving extension
	isTrimmed := false
	if request.Name != "" && len(request.Name) > 255 {
		// Get the file extension
		ext := filepath.Ext(request.Name)

		// Calculate the base name (filename without extension)
		baseName := strings.TrimSuffix(request.Name, ext)

		// Calculate how much we need to trim from the base name
		maxBaseNameLength := 255 - len(ext)

		if maxBaseNameLength > 0 {
			// Trim the base name if it's too long
			if len(baseName) > maxBaseNameLength {
				baseName = baseName[:maxBaseNameLength]
			}
			// Reconstruct the filename with preserved extension
			request.Name = baseName + ext
		} else {
			// If extension itself is >= 255 characters, just truncate the whole name
			request.Name = request.Name[:255]
		}

		isTrimmed = true
	}

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
	file, err := uc.fsService.ReadFileByIDUserID(ctx, fileID, userID.(string))
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error reading file", nil, nil)
		return
	}

	oldDirectoryID := file.DirectoryID.Hex()

	// Update the file using the fsService
	err = uc.fsService.UpdateFile(ctx, fileID, userID.(string), &request)
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		if err.Error() == "directory not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
			return
		}
		if err.Error() == "error fetching directory" {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error fetching directory", nil, nil)
			return
		}
		if err.Error() == "no fields to update" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "no fields to update", nil, nil)
			return
		}
		if err.Error() == "cannot change file extension" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "cannot change file extension", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error updating file", nil, nil)
		return
	}

	// Fetch the file to get its size
	updatedFile, err := uc.fsService.ReadFileByIDUserID(ctx, fileID, userID.(string))
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error reading file", nil, nil)
		return
	}

	newDirectoryID := updatedFile.DirectoryID.Hex()

	if oldDirectoryID != newDirectoryID {
		err := uc.fsService.RecalculateDirectorySizeAndParents(ctx, oldDirectoryID, userID.(string))
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to recalculate directory sizes", nil, nil)
			return
		}
	}

	err = uc.fsService.RecalculateDirectorySizeAndParents(ctx, newDirectoryID, userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to recalculate directory sizes", nil, nil)
		return
	}

	if isTrimmed {
		meta := map[string]any{
			"is_trimmed": true,
		}
		helper.FormatResponse(ctx, "success", http.StatusOK, "file updated successfully", nil, meta)
		return
	}

	// Return success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "file updated successfully", nil, nil)
}

func (uc *UpdateController) validatePayload(ctx *gin.Context, request *dto.UpdateFileDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	return nil
}
