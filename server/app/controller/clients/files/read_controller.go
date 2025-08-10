package files

import (
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"
	"bongaquino/server/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReadController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
	userService *service.UserService
}

// NewReadController initializes a new ReadController
func NewReadController(
	fsService *service.FSService,
	ipfsService *service.IPFSService,
	userService *service.UserService,
) *ReadController {
	return &ReadController{
		fsService:   fsService,
		ipfsService: ipfsService,
		userService: userService,
	}
}

func (rc *ReadController) Handle(ctx *gin.Context) {
	// Load file configuration
	fileConfig := config.LoadFileConfig()

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

	// Check for optional include_chunks query param (default: false)
	includeChunks := false
	if val, ok := ctx.GetQuery("include_chunks"); ok && (val == "true") {
		includeChunks = true
	}

	var chunks any = nil
	if includeChunks {
		// List file chunks
		var err error
		chunks, err = rc.ipfsService.ListFileChunks(file.Hash)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error listing file chunks", nil, nil)
			return
		}
	}

	// Prepare the file details
	fileData := gin.H{
		"id":           file.ID.Hex(),
		"directory_id": file.DirectoryID.Hex(),
		"name":         file.Name,
		"size":         file.Size,
		"hash":         file.Hash,
		"content_type": file.ContentType,
		"chunks":       chunks,
		"access":       file.Access,
		"is_encrypted": file.IsEncrypted,
		"recipients":   nil,
		"created_at":   file.CreatedAt,
		"updated_at":   file.UpdatedAt,
	}
	// Check if file.access is "email"
	if file.Access == fileConfig.EmailAccess {
		// Fetch file access details from fsService
		fileAccessList, _ := rc.fsService.ListFileAccessByFileID(ctx, file.ID.Hex())
		// Add file access details to the file data
		fileData["recipients"] = make([]gin.H, len(fileAccessList))
		// Loop thru file access list to get recipient IDs then use recipient ID to get user email
		for j, fileAccess := range fileAccessList {
			// Get user email from fsService using recipient ID
			user, profile, _, _, _, _ := rc.userService.GetUserInfo(ctx, fileAccess.RecipientID.Hex())
			fileData["recipients"].([]gin.H)[j] = gin.H{
				"id":          user.ID.Hex(),
				"email":       user.Email,
				"first_name":  profile.FirstName,
				"middle_name": profile.MiddleName,
				"last_name":   profile.LastName,
				"suffix":      profile.Suffix,
			}
		}
	}

	// Return the file details
	helper.FormatResponse(ctx, "success", http.StatusOK, "file read successfully", fileData, nil)
}
