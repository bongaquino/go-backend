package users

import (
	"net/http"

	"koneksi/services/iam/app/helpers"
	"koneksi/services/iam/app/repositories"

	"github.com/gin-gonic/gin"
)

// RequestTokenController handles user registration
type RequestTokenController struct {
	userRepo *repositories.UserRepository
}

// NewRequestTokenController initializes a new RequestTokenController
func NewRequestTokenController(userRepo *repositories.UserRepository) *RequestTokenController {
	return &RequestTokenController{}
}

// Handle processes the user registration request
func (rc *RequestTokenController) Handle(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	// Validate the payload
	if err := rc.validatePayload(c, &request); err != nil {
		return
	}
}

// validatePayload validates the incoming request payload
func (rc *RequestTokenController) validatePayload(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindJSON(request); err != nil {
		helpers.FormatResponse(c, "error", http.StatusBadRequest, "Invalid input", nil, nil)
		return err
	}
	return nil
}
