package organizations

import (
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateController struct {
	orgService *service.OrganizationService
}

// NewUpdateController initializes a new UpdateController
func NewUpdateController(orgService *service.OrganizationService) *UpdateController {
	return &UpdateController{
		orgService: orgService,
	}
}

// Handle handles the health check request
func (uc *UpdateController) Handle(ctx *gin.Context) {
	// Get orgID from path parameters
	orgID := ctx.Param("orgID")
	if orgID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "orgID is required", nil, nil)
		return
	}

	// Get request body
	var request dto.UpdateOrgDTO
	if err := uc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Update the organization using the service
	updatedOrg, err := uc.orgService.UpdateOrg(ctx, orgID, &request)
	if err != nil {
		if err.Error() == "organization not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.FormatResponse(ctx, "success", http.StatusOK, "organization updated successfully", updatedOrg, nil)
}

func (uc *UpdateController) validatePayload(ctx *gin.Context, request *dto.UpdateOrgDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	return nil
}
