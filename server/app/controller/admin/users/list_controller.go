package users

import (
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	// Get pagination parameters from query params
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid page parameter", nil, err)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid limit parameter", nil, err)
		return
	}

	users, err := lc.userService.ListUsers(ctx.Request.Context(), pageInt, limitInt)
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
