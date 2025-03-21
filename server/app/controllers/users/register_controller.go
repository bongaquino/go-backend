package users

import (
	"net/http"

	"koneksi/server/app/helpers"

	"github.com/gin-gonic/gin"
)

// RegisterController handles user registration
type RegisterController struct {
	userService *providers.UserService
}

// NewRegisterController initializes a new RegisterController
func NewRegisterController(userService *providers.UserService) *RegisterController {
	return &RegisterController{
		userService: userService,
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

	// Call the UserService to handle registration
	user, profile, userRole, err := rc.userService.RegisterUser(c.Request.Context(), &request)
	if err != nil {
		helpers.FormatResponse(c, "error", http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	// Respond with success
	helpers.FormatResponse(c, "success", http.StatusCreated, "User registered successfully", gin.H{
		"user":      user,
		"profile":   profile,
		"user_role": userRole,
	}, nil)
}

// validatePayload validates the incoming request payload
func (rc *RegisterController) validatePayload(c *gin.Context, request any) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helpers.FormatResponse(c, "error", http.StatusBadRequest, "Invalid input", nil, nil)
		return err
	}
	return nil
}
