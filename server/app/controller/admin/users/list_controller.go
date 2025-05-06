package user

import (
	"koneksi/server/app/service"

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
func (c *ListController) Handle(ctx *gin.Context) {
	users, err := c.userService.ListUsers(ctx.Request.Context(), 1, 10)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch users", "details": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"users": users})
}
