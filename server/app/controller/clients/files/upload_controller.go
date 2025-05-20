package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
	userService *service.UserService
}

// NewUploadController initializes a new UploadController
func NewUploadController(fsService *service.FSService,
	ipfsService *service.IPFSService,
	userService *service.UserService,
) *UploadController {
	return &UploadController{
		fsService:   fsService,
		ipfsService: ipfsService,
		userService: userService,
	}
}

func (uc *UploadController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Extract directory ID from the query parameters
	directoryID := ctx.Query("directory_id")
	if directoryID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "directory_id is required", nil, nil)
		return
	}

	// Handle file upload
	file, err := ctx.FormFile("file")
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "failed to get uploaded file", nil, nil)
		return
	}

	// Get file metadata
	fileName := file.Filename
	fileSize := file.Size
	fileType := file.Header.Get("Content-Type")

	// Get user limits
	userLimit, err := uc.userService.GetUserLimits(ctx, userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to get user limits", nil, nil)
		return
	}

	// Check if the user has reached their upload limit
	if userLimit.BytesUsage+fileSize > userLimit.BytesLimit {
		helper.FormatResponse(ctx, "error", http.StatusForbidden, "upload limit reached", nil, nil)
		return
	}

	// Save the uploaded file to a temporary location
	destination := "/tmp/uploads/" + file.Filename
	err = ctx.SaveUploadedFile(file, destination)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to save uploaded file", nil, nil)
		return
	}

	// Upload file to IPFS using the IPFS service
	// _, err = uc.ipfsService.UploadFile(destination)
	// if err != nil {
	// 	helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to upload file to IPFS", nil, nil)
	// 	return
	// }

	// Add file metadata to the database using the FS service

	// Update user limits

	helper.FormatResponse(ctx, "success", http.StatusOK, "file uploaded successfully", gin.H{
		"filename": fileName,
		"size":     fileSize,
		"type":     fileType,
	}, nil)
}
