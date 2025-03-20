package middleware

import (
	"net/http"

	"koneksi/server/app/helpers"
	"koneksi/server/app/repositories"

	"github.com/gin-gonic/gin"
)

type VerifiedMiddleware struct {
	HandleVerified gin.HandlerFunc
}

func NewVerifiedMiddleware(userRepo *repositories.UserRepository) *VerifiedMiddleware {
	return &VerifiedMiddleware{
		HandleVerified: func(c *gin.Context) {
			// Retrieve the user ID from the context
			email, exists := c.Get("email")
			if !exists {
				helpers.FormatResponse(c, "error", http.StatusUnauthorized, "User ID not found in context", nil, nil)
				c.Abort()
				return
			}

			// Check if the user is verified using UserRepository
			user, err := userRepo.ReadUserByEmail(c.Request.Context(), email.(string))
			if err != nil {
				helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Failed to retrieve user", nil, nil)
				c.Abort()
				return
			}

			if user == nil || !user.IsVerified {
				helpers.FormatResponse(c, "error", http.StatusForbidden, "User is not verified", nil, nil)
				c.Abort()
				return
			}

			// Continue to the next middleware
			c.Next()
		},
	}
}
