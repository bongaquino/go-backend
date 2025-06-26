package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/service"
	"koneksi/server/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShareController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
	userService *service.UserService
}

// NewShareController initializes a new ShareController
func NewShareController(
	fsService *service.FSService,
	ipfsService *service.IPFSService,
	userService *service.UserService,
) *ShareController {
	return &ShareController{
		fsService:   fsService,
		ipfsService: ipfsService,
		userService: userService,
	}
}

func (sc *ShareController) Handle(ctx *gin.Context) {
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

	// Get the access type from the query parameters (default: "private")
	accessType := ctx.DefaultQuery("access", fileConfig.DefaultAccess)

	// Check if access type is valid by checking against the allowed options
	if !helper.Contains(fileConfig.AccessOptions, accessType) {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid access type", nil, nil)
		return
	}

	// Delete all existing file access records for the file ID (if any)
	if err := sc.fsService.DeleteFileAccessByFileID(ctx, fileID); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error deleting existing file access records", nil, nil)
		return
	}

	// Validate the access type and prepare the request body if needed
	var requestBody map[string]any
	var responseBody map[string]any
	switch accessType {
	case fileConfig.PublicAccess, fileConfig.PrivateAccess:
		// No body needed
	case fileConfig.TemporaryAccess:
		// Verify if duration is provided in the request body
		requestBody = make(map[string]any)
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "duration is required for temporary access", nil, nil)
			return
		}
		if _, ok := requestBody["duration"]; !ok || requestBody["duration"] == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "duration is required for temporary access", nil, nil)
			return
		}
		// Generate a temporary file key
		fileKey, err := helper.GenerateFileKey(fileID)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to generate temporary file key", nil, nil)
			return
		}
		// Save the file key and file ID in Redis
		durationVal, ok := requestBody["duration"].(float64)
		if !ok {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "duration must be a number", nil, nil)
			return
		}
		duration := time.Duration(int(durationVal)) * time.Second
		if err := sc.fsService.SetTemporaryFileKey(ctx, fileKey, fileID, duration); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to save temporary file key", nil, nil)
			return
		}
		responseBody = map[string]any{
			"file_key": fileKey,
			"duration": duration.String(),
		}
	case fileConfig.PasswordAccess:
		// Verify if password is provided in the request body
		requestBody = make(map[string]any)
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "password is required for password-protected access", nil, nil)
			return
		}
		if _, ok := requestBody["password"]; !ok || requestBody["password"] == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "password is required for password-protected access", nil, nil)
			return
		}
		// Extract the password from the request body
		password, ok := requestBody["password"].(string)
		if !ok || password == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "password must be a non-empty string", nil, nil)
			return
		}
		// Validate the password
		_, err := helper.ValidatePassword(password)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid password format", nil, nil)
			return
		}
		// Hash the password
		hashedPassword, err := helper.Hash(password)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to hash password", nil, nil)
			return
		}
		// Create the file access record
		fileObjID, err := primitive.ObjectIDFromHex(fileID)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
			return
		}
		ownerObjID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid user ID format", nil, nil)
			return
		}
		fileAccess := &model.FileAccess{
			FileID:      fileObjID,
			OwnerID:     ownerObjID,
			RecipientID: nil,
			Password:    &hashedPassword,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := sc.fsService.CreateFileAccess(ctx, fileAccess); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to create file access", nil, nil)
			return
		}
	case fileConfig.EmailAccess:
		requestBody = make(map[string]any)
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "at least one email is required for email access", nil, nil)
			return
		}
		emails, ok := requestBody["emails"].([]any)
		if !ok || len(emails) == 0 {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "at least one email is required for email access", nil, nil)
			return
		}
		// Check if all emails exist in the user service using UserExists method
		for _, email := range emails {
			exists, err := sc.userService.UserExists(ctx, email.(string))
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to verify emails", nil, nil)
				return
			}
			if !exists {
				helper.FormatResponse(ctx, "error", http.StatusBadRequest, "one or more provided emails are invalid", nil, nil)
				return
			}
		}
		// Create file access records for each email
		fileObjID, err := primitive.ObjectIDFromHex(fileID)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
			return
		}
		ownerObjID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid user ID format", nil, nil)
			return
		}
		for _, email := range emails {
			// Get the user ID for the email
			user, _, err := sc.userService.GetUserProfileByEmail(ctx, email.(string))
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to get user by email", nil, nil)
				return
			}
			// Skip if the user ID is the same as the owner ID
			if user.ID.Hex() == ownerObjID.Hex() {
				continue
			}
			// Create the file access record
			fileAccess := &model.FileAccess{
				FileID:      fileObjID,
				OwnerID:     ownerObjID,
				RecipientID: &user.ID,
				Password:    nil, // No password for email access
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := sc.fsService.CreateFileAccess(ctx, fileAccess); err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to create file access", nil, nil)
				return
			}
		}
	}

	// Check if the file ID is in valid format
	if fileID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
		return
	}
	if _, err := primitive.ObjectIDFromHex(fileID); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
		return
	}

	// Update the file's access type using the FS service
	err := sc.fsService.UpdateFileAccess(ctx, fileID, userID.(string), accessType)
	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error updating file access", nil, nil)
		return
	}

	// Return success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "file shared successfully", responseBody, nil)
}
