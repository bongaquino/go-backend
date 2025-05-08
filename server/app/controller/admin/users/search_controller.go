package users

import (
	"koneksi/server/app/helper"
	"koneksi/server/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	userService *service.UserService
}

// NewSearchController initializes a new SearchController
func NewSearchController(userService *service.UserService) *SearchController {
	return &SearchController{
		userService: userService,
	}
}

// Handle handles the health check request
func (lc *SearchController) Handle(ctx *gin.Context) {
	// Get email from query params
	email := ctx.Query("email")

	// Validate email
	if email == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "email is required", nil, nil)
		return
	}

	users, err := lc.userService.SearchUsers(ctx, email)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to fetch users", nil, err)
		return
	}

	// Exclude sensitive fields from the response
	for i := range users {
		users[i].Password = "REDACTED"
		users[i].OtpSecret = "REDACTED"
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, users, nil)
}
