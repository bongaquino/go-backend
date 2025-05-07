package users

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateController struct {
	userService *service.UserService
}

// NewUpdateController initializes a new UpdateController
func NewUpdateController(userService *service.UserService) *UpdateController {
	return &UpdateController{
		userService: userService,
	}
}

// Handle handles the health check request
func (lc *UpdateController) Handle(ctx *gin.Context) {
	// Get userID from path parameters
	userID := ctx.Param("userID")
	if userID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "userID is required", nil, nil)
		return
	}

	user, profile, err := lc.userService.GetUserProfile(ctx.Request.Context(), userID)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to fetch user", nil, err)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"user":    user,
		"profile": profile,
	}, nil)
}
