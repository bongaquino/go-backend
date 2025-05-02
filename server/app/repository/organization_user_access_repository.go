package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"koneksi-backend/server/app/model"
)

type OrganizationUserAccessRepository struct {
	collection *mongo.Collection
}

func NewOrganizationUserAccessRepository(db *mongo.Database) *OrganizationUserAccessRepository {
	return &OrganizationUserAccessRepository{
		collection: db.Collection("organization_user_access"),
	}
}

func (r *OrganizationUserAccessRepository) Create(ctx context.Context, access *model.OrganizationUserAccess) error {
	_, err := r.collection.InsertOne(ctx, access)
	return err
}

func (r *OrganizationUserAccessRepository) FindByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]model.OrganizationUserAccess, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"organization_id": organizationID})
	if err != nil {
		return nil, err
	}
	var accesses []model.OrganizationUserAccess
	if err := cursor.All(ctx, &accesses); err != nil {
		return nil, err
	}
	return accesses, nil
}

func (r *OrganizationUserAccessRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *OrganizationUserAccessRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}