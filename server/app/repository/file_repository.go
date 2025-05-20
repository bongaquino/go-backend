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

type FileRepository struct {
	collection *mongoDriver.Collection
}

func NewFileRepository(mongoProvider *provider.MongoProvider) *FileRepository {
	db := mongoProvider.GetDB()
	return &FileRepository{
		collection: db.Collection("files"),
	}
}

func (r *FileRepository) ListByUserID(ctx context.Context, userID string) ([]*model.File, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid user ID format", logger.Error(err))
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		logger.Log.Error("error listing files by userID", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []*model.File
	for cursor.Next(ctx) {
		var file model.File
		if err := cursor.Decode(&file); err != nil {
			logger.Log.Error("error decoding file", logger.Error(err))
			return nil, err
		}
		files = append(files, &file)
	}

	if err := cursor.Err(); err != nil {
		logger.Log.Error("cursor error", logger.Error(err))
		return nil, err
	}

	return files, nil
}

func (r *FileRepository) ListByDirectoryID(ctx context.Context, fileID string) ([]*model.File, error) {
	// Convert fileID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		logger.Log.Error("invalid file ID format", logger.Error(err))
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"file_id": objectID})
	if err != nil {
		logger.Log.Error("error listing files by fileID", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	var files []*model.File
	for cursor.Next(ctx) {
		var file model.File
		if err := cursor.Decode(&file); err != nil {
			logger.Log.Error("error decoding file", logger.Error(err))
			return nil, err
		}
		files = append(files, &file)
	}
	if err := cursor.Err(); err != nil {
		logger.Log.Error("cursor error", logger.Error(err))
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) ListByDirectoryIDUserID(ctx context.Context, directoryID string, userID string) ([]*model.File, error) {
	// Convert directoryID to ObjectID
	directoryObjectID, err := primitive.ObjectIDFromHex(directoryID)
	if err != nil {
		logger.Log.Error("invalid directory ID format", logger.Error(err))
		return nil, err
	}

	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid user ID format", logger.Error(err))
		return nil, err
	}

	// Create a filter to find files by directory ID and user ID
	filter := bson.M{
		"directory_id": directoryObjectID,
		"user_id":      userObjectID,
		"is_deleted":   false,
	}

	// Find files in the collection
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		logger.Log.Error("error listing files by directory ID and user ID", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	// Initialize a slice to hold the files
	var files []*model.File

	// Iterate through the cursor and decode each file
	for cursor.Next(ctx) {
		var file model.File
		if err := cursor.Decode(&file); err != nil {
			logger.Log.Error("error decoding file", logger.Error(err))
			return nil, err
		}
		files = append(files, &file)
	}

	// Check for any errors during iteration
	if err := cursor.Err(); err != nil {
		logger.Log.Error("cursor error", logger.Error(err))
		return nil, err
	}

	// Return the list of files
	return files, nil
}

func (r *FileRepository) Create(ctx context.Context, file *model.File) error {
	file.ID = primitive.NewObjectID()
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		logger.Log.Error("error creating file", logger.Error(err))
		return err
	}
	return nil
}

func (r *FileRepository) Read(ctx context.Context, id string) (*model.File, error) {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var file model.File
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&file)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading file", logger.Error(err))
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) Update(ctx context.Context, id string, update bson.M) error {
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
		logger.Log.Error("error updating file by userID", logger.Error(err))
		return err
	}
	return nil
}
