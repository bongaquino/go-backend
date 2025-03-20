package repositories

import (
	"context"
	"time"

	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// UserRoleRepository handles database operations for the UserRole model.
type UserRoleRepository struct {
	collection *mongoDriver.Collection
}

// NewUserRoleRepository initializes a new UserRoleRepository.
func NewUserRoleRepository(mongoService *services.MongoService) *UserRoleRepository {
	db := mongoService.GetDB()
	return &UserRoleRepository{
		collection: db.Collection("user_roles"),
	}
}

// CreateUserRole inserts a new user-role relationship into the database.
func (r *UserRoleRepository) CreateUserRole(ctx context.Context, userRole *models.UserRole) error {
	// Generate a new ObjectID for the userRole
	userRole.ID = primitive.NewObjectID()

	// Set timestamps
	userRole.CreatedAt = time.Now()
	userRole.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, userRole)
	if err != nil {
		logger.Log.Error("error creating user role", logger.Error(err))
		return err
	}
	return nil
}

// ReadUserRoles retrieves all roles associated with a user.
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

// ReadUsersByRole retrieves all users associated with a role.
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

// DeleteUserRole removes a specific user-role relationship.
func (r *UserRoleRepository) DeleteUserRole(ctx context.Context, userID, roleID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID, "role_id": roleID})
	if err != nil {
		logger.Log.Error("error deleting user role", logger.Error(err))
		return err
	}
	return nil
}
