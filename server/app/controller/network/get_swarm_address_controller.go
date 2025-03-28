package network

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// GetSwarmAddressController handles fetching swarm addresses from the IPFS network
type GetSwarmAddressController struct {
	ipfsService *service.IPFSService
}

// NewGetSwarmAddressController initializes a new GetSwarmAddressController
func NewGetSwarmAddressController(ipfsService *service.IPFSService) *GetSwarmAddressController {
	return &GetSwarmAddressController{
		ipfsService: ipfsService,
	}
}

// Handle processes the request to fetch swarm addresses
func (gsc *GetSwarmAddressController) Handle(c *gin.Context) {
	// Fetch the number of peers from the IPFS service
	numPeers, err := gsc.ipfsService.GetSwarmPeers()
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with the number of peers
	helper.FormatResponse(c, "success", http.StatusOK, "Swarm addresses fetched successfully", gin.H{
		"num_peers": numPeers,
	}, nil)
}
