package repository

import (
	"context"
	"time"

	"bongaquino/server/app/model"
	"bongaquino/server/app/provider"
	"bongaquino/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type PolicyPermissionRepository struct {
	collection *mongoDriver.Collection
}

func NewPolicyPermissionRepository(mongoProvider *provider.MongoProvider) *PolicyPermissionRepository {
	db := mongoProvider.GetDB()
	return &PolicyPermissionRepository{
		collection: db.Collection("policy_permission"),
	}
}

func (r *PolicyPermissionRepository) Create(ctx context.Context, policyPermission *model.PolicyPermission) error {
	policyPermission.ID = primitive.NewObjectID()
	policyPermission.CreatedAt = time.Now()
	policyPermission.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, policyPermission)
	if err != nil {
		logger.Log.Error("error creating policy permission", logger.Error(err))
		return err
	}
	return nil
}

func (r *PolicyPermissionRepository) ReadByPolicyID(ctx context.Context, policyID string) ([]model.PolicyPermission, error) {
	var results []model.PolicyPermission

	cursor, err := r.collection.Find(ctx, bson.M{"policy_id": policyID})
	if err != nil {
		logger.Log.Error("error retrieving policy permissions", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		logger.Log.Error("error decoding policy permissions", logger.Error(err))
		return nil, err
	}

	return results, nil
}

func (r *PolicyPermissionRepository) ReadByPermissionID(ctx context.Context, permissionID string) ([]model.PolicyPermission, error) {
	var results []model.PolicyPermission

	cursor, err := r.collection.Find(ctx, bson.M{"permission_id": permissionID})
	if err != nil {
		logger.Log.Error("error retrieving permission policies", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		logger.Log.Error("error decoding permission policies", logger.Error(err))
		return nil, err
	}

	return results, nil
}

func (r *PolicyPermissionRepository) ReadByPolicyIDPermissionID(ctx context.Context, policyID string, permissionID string) (*model.PolicyPermission, error) {
	// Convert policyID to ObjectID
	policyObjectID, err := primitive.ObjectIDFromHex(policyID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	// Convert permissionID to ObjectID
	permissionObjectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var result model.PolicyPermission
	err = r.collection.FindOne(ctx, bson.M{"policy_id": policyObjectID, "permission_id": permissionObjectID}).Decode(&result)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error retrieving policy permission", logger.Error(err))
		return nil, err
	}
	return &result, nil
}

func (r *PolicyPermissionRepository) Delete(ctx context.Context, policyID, permissionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"policy_id": policyID, "permission_id": permissionID})
	if err != nil {
		logger.Log.Error("error deleting policy permission", logger.Error(err))
		return err
	}
	return nil
}
