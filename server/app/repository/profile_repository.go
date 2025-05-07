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

type ProfileRepository struct {
	collection *mongoDriver.Collection
}

func NewProfileRepository(mongoProvider *provider.MongoProvider) *ProfileRepository {
	db := mongoProvider.GetDB()
	return &ProfileRepository{
		collection: db.Collection("profiles"),
	}
}

func (r *ProfileRepository) Create(ctx context.Context, profile *model.Profile) error {
	profile.ID = primitive.NewObjectID()

	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, profile)
	if err != nil {
		logger.Log.Error("error creating profile", logger.Error(err))
		return err
	}
	return nil
}

func (r *ProfileRepository) ReadByUserID(ctx context.Context, userID string) (*model.Profile, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var profile model.Profile
	err = r.collection.FindOne(ctx, bson.M{"user_id": objectID}).Decode(&profile)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading profile by user ID", logger.Error(err))
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) UpdateByUserID(ctx context.Context, userID string, update bson.M) error {
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
		logger.Log.Error("error updating profile", logger.Error(err))
		return err
	}
	return nil
}

func (r *ProfileRepository) Delete(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Log.Error("error deleting profile", logger.Error(err))
		return err
	}
	return nil
}
