// server/app/controller/users/verify_account_controller.go
package users

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// VVerifyAccountController handles verifying user accounts
type VerifyAccountController struct {
	userService *service.UserService
}

// NewVerifyAccountController initializes a new VerifyAccountController
func NewVerifyAccountController(userService *service.UserService) *VerifyAccountController {
	return &VerifyAccountController{
		userService: userService,
	}
}

func (vac *VerifyAccountController) Handle(c *gin.Context) {
	var request struct {
		Email              string `json:"email" binding:"required,email"`
		Token              string `json:"token" binding:"required"`
	}

	// Validate the request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return
	}

	// Verify token using the UserService
	err := vac.userService.VerifyUserAccount(c.Request.Context(), request.Email, request.Token)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helper.FormatResponse(c, "success", http.StatusOK, "account verified successfully", nil, nil)
}