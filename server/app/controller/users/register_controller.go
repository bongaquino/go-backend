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
	emailService *service.EmailService
}

func NewRegisterController(userService *service.UserService, tokenService *service.TokenService, emailService *service.EmailService) *RegisterController {
	return &RegisterController{
		userService:  userService,
		tokenService: tokenService,
		emailService: emailService,
	}
}

func (rc *RegisterController) Handle(ctx *gin.Context) {
	var request dto.CreateUser

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

	// Add user role to the request
	request.Role = "user"

	// Register the user
	user, profile, userRole, roleName, err := rc.userService.RegisterUser(ctx.Request.Context(), &request)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Generate a verification code
	code, err := rc.userService.GenerateVerificationCode(ctx.Request.Context(), user.ID.Hex())
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Send the verification email
	err = rc.emailService.SendVerificationCode(user.Email, code)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to send verification email", nil, nil)
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

func (rc *RegisterController) validatePayload(ctx *gin.Context, request *dto.CreateUser) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	// Check if new passwords pass validation
	isValid, validationErr := helper.ValidatePassword(request.Password)
	if !isValid || validationErr != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, validationErr.Error(), nil, nil)
		return validationErr
	}
	return nil
}
