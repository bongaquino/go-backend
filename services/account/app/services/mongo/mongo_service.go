package mongo

import (
	"context"
	"time"

	"koneksi/orchestrator/config"
	"koneksi/orchestrator/core/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoService initializes a new MongoService
func NewMongoService() *MongoService {
	mongoConfig := config.LoadMongoConfig()

	clientOptions := options.Client().ApplyURI(mongoConfig.GetMongoUri())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Log.Fatal("database connection error", logger.Error(err))
	}

	return &MongoService{
		client: client,
		db:     client.Database(mongoConfig.MongoDatabase),
	}
}

// GetDB retrieves the MongoDB database instance
func (m *MongoService) GetDB() *mongo.Database {
	if m.db == nil {
		logger.Log.Fatal("database not initialized")
	}
	return m.db
}
