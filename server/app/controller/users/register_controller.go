package users

import (
	"net/http"

	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	userService *service.UserService
}

func NewRegisterController(userService *service.UserService) *RegisterController {
	return &RegisterController{
		userService: userService,
	}
}

func (rc *RegisterController) Handle(c *gin.Context) {
	var request dto.NewUser

	if err := rc.validatePayload(c, &request); err != nil {
		return
	}

	user, profile, userRole, err := rc.userService.RegisterUser(c.Request.Context(), &request)
	if err != nil {
		helper.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.FormatResponse(c, "success", http.StatusCreated, "user registered successfully", gin.H{
		"user":      user,
		"profile":   profile,
		"user_role": userRole,
	}, nil)
}

func (rc *RegisterController) validatePayload(c *gin.Context, request any) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(c, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
