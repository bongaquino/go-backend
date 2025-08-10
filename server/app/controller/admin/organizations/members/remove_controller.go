package members

import (
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RemoveController struct {
	orgService *service.OrganizationService
}

// NewRemoveController initializes a new RemoveController
func NewRemoveController(orgService *service.OrganizationService) *RemoveController {
	return &RemoveController{
		orgService: orgService,
	}
}

// Handle handles the health check request
func (ac *RemoveController) Handle(ctx *gin.Context) {
	// Get orgID from path parameters
	orgID := ctx.Param("orgID")
	if orgID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "orgID is required", nil, nil)
		return
	}

	// Get userID from path parameters
	userID := ctx.Param("userID")
	if userID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "userID is required", nil, nil)
		return
	}

	// Remove the user to the organization using the service
	if err := ac.orgService.RemoveMember(ctx, orgID, userID); err != nil {
		if err.Error() == "organization not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "organization not found", nil, nil)
			return
		}
		if err.Error() == "user not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "user not found", nil, nil)
			return
		}
		if err.Error() == "role not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "role not found", nil, nil)
			return
		}
		if err.Error() == "user is not a member" {
			helper.FormatResponse(ctx, "error", http.StatusConflict, "user is not a member", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update member role", nil, err)
		return
	}

	// Return the response
	helper.FormatResponse(ctx, "success", http.StatusOK, "member removed successfully", nil, nil)
}
