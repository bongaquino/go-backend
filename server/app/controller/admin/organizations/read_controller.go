package organizations

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReadController struct {
	orgService *service.OrganizationService
}

// NewReadController initializes a new ReadController
func NewReadController(orgService *service.OrganizationService) *ReadController {
	return &ReadController{
		orgService: orgService,
	}
}

// Handle handles the health check request
func (lc *ReadController) Handle(ctx *gin.Context) {
	// Get orgID from path parameters
	orgID := ctx.Param("orgID")
	if orgID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "orgID is required", nil, nil)
		return
	}

	// Read the organization using the service
	org, err := lc.orgService.ReadOrg(ctx, orgID)

	// If err is not found, return a 404 error
	if err != nil {
		if err.Error() == "organization not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "organization not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to fetch organization", nil, err)
		return
	}

	helper.FormatResponse(ctx, "success", http.StatusOK, "organization updated successfully", org, nil)
}
