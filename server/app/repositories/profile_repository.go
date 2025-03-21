package repositories

import (
	"context"
	"time"

	"koneksi/server/app/models"
	"koneksi/server/app/providers"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepository struct {
	collection *mongoDriver.Collection
}

func NewProfileRepository(mongoProvider *providers.MongoProvider) *ProfileRepository {
	db := mongoProvider.GetDB()
	return &ProfileRepository{
		collection: db.Collection("profiles"),
	}
}

func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
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

func (r *ProfileRepository) ReadProfileByUserID(ctx context.Context, userID string) (*models.Profile, error) {
	var profile models.Profile
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&profile)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading profile by user ID", logger.Error(err))
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, userID string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating profile", logger.Error(err))
		return err
	}
	return nil
}

func (r *ProfileRepository) DeleteProfile(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Log.Error("error deleting profile", logger.Error(err))
		return err
	}
	return nil
}
