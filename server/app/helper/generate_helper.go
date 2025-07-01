package helper

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"koneksi/server/config"
	"strings"
)

// GenerateClientID creates a secure client ID in the format: id_<clean_base64>
func GenerateClientID() (string, error) {
	randomBytes, err := generateRandomBytes(32)
	if err != nil {
		return "", err
	}

	appConfig := config.LoadAppConfig()
	h := hmac.New(sha256.New, []byte(appConfig.AppKey))
	h.Write(randomBytes)
	hashed := h.Sum(nil)

	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed)
	cleaned := removeDashesAndUnderscores(encoded)
	return fmt.Sprintf("id_%s", cleaned), nil
}

// GenerateClientSecret creates a secure client secret in the format: sk_<clean_base64>
func GenerateClientSecret() (string, error) {
	randomBytes, err := generateRandomBytes(32)
	if err != nil {
		return "", err
	}

	appConfig := config.LoadAppConfig()
	h := hmac.New(sha256.New, []byte(appConfig.AppKey))
	h.Write(randomBytes)
	hashed := h.Sum(nil)

	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed)
	cleaned := removeDashesAndUnderscores(encoded)
	return fmt.Sprintf("sk_%s", cleaned), nil
}

// removeDashesAndUnderscores removes '-' and '_' characters from the input string
func removeDashesAndUnderscores(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "-", ""), "_", "")
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

// GenerateFileKey generates a secure file key
func GenerateFileKey(fileID string) (string, error) {
	randomBytes, err := generateRandomBytes(16)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	appConfig := config.LoadAppConfig()
	h := hmac.New(sha256.New, []byte(appConfig.AppKey))
	h.Write(randomBytes)
	hashed := h.Sum(nil)

	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed)
	cleaned := removeDashesAndUnderscores(encoded)

	return fmt.Sprintf("%s_%s", fileID, cleaned), nil
}

// GenerateSalt generates a secure random salt
func GenerateSalt() (string, error) {
	randomBytes, err := generateRandomBytes(16)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	cleaned := removeDashesAndUnderscores(encoded)

	return cleaned, nil
}

// GenerateNonce generates a secure random nonce
func GenerateNonce() (string, error) {
	randomBytes, err := generateRandomBytes(12)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	cleaned := removeDashesAndUnderscores(encoded)

	return cleaned, nil
}
