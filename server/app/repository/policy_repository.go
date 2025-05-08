package repository

import (
	"context"
	"time"

	"koneksi/server/app/model"
	"koneksi/server/app/provider"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type PolicyRepository struct {
	collection *mongoDriver.Collection
}

func NewPolicyRepository(mongoProvider *provider.MongoProvider) *PolicyRepository {
	db := mongoProvider.GetDB()
	return &PolicyRepository{
		collection: db.Collection("policies"),
	}
}

func (r *PolicyRepository) List(ctx context.Context) ([]*model.Policy, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Log.Error("error listing policies", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var policies []*model.Policy
	for cursor.Next(ctx) {
		var policy model.Policy
		if err := cursor.Decode(&policy); err != nil {
			logger.Log.Error("error decoding policy", logger.Error(err))
			return nil, err
		}
		policies = append(policies, &policy)
	}

	if err := cursor.Err(); err != nil {
		logger.Log.Error("error iterating over policies", logger.Error(err))
		return nil, err
	}

	return policies, nil
}

func (r *PolicyRepository) Create(ctx context.Context, policy *model.Policy) error {
	policy.ID = primitive.NewObjectID()
	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, policy)
	if err != nil {
		logger.Log.Error("error creating policy", logger.Error(err))
		return err
	}
	return nil
}

func (r *PolicyRepository) ReadByName(ctx context.Context, name string) (*model.Policy, error) {
	var policy model.Policy
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&policy)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading policy by name", logger.Error(err))
		return nil, err
	}
	return &policy, nil
}

func (r *PolicyRepository) Update(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating policy", logger.Error(err))
		return err
	}
	return nil
}

func (r *PolicyRepository) Delete(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logger.Log.Error("error deleting policy", logger.Error(err))
		return err
	}
	return nil
}
