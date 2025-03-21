package providers

import (
	"context"
	"time"

	"koneksi/server/config"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProvider struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoProvider initializes a new MongoProvider
func NewMongoProvider() *MongoProvider {
	mongoConfig := config.LoadMongoConfig()

	clientOptions := options.Client().ApplyURI(mongoConfig.GetMongoUri())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Log.Fatal("database connection error", logger.Error(err))
	}

	return &MongoProvider{
		client: client,
		db:     client.Database(mongoConfig.MongoDatabase),
	}
}

// GetDB retrieves the MongoDB database instance
func (m *MongoProvider) GetDB() *mongo.Database {
	if m.db == nil {
		logger.Log.Fatal("database not initialized")
	}
	return m.db
}
