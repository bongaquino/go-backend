package service

import (
	"context"
	"errors"
	"koneksi/server/app/helper"
	"koneksi/server/app/provider"
	"koneksi/server/app/repository"
)

type TokenService struct {
	userRepo    *repository.UserRepository
	jwtProvider *provider.JWTProvider
	mfaService  *MFAService
}

func NewTokenService(userRepo *repository.UserRepository, jwtProvider *provider.JWTProvider, mfaService *MFAService) *TokenService {
	return &TokenService{
		userRepo:    userRepo,
		jwtProvider: jwtProvider,
		mfaService:  mfaService,
	}
}

// AuthenticateUser validates user credentials and generates tokens
func (ts *TokenService) AuthenticateUser(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error) {
	user, err := ts.userRepo.ReadUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", errors.New("invalid credentials")
	}

	if !helper.CheckHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
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
