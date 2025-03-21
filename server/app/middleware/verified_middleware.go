package middleware

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/repository"

	"github.com/gin-gonic/gin"
)

type VerifiedMiddleware struct {
	HandleVerified gin.HandlerFunc
}

func NewVerifiedMiddleware(userRepo *repository.UserRepository) *VerifiedMiddleware {
	return &VerifiedMiddleware{
		HandleVerified: func(c *gin.Context) {
			// Retrieve the userID from the context
			userIDValue, exists := c.Get("userID")
			if !exists {
				helper.FormatResponse(c, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
				c.Abort()
				return
			}

			// Ensure the userID is a string
			userID, ok := userIDValue.(string)
			if !ok {
				helper.FormatResponse(c, "error", http.StatusInternalServerError, "invalid user ID format in context", nil, nil)
				c.Abort()
				return
			}

			// Check if the user is verified using UserRepository
			user, err := userRepo.ReadUserByID(c.Request.Context(), userID)
			if err != nil {
				helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to retrieve user", nil, nil)
				c.Abort()
				return
			}

			if user == nil || !user.IsVerified {
				helper.FormatResponse(c, "error", http.StatusForbidden, "user is not verified", nil, nil)
				c.Abort()
				return
			}

			// Continue to the next middleware
			c.Next()
		},
	}
}
