package middleware

import (
	"net/http"

	"koneksi/server/app/helpers"
	"koneksi/server/app/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerifiedMiddleware struct {
	HandleVerified gin.HandlerFunc
}

func NewVerifiedMiddleware(userRepo *repositories.UserRepository) *VerifiedMiddleware {
	return &VerifiedMiddleware{
		HandleVerified: func(c *gin.Context) {
			// Retrieve the userID from the context
			userIDValue, exists := c.Get("userID")
			if !exists {
				helpers.FormatResponse(c, "error", http.StatusUnauthorized, "User ID not found in context", nil, nil)
				c.Abort()
				return
			}

			// Ensure the userID is a string
			userID, ok := userIDValue.(string)
			if !ok {
				helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Invalid user ID format in context", nil, nil)
				c.Abort()
				return
			}

			// Convert userID to primitive.ObjectID
			objectID, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Invalid user ID format", nil, nil)
				c.Abort()
				return
			}

			// Check if the user is verified using UserRepository
			user, err := userRepo.ReadUserByID(c.Request.Context(), objectID)
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
