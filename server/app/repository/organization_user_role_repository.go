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
		collection: mongoProvider.GetDB().Collection("organization_user_access"),
	}
}

func (r *OrganizationUserRoleRepository) Create(ctx context.Context, orgUserRole *model.OrganizationUserRole) error {
	orgUserRole.ID = primitive.NewObjectID()
	orgUserRole.CreatedAt = time.Now()
	orgUserRole.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, orgUserRole)
	if err != nil {
		logger.Log.Error("error creating organization user access", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationUserRoleRepository) FindByOrganizationID(ctx context.Context, organizationID string) ([]model.OrganizationUserRole, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"organization_id": organizationID})
	if err != nil {
		logger.Log.Error("error finding organization user access by organization ID", logger.Error(err))
		return nil, err
	}
	var orgUserRole []model.OrganizationUserRole
	if err := cursor.All(ctx, &orgUserRole); err != nil {
		logger.Log.Error("error decoding organization user access", logger.Error(err))
		return nil, err
	}
	return orgUserRole, nil
}

func (r *OrganizationUserRoleRepository) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating organization user access", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationUserRoleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Log.Error("error deleting organization user access", logger.Error(err))
		return err
	}
	return nil
}
