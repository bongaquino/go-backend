package service

import (
	"context"
	"errors"
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"
	"time"
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

// RegisterUser registers a new user
func (us *UserService) RegisterUser(ctx context.Context, request *dto.RegisterUser) (*model.User, *model.Profile, *model.UserRole, error) {
	existingUser, err := us.userRepo.ReadUserByEmail(ctx, request.Email)
	if err != nil {
		logger.Log.Error("error checking existing user", logger.Error(err))
		return nil, nil, nil, errors.New("internal server error")
	}
	if existingUser != nil {
		return nil, nil, nil, errors.New("email already exists")
	}

	user := &model.User{
		Email:      request.Email,
		Password:   request.Password,
		IsVerified: true,
	}
	if err := us.userRepo.CreateUser(ctx, user); err != nil {
		logger.Log.Error("error creating user", logger.Error(err))
		return nil, nil, nil, errors.New("failed to create user")
	}

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

// ChangePassword changes the user's password
func (us *UserService) ChangePassword(ctx context.Context, userID string, request *dto.ChangePasswordDTO) error {
	// Fetch the user from the repository
	user, err := us.userRepo.ReadUser(ctx, userID)
	if err != nil {
		logger.Log.Error("error fetching user by ID", logger.Error(err))
		return errors.New("failed to retrieve user")
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify if the old password is the same as the new password
	if request.OldPassword == request.NewPassword {
		return errors.New("new password must be different from the old password")
	}

	// Verify the old password
	if !helper.CheckHash(request.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := helper.Hash(request.NewPassword)
	if err != nil {
		logger.Log.Error("error hashing new password", logger.Error(err))
		return errors.New("failed to hash new password")
	}

	// Update the user's password in the repository
	update := map[string]any{
		"password":  hashedPassword,
		"updatedAt": time.Now(),
	}
	if err := us.userRepo.UpdateUser(ctx, user.Email, update); err != nil {
		logger.Log.Error("error updating user password", logger.Error(err))
		return errors.New("failed to update password")
	}

	return nil
}
