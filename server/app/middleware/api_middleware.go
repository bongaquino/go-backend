package middleware

import (
	"net/http"
	"time"

	"koneksi/server/app/helper"
	"koneksi/server/app/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type APIMiddleware struct {
	Handle gin.HandlerFunc
}

func NewAPIMiddleware(svcAccRepo *repository.ServiceAccountRepository) *APIMiddleware {
	return &APIMiddleware{
		Handle: func(ctx *gin.Context) {
			// Retrieve ClientID from the request header
			clientID := ctx.GetHeader("Client-ID")
			if clientID == "" {
				helper.FormatResponse(ctx, "error", http.StatusBadRequest, "Client-ID header is required", nil, nil)
				ctx.Abort()
				return
			}

			// Retrieve the ClientSecret from the request header
			clientSecret := ctx.GetHeader("Client-Secret")
			if clientSecret == "" {
				helper.FormatResponse(ctx, "error", http.StatusBadRequest, "Client-Secret header is required", nil, nil)
				ctx.Abort()
				return
			}

			// Read the ClientID and ClientSecret using ServiceAccountRepository
			serviceAccount, err := svcAccRepo.ReadByClientID(ctx.Request.Context(), clientID)
			if err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to validate credentials", nil, nil)
				ctx.Abort()
				return
			}
			if serviceAccount == nil {
				helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "invalid credentials", nil, nil)
				ctx.Abort()
				return
			}

			// Check if the ClientSecret is valid
			if !helper.CheckHash(clientSecret, serviceAccount.ClientSecret) {
				helper.FormatResponse(ctx, "error", http.StatusUnauthorized, "invalid credentials", nil, nil)
				ctx.Abort()
				return
			}

			// Update service account last used timestamp
			serviceAccount.LastUsedAt = time.Now()

			// Update the service account in the repository
			if err := svcAccRepo.UpdateByClientID(ctx.Request.Context(), clientID, bson.M{
				"last_used_at": serviceAccount.LastUsedAt,
			}); err != nil {
				helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to update service account", nil, nil)
				ctx.Abort()
				return
			}

			// Set the client ID in the context
			ctx.Set("clientID", clientID)

			// Set the user ID in the context
			ctx.Set("userID", serviceAccount.UserID.Hex())

			// Continue to the next middleware
			ctx.Next()
		},
	}
}
