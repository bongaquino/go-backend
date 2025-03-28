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

// NewIPFSProvider initializes a new IPFSProvider using the IPFSConfig
func NewIPFSProvider() *IPFSProvider {
	ipfsConfig := config.LoadIPFSConfig()

	return &IPFSProvider{
		baseURL: ipfsConfig.IpfsNodeURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSwarmAddrs calls the IPFS API to get swarm addresses and returns the number of peers
func (p *IPFSProvider) GetSwarmAddrs() (int, error) {
	url := fmt.Sprintf("%s/api/v0/swarm/addrs", p.baseURL)

	// Make the HTTP request
	resp, err := p.client.Post(url, "application/json", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to call IPFS API: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body
	var result struct {
		Addrs map[string][]string `json:"Addrs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Count the number of peers
	numPeers := len(result.Addrs)
	return numPeers, nil
}
