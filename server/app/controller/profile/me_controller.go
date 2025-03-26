package profile

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MeController handles health-related endpoints
type MeController struct {
	userService *service.UserService
}

// NewMeController initializes a new MeController
func NewMeController(userService *service.UserService) *MeController {
	return &MeController{
		userService: userService,
	}
}

// Handles the health check endpoint
func (hc *MeController) Handle(ctx *gin.Context) {
	// Extract user ID from the context
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
		return
	}

	// Fetch the user profile
	user, profile, err := hc.userService.GetUserProfile(ctx.Request.Context(), userID.(string))
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Sanitize the user object by removing sensitive fields
	sanitizedUser := gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"is_verified":    user.IsVerified,
		"is_mfa_enabled": user.IsMFAEnabled,
	}

	// Sanitize the profile object by removing sensitive fields
	sanitizedProfile := gin.H{
		"first_name": profile.FirstName,
		"last_name":  profile.LastName,
	}

	// Return the user profile
	helper.FormatResponse(ctx, "success", http.StatusOK, "user profile retrieved successfully", gin.H{
		"user":    sanitizedUser,
		"profile": sanitizedProfile,
	}, nil)
}
