package service

import (
	"context"
	"errors"
	"koneksi/server/app/helper"
	"koneksi/server/app/provider"
	"koneksi/server/app/repository"
)

type TokenService struct {
	userRepo   *repository.UserRepository
	jwtService *provider.JwtProvider
}

func NewTokenService(userRepo *repository.UserRepository, jwtService *provider.JwtProvider) *TokenService {
	return &TokenService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// AuthenticateUser validates user credentials and generates tokens
func (ts *TokenService) AuthenticateUser(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	user, err := ts.userRepo.ReadUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", errors.New("invalid credentials")
	}

	if !helper.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err = ts.jwtService.GenerateTokens(user.ID.Hex(), &user.Email, nil)
	if err != nil {
		return "", "", errors.New("could not generate tokens")
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens validates the refresh token and generates new tokens
func (ts *TokenService) RefreshTokens(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	claims, err := ts.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	user, err := ts.userRepo.ReadUserByEmail(ctx, *claims.Email)
	if err != nil || user == nil {
		return "", "", errors.New("user no longer exists")
	}

	accessToken, newRefreshToken, err = ts.jwtService.GenerateTokens(user.ID.Hex(), &user.Email, nil)
	if err != nil {
		return "", "", errors.New("could not generate new tokens")
	}

	return accessToken, newRefreshToken, nil
}

// RevokeToken revokes the refresh token
func (ts *TokenService) RevokeToken(ctx context.Context, refreshToken string) error {
	// Validate the refresh token
	claims, err := ts.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return errors.New("invalid or expired refresh token")
	}

	// Check if the user exists
	user, err := ts.userRepo.ReadUserByEmail(ctx, *claims.Email)
	if err != nil || user == nil {
		return errors.New("user no longer exists")
	}

	// Revoke the refresh token (e.g., remove it from Redis or mark it as invalid)
	err = ts.jwtService.RevokeRefreshToken(user.ID.Hex())
	if err != nil {
		return errors.New("could not revoke refresh token")
	}

	return nil
}
