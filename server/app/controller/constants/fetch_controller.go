package constants

import (
	"net/http"

	"koneksi/server/app/helper"

	"github.com/gin-gonic/gin"
)

// FetchController handles health-related endpoints
type FetchController struct{}

// NewFetchController initializes a new FetchController
func NewFetchController() *FetchController {
	return &FetchController{}
}

// Handles the health check endpoint
func (hc *FetchController) Handle(ctx *gin.Context) {
	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"roles":       "roles",
		"policies":    "policies",
		"permissions": "permissions",
	}, nil)
}
