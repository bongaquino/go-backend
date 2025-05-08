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

type OrganizationUserAccessRepository struct {
	collection *mongo.Collection
}

func NewOrganizationUserAccessRepository(mongoProvider *provider.MongoProvider) *OrganizationUserAccessRepository {
	return &OrganizationUserAccessRepository{
		collection: mongoProvider.GetDB().Collection("organization_user_access"),
	}
}

func (r *OrganizationUserAccessRepository) Create(ctx context.Context, orgUserAccess *model.OrganizationUserAccess) error {
	orgUserAccess.ID = primitive.NewObjectID()
	orgUserAccess.CreatedAt = time.Now()
	orgUserAccess.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, orgUserAccess)
	if err != nil {
		logger.Log.Error("error creating organization user access", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationUserAccessRepository) FindByOrganizationID(ctx context.Context, organizationID string) ([]model.OrganizationUserAccess, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"organization_id": organizationID})
	if err != nil {
		logger.Log.Error("error finding organization user access by organization ID", logger.Error(err))
		return nil, err
	}
	var accesses []model.OrganizationUserAccess
	if err := cursor.All(ctx, &accesses); err != nil {
		logger.Log.Error("error decoding organization user access", logger.Error(err))
		return nil, err
	}
	return accesses, nil
}

func (r *OrganizationUserAccessRepository) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating organization user access", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationUserAccessRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Log.Error("error deleting organization user access", logger.Error(err))
		return err
	}
	return nil
}
