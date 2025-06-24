package directories

import (
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
}

// NewCreateController initializes a new CreateController
func NewCreateController(fsService *service.FSService, ipfsService *service.IPFSService) *CreateController {
	return &CreateController{
		fsService:   fsService,
		ipfsService: ipfsService,
	}
}

func (cc *CreateController) Handle(ctx *gin.Context) {
	// Validate the request payload
	var request dto.CreateDirectoryDTO
	if err := cc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Limit directory name length to 255 characters
	isTrimmed := false
	if request.Name != "" && len(request.Name) > 255 {
		request.Name = request.Name[:255]
		isTrimmed = true
	}

	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	// Convert userID to ObjectID
	userIDStr := userID.(string)
	objectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid user ID format", nil, nil)
		return
	}

	// Convert directoryID to ObjectID if provided
	var dirObjectID *primitive.ObjectID
	if request.DirectoryID != "" {
		tmpID, err := primitive.ObjectIDFromHex(request.DirectoryID)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid directory ID format", nil, nil)
			return
		}
		dirObjectID = &tmpID
	} else {
		rootDir, _, _, err := cc.fsService.ReadRootDirectory(ctx, userID.(string))
        if err != nil {
            helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to get root directory", nil, nil)
            return
        }
		dirObjectID = &rootDir.ID
	}

	// Build directory model
	directory := &model.Directory{
		UserID:      objectID,
		DirectoryID: dirObjectID,
		Name:        request.Name,
		Size:        0,
		IsDeleted:   false,
	}

	// Create the directory using the fsService
	if err := cc.fsService.CreateDirectory(ctx, directory); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to create directory", directory, nil)
		return
	}

	// Prepare the response data
	response := gin.H{
		"directory": gin.H{
			"id":        directory.ID.Hex(),
			"name":      directory.Name,
			"size":      directory.Size,
			"createdAt": directory.CreatedAt,
			"updatedAt": directory.UpdatedAt,
		},
	}

	if isTrimmed {
		meta := map[string]interface{}{
			"is_trimmed": true,
		}
		helper.FormatResponse(ctx, "success", http.StatusOK, "directory created successfully", response, meta)
		return
	}

	// Return success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "directory created successfully", response, nil)
}

func (cc *CreateController) validatePayload(ctx *gin.Context, request *dto.CreateDirectoryDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	return nil
}
