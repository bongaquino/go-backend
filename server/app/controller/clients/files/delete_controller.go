package files

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteController struct {
	ipfsService *service.IPFSService
}

// NewDeleteController initializes a new DeleteController
func NewDeleteController(ipfsService *service.IPFSService) *DeleteController {
	return &DeleteController{
		ipfsService: ipfsService,
	}
}

func (rc *DeleteController) Handle(ctx *gin.Context) {
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
