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

type UserRoleRepository struct {
	collection *mongoDriver.Collection
}

func NewUserRoleRepository(mongoProvider *providers.MongoProvider) *UserRoleRepository {
	db := mongoProvider.GetDB()
	return &UserRoleRepository{
		collection: db.Collection("user_roles"),
	}
}

func (r *UserRoleRepository) CreateUserRole(ctx context.Context, userRole *models.UserRole) error {
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

func (r *UserRoleRepository) ReadUserRoles(ctx context.Context, userID string) ([]models.UserRole, error) {
	var results []models.UserRole

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
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

func (r *UserRoleRepository) ReadUsersByRole(ctx context.Context, roleID string) ([]models.UserRole, error) {
	var results []models.UserRole

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

func (r *UserRoleRepository) DeleteUserRole(ctx context.Context, userID, roleID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID, "role_id": roleID})
	if err != nil {
		logger.Log.Error("error deleting user role", logger.Error(err))
		return err
	}
	return nil
}
