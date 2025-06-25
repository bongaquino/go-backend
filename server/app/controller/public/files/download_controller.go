package files

import (
	"bufio"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"koneksi/server/config"
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
	// Load file configuration
	fileConfig := config.LoadFileConfig()

	// Extract user ID from the context
	fileID := ctx.Param("fileID")
	if fileID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
		return
	}
	if _, err := primitive.ObjectIDFromHex(fileID); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
		return
	}

	// Read file by ID
	file, err := dc.fsService.ReadFileByID(ctx, fileID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "error reading file"
		if err.Error() == "file not found" {
			status = http.StatusNotFound
			message = "file not found"
		}
		helper.FormatResponse(ctx, "error", status, message, nil, nil)
		return
	}

	// Get file hash
	fileHash := file.Hash
	if fileHash == "" {
		helper.FormatResponse(ctx, "error", http.StatusNotFound, "file hash not found", nil, nil)
		return
	}

	// Get file access
	fileAccess := file.Access
	switch fileAccess {
	// If file is private, no access is allowed
	case fileConfig.PrivateAccess:
		helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
		return
	// If file is temporary, validate the file key
	case fileConfig.TemporaryAccess:
		fileKey := ctx.Query("key")
		if fileKey == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file key is required for temporary access", nil, nil)
			return
		}
		fileIDFromKey, err := dc.fsService.GetTemporaryFileKey(ctx, fileKey)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid or expired file key", nil, nil)
			return
		}
		if fileIDFromKey != fileID {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
	// If file is password-protected, validate the password
	case fileConfig.PasswordAccess:
		password := ctx.Query("password")
		if password == "" {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "password is required for password-protected access", nil, nil)
			return
		}
		// Check the password against stored hash
		fileAccess, err := dc.fsService.ReadFileAccessByFileID(ctx, fileID)
		if err != nil {
			if err.Error() == "file access not found" {
				helper.FormatResponse(ctx, "error", http.StatusNotFound, "file access not found", nil, nil)
				return
			}
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to read file access", nil, nil)
			return
		}
		hashedPassword := fileAccess.Password
		isHashValid := helper.CheckHash(password, hashedPassword)
		if !isHashValid {
			helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid password", nil, nil)
			return
		}
	}
	if ctx.DefaultQuery("stream", "false") == "true" {
		url := dc.ipfsService.GetFileURL(fileHash)
		resp, err := dc.ipfsService.GetHTTPClient().Get(url)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error streaming file from IPFS", nil, nil)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "unexpected status code from IPFS node", nil, nil)
			return
		}

		// Set headers, but skip Content-Length to enable chunked transfer
		ctx.Header("Content-Disposition", "attachment; filename="+file.Name)
		ctx.Header("Content-Type", file.ContentType)
		ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "0")

		reader := bufio.NewReader(resp.Body)
		ctx.Status(http.StatusOK)

		buf := make([]byte, 32*1024) // 32KB buffer
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				if _, writeErr := ctx.Writer.Write(buf[:n]); writeErr != nil {
					break
				}
				ctx.Writer.Flush()
			}
			if err != nil {
				break
			}
		}
		return
	}

	// Non-stream mode: download and send full file
	fileContent, err := dc.ipfsService.DownloadFile(fileHash)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error downloading file from IPFS", nil, nil)
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename="+file.Name)
	ctx.Header("Content-Length", strconv.Itoa(len(fileContent)))
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	ctx.Data(http.StatusOK, file.ContentType, fileContent)
}
