package peers

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FetchController struct {
	ipfsService *service.IPFSService
}

// NewFetchController initializes a new FetchController
func NewFetchController(ipfsService *service.IPFSService) *FetchController {
	return &FetchController{
		ipfsService: ipfsService,
	}
}

func (pc *FetchController) Handle(ctx *gin.Context) {
	/// Fetch the number of peers and their details from the IPFS service
	numPeers, peers, err := pc.ipfsService.GetSwarmPeers()
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
