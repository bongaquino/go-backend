package helper

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"koneksi/server/config"
	"strings"
)

// GenerateClientID creates a secure client ID using random bytes and HMAC with appKey
func GenerateClientID() (string, error) {
	randomBytes, err := generateRandomBytes(32)
	if err != nil {
		return "", err
	}

	appConfig := config.LoadAppConfig()
	h := hmac.New(sha512.New, []byte(appConfig.AppKey))
	h.Write(randomBytes)
	hashed := h.Sum(nil)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed), nil
}

// GenerateClientSecret creates a secure client secret using random bytes and HMAC with appKey
func GenerateClientSecret() (string, error) {
	randomBytes, err := generateRandomBytes(64)
	if err != nil {
		return "", err
	}

	appConfig := config.LoadAppConfig()
	h := hmac.New(sha512.New, []byte(appConfig.AppKey))
	h.Write(randomBytes)
	hashed := h.Sum(nil)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed), nil
}

// generateRandomBytes securely generates a random byte slice of specified length
func generateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// GenerateCode generates a secure random reset code
func GenerateCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset code: %w", err)
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

// GenerateNumericCode generates a secure random numeric code
func GenerateNumericCode(length int) (string, error) {
	digits := make([]byte, length)
	_, err := rand.Read(digits)
	if err != nil {
		return "", fmt.Errorf("failed to generate random digit: %w", err)
	}
	for i := range digits {
		digits[i] = byte(digits[i]%10 + '0') // Convert byte to ASCII digit
	}
	return string(digits), nil
}
