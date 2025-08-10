package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"koneksi/server/app/helper"
	"koneksi/server/app/provider"
	"koneksi/server/app/repository"
)

type TokenService struct {
	userRepo      *repository.UserRepository
	jwtProvider   *provider.JWTProvider
	mfaService    *MFAService
	redisProvider *provider.RedisProvider
}

func NewTokenService(userRepo *repository.UserRepository, jwtProvider *provider.JWTProvider, mfaService *MFAService, redisProvider *provider.RedisProvider) *TokenService {
	return &TokenService{
		userRepo:      userRepo,
		jwtProvider:   jwtProvider,
		mfaService:    mfaService,
		redisProvider: redisProvider,
	}
}

// AuthenticateUser validates user credentials and generates tokens
func (ts *TokenService) AuthenticateUser(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error) {
	failedLoginAttemptsKey := fmt.Sprintf("failed_login_attempts:%s", email)

	user, err := ts.userRepo.ReadByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", errors.New("invalid credentials")
	}

	// Check if user account is locked due to too many failed login attempts
	if user.IsLocked {
		return "", "", errors.New("account locked due to multiple failed login attempts")
	}

	if !helper.CheckHash(password, user.Password) {
		// Get the current number of failed attempts
		attemptsStr, err := ts.redisProvider.Get(ctx, failedLoginAttemptsKey)
		if attemptsStr == "0" || err != nil {
			// If no attempts exist or an error occurred, initialize it to 1 with a 24-hour expiration
			err = ts.redisProvider.Set(ctx, failedLoginAttemptsKey, "1", 24*time.Hour)
			if err != nil {
				return "", "", fmt.Errorf("failed to increment failed attempts: %w", err)
			}
		} else {
			// Convert the string to an integer
			attempts, err := strconv.Atoi(attemptsStr)
			if err != nil {
				return "", "", fmt.Errorf("failed to convert failed attempts to integer: %w", err)
			}

			if attempts >= 5 {
				update := map[string]any{
					"is_locked": true,
				}

				if err := ts.userRepo.UpdateByEmail(ctx, email, update); err != nil {
					return "", "", fmt.Errorf("failed to update user lock status: %w", err)
				}
			} else {
				// Increment the failed attempt count
				attempts++

				// Update the failed attempt count in Redis with a 24-hour expiration
				err = ts.redisProvider.Set(ctx, failedLoginAttemptsKey, strconv.Itoa(attempts), 24*time.Hour)
				if err != nil {
					return "", "", fmt.Errorf("failed to update failed attempts: %w", err)
				}
			}
		}

		return "", "", errors.New("invalid credentials")
	}

	// Reset the failed attempt counter on successful login
	err = ts.redisProvider.Del(ctx, failedLoginAttemptsKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to reset failed attempts: %w", err)
	}

	accessToken, refreshToken, err = ts.jwtProvider.GenerateTokens(user.ID.Hex(), &user.Email, nil)
	if err != nil {
		return "", "", errors.New("failed to generate tokens")
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens validates the refresh token and generates new tokens
func (ts *TokenService) RefreshTokens(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	claims, err := ts.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	user, err := ts.userRepo.ReadByEmail(ctx, *claims.Email)
	if err != nil || user == nil {
		return "", "", errors.New("user no longer exists")
	}

	// Check if user account is locked due to too many failed login attempts
	if user.IsLocked {
		return "", "", errors.New("account locked due to multiple failed login attempts")
	}

	accessToken, newRefreshToken, err = ts.jwtProvider.GenerateTokens(user.ID.Hex(), &user.Email, nil)
	if err != nil {
		return "", "", errors.New("failed to generate tokens")
	}

	return accessToken, newRefreshToken, nil
}

// RevokeToken revokes the refresh token
func (ts *TokenService) RevokeToken(ctx context.Context, refreshToken string) error {
	// Validate the refresh token
	claims, err := ts.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return errors.New("invalid or expired refresh token")
	}

	// Check if the user exists
	user, err := ts.userRepo.ReadByEmail(ctx, *claims.Email)
	if err != nil || user == nil {
		return errors.New("user no longer exists")
	}

	// Revoke the refresh token (e.g., remove it from Redis or mark it as invalid)
	err = ts.jwtProvider.RevokeRefreshToken(user.ID.Hex())
	if err != nil {
		return errors.New("failed to revoke token")
	}

	return nil
}

// AuthenticateLoginCode validates the login code and generates tokens
func (ts *TokenService) AuthenticateLoginCode(ctx context.Context, loginCode, otp string) (accessToken string, refreshToken string, err error) {
	userID, err := ts.mfaService.VerifyLoginCode(ctx, loginCode, otp)
	if err != nil {
		return "", "", errors.New("invalid login code or OTP")
	}

	// Check if user exists
	user, err := ts.userRepo.Read(ctx, userID)
	if err != nil || user == nil {
		return "", "", errors.New("user no longer exists")
	}

	// Check if user account is locked due to too many failed login attempts
	if user.IsLocked {
		return "", "", errors.New("account locked due to multiple failed login attempts")
	}

	// Generate tokens
	accessToken, refreshToken, err = ts.jwtProvider.GenerateTokens(userID, nil, nil)
	if err != nil {
		return "", "", errors.New("failed to generate tokens")
	}

	return accessToken, refreshToken, nil
}
