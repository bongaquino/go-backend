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
	file, err := rc.fsService.ReadFileByID(ctx, fileID)
	if err != nil {
		fmt.Println("Error reading file:", err)
		helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
		return
	}

	// Return the file details
	helper.FormatResponse(ctx, "success", http.StatusOK, "file read successfully", gin.H{
		"id":           file.ID.Hex(),
		"access":       file.Access,
		"size":         file.Size,
		"is_encrypted": file.IsEncrypted,
	}, nil)
}
