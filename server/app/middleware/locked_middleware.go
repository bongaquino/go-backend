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

// Unlike other routes that already has middleware to check for authentication and get the user ID, this middleware should be used on routes that do not require authentication but still need to check if the account is locked in a secure manner.
func NewLockedMiddleware(userRepo *repository.UserRepository) *LockedMiddleware {
	return &LockedMiddleware{
		Handle: func(ctx *gin.Context) {
			// Retrieve userID from the context (assumes it's set by a previous middleware)
			userIDValue, exists := ctx.Get("userID")
			if !exists {
				helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "userID not found in context", nil, nil)
				ctx.Abort()
				return
			}

			// Ensure the userID is a string
			userID, ok := userIDValue.(string)
			if !ok {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "invalid user ID format in context", nil, nil)
				ctx.Abort()
				return
			}

			// Check if the user is verified using UserRepository
			user, err := userRepo.Read(ctx.Request.Context(), userID)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to retrieve user", nil, nil)
				ctx.Abort()
				return
			}

			if user != nil && user.IsLocked {
				helper.FormatResponse(ctx, "error", http.StatusForbidden, "account locked due to multiple failed login attempts", nil, nil)
				ctx.Abort()
				return
			}

			// Continue to the next middleware
			ctx.Next()
		},
	}
}
