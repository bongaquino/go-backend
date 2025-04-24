package service

import (
	"context"
	"errors"
	"fmt"
	"koneksi/server/app/dto"
	"koneksi/server/app/helper"
	"koneksi/server/app/model"
	"koneksi/server/app/provider"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"
	"time"
)

type UserService struct {
	userRepo      *repository.UserRepository
	profileRepo   *repository.ProfileRepository
	roleRepo      *repository.RoleRepository
	userRoleRepo  *repository.UserRoleRepository
	redisProvider *provider.RedisProvider
}

func NewUserService(userRepo *repository.UserRepository,
	profileRepo *repository.ProfileRepository,
	roleRepo *repository.RoleRepository,
	userRoleRepo *repository.UserRoleRepository,
	redisProvider *provider.RedisProvider,
) *UserService {
	return &UserService{
		userRepo:      userRepo,
		profileRepo:   profileRepo,
		roleRepo:      roleRepo,
		userRoleRepo:  userRoleRepo,
		redisProvider: redisProvider,
	}
}

// UserExists checks if a user with the given email already exists
func (us *UserService) UserExists(ctx context.Context, email string) (bool, error) {
	// Query the repository to check if the user exists
	user, err := us.userRepo.ReadUserByEmail(ctx, email)
	if err != nil {
		logger.Log.Error("error checking if user exists", logger.Error(err))
		return false, errors.New("internal server error")
	}

	// Return true if the user exists, false otherwise
	return user != nil, nil
}

// RegisterUser registers a new user
func (us *UserService) RegisterUser(ctx context.Context, request *dto.RegisterUser) (*model.User, *model.Profile, *model.UserRole, string, error) {
	user := &model.User{
		Email:      request.Email,
		Password:   request.Password,
		IsVerified: true,
	}
	if err := us.userRepo.CreateUser(ctx, user); err != nil {
		logger.Log.Error("error creating user", logger.Error(err))
		return nil, nil, nil, "", errors.New("failed to create user")
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
		return nil, nil, nil, "", errors.New("failed to create profile")
	}

	userRole, err := us.roleRepo.ReadRoleByName(ctx, "user")
	if err != nil {
		logger.Log.Error("error retrieving default role", logger.Error(err))
		return nil, nil, nil, "", errors.New("failed to assign default role")
	}
	if userRole == nil {
		return nil, nil, nil, "", errors.New("default role not found")
	}

	userRoleAssignment := &model.UserRole{
		UserID: user.ID,
		RoleID: userRole.ID,
	}
	if err := us.userRoleRepo.CreateUserRole(ctx, userRoleAssignment); err != nil {
		logger.Log.Error("error assigning default role", logger.Error(err))
		return nil, nil, nil, "", errors.New("failed to assign default role")
	}

	return user, profile, userRoleAssignment, userRole.Name, nil
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
	if err := us.userRepo.UpdateUser(ctx, user.ID.Hex(), update); err != nil {
		logger.Log.Error("error updating user password", logger.Error(err))
		return errors.New("failed to update password")
	}

	return nil
}

func (us *UserService) GeneratePasswordResetCode(ctx context.Context, email string) (string, error) {
	// Check if the user exists
	user, err := us.userRepo.ReadUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", fmt.Errorf("user not found")
	}

	// Construct the Redis key
	key := fmt.Sprintf("password_reset:%s", email)

	// Check if a reset code already exists in Redis
	existingCode, err := us.redisProvider.Get(ctx, key)
	if err == nil && existingCode != "" {
		return "", fmt.Errorf("password reset already pending")
	}

	// Generate a random reset code using the helper
	resetCode, err := helper.GenerateCode(6) // 6 bytes (~12 hex characters)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset code")
	}

	// Store the reset code in Redis with a 15-minute expiration
	err = us.redisProvider.Set(ctx, key, resetCode, 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to store reset code")
	}

	return resetCode, nil
}

func (us *UserService) ResetPassword(ctx context.Context, email, resetCode, newPassword string) error {
	// Retrieve the user by email from the database
	user, err := us.userRepo.ReadUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Construct the Redis key
	key := fmt.Sprintf("password_reset:%s", email)

	// Retrieve the stored reset code from Redis
	storedCode, err := us.redisProvider.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to retrieve reset code")
	}

	// Compare the stored code with the provided code
	if storedCode != resetCode {
		return fmt.Errorf("invalid reset code")
	}

	// Check if new password is not the same as the old one
	if helper.CheckHash(newPassword, user.Password)  {
		return fmt.Errorf("new password must be different from the old one")
	}

	// Delete the reset code from Redis to prevent reuse
	err = us.redisProvider.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete reset code")
	}

	// Hash the new password
	hashedPassword, err := helper.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	// Update the user's password in the database
	err = us.userRepo.UpdateUserByEmail(ctx, email, map[string]any{
		"password": hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update password")
	}

	return nil
}

func (us *UserService) GetUserProfile(ctx context.Context, userID string) (*model.User, *model.Profile, error) {
	user, err := us.userRepo.ReadUser(ctx, userID)
	if err != nil {
		logger.Log.Error("error fetching user by ID", logger.Error(err))
		return nil, nil, errors.New("failed to retrieve user")
	}
	if user == nil {
		return nil, nil, errors.New("user not found")
	}

	profile, err := us.profileRepo.ReadProfileByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error("error fetching profile by user ID", logger.Error(err))
		return nil, nil, errors.New("failed to retrieve profile")
	}
	if profile == nil {
		return nil, nil, errors.New("profile not found")
	}

	return user, profile, nil
}

func (us *UserService) GetUserProfileByEmail(ctx context.Context, email string) (*model.User, *model.Profile, error) {
	user, err := us.userRepo.ReadUserByEmail(ctx, email)
	if err != nil {
		logger.Log.Error("error fetching user by email", logger.Error(err))
		return nil, nil, errors.New("failed to retrieve user")
	}
	if user == nil {
		return nil, nil, errors.New("user not found")
	}

	profile, err := us.profileRepo.ReadProfileByUserID(ctx, user.ID.Hex())
	if err != nil {
		logger.Log.Error("error fetching profile by user ID", logger.Error(err))
		return nil, nil, errors.New("failed to retrieve profile")
	}
	if profile == nil {
		return nil, nil, errors.New("profile not found")
	}

	return user, profile, nil
}

func (us *UserService) ValidatePassword(ctx context.Context, userID, password string) (bool, error) {
	// Retrieve the user from the database
	user, err := us.userRepo.ReadUser(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	// Compare the provided password with the stored hash
	isValid := helper.CheckHash(password, user.Password)
	return isValid, nil
}
