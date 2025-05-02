package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"koneksi-backend/server/app/model"
)

type OrganizationRepository struct {
	collection *mongo.Collection
}

func NewOrganizationRepository(db *mongo.Database) *OrganizationRepository {
	return &OrganizationRepository{
		collection: db.Collection("organization"),
	}
}

func (r *OrganizationRepository) Create(ctx context.Context, organization *model.Organization) error {
	_, err := r.collection.InsertOne(ctx, organization)
	return err
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Organization, error) {
	var organization model.Organization
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *OrganizationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}