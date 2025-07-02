package files

import (
	"bufio"
	"crypto/cipher"
	"encoding/base64"
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

	// Extract file ID from the context
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
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "file hash not found", nil, nil)
		return
	}

	// Check if file key is passed from the query parameters
	var fileIDFromKey string
	fileKey := ctx.Query("key")
	if fileKey != "" {
		// Check if the file ID is cached
		fileIDFromKey, _ = dc.fsService.GetTemporaryFileKey(ctx, fileKey)
	}

	// If the fileID is NOT found in the cache, perform access checks
	if fileIDFromKey == "" {
		// Access-restricted logic
		switch file.Access {
		case fileConfig.PrivateAccess, fileConfig.EmailAccess:
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "file not found", nil, nil)
			return
		case fileConfig.PasswordAccess:
			// If the file is password-protected, check for the password in the header
			password := ctx.GetHeader("password")
			if password == "" {
				helper.FormatResponse(ctx, "error", http.StatusBadRequest, "password is required for password-protected access", nil, nil)
				return
			}

			// Read file access by file ID to verify the password
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
			if hashedPassword == nil || !helper.CheckHash(password, *hashedPassword) {
				helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid password", nil, nil)
				return
			}
		}
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
	stream := ctx.DefaultQuery("stream", "false")
	if stream == "true" && isEncrypted {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "stream mode is not supported for encrypted downloads", nil, nil)
		return
	}

	if stream == "true" {
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

	// Non-stream mode: download and send full file, decrypt if needed
	fileContent, err := dc.ipfsService.DownloadFile(fileHash)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "error downloading file from IPFS", nil, nil)
		return
	}
	if isEncrypted && ctx.GetHeader("passphrase") != "" {
		plaintext, decErr := aesGCM.Open(nil, nonceBytes, fileContent, nil)
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
	ctx.Header("Content-Length", strconv.Itoa(len(fileContent)))
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	ctx.Data(http.StatusOK, file.ContentType, fileContent)
}
