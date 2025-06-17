package files

import (
	"bufio"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
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
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "user ID not found in context", nil, nil)
		return
	}

	fileID := ctx.Param("fileID")
	if fileID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file ID is required", nil, nil)
		return
	}
	if _, err := primitive.ObjectIDFromHex(fileID); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid file ID format", nil, nil)
		return
	}

	file, err := dc.fsService.ReadFileByIDUserID(ctx, fileID, userID.(string))
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

	fileHash := file.Hash
	if fileHash == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file hash is required for download", nil, nil)
		return
	}

	// Change default value of stream query param to "true"
	if ctx.DefaultQuery("stream", "true") == "true" {
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
