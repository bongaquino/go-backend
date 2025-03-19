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
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTService initializes a new JWTService
func NewJWTService() *JWTService {
	jwtConfig := config.LoadJwtConfig()

	if jwtConfig.JwtSecret == "" {
		logger.Log.Fatal("JWT secret key is missing in environment variables")
	}

	return &JWTService{
		secretKey:     jwtConfig.JwtSecret,
		tokenDuration: time.Duration(jwtConfig.JwtExpiration) * time.Minute,
	}
}

// Claims structure for JWT
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func (j *JWTService) GenerateToken(userID string, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
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
