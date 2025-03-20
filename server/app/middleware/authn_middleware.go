package middleware

import (
	"net/http"
	"strings"

	"koneksi/server/app/helpers"
	"koneksi/server/app/services"

	"github.com/gin-gonic/gin"
)

type AuthnMiddleware struct {
	HandleAuth gin.HandlerFunc
}

func NewAuthnMiddleware(jwtService *services.JWTService) *AuthnMiddleware {
	return &AuthnMiddleware{
		HandleAuth: func(c *gin.Context) {
			// Get the Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				helpers.FormatResponse(c, "error", http.StatusUnauthorized, "Authorization header required", nil, nil)
				c.Abort()
				return
			}

			// Extract the token from the Authorization header
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				helpers.FormatResponse(c, "error", http.StatusUnauthorized, "Invalid authorization header", nil, nil)
				c.Abort()
				return
			}

			// Validate the token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				helpers.FormatResponse(c, "error", http.StatusUnauthorized, "Invalid or expired access token", nil, nil)
				c.Abort()
				return
			}

			// Set the user ID and email in the context
			c.Set("userID", claims.Sub)
			c.Set("email", claims.Email)

			// Continue to the next middleware
			c.Next()
		},
	}
}
