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
	user, err := ts.userRepo.ReadUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", errors.New("invalid credentials")
	}

	// Check if the account is locked
	isLocked, err := ts.isAccountLocked(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to check account status: %w", err)
	}
	if isLocked {
		return "", "", errors.New("your account has been locked due to too many failed login attempts. Please reset your password or contact support")
	}

	if !helper.CheckHash(password, user.Password) {
		// Increment the failed attempt counter
		err = ts.incrementFailedAttempts(ctx, email)
		if err != nil {
			return "", "", fmt.Errorf("failed to increment failed attempts: %w", err)
		}

		return "", "", errors.New("invalid credentials")
	}

	// Reset the failed attempt counter on successful login
	ts.resetFailedAttempts(ctx, email)

	accessToken, refreshToken, err = ts.jwtProvider.GenerateTokens(user.ID.Hex(), &user.Email, nil)
	if err != nil {
		return "", "", errors.New("failed to generate tokens")
	}

	return accessToken, refreshToken, nil
}

func (ts *TokenService) isAccountLocked(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("lockout:%s", email)
	isLocked, err := ts.redisProvider.Get(ctx, key)
	if isLocked != "1" || err != nil {
		return false, nil
	}

	return true, nil
}

// IncrementFailedAttempts increments the failed login attempts count for a given email
func (ts *TokenService) incrementFailedAttempts(ctx context.Context, email string) error {
	key := fmt.Sprintf("failed_login_attempts:%s", email)

	// Get the current number of failed attempts
	attemptsStr, err := ts.redisProvider.Get(ctx, key)
	if attemptsStr == "0" || err != nil {
		// If no attempts exist or an error occurred, initialize it to 1 with a 24-hour expiration
		err = ts.redisProvider.Set(ctx, key, "1", 24*time.Hour)
		if err != nil {
			return err
		}

		return nil
	}

	// Convert the current number of failed attempts to an integer
	attempts := 0
	attempts, err = strconv.Atoi(attemptsStr)
	if err != nil {
		return err
	}

	// Increment the failed attempt counter
	attempts++

	// Update the failed attempt count in Redis with a 24-hour expiration
	err = ts.redisProvider.Set(ctx, key, strconv.Itoa(attempts), 24*time.Hour)
	if err != nil {
		return err
	}

	// Lock the account if there are 5 or more failed attempts
	if attempts >= 5 {
		err = ts.lockAccount(ctx, email)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResetFailedAttempts resets the failed login attempts count for a given email
func (ts *TokenService) resetFailedAttempts(ctx context.Context, email string) error {
	key := fmt.Sprintf("failed_login_attempts:%s", email)
	err := ts.redisProvider.Del(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

// lockAccount locks the account for 24 hours by setting a key in Redis with an expiration time of 24 hours
func (ts *TokenService) lockAccount(ctx context.Context, email string) error {
	lockKey := fmt.Sprintf("lockout:%s", email)

	err := ts.redisProvider.Set(ctx, lockKey, true, 24*time.Hour)
	if err != nil {
		return err
	}

	return nil
}

// RefreshTokens validates the refresh token and generates new tokens
func (ts *TokenService) RefreshTokens(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	claims, err := ts.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	user, err := ts.userRepo.ReadUserByEmail(ctx, *claims.Email)
	if err != nil || user == nil {
		return "", "", errors.New("user no longer exists")
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
	user, err := ts.userRepo.ReadUserByEmail(ctx, *claims.Email)
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

	// Generate tokens
	accessToken, refreshToken, err = ts.jwtProvider.GenerateTokens(userID, nil, nil)
	if err != nil {
		return "", "", errors.New("failed to generate tokens")
	}

	return accessToken, refreshToken, nil
}
