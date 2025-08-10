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

type PermissionRepository struct {
	collection *mongoDriver.Collection
}

func NewPermissionRepository(mongoProvider *provider.MongoProvider) *PermissionRepository {
	db := mongoProvider.GetDB()
	return &PermissionRepository{
		collection: db.Collection("permissions"),
	}
}

func (r *PermissionRepository) List(ctx context.Context) ([]*model.Permission, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Log.Error("error listing permissions", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*model.Permission
	for cursor.Next(ctx) {
		var policy model.Permission
		if err := cursor.Decode(&policy); err != nil {
			logger.Log.Error("error decoding policy", logger.Error(err))
			return nil, err
		}
		permissions = append(permissions, &policy)
	}

	if err := cursor.Err(); err != nil {
		logger.Log.Error("error iterating over permissions", logger.Error(err))
		return nil, err
	}

	return permissions, nil
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
