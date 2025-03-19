package database

import (
	"context"
	"fmt"
	"koneksi/services/account/app/models"
	"koneksi/services/account/app/services/mongo"
	"koneksi/services/account/core/logger"

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
					Keys:    map[string]any{"email": 1},
					Options: mongoOptions.Index().SetUnique(true).SetName("unique_email"),
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
		// Add more collections and their indexes here
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
	collectionNames, err := db.ListCollectionNames(ctx, map[string]interface{}{"name": name})
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
