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
	userRepo    *repositories.UserRepository
	profileRepo *repositories.ProfileRepository
}

// NewRegisterController initializes a new RegisterController
func NewRegisterController(userRepo *repositories.UserRepository, profileRepo *repositories.ProfileRepository) *RegisterController {
	return &RegisterController{
		userRepo:    userRepo,
		profileRepo: profileRepo,
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

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.FormatResponse(c, "error", http.StatusBadRequest, "Invalid input", nil, nil)
		return
	}

	// Check if the email already exists
	existingUser, err := rc.userRepo.ReadUserByUsername(c.Request.Context(), request.Email)
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
		Email:    request.Email,
		Password: request.Password,
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

	// Respond with success
	helpers.FormatResponse(c, "success", http.StatusCreated, "User registered successfully", nil, nil)
}
