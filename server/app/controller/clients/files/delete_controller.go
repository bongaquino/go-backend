package files

import (
	"fmt"
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

func (rc *DeleteController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Get the file ID from the URL parameters
	fileID := ctx.Param("fileID")
	if fileID == ":file" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
		return
	}

	fmt.Println("userID", userID)
	fmt.Println("fileID", fileID)
}
