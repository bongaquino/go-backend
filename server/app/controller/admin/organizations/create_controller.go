package organizations

import (
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateController struct {
	orgService *service.OrganizationService
}

// NewCreateController initializes a new CreateController
func NewCreateController(orgService *service.OrganizationService) *CreateController {
	return &CreateController{
		orgService: orgService,
	}
}

// Handle handles the health check request
func (cc *CreateController) Handle(ctx *gin.Context) {
	var request dto.CreateOrgDTO
	// Bind the request body to the CreateOrgDTO struct
	if err := cc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Create the organization using the service
	if org, err := cc.orgService.CreateOrg(ctx, &request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	} else {
		// Respond with success and include the org
		helper.FormatResponse(ctx, "success", http.StatusOK, "organization created successfully", org, nil)
	}
}

func (cc *CreateController) validatePayload(ctx *gin.Context, request *dto.CreateOrgDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
