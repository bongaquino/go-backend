package repositories

import (
	"context"
	"time"

	"koneksi/services/iam/app/helpers"
	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongoDriver.Collection
}

// NewUserRepository initializes a new UserRepository
func NewUserRepository(mongoService *services.MongoService) *UserRepository {
	db := mongoService.GetDB()
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Hash the password before storing it
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		logger.Log.Error("error hashing password", logger.Error(err))
		return err
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Log.Error("error creating user", logger.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) ReadUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading user by username", logger.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, username string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating user", logger.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, username string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		logger.Log.Error("error deleting user", logger.Error(err))
		return err
	}
	return nil
}
