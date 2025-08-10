package repository

import (
	"context"
	"time"

	"bongaquino/server/app/model"
	"bongaquino/server/app/provider"
	"bongaquino/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type RolePermissionRepository struct {
	collection *mongoDriver.Collection
}

func NewRolePermissionRepository(mongoProvider *provider.MongoProvider) *RolePermissionRepository {
	db := mongoProvider.GetDB()
	return &RolePermissionRepository{
		collection: db.Collection("role_permission"),
	}
}

func (r *RolePermissionRepository) Create(ctx context.Context, rolePermission *model.RolePermission) error {
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

func (r *RolePermissionRepository) ReadByRoleID(ctx context.Context, roleID string) ([]model.RolePermission, error) {
	var results []model.RolePermission

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

func (r *RolePermissionRepository) ReadByPermissionID(ctx context.Context, permissionID string) ([]model.RolePermission, error) {
	var results []model.RolePermission

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

func (r *RolePermissionRepository) Delete(ctx context.Context, roleID, permissionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"role_id": roleID, "permission_id": permissionID})
	if err != nil {
		logger.Log.Error("error deleting role permission", logger.Error(err))
		return err
	}
	return nil
}
