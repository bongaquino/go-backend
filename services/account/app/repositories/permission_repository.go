package repositories

import (
	"context"
	"time"

	"koneksi/services/account/app/models"
	"koneksi/services/account/app/services/mongo"
	"koneksi/services/account/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// PermissionRepository handles database operations for the Permission model.
type PermissionRepository struct {
	collection *mongoDriver.Collection
}

// NewPermissionRepository initializes a new PermissionRepository.
func NewPermissionRepository(mongoService *mongo.MongoService) *PermissionRepository {
	db := mongoService.GetDB()
	return &PermissionRepository{
		collection: db.Collection("permissions"),
	}
}

// CreatePermission inserts a new permission into the database.
func (r *PermissionRepository) CreatePermission(ctx context.Context, permission *models.Permission) error {
	permission.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, permission)
	if err != nil {
		logger.Log.Error("error creating permission", logger.Error(err))
		return err
	}
	return nil
}

// ReadPermissionByName retrieves a permission by name.
func (r *PermissionRepository) ReadPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&permission)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading permission by name", logger.Error(err))
		return nil, err
	}
	return &permission, nil
}

// UpdatePermission updates an existing permission.
func (r *PermissionRepository) UpdatePermission(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating permission", logger.Error(err))
		return err
	}
	return nil
}

// DeletePermission removes a permission from the database.
func (r *PermissionRepository) DeletePermission(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logger.Log.Error("error deleting permission", logger.Error(err))
		return err
	}
	return nil
}
