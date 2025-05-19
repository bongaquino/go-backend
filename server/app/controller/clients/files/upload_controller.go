package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	ipfsService *service.IPFSService
}

// NewUploadController initializes a new UploadController
func NewUploadController(ipfsService *service.IPFSService) *UploadController {
	return &UploadController{
		ipfsService: ipfsService,
	}
}

func (uc *UploadController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	// userID, exists := ctx.Get("userID")
	// if !exists {
	// 	helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
	// 	return
	// }

	/// Upload the number of peers and their details from the IPFS service
	numPeers, peers, err := uc.ipfsService.GetSwarmPeers()
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with the number of peers and their details
	helper.FormatResponse(ctx, "success", http.StatusOK, "peers fetched successfully", gin.H{
		"count": numPeers,
		"peers": peers,
	}, nil)
}
