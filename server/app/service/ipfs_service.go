package service

import (
	"fmt"
	"io"
	"koneksi/server/app/provider"
	"net/http"
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

// DownloadFile retrieves a file from IPFS by hash and returns its content
func (s *IPFSService) DownloadFile(hash string) ([]byte, error) {
	// Build the internal download URL
	url := s.ipfsProvider.GetFileURL(hash)

	// Perform the HTTP GET request to the IPFS node
	resp, err := s.ipfsProvider.Client().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from IPFS: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from IPFS node: %d", resp.StatusCode)
	}

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
