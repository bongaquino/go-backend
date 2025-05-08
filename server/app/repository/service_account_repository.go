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

type ServiceAccountRepository struct {
	collection *mongoDriver.Collection
}

func NewServiceAccountRepository(mongoProvider *provider.MongoProvider) *ServiceAccountRepository {
	db := mongoProvider.GetDB()
	return &ServiceAccountRepository{
		collection: db.Collection("service_accounts"),
	}
}

func (r *ServiceAccountRepository) Create(ctx context.Context, account *model.ServiceAccount) error {
	account.ID = primitive.NewObjectID()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, account)
	if err != nil {
		logger.Log.Error("error creating service account", logger.Error(err))
		return err
	}
	return nil
}

func (r *ServiceAccountRepository) ReadByClientID(ctx context.Context, clientID string) (*model.ServiceAccount, error) {
	var account model.ServiceAccount
	err := r.collection.FindOne(ctx, bson.M{"client_id": clientID}).Decode(&account)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading service account by client ID", logger.Error(err))
		return nil, err
	}
	return &account, nil
}

func (r *ServiceAccountRepository) Update(ctx context.Context, clientID string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateOne(ctx, bson.M{"client_id": clientID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating service account", logger.Error(err))
		return err
	}
	return nil
}

func (r *ServiceAccountRepository) Delete(ctx context.Context, clientID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"client_id": clientID})
	if err != nil {
		logger.Log.Error("error deleting service account", logger.Error(err))
		return err
	}
	return nil
}
