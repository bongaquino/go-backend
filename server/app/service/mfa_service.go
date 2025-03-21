package service

import (
	"context"
	"fmt"

	"koneksi/server/app/helper"
	"koneksi/server/app/repository"
)

// MFAService handles MFA-related operations
type MFAService struct {
	userRepo *repository.UserRepository
}

// NewMFAService initializes a new MFAService
func NewMFAService(userRepo *repository.UserRepository) *MFAService {
	return &MFAService{
		userRepo: userRepo,
	}
}

// GenerateOTP generates an OTP secret and QR code for the user
func (ms *MFAService) GenerateOTP(ctx context.Context, userID string) (string, string, error) {
	// Generate OTP secret
	otpSecret, err := helper.GenerateOTPSecret(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate OTP secret: %w", err)
	}

	// Generate QR code
	qrCode, err := helper.GenerateQRCode(userID, otpSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Save hashed OTP secret to the user's record in the database
	err = ms.userRepo.UpdateOTPSecret(ctx, userID, otpSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to save OTP secret: %w", err)
	}

	return otpSecret, qrCode, nil
}
