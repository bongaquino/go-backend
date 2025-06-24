package files

import (
	"bytes"
	"io"
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/service"
	"net/http"
	"path/filepath"
	"strings"

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

	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid user ID format", nil, nil)
		return
	}

	// Get directory_id from request body
	var request dto.UploadFileDTO
	_ = ctx.ShouldBind(&request)

	// Only use directoryID from request body
	directoryID := request.DirectoryID
	if directoryID == "" {
		// Get the user's root directory
		rootDir, _, _, err := uc.fsService.ReadRootDirectory(ctx, userID.(string))
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to get root directory", nil, nil)
			return
		}
		directoryID = rootDir.ID.Hex()
	}

	dirObjID, err := primitive.ObjectIDFromHex(directoryID)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid directory ID format", nil, nil)
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

	// Limit file name length to 255 characters while preserving extension
	isTrimmed := false
	if fileName != "" && len(fileName) > 255 {
		// Get the file extension
		ext := filepath.Ext(fileName)
		
		// Calculate the base name (filename without extension)
		baseName := strings.TrimSuffix(fileName, ext)
		
		// Calculate how much we need to trim from the base name
		maxBaseNameLength := 255 - len(ext)
		
		if maxBaseNameLength > 0 {
			// Trim the base name if it's too long
			if len(baseName) > maxBaseNameLength {
				baseName = baseName[:maxBaseNameLength]
			}
			// Reconstruct the filename with preserved extension
			fileName = baseName + ext
		} else {
			// If extension itself is >= 255 characters, just truncate the whole name
			fileName = fileName[:255]
		}
		
		isTrimmed = true
	}

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

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to open uploaded file", nil, nil)
		return
	}
	defer src.Close()

	var cid string
	var uploadErr error

	// Check for stream query param to decide upload mode
	stream := ctx.Query("stream")
	if stream == "false" {
		// Non-streaming mode: read the entire file into memory first.
		ctx.Writer.Header().Set("X-Upload-Mode", "non-stream")
		fileBytes, err := io.ReadAll(src)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to read file content for non-streaming upload", nil, nil)
			return
		}
		cid, uploadErr = uc.ipfsService.UploadFile(fileName, bytes.NewReader(fileBytes))
	} else {
		// Streaming mode (default): pass the file stream directly.
		ctx.Writer.Header().Set("X-Upload-Mode", "stream")
		cid, uploadErr = uc.ipfsService.UploadFile(fileName, src)
	}

	if uploadErr != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to upload file to IPFS", nil, nil)
		return
	}

	// Create a new file model
	newFile := &model.File{
		UserID:      userObjID,
		DirectoryID: &dirObjID,
		Name:        fileName,
		Hash:        cid,
		Size:        fileSize,
		ContentType: fileType,
		IsDeleted:   false,
	}

	// Save the file metadata to the database
	err = uc.fsService.CreateFile(ctx, newFile)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to save file metadata", nil, nil)
		return
	}

	err = uc.fsService.RecalculateDirectorySizeAndParents(ctx, directoryID, userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to recalculate directory sizes", nil, nil)
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

	if isTrimmed {
		meta := map[string]interface{}{
			"is_trimmed": true,
		}
		helper.FormatResponse(ctx, "success", http.StatusOK, "file uploaded successfully", gin.H{
			"directory_id": directoryID,
			"file_id":      newFile.ID.Hex(),
			"name":         fileName,
			"hash":         cid,
			"size":         fileSize,
			"content_type": fileType,
		}, meta)
		return
	}

	// Return response
	helper.FormatResponse(ctx, "success", http.StatusOK, "file uploaded successfully", gin.H{
		"directory_id": directoryID,
		"file_id":      newFile.ID.Hex(),
		"name":         fileName,
		"hash":         cid,
		"size":         fileSize,
		"content_type": fileType,
	}, nil)
}
