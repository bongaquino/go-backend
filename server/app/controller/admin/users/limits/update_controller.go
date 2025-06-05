package limits

import (
	"koneksi/server/app/dto"
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
func (uc *UpdateController) Handle(ctx *gin.Context) {
	// Get userID from path parameters
	userID := ctx.Param("userID")
	if userID == "" {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "userID is required", nil, nil)
		return
	}

	// Get request body
	var request dto.UpdateLimitDTO
	if err := uc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Check if user usage is greater than request limit
	limit, err := uc.userService.GetUserLimits(ctx, userID)
	if err != nil {
		if err.Error() == "user limit not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, "user limit not found", nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to get user usage", nil, nil)
		return
	}
	if limit.BytesUsage > request.BytesLimit {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "user usage exceeds new limit", nil, nil)
		return
	}

	// Update limit
	if err := uc.userService.UpdateUserLimit(ctx, userID, request.BytesLimit); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update user limit", nil, nil)
		return
	}

	// Send success response
	helper.FormatResponse(ctx, "success", http.StatusOK, "user limit updated successfully", nil, nil)
}

func (rc *UpdateController) validatePayload(ctx *gin.Context, request *dto.UpdateLimitDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
		return err
	}
	return nil
}
