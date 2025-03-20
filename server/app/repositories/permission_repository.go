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

type PermissionRepository struct {
	collection *mongoDriver.Collection
}

func NewPermissionRepository(mongoService *services.MongoService) *PermissionRepository {
	db := mongoService.GetDB()
	return &PermissionRepository{
		collection: db.Collection("permissions"),
	}
}

func (r *PermissionRepository) CreatePermission(ctx context.Context, permission *models.Permission) error {
	permission.ID = primitive.NewObjectID()

	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, permission)
	if err != nil {
		logger.Log.Error("error creating permission", logger.Error(err))
		return err
	}
	return nil
}

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

func (r *PermissionRepository) UpdatePermission(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating permission", logger.Error(err))
		return err
	}
	return nil
}

func (r *PermissionRepository) DeletePermission(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logger.Log.Error("error deleting permission", logger.Error(err))
		return err
	}
	return nil
}
