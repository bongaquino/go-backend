package members

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Request structure for adding an organization
type UpdateMemberRoleRequest struct {
	RoleID string `json:"role_id" binding:"required"`
}

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
func (ac *UpdateController) Handle(ctx *gin.Context) {
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

	// Bind the request body to the UpdateMemberRoleRequest struct
	var request UpdateMemberRoleRequest
	if err := ac.validatePayload(ctx, &request); err != nil {
		return
	}

	// Update the user to the organization using the service
	if err := ac.orgService.UpdateMember(ctx, orgID, userID, request.RoleID); err != nil {
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
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update member role", nil, err)
		return
	}

	// Return the response
	helper.FormatResponse(ctx, "success", http.StatusOK, "member role updated successfully", gin.H{
		"org_id":  orgID,
		"user_id": userID,
		"role_id": request.RoleID,
	}, nil)
}

func (rc *UpdateController) validatePayload(ctx *gin.Context, request *UpdateMemberRoleRequest) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
