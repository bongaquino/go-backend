package repository

import (
	"context"
	"time"

	"koneksi/server/app/model"
	"koneksi/server/app/provider"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrganizationRepository struct {
	collection *mongo.Collection
}

func NewOrganizationRepository(mongoProvider *provider.MongoProvider) *OrganizationRepository {
	return &OrganizationRepository{
		collection: mongoProvider.GetDB().Collection("organization"),
	}
}

func (r *OrganizationRepository) List(ctx context.Context, page, limit int) ([]model.Organization, error) {
	var orgs []model.Organization

	// Calculate the number of documents to skip
	skip := (page - 1) * limit

	// Set options for pagination
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Log.Error("error fetching orgs", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orgs); err != nil {
		logger.Log.Error("error decoding orgs", logger.Error(err))
		return nil, err
	}

	return orgs, nil
}

func (r *OrganizationRepository) Create(ctx context.Context, organization *model.Organization) error {
	organization.CreatedAt = time.Now()
	organization.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, organization)
	if err != nil {
		logger.Log.Error("error creating organization", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id string) (*model.Organization, error) {
	var organization model.Organization
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&organization)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error finding organization by ID", logger.Error(err))
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating organization", logger.Error(err))
		return err
	}
	return nil
}

func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Log.Error("error deleting organization", logger.Error(err))
		return err
	}
	return nil
}
