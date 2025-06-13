package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Extract directory ID from the query parameters
	directoryID := ctx.Query("directory_id")
	if directoryID == ":directory" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "directory ID is required", nil, nil)
		return
	}

	// Check if user has access to the directory
	isOwner, err := uc.fsService.CheckDirectoryOwnership(ctx, directoryID, userID.(string))
	if err != nil {
		if err.Error() == "directory not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to check directory ownership", nil, nil)
		return
	}
	if !isOwner {
		helper.FormatResponse(ctx, "error", http.StatusForbidden, "access denied to directory", nil, nil)
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

	// Stream the file and upload directly to IPFS
	uploadedFile, err := file.Open()
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to open uploaded file", nil, nil)
		return
	}
	defer uploadedFile.Close()
	cid, err := uc.ipfsService.UploadFile(file.Filename, uploadedFile)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to upload file to IPFS", nil, nil)
		return
	}

	// Convert userID string to primitive.ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "invalid user ID format", nil, nil)
		return
	}

	// Convert directoryID string to primitive.ObjectID
	dirObjID, err := primitive.ObjectIDFromHex(directoryID)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "invalid directory ID format", nil, nil)
		return
	}

	// Build the file model
	fileModel := &model.File{
		UserID:      userObjID,
		DirectoryID: &dirObjID,
		Name:        fileName,
		Hash:        cid,
		Size:        fileSize,
		ContentType: fileType,
		IsDeleted:   false,
	}

	// Save the file metadata to the database
	err = uc.fsService.CreateFile(ctx, fileModel)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to save file metadata", nil, nil)
		return
	}

	// Compute the new usage
	newUsage := userLimit.BytesUsage + fileSize

	// Update user usage
	err = uc.userService.UpdateUserUsage(ctx, userID.(string), newUsage)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update user usage", nil, nil)
		return
	}

	// Return response
	helper.FormatResponse(ctx, "success", http.StatusOK, "file uploaded successfully", gin.H{
		"directory_id": directoryID,
		"file_id":      fileModel.ID.Hex(),
		"name":         fileName,
		"hash":         cid,
		"size":         fileSize,
		"content_type": fileType,
	}, nil)
}
