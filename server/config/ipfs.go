package config

import "koneksi/server/core/env"

// IPFSConfig holds the IPFS configuration
type IPFSConfig struct {
	IpfsNodeURL string
}

func LoadIPFSConfig() *IPFSConfig {
	// Load environment variables
	envVars := env.LoadEnv()

	// Create the configuration from environment variables
	return &IPFSConfig{
		IpfsNodeURL: envVars.IpfsNodeURL,
	}
}
