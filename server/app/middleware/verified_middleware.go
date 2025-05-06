package middleware

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/repository"

	"github.com/gin-gonic/gin"
)

type VerifiedMiddleware struct {
	Handle gin.HandlerFunc
}

func NewVerifiedMiddleware(userRepo *repository.UserRepository) *VerifiedMiddleware {
	return &VerifiedMiddleware{
		Handle: func(ctx *gin.Context) {
			// Retrieve the userID from the context
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
			user, err := userRepo.ReadUser(ctx.Request.Context(), userID)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to retrieve user", nil, nil)
				ctx.Abort()
				return
			}

			if user == nil || !user.IsVerified {
				helper.FormatResponse(ctx, "error", http.StatusForbidden, "user is not verified", nil, nil)
				ctx.Abort()
				return
			}

			// Continue to the next middleware
			ctx.Next()
		},
	}
}
