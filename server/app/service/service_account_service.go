package service

import (
	"context"
	"koneksi/server/app/dto"
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

func (s *ServiceAccountService) CreateServiceAccount(ctx context.Context, request *dto.GenerateServiceAccountDTO) (*model.ServiceAccount, error) {
	// Convert userID string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(*request.UserID)
	if err != nil {
		return nil, err
	}

	// Create service account
	serviceAccount := &model.ServiceAccount{
		UserID:       objectID,
		Name:         request.Name,
		ClientID:     *request.ClientID,
		ClientSecret: *request.ClientSecret,
	}

	err = s.serviceAccountRepo.Create(ctx, serviceAccount)
	if err != nil {
		return nil, err
	}

	return serviceAccount, nil
}

func (s *ServiceAccountService) ListServiceAccounts(ctx context.Context, userID string) ([]*model.ServiceAccount, error) {
	// List service accounts by user ID
	serviceAccounts, err := s.serviceAccountRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return serviceAccounts, nil
}
