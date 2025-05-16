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

type DirectoryRepository struct {
	collection *mongoDriver.Collection
}

func NewDirectoryRepository(mongoProvider *provider.MongoProvider) *DirectoryRepository {
	db := mongoProvider.GetDB()
	return &DirectoryRepository{
		collection: db.Collection("directories"),
	}
}

func (r *DirectoryRepository) ListByUserID(ctx context.Context, userID string) ([]*model.Directory, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid user ID format", logger.Error(err))
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		logger.Log.Error("error listing directories by userID", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var directories []*model.Directory
	for cursor.Next(ctx) {
		var directory model.Directory
		if err := cursor.Decode(&directory); err != nil {
			logger.Log.Error("error decoding directory", logger.Error(err))
			return nil, err
		}
		directories = append(directories, &directory)
	}

	if err := cursor.Err(); err != nil {
		logger.Log.Error("cursor error", logger.Error(err))
		return nil, err
	}

	return directories, nil
}

func (r *DirectoryRepository) Create(ctx context.Context, directory *model.Directory) error {
	directory.ID = primitive.NewObjectID()
	directory.CreatedAt = time.Now()
	directory.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, directory)
	if err != nil {
		logger.Log.Error("error creating directory", logger.Error(err))
		return err
	}
	return nil
}

func (r *DirectoryRepository) Read(ctx context.Context, id string) (*model.Directory, error) {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var directory model.Directory
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&directory)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading directory", logger.Error(err))
		return nil, err
	}
	return &directory, nil
}

func (r *DirectoryRepository) Update(ctx context.Context, id string, update bson.M) error {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	// Set the updated time
	update["updated_at"] = time.Now()

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating directory by userID", logger.Error(err))
		return err
	}
	return nil
}
