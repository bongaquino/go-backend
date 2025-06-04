package repository

import (
	"context"
	"fmt"
	"time"

	"koneksi/server/app/model"
	"koneksi/server/app/provider"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type SettingRepository struct {
	collection *mongoDriver.Collection
}

func NewSettingRepository(mongoProvider *provider.MongoProvider) *SettingRepository {
	db := mongoProvider.GetDB()
	return &SettingRepository{
		collection: db.Collection("settings"),
	}
}

func (r *SettingRepository) Create(ctx context.Context, setting *model.Setting) error {
	setting.ID = primitive.NewObjectID()
	setting.CreatedAt = time.Now()
	setting.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, setting)
	if err != nil {
		logger.Log.Error("error creating settings", logger.Error(err))
		return err
	}
	return nil
}

func (r *SettingRepository) ReadByUserID(ctx context.Context, userID string) (*model.Setting, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var setting model.Setting
	err = r.collection.FindOne(ctx, bson.M{"user_id": objectID}).Decode(&setting)

	fmt.Println(err)

	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading settings by user ID", logger.Error(err))
		return nil, err
	}

	return &setting, nil
}

func (r *SettingRepository) UpdateByUserID(ctx context.Context, userID string, update bson.M) error {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	// Set the updated_at field to the current time
	update["updated_at"] = time.Now()

	_, err = r.collection.UpdateOne(ctx, bson.M{"user_id": objectID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating settings", logger.Error(err))
		return err
	}
	return nil
}

func (r *SettingRepository) Delete(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Log.Error("error deleting settings", logger.Error(err))
		return err
	}
	return nil
}
