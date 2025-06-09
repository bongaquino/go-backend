package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DownloadController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
}

// NewDownloadController initializes a new DownloadController
func NewDownloadController(fsService *service.FSService, ipfsService *service.IPFSService) *DownloadController {
	return &DownloadController{
		fsService:   fsService,
		ipfsService: ipfsService,
	}
}

func (dc *DownloadController) Handle(ctx *gin.Context) {
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

	// Check if the file is owned by the user
	file, err := dc.fsService.ReadFileByIDUserID(ctx, fileID, userID.(string))
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error reading file", nil, nil)
		return
	}

	// Get the file hash
	fileHash := file.Hash
	if fileHash == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file hash is required for download", nil, nil)
		return
	}

	// Download the file from IPFS
	fileContent, err := dc.ipfsService.DownloadFile(fileHash)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error downloading file from IPFS", nil, nil)
		return
	}

	// Set the response headers for file download
	ctx.Header("Content-Disposition", "attachment; filename="+file.Name)
	ctx.Header("Content-Length", strconv.Itoa(len(fileContent)))
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")
	ctx.Data(http.StatusOK, file.ContentType, fileContent)
	helper.FormatResponse(ctx, "success", http.StatusOK, "file downloaded successfully", nil, nil)
}
