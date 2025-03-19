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
			Name: "profiles",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.Profile{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_user_profile"),
				},
			},
		},
		{
			Name: "roles",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.Role{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_roles"),
				},
			},
		},
		{
			Name: "permissions",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.Permission{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_permission_name"),
				},
			},
		},
		{
			Name: "policies",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.Policy{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_policy_name"),
				},
			},
		},
		{
			Name: "policy_permissions",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.PolicyPermission{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_policy_permission"),
				},
			},
		},
		{
			Name: "role_permissions",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.RolePermission{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_role_permission"),
				},
			},
		},
		{
			Name: "user_roles",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.UserRole{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_user_role"),
				},
			},
		},
		{
			Name: "service_accounts",
			Indexes: []mongoDriver.IndexModel{
				{
					Keys:    models.ServiceAccount{}.GetIndexes(),
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_client_id"),
				},
			},
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
