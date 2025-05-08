package database

import (
	"context"
	"fmt"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"
)

// SeedCollections seeds initial data into MongoDB collections
func SeedCollections(
	permissionRepo *repository.PermissionRepository,
	roleRepo *repository.RoleRepository,
	rolePermissionRepo *repository.RolePermissionRepository,
	policyRepo *repository.PolicyRepository,
	policyPermissionRepo *repository.PolicyPermissionRepository,
) {
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
		{"policies", func(ctx context.Context) error { return seedPolicies(ctx, policyRepo) }},
	}

	for _, seeder := range seeders {
		if err := seeder.Seeder(ctx); err != nil {
			logger.Log.Error(fmt.Sprintf("failed to seed collection: %s", seeder.Name), logger.Error(err))
		} else {
			logger.Log.Info(fmt.Sprintf("seeded collection: %s", seeder.Name))
		}
	}
}

// seedPermissions inserts initial permissions using the repository
func seedPermissions(ctx context.Context, permissionRepo *repository.PermissionRepository) error {
	permissions := []model.Permission{
		{Name: "list:users"},
		{Name: "create:user"},
		{Name: "read:user"},
		{Name: "update:user"},
		{Name: "list:organizations"},
		{Name: "create:organization"},
		{Name: "read:organization"},
		{Name: "update:organization"},
		{Name: "list:files"},
		{Name: "upload:file"},
		{Name: "download:file"},
	}

	for _, perm := range permissions {
		existing, err := permissionRepo.ReadByName(ctx, perm.Name)
		if err != nil {
			return err
		}
		if existing == nil {
			if err := permissionRepo.Create(ctx, &perm); err != nil {
				return err
			}
		} else {
			logger.Log.Info(fmt.Sprintf("skipping permission: %s (already exists)", perm.Name))
		}
	}
	return nil
}

// seedRoles inserts initial roles using the repository
func seedRoles(ctx context.Context, roleRepo *repository.RoleRepository) error {
	roles := []model.Role{
		{Name: "system_admin"},
		{Name: "system_user"},
		{Name: "organization_admin"},
		{Name: "organization_user"},
		{Name: "organization_viewer"},
	}

	for _, role := range roles {
		existing, err := roleRepo.ReadByName(ctx, role.Name)
		if err != nil {
			return err
		}
		if existing == nil {
			if err := roleRepo.Create(ctx, &role); err != nil {
				return err
			}
		} else {
			logger.Log.Info(fmt.Sprintf("skipping role: %s (already exists)", role.Name))
		}
	}
	return nil
}

// seedRolePermissions assigns specific permissions to roles using a role-permission map
func seedRolePermissions(
	ctx context.Context,
	roleRepo *repository.RoleRepository,
	permissionRepo *repository.PermissionRepository,
	rolePermissionRepo *repository.RolePermissionRepository,
) error {
	rolePermissionsMap := map[string][]string{
		"system_admin": {
			"list:users", "create:user", "read:user", "update:user",
			"list:organizations", "create:organization", "read:organization", "update:organization",
			"list:files", "upload:file", "download:file",
		},
		"system_user": {
			"list:files", "upload:file", "download:file",
		},
		"organization_admin": {
			"read:organization", "update:organization",
			"list:files", "upload:file", "download:file",
		},
		"organization_user": {
			"list:files", "upload:file", "download:file",
		},
		"organization_viewer": {
			"list:files", "download:file",
		},
	}

	for roleName, permissionNames := range rolePermissionsMap {
		role, err := roleRepo.ReadByName(ctx, roleName)
		if err != nil {
			return err
		}
		if role == nil {
			logger.Log.Warn(fmt.Sprintf("role %s not found, skipping permission seeding", roleName))
			continue
		}

		existingPermissions, err := rolePermissionRepo.ReadByRoleID(ctx, role.ID.Hex())
		if err != nil {
			return err
		}

		existingMap := make(map[string]bool)
		for _, rp := range existingPermissions {
			existingMap[rp.PermissionID.Hex()] = true
		}

		for _, permName := range permissionNames {
			perm, err := permissionRepo.ReadByName(ctx, permName)
			if err != nil {
				return err
			}
			if perm == nil {
				logger.Log.Warn(fmt.Sprintf("skipping role permission seeding: Permission %s not found", permName))
				continue
			}

			if !existingMap[perm.ID.Hex()] {
				rolePermission := model.RolePermission{
					RoleID:       role.ID,
					PermissionID: perm.ID,
				}
				if err := rolePermissionRepo.Create(ctx, &rolePermission); err != nil {
					return err
				}
			} else {
				logger.Log.Info(fmt.Sprintf("skipping role permission: %s -> %s (already exists)", roleName, permName))
			}
		}
	}
	return nil
}

// seedPolicies inserts a default organization-level policy
func seedPolicies(ctx context.Context, policyRepo *repository.PolicyRepository) error {
	defaultPolicy := model.Policy{
		Name: "default_organization_access",
	}

	existing, err := policyRepo.ReadByName(ctx, defaultPolicy.Name)
	if err != nil {
		return err
	}
	if existing == nil {
		if err := policyRepo.Create(ctx, &defaultPolicy); err != nil {
			return err
		}
	} else {
		logger.Log.Info(fmt.Sprintf("skipping policy: %s (already exists)", defaultPolicy.Name))
	}

	return nil
}

// seedPolicyPermissions assigns permissions to the default policy
func seedPolicyPermissions(
	ctx context.Context,
	policyRepo *repository.PolicyRepository,
	permissionRepo *repository.PermissionRepository,
	policyPermissionRepo *repository.PolicyPermissionRepository,
) error {
	defaultPolicy, err := policyRepo.ReadByName(ctx, "default_organization_access")
	if err != nil {
		return err
	}
	if defaultPolicy == nil {
		logger.Log.Warn("default policy not found, skipping policy permission seeding")
		return nil
	}

	permissions := []string{
		"list:users", "create:user", "read:user", "update:user",
		"list:organizations", "create:organization", "read:organization", "update:organization",
	}

	for _, permName := range permissions {
		perm, err := permissionRepo.ReadByName(ctx, permName)
		if err != nil {
			return err
		}
		if perm == nil {
			logger.Log.Warn(fmt.Sprintf("skipping policy permission seeding: Permission %s not found", permName))
			continue
		}

		existing, err := policyPermissionRepo.ReadByPolicyIDPermissionID(ctx, defaultPolicy.ID.Hex(), perm.ID.Hex())
		if err != nil {
			return err
		}
		if existing == nil {
			policyPermission := model.PolicyPermission{
				PolicyID:     defaultPolicy.ID,
				PermissionID: perm.ID,
			}
			if err := policyPermissionRepo.Create(ctx, &policyPermission); err != nil {
				return err
			}
		} else {
			logger.Log.Info(fmt.Sprintf("skipping policy permission: %s -> %s (already exists)", defaultPolicy.Name, permName))
		}
	}

	return nil
}
