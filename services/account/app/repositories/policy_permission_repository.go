package repositories

import (
	"context"

	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services/mongo"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// PolicyPermissionRepository handles database operations for the PolicyPermission model.
type PolicyPermissionRepository struct {
	collection *mongoDriver.Collection
}

// NewPolicyPermissionRepository initializes a new PolicyPermissionRepository.
func NewPolicyPermissionRepository(mongoService *mongo.MongoService) *PolicyPermissionRepository {
	db := mongoService.GetDB()
	return &PolicyPermissionRepository{
		collection: db.Collection("policy_permissions"),
	}
}

// CreatePolicyPermission inserts a new policy-permission relationship into the database.
func (r *PolicyPermissionRepository) CreatePolicyPermission(ctx context.Context, policyPermission *models.PolicyPermission) error {
	_, err := r.collection.InsertOne(ctx, policyPermission)
	if err != nil {
		logger.Log.Error("error creating policy permission", logger.Error(err))
		return err
	}
	return nil
}

// ReadPolicyPermissions retrieves all permissions associated with a policy.
func (r *PolicyPermissionRepository) ReadPolicyPermissions(ctx context.Context, policyID string) ([]models.PolicyPermission, error) {
	var results []models.PolicyPermission

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

// ReadPermissionPolicies retrieves all policies associated with a permission.
func (r *PolicyPermissionRepository) ReadPermissionPolicies(ctx context.Context, permissionID string) ([]models.PolicyPermission, error) {
	var results []models.PolicyPermission

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

// DeletePolicyPermission removes a specific policy-permission relationship.
func (r *PolicyPermissionRepository) DeletePolicyPermission(ctx context.Context, policyID, permissionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"policy_id": policyID, "permission_id": permissionID})
	if err != nil {
		logger.Log.Error("error deleting policy permission", logger.Error(err))
		return err
	}
	return nil
}
