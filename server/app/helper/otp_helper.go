package helper

import (
	"fmt"
	"koneksi/server/config"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateOTPSecret generates a TOTP secret for the user
func GenerateOTPSecret(userID string) (string, error) {
	appConfig := config.LoadAppConfig()

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      appConfig.AppName,
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

// VerifyOTP verifies the provided OTP against the stored secret
func VerifyOTP(secret, otp string) bool {
	return totp.Validate(otp, secret)
}
