package repositories

import (
	"context"
	"time"

	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services/mongo"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// PolicyRepository handles database operations for the Policy model.
type PolicyRepository struct {
	collection *mongoDriver.Collection
}

// NewPolicyRepository initializes a new PolicyRepository.
func NewPolicyRepository(mongoService *mongo.MongoService) *PolicyRepository {
	db := mongoService.GetDB()
	return &PolicyRepository{
		collection: db.Collection("policies"),
	}
}

// CreatePolicy inserts a new policy into the database.
func (r *PolicyRepository) CreatePolicy(ctx context.Context, policy *models.Policy) error {
	policy.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, policy)
	if err != nil {
		logger.Log.Error("error creating policy", logger.Error(err))
		return err
	}
	return nil
}

// ReadPolicyByName retrieves a policy by name.
func (r *PolicyRepository) ReadPolicyByName(ctx context.Context, name string) (*models.Policy, error) {
	var policy models.Policy
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

// UpdatePolicy updates an existing policy.
func (r *PolicyRepository) UpdatePolicy(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating policy", logger.Error(err))
		return err
	}
	return nil
}

// DeletePolicy removes a policy from the database.
func (r *PolicyRepository) DeletePolicy(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logger.Log.Error("error deleting policy", logger.Error(err))
		return err
	}
	return nil
}
