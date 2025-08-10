package repository

import (
	"context"
	"time"

	"bongaquino/server/app/model"
	"bongaquino/server/app/provider"
	"bongaquino/server/core/logger"

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

func (r *DirectoryRepository) ListByDirectoryIDUserID(ctx context.Context, directoryID string, userID string) ([]*model.Directory, error) {
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

	filter := bson.M{
		"directory_id": directoryObjectID,
		"user_id":      userObjectID,
		"is_deleted":   false,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		logger.Log.Error("error listing directories by directory ID and user ID", logger.Error(err))
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

func (r *DirectoryRepository) ReadByIDUserID(ctx context.Context, id string, userID string) (*model.Directory, error) {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid user ID format", logger.Error(err))
		return nil, err
	}

	// Find the directory by ID and userID
	var directory model.Directory
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "user_id": userObjectID, "is_deleted": false}).Decode(&directory)

	// Check if the directory was found
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading directory", logger.Error(err))
		return nil, err
	}

	// Return the directory
	return &directory, nil
}

func (r *DirectoryRepository) ReadByUserIDName(ctx context.Context, userID string, name string) (*model.Directory, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid user ID format", logger.Error(err))
		return nil, err
	}

	var directory model.Directory
	err = r.collection.FindOne(ctx, bson.M{"user_id": objectID, "name": name, "is_deleted": false}).Decode(&directory)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading directory by userID and name", logger.Error(err))
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

func (r *DirectoryRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid user ID format", logger.Error(err))
		return 0, err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"user_id": objectID, "is_deleted": false})
	if err != nil {
		logger.Log.Error("error counting directories by userID", logger.Error(err))
		return 0, err
	}
	return count, nil
}

// FindAllDescendantIDs returns all descendant directory IDs (as strings) for a given directory, including nested children.
func (r *DirectoryRepository) FindAllDescendantIDs(ctx context.Context, directoryID, userID string) ([]string, error) {
    var result []string
    queue := []string{directoryID}

    for len(queue) > 0 {
        currentID := queue[0]
        queue = queue[1:]

        // Find direct children of currentID
        children, err := r.ListByDirectoryIDUserID(ctx, currentID, userID)
        if err != nil {
            return nil, err
        }
        for _, child := range children {
            childID := child.ID.Hex()
            result = append(result, childID)
            queue = append(queue, childID)
        }
    }

    return result, nil
}