package services

import (
	"context"
	"errors"
	"time"

	"koneksi/services/iam/config"
	"koneksi/services/iam/core/logger"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService handles JWT-related operations
type JWTService struct {
	redisService    *RedisService
	secretKey       string
	tokenDuration   time.Duration
	refreshDuration time.Duration
}

// NewJWTService initializes a new JWTService with Redis dependency
func NewJWTService(redisService *RedisService) *JWTService {
	jwtConfig := config.LoadJwtConfig()

	if jwtConfig.JwtSecret == "" {
		logger.Log.Fatal("JWT secret key is missing in environment variables")
	}

	return &JWTService{
		redisService:    redisService,
		secretKey:       jwtConfig.JwtSecret,
		tokenDuration:   time.Duration(jwtConfig.JwtTokenExpiration) * time.Second,
		refreshDuration: time.Duration(jwtConfig.JwtRefreshExpiration) * time.Second,
	}
}

// Claims structure for JWT
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateTokens creates an access and refresh token for a user
func (j *JWTService) GenerateTokens(userID string, email string) (accessToken string, refreshToken string, err error) {
	// Generate access token
	accessClaims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshClaims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	// Store refresh token in Redis
	ctx := context.Background()
	err = j.redisService.Set(ctx, "refresh_token:"+userID, refreshToken, j.refreshDuration)
	if err != nil {
		logger.Log.Error("Failed to store refresh token in Redis", logger.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken parses and validates a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	// Parse token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken checks if the refresh token is valid and exists in Redis
func (j *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	// Validate JWT signature & structure
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Ensure token type is "refresh"
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type, expected refresh token")
	}

	// Check if the refresh token exists in Redis
	ctx := context.Background()
	storedToken, err := j.redisService.Get(ctx, "refresh_token:"+claims.UserID)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// Ensure the token matches the stored token
	if storedToken != tokenString {
		return nil, errors.New("refresh token mismatch")
	}

	return claims, nil
}

// RevokeRefreshToken removes a refresh token from Redis (logout)
func (j *JWTService) RevokeRefreshToken(userID string) error {
	ctx := context.Background()
	return j.redisService.Del(ctx, "refresh_token:"+userID)
}
