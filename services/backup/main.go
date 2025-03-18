package main

import (
	_ "koneksi/services/backup/docs"
	"koneksi/services/backup/start"
)

// @title Koneksi Orchestrator API
// @version 1.0
// @description This is the API documentation for Koneksi Orchestrator.

// @host localhost:8080
// @BasePath /

func main() {
	// Initialize the server
	start.InitializeKernel()
}
