package service

import (
	"koneksi/server/app/provider"
)

// IPFSService handles business logic related to IPFS
type IPFSService struct {
	ipfsProvider *provider.IPFSProvider
}

// NewIPFSService initializes a new IPFSService
func NewIPFSService(ipfsProvider *provider.IPFSProvider) *IPFSService {
	return &IPFSService{
		ipfsProvider: ipfsProvider,
	}
}

// GetSwarmPeers fetches the number of peers from the IPFS provider
func (s *IPFSService) GetSwarmPeers() (int, error) {
	return s.ipfsProvider.GetSwarmAddrs()
}
