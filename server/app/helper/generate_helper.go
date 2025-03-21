package helper

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// GenerateClientID generates a cryptographically safe client ID.
func GenerateClientID() (string, error) {
	return generateRandomString(32) // 32 bytes (~43 Base64 characters)
}

// GenerateClientSecret generates a cryptographically safe client secret.
func GenerateClientSecret() (string, error) {
	return generateRandomString(64) // 64 bytes (~86 Base64 characters)
}

// generateRandomString generates a secure random string of the given length in bytes.
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}

// GenerateResetCode generates a secure random reset code.
func GenerateResetCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset code: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
