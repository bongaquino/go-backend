package repositories

import (
	"context"
	"time"

	"koneksi/services/account/app/models"
	"koneksi/services/account/app/services/mongo"
	"koneksi/services/account/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// RoleRepository handles database operations for the Role model.
type RoleRepository struct {
	collection *mongoDriver.Collection
}

// NewRoleRepository initializes a new RoleRepository.
func NewRoleRepository(mongoService *mongo.MongoService) *RoleRepository {
	db := mongoService.GetDB()
	return &RoleRepository{
		collection: db.Collection("roles"),
	}
}

// CreateRole inserts a new role into the database.
func (r *RoleRepository) CreateRole(ctx context.Context, role *models.Role) error {
	role.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		logger.Log.Error("error creating role", logger.Error(err))
		return err
	}
	return nil
}

// ReadRoleByName retrieves a role by name.
func (r *RoleRepository) ReadRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading role by name", logger.Error(err))
		return nil, err
	}
	return &role, nil
}

// UpdateRole updates an existing role.
func (r *RoleRepository) UpdateRole(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating role", logger.Error(err))
		return err
	}
	return nil
}

// DeleteRole removes a role from the database.
func (r *RoleRepository) DeleteRole(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logger.Log.Error("error deleting role", logger.Error(err))
		return err
	}
	return nil
}
