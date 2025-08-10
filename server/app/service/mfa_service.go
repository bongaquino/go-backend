package service

import (
	"context"
	"fmt"
	"time"

	"koneksi/server/app/helper"
	"koneksi/server/app/provider"
	"koneksi/server/app/repository"

	"go.mongodb.org/mongo-driver/bson"
)

// MFAService handles MFA-related operations
type MFAService struct {
	userRepo      *repository.UserRepository
	settingRepo   *repository.SettingRepository
	redisProvider *provider.RedisProvider
}

// NewMFAService initializes a new MFAService
func NewMFAService(userRepo *repository.UserRepository, settingRepo *repository.SettingRepository, redisProvider *provider.RedisProvider) *MFAService {
	return &MFAService{
		userRepo:      userRepo,
		settingRepo:   settingRepo,
		redisProvider: redisProvider,
	}
}

// GenerateOTP generates an OTP secret and QR code for the user
func (ms *MFAService) GenerateOTP(ctx context.Context, userID string) (string, string, error) {
	// Generate OTP secret
	otpSecret, err := helper.GenerateOTPSecret(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate OTP secret: %w", err)
	}

	// Encrypt the OTP secret
	encryptedOtpSecret, err := helper.Encrypt(otpSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt OTP secret: %w", err)
	}

	// Generate QR code
	qrCode, err := helper.GenerateQRCode(userID, otpSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Save the encrypted OTP secret to the user's record in the database
	err = ms.userRepo.Update(ctx, userID, bson.M{"otp_secret": encryptedOtpSecret})
	if err != nil {
		return "", "", fmt.Errorf("failed to save OTP secret: %w", err)
	}

	return otpSecret, qrCode, nil
}

func (ms *MFAService) VerifyOTP(ctx context.Context, userID, otp string) (bool, error) {
	// Retrieve the user from the database
	user, err := ms.userRepo.Read(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	// Check if OTP secret is already generated
	if len(user.OtpSecret) == 0 {
		return false, fmt.Errorf("OTP secret not set")
	}

	// Decrypt the stored OTP secret
	decryptedOtpSecret, err := helper.Decrypt(user.OtpSecret)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt OTP secret: %w", err)
	}

	// Verify the OTP using the decrypted secret
	isValid := helper.VerifyOTP(decryptedOtpSecret, otp)
	return isValid, nil
}

func (ms *MFAService) EnableMFA(ctx context.Context, userID string) error {
	update := bson.M{
		"is_mfa_enabled": true,
	}
	err := ms.settingRepo.UpdateByUserID(ctx, userID, update)
	if err != nil {
		return fmt.Errorf("failed to enable MFA: %w", err)
	}

	return nil
}

func (ms *MFAService) DisableMFA(ctx context.Context, userID string) error {
	updateUser := bson.M{
		"otp_secret": "",
	}
	err := ms.userRepo.Update(ctx, userID, updateUser)
	if err != nil {
		return fmt.Errorf("failed to clear OTP secret: %w", err)
	}

	updateSetting := bson.M{
		"is_mfa_enabled": false,
	}
	err = ms.settingRepo.UpdateByUserID(ctx, userID, updateSetting)
	if err != nil {
		return fmt.Errorf("failed to disable MFA: %w", err)
	}

	return nil
}

// Generate login code for the user
func (ms *MFAService) GenerateLoginCode(ctx context.Context, userID string) (string, error) {
	// Generate a login code
	loginCode, err := helper.GenerateCode(6)
	if err != nil {
		return "", fmt.Errorf("failed to generate login code")
	}

	// Construct the Redis key
	key := fmt.Sprintf("login_code:%s", loginCode)

	// Check if a login code already exists in Redis
	existingCode, err := ms.redisProvider.Get(ctx, key)
	if err == nil && existingCode != "" {
		return "", fmt.Errorf("login already pending")
	}

	// Store the login code in Redis with a 3-minute expiration
	err = ms.redisProvider.Set(ctx, key, userID, 3*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to store login code")
	}

	// Return the login code
	return loginCode, nil
}

func (ms *MFAService) VerifyLoginCode(ctx context.Context, loginCode, otp string) (string, error) {
	// Construct the Redis key
	key := fmt.Sprintf("login_code:%s", loginCode)

	// Retrieve the user ID from Redis using the login code
	userID, err := ms.redisProvider.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("invalid login code")
	}

	// Verify the OTP
	isValid, err := ms.VerifyOTP(ctx, userID, otp)
	if err != nil {
		return "", fmt.Errorf("failed to verify OTP")
	}
	if !isValid {
		return "", fmt.Errorf("invalid OTP")
	}

	// Delete the login code from Redis
	err = ms.redisProvider.Del(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to delete login code")
	}

	return userID, nil
}
