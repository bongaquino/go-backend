package repositories

import (
	"context"
	"time"

	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// ProfileRepository handles database operations for the Profile model.
type ProfileRepository struct {
	collection *mongoDriver.Collection
}

// NewProfileRepository initializes a new ProfileRepository.
func NewProfileRepository(mongoService *services.MongoService) *ProfileRepository {
	db := mongoService.GetDB()
	return &ProfileRepository{
		collection: db.Collection("profiles"),
	}
}

// CreateProfile inserts a new profile into the database.
func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, profile)
	if err != nil {
		logger.Log.Error("error creating profile", logger.Error(err))
		return err
	}
	return nil
}

// ReadProfileByUserID retrieves a profile by user ID.
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

// UpdateProfile updates an existing profile.
func (r *ProfileRepository) UpdateProfile(ctx context.Context, userID string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating profile", logger.Error(err))
		return err
	}
	return nil
}

// DeleteProfile removes a profile from the database.
func (r *ProfileRepository) DeleteProfile(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Log.Error("error deleting profile", logger.Error(err))
		return err
	}
	return nil
}
