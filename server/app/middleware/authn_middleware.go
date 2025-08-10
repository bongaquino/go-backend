package middleware

import (
	"net/http"
	"strings"

	"bongaquino/server/app/helper"
	"bongaquino/server/app/provider"

	"github.com/gin-gonic/gin"
)

type AuthnMiddleware struct {
	Handle gin.HandlerFunc
}

func NewAuthnMiddleware(jwtService *provider.JWTProvider) *AuthnMiddleware {
	return &AuthnMiddleware{
		Handle: func(ctx *gin.Context) {
			// Get the Authorization header
			authHeader := ctx.GetHeader("Authorization")
			if authHeader == "" {
				helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "authorization header required", nil, nil)
				ctx.Abort()
				return
			}

			// Extract the token from the Authorization header
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "invalid authorization header", nil, nil)
				ctx.Abort()
				return
			}

			// Validate the token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "invalid or expired access token", nil, nil)
				ctx.Abort()
				return
			}

			// Set the user ID in the context
			ctx.Set("userID", claims.Sub)

			// Continue to the next middleware
			ctx.Next()
		},
	}
}
