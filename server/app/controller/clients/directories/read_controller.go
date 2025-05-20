package directories

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReadController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
}

// NewReadController initializes a new ReadController
func NewReadController(fsService *service.FSService, ipfsService *service.IPFSService) *ReadController {
	return &ReadController{
		fsService:   fsService,
		ipfsService: ipfsService,
	}
}

func (rc *ReadController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Get the directory ID from the URL parameters
	directoryID := ctx.Param("directoryID")

	// Check if the directory ID is not empty
	if directoryID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
		return
	}

	if directoryID == "root" {
		// Use fsService to read the root directory
		directory, subDirectories, files, err := rc.fsService.ReadRootDirectory(ctx, userID.(string))
		if err != nil {
			if err.Error() == "directory not found" {
				helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
				return
			}
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to read root directory", nil, nil)
			return
		}

		// Ensure subDirectories and files are not nil
		if subDirectories == nil {
			subDirectories = []*model.Directory{}
		}
		if files == nil {
			files = []*model.File{}
		}

		// Format the directory data
		directoryData := gin.H{
			"id":        directory.ID.Hex(),
			"name":      directory.Name,
			"size":      directory.Size,
			"createdAt": directory.CreatedAt,
			"updatedAt": directory.UpdatedAt,
		}

		// Format the subdirectories
		subDirectoriesData := make([]gin.H, len(subDirectories))
		for i, subDir := range subDirectories {
			subDirectoriesData[i] = gin.H{
				"id":        subDir.ID.Hex(),
				"name":      subDir.Name,
				"size":      subDir.Size,
				"createdAt": subDir.CreatedAt,
				"updatedAt": subDir.UpdatedAt,
			}
		}

		// Format the files
		filesData := make([]gin.H, len(files))
		for i, file := range files {
			filesData[i] = gin.H{
				"id":          file.ID.Hex(),
				"name":        file.Name,
				"hash":        file.Hash,
				"size":        file.Size,
				"contentType": file.ContentType,
				"createdAt":   file.CreatedAt,
				"updatedAt":   file.UpdatedAt,
			}
		}

		// Prepare the response
		response := gin.H{
			"directory":      directoryData,
			"subdirectories": subDirectoriesData,
			"files":          filesData,
		}

		// Send the response
		helper.FormatResponse(ctx, "success", http.StatusOK, "directory read successfully", response, nil)
	} else {
		// Check if the file ID is in valid format
		if directoryID == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
			return
		}
		if _, err := primitive.ObjectIDFromHex(directoryID); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
			return
		}

		// Use fsService to read the root directory
		directory, subDirectories, files, err := rc.fsService.ReadDirectory(ctx, directoryID, userID.(string))
		if err != nil {
			if err.Error() == "directory not found" {
				helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
				return
			}
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to read directory", nil, nil)
			return
		}

		// Ensure subDirectories and files are not nil
		if subDirectories == nil {
			subDirectories = []*model.Directory{}
		}
		if files == nil {
			files = []*model.File{}
		}

		// Format the directory data
		directoryData := gin.H{
			"id":        directory.ID.Hex(),
			"name":      directory.Name,
			"size":      directory.Size,
			"createdAt": directory.CreatedAt,
			"updatedAt": directory.UpdatedAt,
		}

		// Format the subdirectories
		subDirectoriesData := make([]gin.H, len(subDirectories))
		for i, subDir := range subDirectories {
			subDirectoriesData[i] = gin.H{
				"id":        subDir.ID.Hex(),
				"name":      subDir.Name,
				"size":      subDir.Size,
				"createdAt": subDir.CreatedAt,
				"updatedAt": subDir.UpdatedAt,
			}
		}

		// Format the files
		filesData := make([]gin.H, len(files))
		for i, file := range files {
			filesData[i] = gin.H{
				"id":          file.ID.Hex(),
				"name":        file.Name,
				"hash":        file.Hash,
				"size":        file.Size,
				"contentType": file.ContentType,
				"createdAt":   file.CreatedAt,
				"updatedAt":   file.UpdatedAt,
			}
		}

		// Prepare the response
		response := gin.H{
			"directory":      directoryData,
			"subdirectories": subDirectoriesData,
			"files":          filesData,
		}

		// Send the response
		helper.FormatResponse(ctx, "success", http.StatusOK, "directory read successfully", response, nil)
	}
}
