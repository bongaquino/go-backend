package repositories

import (
	"context"
	"time"

	"koneksi/server/app/models"
	"koneksi/server/app/services"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type RoleRepository struct {
	collection *mongoDriver.Collection
}

func NewRoleRepository(mongoService *services.MongoService) *RoleRepository {
	db := mongoService.GetDB()
	return &RoleRepository{
		collection: db.Collection("roles"),
	}
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *models.Role) error {
	role.ID = primitive.NewObjectID()

	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		logger.Log.Error("error creating role", logger.Error(err))
		return err
	}
	return nil
}

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

func (r *RoleRepository) ReadRoleByID(ctx context.Context, roleID string) (*models.Role, error) {
	var role models.Role
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&role)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading role by ID", logger.Error(err))
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) UpdateRole(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating role", logger.Error(err))
		return err
	}
	return nil
}

func (r *RoleRepository) DeleteRole(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logger.Log.Error("error deleting role", logger.Error(err))
		return err
	}
	return nil
}
