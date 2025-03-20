package repositories

import (
	"context"
	"time"

	"koneksi/server/app/helpers"
	"koneksi/server/app/models"
	"koneksi/server/app/services"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// Generate a new ObjectID for the user
	user.ID = primitive.NewObjectID()

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Hash the password before storing it
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		logger.Log.Error("error hashing password", logger.Error(err))
		return err
	}
	user.Password = hashedPassword

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Log.Error("error creating user", logger.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) ReadUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading user by email", logger.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ReadUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		logger.Log.Error("error reading user by ID", logger.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, email string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": update})
	if err != nil {
		logger.Log.Error("error updating user", logger.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, email string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		logger.Log.Error("error deleting user", logger.Error(err))
		return err
	}
	return nil
}
