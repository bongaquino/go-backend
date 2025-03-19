package database

import (
	"context"
	"fmt"
	"koneksi/services/account/app/models"
	"koneksi/services/account/app/services/mongo"
	"koneksi/services/account/core/logger"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// MigrateCollections creates or updates collections for each model
func MigrateCollections(mongoService *mongo.MongoService) {
	db := mongoService.GetDB()
	ctx := context.Background()

	collections := []struct {
		Name    string
		Indexes []mongo.IndexModel
	}{
		{
			Name: "users",
			Indexes: []mongo.IndexModel{
				{
					Keys:    models.User{}.GetIndexes(),
					Options: options.Index().SetUnique(true),
				},
			},
		},
		{
			Name: "roles",
			Indexes: []mongo.IndexModel{
				{
					Keys:    models.Role{}.GetIndexes(),
					Options: options.Index().SetUnique(true),
				},
			},
		},
		// Add more collections and their indexes here
	}

	for _, collection := range collections {
		err := ensureCollection(db, ctx, collection.Name, collection.Indexes)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to migrate collection: %s", collection.Name), logger.Error(err))
		} else {
			logger.Log.Info(fmt.Sprintf("Migrated collection: %s", collection.Name))
		}
	}
}

// ensureCollection ensures the collection exists and applies indexes
func ensureCollection(db *mongo.Database, ctx context.Context, name string, indexes []mongo.IndexModel) error {
	// Create collection if it doesn't exist
	err := db.CreateCollection(ctx, name)
	if err != nil && !mongo.IsCollectionExistsError(err) {
		return err
	}

	// Apply indexes
	collection := db.Collection(name)
	_, err = collection.Indexes().CreateMany(ctx, indexes)
	return err
}
