package directories

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/service"
	"koneksi/server/config"
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

	// Get the directory ID from the URL parameters
	directoryID := ctx.Param("directoryID")

	// Check if the directory ID is not empty
	if directoryID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "directory ID is required", nil, nil)
		return
	}

	if directoryID == "root" {
		// Use fsService to read the root directory
		directory, subDirectories, files, err := rc.fsService.ReadRootDirectory(ctx, userID.(string))
		if err != nil {
			if err.Error() == "directory not found" {
				helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
				return
			}
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to read root directory", nil, nil)
			return
		}

		// Ensure subDirectories and files are not nil
		if subDirectories == nil {
			subDirectories = []*model.Directory{}
		}
		if files == nil {
			files = []*model.File{}
		}

		// Format the directory data
		directoryData := gin.H{
			"id":         directory.ID.Hex(),
			"name":       directory.Name,
			"size":       directory.Size,
			"created_at": directory.CreatedAt,
			"updated_at": directory.UpdatedAt,
		}

		// Format the subdirectories
		subDirectoriesData := make([]gin.H, len(subDirectories))
		for i, subDir := range subDirectories {
			subDirectoriesData[i] = gin.H{
				"id":         subDir.ID.Hex(),
				"name":       subDir.Name,
				"size":       subDir.Size,
				"created_at": subDir.CreatedAt,
				"updated_at": subDir.UpdatedAt,
			}
		}

		// Format the files
		filesData := make([]gin.H, len(files))
		for i, file := range files {
			filesData[i] = gin.H{
				"id":           file.ID.Hex(),
				"directory_id": file.DirectoryID.Hex(),
				"name":         file.Name,
				"hash":         file.Hash,
				"size":         file.Size,
				"content_type": file.ContentType,
				"access":       file.Access,
				"recipients":   nil,
				"created_at":   file.CreatedAt,
				"updated_at":   file.UpdatedAt,
			}
			// Check if file.access is "email"
			if file.Access == fileConfig.EmailAccess {
				// Fetch file access details from fsService
				fileAccessList, _ := rc.fsService.ListFileAccessByFileID(ctx, file.ID.Hex())
				// Add file access details to the file data
				filesData[i]["recipients"] = make([]gin.H, len(fileAccessList))
				// Loop thru file access list to get recipient IDs then use recipient ID to get user email
				for j, fileAccess := range fileAccessList {
					// Get user email from fsService using recipient ID
					user, profile, _, _, _, _ := rc.userService.GetUserInfo(ctx, fileAccess.RecipientID.Hex())
					filesData[i]["recipients"].([]gin.H)[j] = gin.H{
						"id":          user.ID.Hex(),
						"email":       user.Email,
						"first_name":  profile.FirstName,
						"middle_name": profile.MiddleName,
						"last_name":   profile.LastName,
						"suffix":      profile.Suffix,
					}
				}
			}
		}

		// Prepare the response
		response := gin.H{
			"directory":      directoryData,
			"subdirectories": subDirectoriesData,
			"files":          filesData,
		}

		// Send the response
		helper.FormatResponse(ctx, "success", http.StatusOK, "directory read successfully", response, nil)
	} else {
		// Check if the directory ID is in valid format
		if directoryID == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "directory ID is required", nil, nil)
			return
		}
		if _, err := primitive.ObjectIDFromHex(directoryID); err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid directory ID format", nil, nil)
			return
		}

		// Use fsService to read the root directory
		directory, subDirectories, files, err := rc.fsService.ReadDirectory(ctx, directoryID, userID.(string))
		if err != nil {
			if err.Error() == "directory not found" {
				helper.FormatResponse(ctx, "error", http.StatusNotFound, "directory not found", nil, nil)
				return
			}
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to read directory", nil, nil)
			return
		}

		// Ensure subDirectories and files are not nil
		if subDirectories == nil {
			subDirectories = []*model.Directory{}
		}
		if files == nil {
			files = []*model.File{}
		}

		// Format the directory data
		directoryData := gin.H{
			"id":         directory.ID.Hex(),
			"name":       directory.Name,
			"size":       directory.Size,
			"created_at": directory.CreatedAt,
			"updated_at": directory.UpdatedAt,
		}

		// Format the subdirectories
		subDirectoriesData := make([]gin.H, len(subDirectories))
		for i, subDir := range subDirectories {
			subDirectoriesData[i] = gin.H{
				"id":         subDir.ID.Hex(),
				"name":       subDir.Name,
				"size":       subDir.Size,
				"created_at": subDir.CreatedAt,
				"updated_at": subDir.UpdatedAt,
			}
		}

		// Format the files
		filesData := make([]gin.H, len(files))
		for i, file := range files {
			filesData[i] = gin.H{
				"id":           file.ID.Hex(),
				"directory_id": file.DirectoryID.Hex(),
				"name":         file.Name,
				"hash":         file.Hash,
				"size":         file.Size,
				"content_type": file.ContentType,
				"access":       file.Access,
				"recipients":   nil,
				"created_at":   file.CreatedAt,
				"updated_at":   file.UpdatedAt,
			}
			// Check if file.access is "email"
			if file.Access == fileConfig.EmailAccess {
				// Fetch file access details from fsService
				fileAccessList, _ := rc.fsService.ListFileAccessByFileID(ctx, file.ID.Hex())
				// Add file access details to the file data
				filesData[i]["recipients"] = make([]gin.H, len(fileAccessList))
				// Loop thru file access list to get recipient IDs then use recipient ID to get user email
				for j, fileAccess := range fileAccessList {
					// Get user email from fsService using recipient ID
					user, profile, _, _, _, _ := rc.userService.GetUserInfo(ctx, fileAccess.RecipientID.Hex())
					filesData[i]["recipients"].([]gin.H)[j] = gin.H{
						"id":          user.ID.Hex(),
						"email":       user.Email,
						"first_name":  profile.FirstName,
						"middle_name": profile.MiddleName,
						"last_name":   profile.LastName,
						"suffix":      profile.Suffix,
					}
				}
			}
		}

		// Prepare the response
		response := gin.H{
			"directory":      directoryData,
			"subdirectories": subDirectoriesData,
			"files":          filesData,
		}

		// Send the response
		helper.FormatResponse(ctx, "success", http.StatusOK, "directory read successfully", response, nil)
	}
}
