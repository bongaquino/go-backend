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

type FileAccessRepository struct {
	collection *mongoDriver.Collection
}

func NewFileAccessRepository(mongoProvider *provider.MongoProvider) *FileAccessRepository {
	db := mongoProvider.GetDB()
	return &FileAccessRepository{
		collection: db.Collection("file_access"),
	}
}

func (r *FileAccessRepository) Create(ctx context.Context, fileAccess *model.FileAccess) error {
	fileAccess.ID = primitive.NewObjectID()
	fileAccess.CreatedAt = time.Now()
	fileAccess.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, fileAccess)
	if err != nil {
		logger.Log.Error("error creating file access", logger.Error(err))
		return err
	}
	return nil
}

func (r *FileAccessRepository) ReadByFileID(ctx context.Context, fileID string) (*model.FileAccess, error) {
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var fileAccess model.FileAccess
	err = r.collection.FindOne(ctx, bson.M{"file_id": objectID}).Decode(&fileAccess)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading file access", logger.Error(err))
		return nil, err
	}
	return &fileAccess, nil
}

func (r *FileAccessRepository) ReadByOwnerID(ctx context.Context, ownerID string) (*model.FileAccess, error) {
	objectID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var fileAccess model.FileAccess
	err = r.collection.FindOne(ctx, bson.M{"owner_id": objectID}).Decode(&fileAccess)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading file access", logger.Error(err))
		return nil, err
	}
	return &fileAccess, nil
}

func (r *FileAccessRepository) ReadByRecipientID(ctx context.Context, recipientID string) (*model.FileAccess, error) {
	objectID, err := primitive.ObjectIDFromHex(recipientID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var fileAccess model.FileAccess
	err = r.collection.FindOne(ctx, bson.M{"recipient_id": objectID}).Decode(&fileAccess)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading file access", logger.Error(err))
		return nil, err
	}
	return &fileAccess, nil
}

func (r *FileAccessRepository) Update(ctx context.Context, id string, update bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	update["updated_at"] = time.Now()

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating file access", logger.Error(err))
		return err
	}
	return nil
}

func (r *FileAccessRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.Log.Error("error deleting file access", logger.Error(err))
		return err
	}
	return nil
}

func (r *FileAccessRepository) DeleteByFileID(ctx context.Context, fileID string) error {
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		logger.Log.Error("invalid file ID format", logger.Error(err))
		return err
	}

	_, err = r.collection.DeleteMany(ctx, bson.M{"file_id": objectID})
	if err != nil {
		logger.Log.Error("error deleting file access by file ID", logger.Error(err))
		return err
	}
	return nil
}
