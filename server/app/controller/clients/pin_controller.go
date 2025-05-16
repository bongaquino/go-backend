package clients

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PinController struct {
	ipfsService *service.IPFSService
}

// NewPinController initializes a new PinController
func NewPinController(ipfsService *service.IPFSService) *PinController {
	return &PinController{
		ipfsService: ipfsService,
	}
}

func (pc *PinController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	// userID, exists := ctx.Get("userID")
	// if !exists {
	// 	helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
	// 	return
	// }

	/// Fetch the number of peers and their details from the IPFS service
	numPeers, peers, err := pc.ipfsService.GetSwarmPeers()
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with the number of peers and their details
	helper.FormatResponse(ctx, "success", http.StatusOK, "swarm addresses fetched successfully", gin.H{
		"num_peers": numPeers,
		"peers":     peers,
	}, nil)
}
