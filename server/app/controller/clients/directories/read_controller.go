package directories

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Get the directory ID from the URL parameters
	directoryID := ctx.Param("directoryID")

	if directoryID == "root" {
		// Use fsService to read the root directory
		directory, subDirectories, files, err := rc.fsService.ReadRootDirectory(ctx, userID.(string))
		if err != nil {
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

		// Prepare the response
		response := gin.H{
			"directory": gin.H{
				"id":        directory.ID.Hex(),
				"name":      directory.Name,
				"size":      directory.Size,
				"createdAt": directory.CreatedAt,
				"updatedAt": directory.UpdatedAt,
			},
			"subdirectories": subDirectories,
			"files":          files,
		}

		// Send the response
		helper.FormatResponse(ctx, "success", http.StatusOK, "directory fetched successfully", response, nil)
	} else {
		// Use fsService to read the root directory
		directory, subDirectories, files, err := rc.fsService.ReadDirectory(ctx, directoryID, userID.(string))
		if err != nil {
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

		// Prepare the response
		response := gin.H{
			"directory": gin.H{
				"id":        directory.ID.Hex(),
				"name":      directory.Name,
				"size":      directory.Size,
				"createdAt": directory.CreatedAt,
				"updatedAt": directory.UpdatedAt,
			},
			"subdirectories": subDirectories,
			"files":          files,
		}

		// Send the response
		helper.FormatResponse(ctx, "success", http.StatusOK, "directory fetched successfully", response, nil)
	}
}
