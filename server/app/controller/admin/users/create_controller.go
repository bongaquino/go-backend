package users

import (
	"net/http"

	"bongaquino/server/app/dto"
	"bongaquino/server/app/helper"
	"bongaquino/server/app/service"

	"github.com/gin-gonic/gin"
)

type CreateController struct {
	userService  *service.UserService
	tokenService *service.TokenService
	emailService *service.EmailService
}

func NewCreateController(userService *service.UserService, tokenService *service.TokenService, emailService *service.EmailService) *CreateController {
	return &CreateController{
		userService:  userService,
		tokenService: tokenService,
		emailService: emailService,
	}
}

func (cc *CreateController) Handle(ctx *gin.Context) {
	var request dto.CreateUserDTO

	if err := cc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Check if user already exists
	exists, err := cc.userService.UserExists(ctx.Request.Context(), request.Email)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}
	if exists {
		helper.FormatResponse(ctx, "error", http.StatusConflict, "user already exists", nil, nil)
		return
	}

	// Add user role to the request
	if request.Role == "" {
		request.Role = "system_user"
	}
	request.IsVerified = true

	// Create the user
	user, profile, userRole, roleName, err := cc.userService.CreateUser(ctx.Request.Context(), &request)

	// If err is "role not found", return 404
	if err != nil {
		if err.Error() == "role not found" {
			helper.FormatResponse(ctx, "error", http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	helper.FormatResponse(ctx, "success", http.StatusCreated, "user created successfully", gin.H{
		"user": gin.H{
			"email": user.Email,
		},
		"profile": profile,
		"user_role": gin.H{
			"role_id":   userRole.RoleID,
			"role_name": roleName,
		},
	}, nil)
}

func (cc *CreateController) validatePayload(ctx *gin.Context, request *dto.CreateUserDTO) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid request body", nil, nil)
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
