package users

import (
	"net/http"

	"koneksi/services/iam/app/helpers"
	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/repositories"
	"koneksi/services/iam/core/logger"

	"github.com/gin-gonic/gin"
)

// RegisterController handles user registration
type RegisterController struct {
	userRepo     *repositories.UserRepository
	profileRepo  *repositories.ProfileRepository
	roleRepo     *repositories.RoleRepository
	userRoleRepo *repositories.UserRoleRepository
}

// NewRegisterController initializes a new RegisterController
func NewRegisterController(userRepo *repositories.UserRepository, profileRepo *repositories.ProfileRepository, roleRepo *repositories.RoleRepository, userRoleRepo *repositories.UserRoleRepository) *RegisterController {
	return &RegisterController{
		userRepo:     userRepo,
		profileRepo:  profileRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

// Handle processes the user registration request
func (rc *RegisterController) Handle(c *gin.Context) {
	var request struct {
		FirstName       string  `json:"first_name" binding:"required"`
		MiddleName      *string `json:"middle_name"`
		LastName        string  `json:"last_name" binding:"required"`
		Suffix          *string `json:"suffix"`
		Email           string  `json:"email" binding:"required,email"`
		Password        string  `json:"password" binding:"required,min=8"`
		ConfirmPassword string  `json:"confirm_password" binding:"required,eqfield=Password"`
	}

	// Validate the payload
	if err := rc.validatePayload(c, &request); err != nil {
		return
	}

	// Check if the email already exists
	existingUser, err := rc.userRepo.ReadUserByEmail(c.Request.Context(), request.Email)
	if err != nil {
		logger.Log.Error("error checking existing user", logger.Error(err))
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Internal server error", nil, nil)
		return
	}
	if existingUser != nil {
		helpers.FormatResponse(c, "error", http.StatusConflict, "Email already exists", nil, nil)
		return
	}

	// Create the user
	user := &models.User{
		Email:      request.Email,
		Password:   request.Password,
		IsVerified: true,
	}
	if err := rc.userRepo.CreateUser(c.Request.Context(), user); err != nil {
		logger.Log.Error("error creating user", logger.Error(err))
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Failed to create user", nil, nil)
		return
	}

	// Create the profile
	profile := &models.Profile{
		UserID:     user.ID,
		FirstName:  request.FirstName,
		MiddleName: request.MiddleName,
		LastName:   request.LastName,
		Suffix:     request.Suffix,
	}
	if err := rc.profileRepo.CreateProfile(c.Request.Context(), profile); err != nil {
		logger.Log.Error("error creating profile", logger.Error(err))
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Failed to create profile", nil, nil)
		return
	}

	// Assign the "user" role to the newly registered user
	userRole, err := rc.roleRepo.ReadRoleByName(c.Request.Context(), "user")
	if err != nil {
		logger.Log.Error("error retrieving default role", logger.Error(err))
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Failed to assign default role", nil, nil)
		return
	}
	if userRole == nil {
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Default role not found", nil, nil)
		return
	}

	userRoleAssignment := &models.UserRole{
		UserID: user.ID,
		RoleID: userRole.ID,
	}
	if err := rc.userRoleRepo.CreateUserRole(c.Request.Context(), userRoleAssignment); err != nil {
		logger.Log.Error("error assigning default role", logger.Error(err))
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Failed to assign default role", nil, nil)
		return
	}

	// Respond with success
	helpers.FormatResponse(c, "success", http.StatusCreated, "User registered successfully", gin.H{
		"user":      user,
		"profile":   profile,
		"user_role": userRoleAssignment,
	}, nil)
}

// validatePayload validates the incoming request payload
func (rc *RegisterController) validatePayload(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helpers.FormatResponse(c, "error", http.StatusBadRequest, "Invalid input", nil, nil)
		return err
	}
	return nil
}
