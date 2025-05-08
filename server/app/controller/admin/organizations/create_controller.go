package organizations

import (
	"fmt"
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
func (lc *CreateController) Handle(ctx *gin.Context) {
	var request dto.CreateOrgDTO
	// Bind the request body to the CreateOrgDTO struct
	if err := lc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Create the organization using the service
	if org, err := lc.orgService.CreateOrg(ctx, &request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to create organization", nil, nil)
		return
	} else {
		// Respond with success and include the org
		helper.FormatResponse(ctx, "success", http.StatusOK, nil, org, nil)
	}
}

func (rc *CreateController) validatePayload(ctx *gin.Context, request *dto.CreateOrgDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		fmt.Println("error", err)
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
