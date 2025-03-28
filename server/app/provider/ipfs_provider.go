package provider

import (
	"encoding/json"
	"fmt"
	"koneksi/server/config"
	"net/http"
	"time"
)

// IPFSProvider handles interactions with the IPFS API
type IPFSProvider struct {
	baseURL string
	client  *http.Client
}

// NewIPFSProvider initializes a new IPFSProvider
func NewIPFSProvider() *IPFSProvider {
	ipfsConfig := config.LoadIPFSConfig()

	return &IPFSProvider{
		baseURL: ipfsConfig.IpfsNodeURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSwarmAddrsDetailed calls the IPFS API to get swarm addresses and returns the number of peers and their details
func (p *IPFSProvider) GetSwarmAddrsDetailed() (int, map[string][]string, error) {
	url := fmt.Sprintf("%s/api/v0/swarm/addrs", p.baseURL)

	// Make the HTTP request
	resp, err := p.client.Post(url, "application/json", nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to call IPFS API: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return 0, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body
	var result struct {
		Addrs map[string][]string `json:"Addrs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Count the number of peers
	numPeers := len(result.Addrs)
	return numPeers, result.Addrs, nil
}
