package main

import (
	_ "argo/docs"
	"argo/start"
)

// @title Argo API
// @version 1.0
// @description This is the API documentation for Argo.

// @host localhost:8080
// @BasePath /

func main() {
	// Initialize the server
	start.InitializeKernel()
}
