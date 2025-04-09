package middleware

import (
	"net/http"
	"strings"

	"koneksi/server/app/helper"
	"koneksi/server/app/provider"

	"github.com/gin-gonic/gin"
)

type AuthnMiddleware struct {
	Handle gin.HandlerFunc
}

func NewAuthnMiddleware(jwtService *provider.JWTProvider) *AuthnMiddleware {
	return &AuthnMiddleware{
		Handle: func(c *gin.Context) {
			// Get the Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				helper.FormatResponse(c, "error", http.StatusUnauthorized, "authorization header required", nil, nil)
				c.Abort()
				return
			}

			// Extract the token from the Authorization header
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				helper.FormatResponse(c, "error", http.StatusUnauthorized, "invalid authorization header", nil, nil)
				c.Abort()
				return
			}

			// Validate the token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				helper.FormatResponse(c, "error", http.StatusUnauthorized, "invalid or expired access token", nil, nil)
				c.Abort()
				return
			}

			// Set the user ID in the context
			c.Set("userID", claims.Sub)

			// Continue to the next middleware
			c.Next()
		},
	}
}
