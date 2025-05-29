package users

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReadController struct {
	userService *service.UserService
	orgService  *service.OrganizationService
}

// NewReadController initializes a new ReadController
func NewReadController(userService *service.UserService, orgService *service.OrganizationService) *ReadController {
	return &ReadController{
		userService: userService,
		orgService:  orgService,
	}
}

// Handle handles the health check request
func (rc *ReadController) Handle(ctx *gin.Context) {
	// Get userID from path parameters
	userID := ctx.Param("userID")
	if userID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "userID is required", nil, nil)
		return
	}

	user, profile, role, limit, err := rc.userService.GetUserProfile(ctx.Request.Context(), userID)

	// if err is not found, return a 404 error
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "profile not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "user not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to fetch user", nil, err)
	}

	org, _ := rc.orgService.GetOrganizationByUserID(ctx.Request.Context(), userID)

	// Exclude sensitive fields from the response
	user.Password = "REDACTED"
	user.OtpSecret = "REDACTED"

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"user":    user,
		"profile": profile,
		"role":    role,
		"limit":   limit,
		"org":     org,
	}, nil)
}
