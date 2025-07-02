package files

import (
	"bufio"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
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

	file, err := dc.fsService.ReadFileByID(ctx, fileID)

	if err != nil {
		if err.Error() == "file not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error reading file", nil, nil)
		return
	}

	// Check if the user has access to the file
	isValid := dc.fsService.ValidateFileAccess(ctx, fileID, userID.(string))
	if !isValid {
		helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
		return
	}

	// Get file hash
	fileHash := file.Hash
	if fileHash == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "file hash is required for download", nil, nil)
		return
	}

	isEncrypted := file.IsEncrypted
	var aesGCM cipher.AEAD
	var nonceBytes []byte
	var decryptedSalt, decryptedNonce string
	if isEncrypted {
		passphrase := ctx.GetHeader("passphrase")
		if passphrase != "" {
			// Decrypt the salt if the file is encrypted
			var err error
			decryptedSalt, err = helper.Decrypt(file.Salt)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to decrypt salt", nil, nil)
				return
			}
			// Decrypt the nonce if the file is encrypted
			decryptedNonce, err = helper.Decrypt(file.Nonce)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to decrypt nonce", nil, nil)
				return
			}
			keyBytes, err := helper.DeriveKey(passphrase, decryptedSalt)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to derive key", nil, nil)
				return
			}
			aesGCM, err = helper.CreateAesGcmCipher(keyBytes)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to create AES-GCM cipher", nil, nil)
				return
			}
			nonceBytes, err = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(decryptedNonce)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "invalid nonce encoding", nil, nil)
				return
			}
		}
	}

	// Check if stream mode is enabled and file is encrypted
	stream := ctx.DefaultQuery("stream", "true")
	if stream == "true" && isEncrypted {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "stream mode is not supported for encrypted downloads", nil, nil)
		return
	}

	// Change default value of stream query param to "true"
	if ctx.DefaultQuery("stream", "true") == "true" {
		// Stream mode: stream file content from IPFS (only for unencrypted files)
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
	// For encrypted files or non-stream mode, download the full file, decrypt if needed, then send
	fileContent, err := dc.ipfsService.DownloadFile(fileHash)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error downloading file from IPFS", nil, nil)
		return
	}
	if isEncrypted && ctx.GetHeader("passphrase") != "" {
		plaintext, decErr := aesGCM.Open(nil, nonceBytes, fileContent, nil)
		fmt.Println(decErr)
		if decErr != nil {
			helper.FormatResponse(ctx, "error", http.StatusForbidden, "failed to decrypt file", nil, nil)
			return
		}
		fileContent = plaintext
	}
	if isEncrypted && ctx.GetHeader("passphrase") == "" {
		ctx.Header("Content-Disposition", "attachment; filename="+file.Name)
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Length", strconv.Itoa(len(fileContent)))
		ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "0")
		ctx.Data(http.StatusOK, "application/octet-stream", fileContent)
		return
	}
	ctx.Header("Content-Disposition", "attachment; filename="+file.Name)
	ctx.Header("Content-Type", file.ContentType)
	ctx.Header("Content-Length", strconv.Itoa(len(fileContent)))
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	ctx.Data(http.StatusOK, file.ContentType, fileContent)
}
