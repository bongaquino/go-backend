package database

import (
	"context"
	"fmt"
	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/repositories"
	"koneksi/services/iam/core/logger"
)

// SeedCollections seeds initial data into MongoDB collections
func SeedCollections(permissionRepo *repositories.PermissionRepository, roleRepo *repositories.RoleRepository, rolePermissionRepo *repositories.RolePermissionRepository) {
	ctx := context.Background()

	seeders := []struct {
		Name   string
		Seeder func(context.Context) error
	}{
		{"permissions", func(ctx context.Context) error { return seedPermissions(ctx, permissionRepo) }},
		{"roles", func(ctx context.Context) error { return seedRoles(ctx, roleRepo) }},
		{"role_permissions", func(ctx context.Context) error {
			return seedRolePermissions(ctx, roleRepo, permissionRepo, rolePermissionRepo)
		}},
	}

	for _, seeder := range seeders {
		if err := seeder.Seeder(ctx); err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to seed collection: %s", seeder.Name), logger.Error(err))
		} else {
			logger.Log.Info(fmt.Sprintf("Seeded collection: %s", seeder.Name))
		}
	}
}

// seedPermissions inserts initial permissions using the repository
func seedPermissions(ctx context.Context, permissionRepo *repositories.PermissionRepository) error {
	permissions := []models.Permission{
		{Name: "upload_files"},
		{Name: "download_files"},
		{Name: "list_files"},
	}

	for _, perm := range permissions {
		existing, err := permissionRepo.ReadPermissionByName(ctx, perm.Name)
		if err != nil {
			return err
		}
		if existing == nil {
			if err := permissionRepo.CreatePermission(ctx, &perm); err != nil {
				return err
			}
		} else {
			logger.Log.Info(fmt.Sprintf("Skipping permission: %s (already exists)", perm.Name))
		}
	}
	return nil
}

// seedRoles inserts initial roles using the repository
func seedRoles(ctx context.Context, roleRepo *repositories.RoleRepository) error {
	roles := []models.Role{
		{Name: "user"},
	}

	for _, role := range roles {
		existing, err := roleRepo.ReadRoleByName(ctx, role.Name)
		if err != nil {
			return err
		}
		if existing == nil {
			if err := roleRepo.CreateRole(ctx, &role); err != nil {
				return err
			}
		} else {
			logger.Log.Info(fmt.Sprintf("Skipping role: %s (already exists)", role.Name))
		}
	}
	return nil
}

// seedRolePermissions assigns all permissions to the "user" role using repositories
func seedRolePermissions(ctx context.Context, roleRepo *repositories.RoleRepository, permissionRepo *repositories.PermissionRepository, rolePermissionRepo *repositories.RolePermissionRepository) error {
	// Find the "user" role
	userRole, err := roleRepo.ReadRoleByName(ctx, "user")
	if err != nil {
		return err
	}
	if userRole == nil {
		return fmt.Errorf("user role not found")
	}

	// Get all permissions
	permissions := []string{"upload_files", "download_files", "list_files"}
	for _, permName := range permissions {
		perm, err := permissionRepo.ReadPermissionByName(ctx, permName)
		if err != nil {
			return err
		}
		if perm == nil {
			logger.Log.Warn(fmt.Sprintf("Skipping role permission seeding: Permission %s not found", permName))
			continue
		}

		// Check if the role-permission already exists
		existingPermissions, err := rolePermissionRepo.ReadRolePermissions(ctx, userRole.ID)
		if err != nil {
			return err
		}
		alreadyExists := false
		for _, rp := range existingPermissions {
			if rp.PermissionID == perm.ID {
				alreadyExists = true
				break
			}
		}

		if !alreadyExists {
			rolePermission := models.RolePermission{
				RoleID:       userRole.ID,
				PermissionID: perm.ID,
			}
			if err := rolePermissionRepo.CreateRolePermission(ctx, &rolePermission); err != nil {
				return err
			}
		} else {
			logger.Log.Info(fmt.Sprintf("Skipping role permission: %s -> %s (already exists)", userRole.Name, perm.Name))
		}
	}
	return nil
}
