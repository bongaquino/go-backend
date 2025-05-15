package repository

import (
	"context"
	"time"

	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/provider"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type ServiceAccountRepository struct {
	collection *mongoDriver.Collection
}

func NewServiceAccountRepository(mongoProvider *provider.MongoProvider) *ServiceAccountRepository {
	db := mongoProvider.GetDB()
	return &ServiceAccountRepository{
		collection: db.Collection("service_accounts"),
	}
}

func (r *ServiceAccountRepository) ListByUserID(ctx context.Context, userID string) ([]*model.ServiceAccount, error) {
	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var accounts []*model.ServiceAccount
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		logger.Log.Error("error reading service accounts by user ID", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var account model.ServiceAccount
		if err := cursor.Decode(&account); err != nil {
			logger.Log.Error("error decoding service account", logger.Error(err))
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	if err := cursor.Err(); err != nil {
		logger.Log.Error("cursor error", logger.Error(err))
		return nil, err
	}

	return accounts, nil
}

func (r *ServiceAccountRepository) Create(ctx context.Context, account *model.ServiceAccount) error {
	account.ID = primitive.NewObjectID()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	hashedSecret, err := helper.Hash(account.ClientSecret)
	if err != nil {
		logger.Log.Error("error hashing client secret", logger.Error(err))
		return err
	}
	account.ClientSecret = hashedSecret

	_, err = r.collection.InsertOne(ctx, account)
	if err != nil {
		logger.Log.Error("error creating service account", logger.Error(err))
		return err
	}
	return nil
}

func (r *ServiceAccountRepository) ReadByClientID(ctx context.Context, clientID string) (*model.ServiceAccount, error) {
	// Convert clientID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return nil, err
	}

	var account model.ServiceAccount
	err = r.collection.FindOne(ctx, bson.M{"client_id": objectID}).Decode(&account)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading service account by client ID", logger.Error(err))
		return nil, err
	}
	return &account, nil
}

func (r *ServiceAccountRepository) UpdateByClientID(ctx context.Context, clientID string, update bson.M) error {
	// Convert clientID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}

	// Set the updated time
	update["updated_at"] = time.Now()

	_, err = r.collection.UpdateOne(ctx, bson.M{"client_id": objectID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating service account", logger.Error(err))
		return err
	}
	return nil
}

func (r *ServiceAccountRepository) DeleteByClientID(ctx context.Context, clientID string) error {
	// Convert clientID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		logger.Log.Error("invalid ID format", logger.Error(err))
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"client_id": objectID})
	if err != nil {
		logger.Log.Error("error deleting service account", logger.Error(err))
		return err
	}
	return nil
}
