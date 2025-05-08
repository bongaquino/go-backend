package constants

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// FetchController handles health-related endpoints
type FetchController struct {
	userService *service.UserService
	orgService  *service.OrganizationService
}

// NewFetchController initializes a new FetchController
func NewFetchController(userService *service.UserService, orgService *service.OrganizationService) *FetchController {
	return &FetchController{
		userService: userService,
		orgService:  orgService,
	}
}

// Handles the health check endpoint
func (fc *FetchController) Handle(ctx *gin.Context) {

	roles, err := fc.userService.ListRoles(ctx)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	policies, err := fc.orgService.ListPolicies(ctx)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"roles":       roles,
		"policies":    policies,
		"permissions": "permissions",
	}, nil)
}
