package repository

import (
	"context"
	"time"

	"koneksi/server/app/model"
	"koneksi/server/app/provider"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type PermissionRepository struct {
	collection *mongoDriver.Collection
}

func NewPermissionRepository(mongoProvider *provider.MongoProvider) *PermissionRepository {
	db := mongoProvider.GetDB()
	return &PermissionRepository{
		collection: db.Collection("permissions"),
	}
}

func (r *PermissionRepository) Create(ctx context.Context, permission *model.Permission) error {
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

func (r *PermissionRepository) ReadByName(ctx context.Context, name string) (*model.Permission, error) {
	var permission model.Permission
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

func (r *PermissionRepository) Update(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating permission", logger.Error(err))
		return err
	}
	return nil
}

func (r *PermissionRepository) Delete(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logger.Log.Error("error deleting permission", logger.Error(err))
		return err
	}
	return nil
}
