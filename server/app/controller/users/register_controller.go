package users

import (
	"net/http"

	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewRegisterController(userService *service.UserService, tokenService *service.TokenService) *RegisterController {
	return &RegisterController{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (rc *RegisterController) Handle(ctx *gin.Context) {
	var request dto.RegisterUser

	if err := rc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Check if user already exists
	exists, err := rc.userService.UserExists(ctx.Request.Context(), request.Email)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}
	if exists {
		helper.FormatResponse(ctx, "error", http.StatusConflict, "user already exists", nil, nil)
		return
	}

	// Register the user
	user, profile, userRole, roleName, err := rc.userService.RegisterUser(ctx.Request.Context(), &request)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := rc.tokenService.AuthenticateUser(ctx.Request.Context(), user.Email, request.Password)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	helper.FormatResponse(ctx, "success", http.StatusCreated, "user registered successfully", gin.H{
		"user": gin.H{
			"email": user.Email,
		},
		"profile": profile,
		"user_role": gin.H{
			"role_id":   userRole.RoleID,
			"role_name": roleName,
		},
		"tokens": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}, nil)
}

func (rc *RegisterController) validatePayload(ctx *gin.Context, request any) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
