package main

import (
	_ "koneksi/services/iam/docs"
	"koneksi/services/iam/start"
)

// @title Koneksi IAM Service
// @version 1.0
// @description This is the API documentation for Koneksi IAM Service.

// @host localhost:8080
// @BasePath /

func main() {
	// Initialize the server
	start.InitializeKernel()
}
