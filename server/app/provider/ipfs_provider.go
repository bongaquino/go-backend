package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"koneksi/server/config"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// IPFSProvider handles interactions with the IPFS API
type IPFSProvider struct {
	nodeURL     string
	downloadURL string
	client      *http.Client
}

// NewIPFSProvider initializes a new IPFSProvider
func NewIPFSProvider() *IPFSProvider {
	ipfsConfig := config.LoadIPFSConfig()

	return &IPFSProvider{
		nodeURL:     ipfsConfig.IPFSNodeURL,
		downloadURL: ipfsConfig.IPFSDownloadURL,
		client: &http.Client{
			Timeout: 0,
		},
	}
}

// GetSwarmAddrsDetailed calls the IPFS API to get swarm addresses and returns the number of peers and their details
func (p *IPFSProvider) GetSwarmAddrsDetailed() (int, map[string][]string, error) {
	url := fmt.Sprintf("%s/api/v0/swarm/addrs", p.nodeURL)

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

// Pin uploads a file to IPFS and pins it
func (p *IPFSProvider) Pin(filename string, file io.Reader) (string, error) {
	// Build the URL for the IPFS API
	url := fmt.Sprintf("%s/api/v0/add?pin=true", p.nodeURL)

	// Use io.Pipe to stream the multipart form data
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Write the multipart form in a goroutine
	go func() {
		defer pw.Close()
		defer writer.Close()
		part, err := writer.CreateFormFile("file", filepath.Base(filename))
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to create form file: %w", err))
			return
		}
		if _, err = io.Copy(part, file); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to copy file content: %w", err))
			return
		}
	}()

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := p.client.Do(req)
	fmt.Println("resp", resp)
	if err != nil {
		return "", fmt.Errorf("failed to call IPFS API: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body
	var result struct {
		Hash string `json:"Hash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Hash == "" {
		return "", fmt.Errorf("empty hash in response")
	}

	return result.Hash, nil
}

// GetFileURL returns the public URL to access a pinned file using its IPFS hash
func (p *IPFSProvider) GetFileURL(hash string) string {
	return fmt.Sprintf("%s/ipfs/%s", p.downloadURL, hash)
}

func (p *IPFSProvider) GetInternalNodeURL() string {
	return p.nodeURL
}

func (p *IPFSProvider) Client() *http.Client {
	return p.client
}
