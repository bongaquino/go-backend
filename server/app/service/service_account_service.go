package service

import (
	"context"
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceAccountService struct {
	serviceAccountRepo *repository.ServiceAccountRepository
	userRepo           *repository.UserRepository
	limitRepo          *repository.LimitRepository
}

func NewServiceAccountService(
	serviceAccountRepo *repository.ServiceAccountRepository,
	userRepo *repository.UserRepository,
	limitRepo *repository.LimitRepository,
) *ServiceAccountService {
	return &ServiceAccountService{
		serviceAccountRepo: serviceAccountRepo,
		userRepo:           userRepo,
		limitRepo:          limitRepo,
	}
}

func GenerateClientCredentials() (string, string, error) {
	clientID, err := helper.GenerateClientID()
	if err != nil {
		return "", "", err
	}

	clientSecret, err := helper.GenerateClientSecret()
	if err != nil {
		return "", "", err
	}

	return clientID, clientSecret, nil
}

func (s *ServiceAccountService) CreateServiceAccount(ctx context.Context, userID, clientID, clientSecret string) (string, string, error) {
	// Convert userID string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", "", err
	}

	// Create service account
	serviceAccount := &model.ServiceAccount{
		UserID:       objectID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	err = s.serviceAccountRepo.Create(ctx, serviceAccount)
	if err != nil {
		return "", "", err
	}

	return clientID, clientSecret, nil
}
