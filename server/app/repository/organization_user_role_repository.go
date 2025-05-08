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

type OrganizationUserRoleRepository struct {
	collection *mongo.Collection
}

func NewOrganizationUserRoleRepository(mongoProvider *provider.MongoProvider) *OrganizationUserRoleRepository {
	return &OrganizationUserRoleRepository{
		collection: mongoProvider.GetDB().Collection("organization_user_role"),
	}
}

func (r *OrganizationUserRoleRepository) Create(ctx context.Context, orgUserRole *model.OrganizationUserRole) error {
	orgUserRole.ID = primitive.NewObjectID()
	orgUserRole.CreatedAt = time.Now()
	orgUserRole.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, orgUserRole)
	if err != nil {
		logger.Log.Error("error creating organization user role", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationUserRoleRepository) ReadByUserID(ctx context.Context, userID string) ([]model.OrganizationUserRole, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		logger.Log.Error("error finding organization user role by user ID", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var orgUserRole []model.OrganizationUserRole
	if err := cursor.All(ctx, &orgUserRole); err != nil {
		logger.Log.Error("error decoding organization user role", logger.Error(err))
		return nil, err
	}
	return orgUserRole, nil
}

func (r *OrganizationUserRoleRepository) ReadByOrganizationID(ctx context.Context, organizationID string) ([]model.OrganizationUserRole, error) {
	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"organization_id": objectID})
	if err != nil {
		logger.Log.Error("error finding organization user role by organization ID", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var orgUserRole []model.OrganizationUserRole
	if err := cursor.All(ctx, &orgUserRole); err != nil {
		logger.Log.Error("error decoding organization user role", logger.Error(err))
		return nil, err
	}
	return orgUserRole, nil
}

func (r *OrganizationUserRoleRepository) ReadByUserIDOrganizationID(ctx context.Context, userID, organizationID string) (*model.OrganizationUserRole, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	organizationObjectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return nil, err
	}

	var orgUserRole model.OrganizationUserRole
	err = r.collection.FindOne(ctx, bson.M{"user_id": userObjectID, "organization_id": organizationObjectID}).Decode(&orgUserRole)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error finding organization user role by user ID and organization ID", logger.Error(err))
		return nil, err
	}
	return &orgUserRole, nil
}

func (r *OrganizationUserRoleRepository) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating organization user role", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationUserRoleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Log.Error("error deleting organization user role", logger.Error(err))
		return err
	}
	return nil
}
