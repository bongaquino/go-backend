package service

import (
	"context"
	"errors"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"
)

type UserService struct {
	userRepo     *repository.UserRepository
	profileRepo  *repository.ProfileRepository
	roleRepo     *repository.RoleRepository
	userRoleRepo *repository.UserRoleRepository
}

func NewUserService(userRepo *repository.UserRepository, profileRepo *repository.ProfileRepository, roleRepo *repository.RoleRepository, userRoleRepo *repository.UserRoleRepository) *UserService {
	return &UserService{
		userRepo:     userRepo,
		profileRepo:  profileRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (us *UserService) RegisterUser(ctx context.Context, request *struct {
	FirstName       string
	MiddleName      *string
	LastName        string
	Suffix          *string
	Email           string
	Password        string
	ConfirmPassword string
}) (*model.User, *model.Profile, *model.UserRole, error) {
	// Check if the email already exists
	existingUser, err := us.userRepo.ReadUserByEmail(ctx, request.Email)
	if err != nil {
		logger.Log.Error("error checking existing user", logger.Error(err))
		return nil, nil, nil, errors.New("internal server error")
	}
	if existingUser != nil {
		return nil, nil, nil, errors.New("email already exists")
	}

	// Create the user
	user := &model.User{
		Email:      request.Email,
		Password:   request.Password, // Ensure password is hashed before saving
		IsVerified: true,
	}
	if err := us.userRepo.CreateUser(ctx, user); err != nil {
		logger.Log.Error("error creating user", logger.Error(err))
		return nil, nil, nil, errors.New("failed to create user")
	}

	// Create the profile
	profile := &model.Profile{
		UserID:     user.ID,
		FirstName:  request.FirstName,
		MiddleName: request.MiddleName,
		LastName:   request.LastName,
		Suffix:     request.Suffix,
	}
	if err := us.profileRepo.CreateProfile(ctx, profile); err != nil {
		logger.Log.Error("error creating profile", logger.Error(err))
		return nil, nil, nil, errors.New("failed to create profile")
	}

	// Assign the "user" role to the newly registered user
	userRole, err := us.roleRepo.ReadRoleByName(ctx, "user")
	if err != nil {
		logger.Log.Error("error retrieving default role", logger.Error(err))
		return nil, nil, nil, errors.New("failed to assign default role")
	}
	if userRole == nil {
		return nil, nil, nil, errors.New("default role not found")
	}

	userRoleAssignment := &model.UserRole{
		UserID: user.ID,
		RoleID: userRole.ID,
	}
	if err := us.userRoleRepo.CreateUserRole(ctx, userRoleAssignment); err != nil {
		logger.Log.Error("error assigning default role", logger.Error(err))
		return nil, nil, nil, errors.New("failed to assign default role")
	}

	return user, profile, userRoleAssignment, nil
}
