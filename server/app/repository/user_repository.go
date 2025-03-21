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

type UserRepository struct {
	collection *mongoDriver.Collection
}

func NewUserRepository(mongoProvider *provider.MongoProvider) *UserRepository {
	db := mongoProvider.GetDB()
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashedPassword, err := helper.Hash(user.Password)
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

func (r *UserRepository) ReadUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
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

func (r *UserRepository) ReadUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
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

func (r *UserRepository) UpdateOTPSecret(ctx context.Context, userID, otpSecret string) error {
	hashedSecret, err := helper.Hash(otpSecret)
	if err != nil {
		logger.Log.Error("error hashing OTP secret", logger.Error(err))
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"otp_secret": hashedSecret}})
	if err != nil {
		logger.Log.Error("error saving OTP secret", logger.Error(err))
		return err
	}
	return nil
}
