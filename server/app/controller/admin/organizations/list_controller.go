package organizations

import (
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListController struct {
	orgService *service.OrganizationService
}

// NewListController initializes a new ListController
func NewListController(orgService *service.OrganizationService) *ListController {
	return &ListController{
		orgService: orgService,
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

	users, err := lc.orgService.ListOrgs(ctx.Request.Context(), pageInt, limitInt)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to fetch users", nil, err)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, nil, users, nil)
}
