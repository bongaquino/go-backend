package repositories

import (
	"context"

	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services/mongo"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// RolePermissionRepository handles database operations for the RolePermission model.
type RolePermissionRepository struct {
	collection *mongoDriver.Collection
}

// NewRolePermissionRepository initializes a new RolePermissionRepository.
func NewRolePermissionRepository(mongoService *mongo.MongoService) *RolePermissionRepository {
	db := mongoService.GetDB()
	return &RolePermissionRepository{
		collection: db.Collection("role_permissions"),
	}
}

// CreateRolePermission inserts a new role-permission relationship into the database.
func (r *RolePermissionRepository) CreateRolePermission(ctx context.Context, rolePermission *models.RolePermission) error {
	_, err := r.collection.InsertOne(ctx, rolePermission)
	if err != nil {
		logger.Log.Error("error creating role permission", logger.Error(err))
		return err
	}
	return nil
}

// ReadRolePermissions retrieves all permissions associated with a role.
func (r *RolePermissionRepository) ReadRolePermissions(ctx context.Context, roleID string) ([]models.RolePermission, error) {
	var results []models.RolePermission

	cursor, err := r.collection.Find(ctx, bson.M{"role_id": roleID})
	if err != nil {
		logger.Log.Error("error retrieving role permissions", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		logger.Log.Error("error decoding role permissions", logger.Error(err))
		return nil, err
	}

	return results, nil
}

// ReadPermissionRoles retrieves all roles associated with a permission.
func (r *RolePermissionRepository) ReadPermissionRoles(ctx context.Context, permissionID string) ([]models.RolePermission, error) {
	var results []models.RolePermission

	cursor, err := r.collection.Find(ctx, bson.M{"permission_id": permissionID})
	if err != nil {
		logger.Log.Error("error retrieving permission roles", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		logger.Log.Error("error decoding permission roles", logger.Error(err))
		return nil, err
	}

	return results, nil
}

// DeleteRolePermission removes a specific role-permission relationship.
func (r *RolePermissionRepository) DeleteRolePermission(ctx context.Context, roleID, permissionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"role_id": roleID, "permission_id": permissionID})
	if err != nil {
		logger.Log.Error("error deleting role permission", logger.Error(err))
		return err
	}
	return nil
}
