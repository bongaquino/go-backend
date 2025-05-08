package repository

import (
	"context"
	"time"

	"koneksi/server/app/model"
	"koneksi/server/app/provider"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccessRepository struct {
	collection *mongo.Collection
}

func NewAccessRepository(mongoProvider *provider.MongoProvider) *AccessRepository {
	return &AccessRepository{
		collection: mongoProvider.GetDB().Collection("access"),
	}
}

func (r *AccessRepository) Create(ctx context.Context, access *model.Access) error {
	access.ID = primitive.NewObjectID()
	access.CreatedAt = time.Now()
	access.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, access)
	if err != nil {
		logger.Log.Error("error creating access", logger.Error(err))
		return err
	}
	return nil
}

func (r *AccessRepository) ReadByName(ctx context.Context, name string) (*model.Access, error) {
	var access model.Access
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&access)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error finding access by name", logger.Error(err))
		return nil, err
	}
	return &access, nil
}

func (r *AccessRepository) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating access", logger.Error(err))
		return err
	}
	return nil
}

func (r *AccessRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Log.Error("error deleting access", logger.Error(err))
		return err
	}
	return nil
}
