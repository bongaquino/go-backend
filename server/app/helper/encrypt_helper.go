package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"koneksi/server/config"
)

// Encrypt encrypts the given
func Encrypt(data string) (string, error) {
	appConfig := config.LoadAppConfig()

	block, err := aes.NewCipher([]byte(appConfig.AppKey))
	if err != nil {
		return "", err
	}

	plaintext := []byte(data)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// Generate a random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Encode the ciphertext to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given base64-encoded ciphertext using the provided key
func Decrypt(encryptedData string) (string, error) {
	appConfig := config.LoadAppConfig()

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(appConfig.AppKey))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
