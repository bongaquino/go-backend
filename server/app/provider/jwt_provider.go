package provider

import (
	"context"
	"errors"
	"time"

	"bongaquino/server/config"
	"bongaquino/server/core/logger"

	"github.com/golang-jwt/jwt/v4"
)

// JWTProvider handles JWT-related operations
type JWTProvider struct {
	redisProvider   *RedisProvider
	secretKey       string
	tokenDuration   time.Duration
	refreshDuration time.Duration
}

// NewJWTProvider initializes a new JWTProvider with Redis dependency
func NewJWTProvider(redisProvider *RedisProvider) *JWTProvider {
	jwtConfig := config.LoadJWTConfig()

	if jwtConfig.JWTSecret == "" {
		logger.Log.Fatal("JWT secret key is missing in environment variables")
	}

	return &JWTProvider{
		redisProvider:   redisProvider,
		secretKey:       jwtConfig.JWTSecret,
		tokenDuration:   time.Duration(jwtConfig.JWTTokenExpiration) * time.Second,
		refreshDuration: time.Duration(jwtConfig.JWTRefreshExpiration) * time.Second,
	}
}

// Claims structure for JWT
type Claims struct {
	Sub      string  `json:"sub"`
	Email    *string `json:"email,omitempty"`
	ClientId *string `json:"client_id,omitempty"`
	Scope    string  `json:"scope"`
	jwt.RegisteredClaims
}

// GenerateTokens creates an access and refresh token for a user
func (j *JWTProvider) GenerateTokens(userID string, email, clientID *string) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessClaims := Claims{
		Sub:      userID,
		Email:    email,
		ClientId: clientID,
		Scope:    "access",
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
		Sub:      userID,
		Email:    email,
		ClientId: clientID,
		Scope:    "refresh",
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
	err = j.redisProvider.Set(ctx, "refresh_token:"+userID, refreshToken, j.refreshDuration)
	if err != nil {
		logger.Log.Error("failed to store refresh token in Redis", logger.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken parses and validates a JWT token
func (j *JWTProvider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateRefreshToken checks if the refresh token is valid and exists in Redis
func (j *JWTProvider) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Scope != "refresh" {
		return nil, errors.New("invalid token type, expected refresh token")
	}

	ctx := context.Background()
	storedToken, err := j.redisProvider.Get(ctx, "refresh_token:"+claims.Sub)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	if storedToken != tokenString {
		return nil, errors.New("refresh token mismatch")
	}

	return claims, nil
}

// RevokeRefreshToken removes a refresh token from Redis
func (j *JWTProvider) RevokeRefreshToken(userID string) error {
	ctx := context.Background()
	return j.redisProvider.Del(ctx, "refresh_token:"+userID)
}
