package members

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Request structure for adding an organization
type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

type AddController struct {
	orgService *service.OrganizationService
}

// NewAddController initializes a new AddController
func NewAddController(orgService *service.OrganizationService) *AddController {
	return &AddController{
		orgService: orgService,
	}
}

// Handle handles the health check request
func (ac *AddController) Handle(ctx *gin.Context) {
	// Get orgID from path parameters
	orgID := ctx.Param("orgID")
	if orgID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "orgID is required", nil, nil)
		return
	}

	// Bind the request body to the AddMemberRequest struct
	var request AddMemberRequest
	if err := ac.validatePayload(ctx, &request); err != nil {
		return
	}

	// Add the user to the organization using the service
	if err := ac.orgService.AddMember(ctx, orgID, request.UserID, request.RoleID); err != nil {
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
		if err.Error() == "user is already a member" {
			helper.FormatResponse(ctx, "error", http.StatusConflict, "user is already a member", nil, nil)
			return
		}
		if err.Error() == "user is already a member of another organization" {
			helper.FormatResponse(ctx, "error", http.StatusConflict, "user is already a member of another organization", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to add member", nil, err)
		return
	}

	// Return the response
	helper.FormatResponse(ctx, "success", http.StatusOK, "user added successfully", gin.H{
		"org_id":  orgID,
		"user_id": request.UserID,
		"role_id": request.RoleID,
	}, nil)
}

func (rc *AddController) validatePayload(ctx *gin.Context, request *AddMemberRequest) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
