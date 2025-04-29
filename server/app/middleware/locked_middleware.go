package middleware

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/repository"

	"github.com/gin-gonic/gin"
)

type LockedMiddleware struct {
	Handle gin.HandlerFunc
}

// TODO: Implement the middleware to check if the user account is locked
// Unlike other routes that already has middleware to check for authentication and get the user ID, this middleware should be used on routes that do not require authentication but still need to check if the account is locked in a secure manner.
func NewLockedMiddleware(userRepo *repository.UserRepository) *LockedMiddleware {
	return &LockedMiddleware{
		Handle: func(c *gin.Context) {
			// Retrieve userID from the context (assumes it's set by a previous middleware)
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
			user, err := userRepo.ReadUser(c.Request.Context(), userID)
			if err != nil {
				helper.FormatResponse(c, "error", http.StatusInternalServerError, "failed to retrieve user", nil, nil)
				c.Abort()
				return
			}

			if user != nil && user.IsLocked {
				helper.FormatResponse(c, "error", http.StatusForbidden, "account locked due to multiple failed login attempts", nil, nil)
				c.Abort()
				return
			}

			// Continue to the next middleware
			c.Next()
		},
	}
}
