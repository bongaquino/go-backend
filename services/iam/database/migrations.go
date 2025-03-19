package database

import (
	"context"
	"fmt"
	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services/mongo"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

// MigrateCollections creates or updates collections for each model
func MigrateCollections(mongoService *mongo.MongoService) {
	db := mongoService.GetDB()
	ctx := context.Background()

	collections := []struct {
		Name    string
		Indexes []mongoDriver.IndexModel
	}{
		{
			Name: "users",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_email"),
				},
			},
		},
		{
			Name:    "profiles",
			Indexes: generateIndexModels(models.Profile{}.GetIndexes(), "unique_user_profile"),
		},
		{
			Name:    "roles",
			Indexes: generateIndexModels(models.Role{}.GetIndexes(), "unique_roles"),
		},
		{
			Name:    "permissions",
			Indexes: generateIndexModels(models.Permission{}.GetIndexes(), "unique_permission_name"),
		},
		{
			Name:    "policies",
			Indexes: generateIndexModels(models.Policy{}.GetIndexes(), "unique_policy_name"),
		},
		{
			Name:    "policy_permissions",
			Indexes: generateIndexModels(models.PolicyPermission{}.GetIndexes(), "unique_policy_permission"),
		},
		{
			Name:    "role_permissions",
			Indexes: generateIndexModels(models.RolePermission{}.GetIndexes(), "unique_role_permission"),
		},
		{
			Name:    "user_roles",
			Indexes: generateIndexModels(models.UserRole{}.GetIndexes(), "unique_user_role"),
		},
		{
			Name:    "service_accounts",
			Indexes: generateIndexModels(models.ServiceAccount{}.GetIndexes(), "unique_client_id"),
		},
	}

	for _, collection := range collections {
		if err := ensureCollection(db, ctx, collection.Name, collection.Indexes); err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to migrate collection: %s", collection.Name), logger.Error(err))
		} else {
			logger.Log.Info(fmt.Sprintf("Migrated collection: %s", collection.Name))
		}
	}
}

// generateIndexModels converts multiple indexes from GetIndexes() into []mongoDriver.IndexModel
func generateIndexModels(indexes []bson.D, baseIndexName string) []mongoDriver.IndexModel {
	var indexModels []mongoDriver.IndexModel
	for i, index := range indexes {
		indexName := fmt.Sprintf("%s_%d", baseIndexName, i+1) // Ensure unique names
		indexModels = append(indexModels, mongoDriver.IndexModel{
			Keys:    index,
			Options: mongoOptions.Index().SetUnique(true).SetName(indexName),
		})
	}
	return indexModels
}

// ensureCollection ensures the collection exists and applies indexes
func ensureCollection(db *mongoDriver.Database, ctx context.Context, name string, indexes []mongoDriver.IndexModel) error {
	// Check if collection exists before trying to create it
	collectionNames, err := db.ListCollectionNames(ctx, bson.M{"name": name})
	if err != nil {
		return err
	}
	if len(collectionNames) == 0 {
		err = db.CreateCollection(ctx, name)
		if err != nil {
			return fmt.Errorf("failed to create collection %s: %w", name, err)
		}
	}

	// Apply indexes
	collection := db.Collection(name)
	if len(indexes) > 0 {
		_, err = collection.Indexes().CreateMany(ctx, indexes)
		if err != nil {
			return fmt.Errorf("failed to create indexes on %s: %w", name, err)
		}
	}
	return nil
}
