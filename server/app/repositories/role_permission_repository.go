package repositories

import (
	"context"
	"time"

	"koneksi/server/app/models"
	"koneksi/server/app/services"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type RolePermissionRepository struct {
	collection *mongoDriver.Collection
}

func NewRolePermissionRepository(mongoService *services.MongoService) *RolePermissionRepository {
	db := mongoService.GetDB()
	return &RolePermissionRepository{
		collection: db.Collection("role_permissions"),
	}
}

func (r *RolePermissionRepository) CreateRolePermission(ctx context.Context, rolePermission *models.RolePermission) error {
	rolePermission.ID = primitive.NewObjectID()

	rolePermission.CreatedAt = time.Now()
	rolePermission.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, rolePermission)
	if err != nil {
		logger.Log.Error("error creating role permission", logger.Error(err))
		return err
	}
	return nil
}

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

func (r *RolePermissionRepository) DeleteRolePermission(ctx context.Context, roleID, permissionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"role_id": roleID, "permission_id": permissionID})
	if err != nil {
		logger.Log.Error("error deleting role permission", logger.Error(err))
		return err
	}
	return nil
}
