package middleware

import (
	"context"
	"fmt"
	"net/http"

	"koneksi/server/app/helpers"
	"koneksi/server/app/repositories"

	"slices"

	"github.com/gin-gonic/gin"
)

type AuthzMiddleware struct {
	userRoleRepository *repositories.UserRoleRepository
	roleRepository     *repositories.RoleRepository
}

func NewAuthzMiddleware(userRoleRepository *repositories.UserRoleRepository, roleRepository *repositories.RoleRepository) *AuthzMiddleware {
	return &AuthzMiddleware{
		userRoleRepository: userRoleRepository,
		roleRepository:     roleRepository,
	}
}

func (m *AuthzMiddleware) HandleAuthz(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve userID from the context (assumes it's set by a previous middleware)
		userID, exists := c.Get("userID")
		if !exists {
			helpers.FormatResponse(c, "error", http.StatusUnauthorized, "User ID not found in context", nil, nil)
			c.Abort()
			return
		}

		// Fetch roles from the database using the userID
		roles, err := m.getUserRoles(c.Request.Context(), userID.(string))
		if err != nil {
			helpers.FormatResponse(c, "error", http.StatusInternalServerError, "Failed to retrieve user roles", nil, nil)
			c.Abort()
			return
		}

		// Check if the user has at least one of the required roles
		hasRole := false
		for _, requiredRole := range requiredRoles {
			if slices.Contains(roles, requiredRole) {
				hasRole = true
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			helpers.FormatResponse(c, "error", http.StatusForbidden, "User does not have the required role", nil, nil)
			c.Abort()
			return
		}

		// Continue to the next middleware
		c.Next()
	}
}

// getUserRoles fetches roles for a given userID from the database
func (m *AuthzMiddleware) getUserRoles(ctx context.Context, userID string) ([]string, error) {
	// Fetch user roles from the UserRoleRepository
	userRoles, err := m.userRoleRepository.ReadUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	fmt.Println(userRoles)

	// Initialize a slice to store role names
	var roles []string

	// Iterate through user roles and fetch role names using RoleRepository
	for _, userRole := range userRoles {
		// Fetch the role by RoleID
		role, err := m.roleRepository.ReadRoleByID(ctx, userRole.RoleID.Hex())
		if err != nil {
			return nil, err
		}
		if role != nil {
			roles = append(roles, role.Name)
		}
	}

	return roles, nil
}
