package directories

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReadController struct {
	fsService   *service.FSService
	ipfsService *service.IPFSService
}

// NewReadController initializes a new ReadController
func NewReadController(fsService *service.FSService, ipfsService *service.IPFSService) *ReadController {
	return &ReadController{
		fsService:   fsService,
		ipfsService: ipfsService,
	}
}

func (rc *ReadController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	// userID, exists := ctx.Get("userID")
	// if !exists {
	// 	helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
	// 	return
	// }

	numPeers, peers, err := rc.ipfsService.GetSwarmPeers()
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.FormatResponse(ctx, "success", http.StatusOK, "peers fetched successfully", gin.H{
		"count": numPeers,
		"peers": peers,
	}, nil)
}
