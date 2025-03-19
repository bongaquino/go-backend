package services

import (
	"errors"
	"time"

	"koneksi/services/iam/config"
	"koneksi/services/iam/core/logger"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService handles JWT-related operations
type JWTService struct {
	secretKey       string
	tokenDuration   time.Duration
	refreshDuration time.Duration
}

// NewJWTService initializes a new JWTService
func NewJWTService() *JWTService {
	jwtConfig := config.LoadJwtConfig()

	if jwtConfig.JwtSecret == "" {
		logger.Log.Fatal("JWT secret key is missing in environment variables")
	}

	return &JWTService{
		secretKey:       jwtConfig.JwtSecret,
		tokenDuration:   time.Duration(jwtConfig.JwtTokenExpiration) * time.Minute,   // Access token expiration
		refreshDuration: time.Duration(jwtConfig.JwtRefreshExpiration) * time.Minute, // Refresh token expiration
	}
}

// Claims structure for JWT
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // New claim to differentiate token types
	jwt.RegisteredClaims
}

// GenerateTokens creates an access and refresh token for a user
func (j *JWTService) GenerateTokens(userID string, email string) (accessToken string, refreshToken string, err error) {
	// Generate access token (short-lived)
	accessClaims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "access", // Identifies this as an access token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (long-lived)
	refreshClaims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "refresh", // Identifies this as a refresh token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken parses and validates a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

// ValidateRefreshToken ensures the token is a refresh token
func (j *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Ensure this is a refresh token
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type, expected refresh token")
	}

	return claims, nil
}
