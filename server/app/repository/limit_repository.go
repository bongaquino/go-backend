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

type LimitRepository struct {
	collection *mongoDriver.Collection
}

func NewLimitRepository(mongoProvider *provider.MongoProvider) *LimitRepository {
	db := mongoProvider.GetDB()
	return &LimitRepository{
		collection: db.Collection("policies"),
	}
}

func (r *LimitRepository) Create(ctx context.Context, limit *model.Limit) error {
	limit.ID = primitive.NewObjectID()
	limit.CreatedAt = time.Now()
	limit.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, limit)
	if err != nil {
		logger.Log.Error("error creating limit", logger.Error(err))
		return err
	}
	return nil
}

func (r *LimitRepository) Read(ctx context.Context, id string) (*model.Limit, error) {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var limit model.Limit
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&limit)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading limit", logger.Error(err))
		return nil, err
	}
	return &limit, nil
}

func (r *LimitRepository) ReadByUserID(ctx context.Context, userID string) (*model.Limit, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var limit model.Limit
	err = r.collection.FindOne(ctx, bson.M{"user_id": objectID}).Decode(&limit)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading limit by userID", logger.Error(err))
		return nil, err
	}
	return &limit, nil
}

func (r *LimitRepository) ReadByOrganizationID(ctx context.Context, orgID string) (*model.Limit, error) {
	// Convert orgID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var limit model.Limit
	err = r.collection.FindOne(ctx, bson.M{"organization_id": objectID}).Decode(&limit)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading limit by orgID", logger.Error(err))
		return nil, err
	}
	return &limit, nil
}

func (r *LimitRepository) Update(ctx context.Context, name string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating policy", logger.Error(err))
		return err
	}
	return nil
}

func (r *LimitRepository) UpdateByUserID(ctx context.Context, userID string, update bson.M) error {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	// Set the updated time
	update["updated_at"] = time.Now()

	_, err = r.collection.UpdateOne(ctx, bson.M{"user_id": objectID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating limit by userID", logger.Error(err))
		return err
	}
	return nil
}

func (r *LimitRepository) UpdateByOrganizationID(ctx context.Context, orgID string, update bson.M) error {
	// Convert orgID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	// Set the updated time
	update["updated_at"] = time.Now()

	_, err = r.collection.UpdateOne(ctx, bson.M{"organization_id": objectID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating limit by orgID", logger.Error(err))
		return err
	}
	return nil
}
