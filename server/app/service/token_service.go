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
	// Check if user exists
	user, err := ts.userRepo.ReadUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", errors.New("invalid credentials")
	}

	// Verify password
	if !helper.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, refreshToken, err = ts.jwtService.GenerateTokens(user.ID.Hex(), &user.Email, nil)
	if err != nil {
		return "", "", errors.New("could not generate tokens")
	}

	return accessToken, refreshToken, nil
}
