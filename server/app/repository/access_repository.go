package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"koneksi-backend/server/app/model"
)

type AccessRepository struct {
	collection *mongo.Collection
}

func NewAccessRepository(db *mongo.Database) *AccessRepository {
	return &AccessRepository{
		collection: db.Collection("access"),
	}
}

func (r *AccessRepository) Create(ctx context.Context, access *model.Access) error {
	_, err := r.collection.InsertOne(ctx, access)
	return err
}

func (r *AccessRepository) FindByName(ctx context.Context, name string) (*model.Access, error) {
	var access model.Access
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&access)
	if err != nil {
		return nil, err
	}
	return &access, nil
}

func (r *AccessRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *AccessRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}