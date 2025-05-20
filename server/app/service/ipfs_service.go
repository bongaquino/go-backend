package service

import (
	"io"
	"koneksi/server/app/provider"
)

// PeerDetails represents detailed information about a peer
type PeerDetails struct {
	ID        string   `json:"id"`
	Addresses []string `json:"addresses"`
}

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

// GetSwarmPeers fetches the number of peers and their details from the IPFS provider
func (s *IPFSService) GetSwarmPeers() (int, []PeerDetails, error) {
	numPeers, addrs, err := s.ipfsProvider.GetSwarmAddrsDetailed()
	if err != nil {
		return 0, nil, err
	}

	// Convert the map to a slice of PeerDetails
	var peers []PeerDetails
	for id, addresses := range addrs {
		peers = append(peers, PeerDetails{
			ID:        id,
			Addresses: addresses,
		})
	}

	return numPeers, peers, nil
}

// UploadFile uploads a file to IPFS and pins it
func (s *IPFSService) UploadFile(filename string, reader io.Reader) (string, error) {
	return s.ipfsProvider.Pin(filename, reader)
}

// GetFileURL returns the public URL to access a pinned file using its IPFS hash
func (s *IPFSService) GetFileURL(hash string) string {
	return s.ipfsProvider.GetFileURL(hash)
}
