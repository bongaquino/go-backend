package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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

// GenerateOTPSecret generates a TOTP secret for the user
func GenerateOTPSecret(userID string) (string, error) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "KoneksiApp",
		AccountName: userID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP secret: %w", err)
	}
	return secret.Secret(), nil
}

// GenerateQRCode generates a QR code URL for the TOTP secret
func GenerateQRCode(userID, secret string) (string, error) {
	key, err := otp.NewKeyFromURL(fmt.Sprintf("otpauth://totp/KoneksiApp:%s?secret=%s&issuer=KoneksiApp", userID, secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	return key.URL(), nil
}
