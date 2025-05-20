package files

import (
	"fmt"
	"koneksi/server/app/helper"
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

	// Read the file from the FS service
	file, err := rc.fsService.ReadFileByIDUserID(ctx, fileID, userID.(string))
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error reading file", nil, nil)
		return
	}

	// Get the file details
	hash := file.Hash

	// Generate file URL
	fileURL := rc.ipfsService.GetFileURL(hash)
	if fileURL == "" {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error generating file URL", nil, nil)
		return
	}

	fmt.Println("file", file)

	// Return the file details
	helper.FormatResponse(ctx, "success", http.StatusOK, "file read successfully", gin.H{
		"id":           file.ID.Hex(),
		"directory_id": file.DirectoryID.Hex(),
		"name":         file.Name,
		"size":         file.Size,
		"hash":         file.Hash,
		"content_type": file.ContentType,
		"ipfs_url":     fileURL,
	}, nil)
}
