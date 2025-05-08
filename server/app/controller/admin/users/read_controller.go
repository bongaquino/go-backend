package users

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReadController struct {
	userService *service.UserService
}

// NewReadController initializes a new ReadController
func NewReadController(userService *service.UserService) *ReadController {
	return &ReadController{
		userService: userService,
	}
}

// Handle handles the health check request
func (lc *ReadController) Handle(ctx *gin.Context) {
	// Get userID from path parameters
	userID := ctx.Param("userID")
	if userID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "userID is required", nil, nil)
		return
	}

	user, profile, err := lc.userService.GetUserProfile(ctx.Request.Context(), userID)

	// if err is not found, return a 404 error
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "profile not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "user not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to fetch user", nil, err)
	}

	// Exclude sensitive fields from the response
	user.Password = "REDACTED"
	user.OtpSecret = "REDACTED"

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, gin.H{
		"user":    user,
		"profile": profile,
	}, nil)
}
