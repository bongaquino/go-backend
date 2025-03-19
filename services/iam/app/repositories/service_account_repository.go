package repositories

import (
	"context"
	"time"

	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// ServiceAccountRepository handles database operations for the ServiceAccount model.
type ServiceAccountRepository struct {
	collection *mongoDriver.Collection
}

// NewServiceAccountRepository initializes a new ServiceAccountRepository.
func NewServiceAccountRepository(mongoService *services.MongoService) *ServiceAccountRepository {
	db := mongoService.GetDB()
	return &ServiceAccountRepository{
		collection: db.Collection("service_accounts"),
	}
}

// CreateServiceAccount inserts a new service account into the database.
func (r *ServiceAccountRepository) CreateServiceAccount(ctx context.Context, account *models.ServiceAccount) error {
	// Generate a new ObjectID for the account
	account.ID = primitive.NewObjectID()

	// Set timestamps
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, account)
	if err != nil {
		logger.Log.Error("error creating service account", logger.Error(err))
		return err
	}
	return nil
}

// ReadServiceAccountByClientID retrieves a service account by its client ID.
func (r *ServiceAccountRepository) ReadServiceAccountByClientID(ctx context.Context, clientID string) (*models.ServiceAccount, error) {
	var account models.ServiceAccount
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

// UpdateServiceAccount updates an existing service account.
func (r *ServiceAccountRepository) UpdateServiceAccount(ctx context.Context, clientID string, update bson.M) error {
	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateOne(ctx, bson.M{"client_id": clientID}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating service account", logger.Error(err))
		return err
	}
	return nil
}

// DeleteServiceAccount removes a service account from the database.
func (r *ServiceAccountRepository) DeleteServiceAccount(ctx context.Context, clientID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"client_id": clientID})
	if err != nil {
		logger.Log.Error("error deleting service account", logger.Error(err))
		return err
	}
	return nil
}
