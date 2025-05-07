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

type UserRoleRepository struct {
	collection *mongoDriver.Collection
}

func NewUserRoleRepository(mongoProvider *provider.MongoProvider) *UserRoleRepository {
	db := mongoProvider.GetDB()
	return &UserRoleRepository{
		collection: db.Collection("user_roles"),
	}
}

func (r *UserRoleRepository) Create(ctx context.Context, userRole *model.UserRole) error {
	userRole.ID = primitive.NewObjectID()

	userRole.CreatedAt = time.Now()
	userRole.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, userRole)
	if err != nil {
		logger.Log.Error("error creating user role", logger.Error(err))
		return err
	}
	return nil
}

func (r *UserRoleRepository) ReadByUserID(ctx context.Context, userID string) ([]model.UserRole, error) {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	// Initialize an empty slice to hold the results
	var results []model.UserRole

	// Use the Find method to retrieve user roles
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		logger.Log.Error("error retrieving user roles", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		logger.Log.Error("error decoding user roles", logger.Error(err))
		return nil, err
	}

	return results, nil
}

func (r *UserRoleRepository) ReadByRoleID(ctx context.Context, roleID string) ([]model.UserRole, error) {
	var results []model.UserRole

	cursor, err := r.collection.Find(ctx, bson.M{"role_id": roleID})
	if err != nil {
		logger.Log.Error("error retrieving users by role", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		logger.Log.Error("error decoding users by role", logger.Error(err))
		return nil, err
	}

	return results, nil
}

func (r *UserRoleRepository) Update(ctx context.Context, userID string, update bson.M) error {
	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	// Set the updated_at field to the current time
	update["updated_at"] = time.Now()

	// Update the user role in the database
	_, err = r.collection.UpdateOne(ctx, bson.M{"user_id": userObjectID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating user role", logger.Error(err))
		return err
	}
	return nil
}

func (r *UserRoleRepository) Delete(ctx context.Context, userID, roleID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID, "role_id": roleID})
	if err != nil {
		logger.Log.Error("error deleting user role", logger.Error(err))
		return err
	}
	return nil
}
