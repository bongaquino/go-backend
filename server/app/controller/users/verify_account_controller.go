// server/app/controller/users/verify_account_controller.go
package users

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// VerifyAccountController handles verifying user accounts
type VerifyAccountController struct {
	userService *service.UserService
}

// NewVerifyAccountController initializes a new VerifyAccountController
func NewVerifyAccountController(userService *service.UserService) *VerifyAccountController {
	return &VerifyAccountController{
		userService: userService,
	}
}

func (vac *VerifyAccountController) Handle(ctx *gin.Context) {
	var request struct {
		VerificationCode string `json:"verification_code" binding:"required"`
	}

	// Validate the request payload
	if err := ctx.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// userID email from the user token
	userID, exists := ctx.Get("userID")
	if !exists {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "unauthorized", nil, nil)
		return
	}

	// Verify code using the UserService
	err := vac.userService.VerifyUserAccount(ctx.Request.Context(), userID.(string), request.VerificationCode)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(ctx, "success", http.StatusOK, "account verified successfully", nil, nil)
}
