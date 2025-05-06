package user

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListController handles health-related endpoints
type ListController struct {
	userService *service.UserService
}

// NewListController initializes a new ListController
func NewListController(userService *service.UserService) *ListController {
	return &ListController{
		userService: userService,
	}
}

// Handle handles the health check request
func (lc *ListController) Handle(ctx *gin.Context) {
	users, err := lc.userService.ListUsers(ctx.Request.Context(), 1, 10)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch users", "details": err.Error()})
		return
	}
	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, users, nil)
}
