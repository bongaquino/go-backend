package database

import (
	"context"
	"fmt"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
		{"policy_permissions", func(ctx context.Context) error {
			return seedPolicyPermissions(ctx, policyRepo, permissionRepo, policyPermissionRepo)
		}},
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
		{Name: "user:browse"},
		{Name: "user:read"},
		{Name: "user:add"},
		{Name: "user:edit"},
		{Name: "user:delete"},
		{Name: "organization:browse"},
		{Name: "organization:read"},
		{Name: "organization:edit"},
		{Name: "organization:add"},
		{Name: "organization:delete"},
		{Name: "directory:browse"},
		{Name: "directory:read"},
		{Name: "directory:edit"},
		{Name: "directory:add"},
		{Name: "directory:delete"},
		{Name: "file:upload"},
		{Name: "file:download"},
		{Name: "file:read"},
		{Name: "file:edit"},
		{Name: "file:delete"},
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
			"user:browse", "user:add", "user:read", "user:edit", "user:delete",
			"organization:browse", "organization:add", "organization:read", "organization:edit", "organization:delete",
			"directory:browse", "directory:add", "directory:read", "directory:edit", "directory:delete",
			"file:upload", "file:download", "file:read", "file:edit", "file:delete",
		},
		"system_user": {
			"directory:browse", "directory:add", "directory:read", "directory:edit", "directory:delete",
			"file:upload", "file:download", "file:read", "file:edit", "file:delete",
		},
		"organization_admin": {
			"organization:browse", "organization:add", "organization:read", "organization:edit", "organization:delete",
			"directory:browse", "directory:add", "directory:read", "directory:edit", "directory:delete",
			"file:upload", "file:download", "file:read", "file:edit", "file:delete",
		},
		"organization_user": {
			"directory:browse", "directory:add", "directory:read", "directory:edit", "directory:delete",
			"file:upload", "file:download", "file:read", "file:edit", "file:delete",
		},
		"organization_viewer": {
			"organization:browse", "organization:read",
			"directory:browse", "directory:read",
			"file:download", "file:read",
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

// seedPolicies inserts default policies
func seedPolicies(ctx context.Context, policyRepo *repository.PolicyRepository) error {
	policies := []model.Policy{
		{Name: "default_organization_policy"},
		{Name: "default_service_account_policy"},
	}

	for _, policy := range policies {
		existing, err := policyRepo.ReadByName(ctx, policy.Name)
		if err != nil {
			return err
		}
		if existing == nil {
			if err := policyRepo.Create(ctx, &policy); err != nil {
				return err
			}
		} else {
			logger.Log.Info(fmt.Sprintf("skipping policy: %s (already exists)", policy.Name))
		}
	}

	return nil
}

// seedPolicyPermissions assigns permissions to the policies
func seedPolicyPermissions(
	ctx context.Context,
	policyRepo *repository.PolicyRepository,
	permissionRepo *repository.PermissionRepository,
	policyPermissionRepo *repository.PolicyPermissionRepository,
) error {
	// default_organization_policy → org perms only
	orgPolicy, err := policyRepo.ReadByName(ctx, "default_organization_policy")
	if err != nil {
		return err
	}
	if orgPolicy != nil {
		orgPerms := []string{
			"organization:browse", "organization:add", "organization:read", "organization:edit", "organization:delete",
		}
		if err := assignPolicyPermissions(ctx, orgPolicy.ID.Hex(), orgPerms, permissionRepo, policyPermissionRepo); err != nil {
			return err
		}
	} else {
		logger.Log.Warn("default_organization_policy not found, skipping")
	}

	// default_service_account_policy → directory + file perms
	servicePolicy, err := policyRepo.ReadByName(ctx, "default_service_account_policy")
	if err != nil {
		return err
	}
	if servicePolicy != nil {
		dirFilePerms := []string{
			"directory:browse", "directory:add", "directory:read", "directory:edit", "directory:delete",
			"file:upload", "file:download", "file:read", "file:edit", "file:delete",
		}
		if err := assignPolicyPermissions(ctx, servicePolicy.ID.Hex(), dirFilePerms, permissionRepo, policyPermissionRepo); err != nil {
			return err
		}
	} else {
		logger.Log.Warn("default_service_account_policy not found, skipping")
	}

	return nil
}

func assignPolicyPermissions(
	ctx context.Context,
	policyID string,
	permNames []string,
	permissionRepo *repository.PermissionRepository,
	policyPermissionRepo *repository.PolicyPermissionRepository,
) error {
	for _, permName := range permNames {
		perm, err := permissionRepo.ReadByName(ctx, permName)
		if err != nil {
			return err
		}
		if perm == nil {
			logger.Log.Warn(fmt.Sprintf("permission %s not found", permName))
			continue
		}

		existing, err := policyPermissionRepo.ReadByPolicyIDPermissionID(ctx, policyID, perm.ID.Hex())
		if err != nil {
			return err
		}
		if existing == nil {
			objID, err := primitive.ObjectIDFromHex(policyID)
			if err != nil {
				return fmt.Errorf("invalid policyID hex: %w", err)
			}
			if err := policyPermissionRepo.Create(ctx, &model.PolicyPermission{
				PolicyID:     objID,
				PermissionID: perm.ID,
			}); err != nil {
				return err
			}
		} else {
			logger.Log.Info("policy permission already exists")
		}
	}
	return nil
}
