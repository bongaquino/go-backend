package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"bongaquino/server/config"

	"golang.org/x/crypto/pbkdf2"
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

// DeriveKey derives a key from the passphrase and salt using PBKDF2
func DeriveKey(passphrase, salt string) ([]byte, error) {
	// Use PBKDF2 to derive a key from the passphrase and salt
	key := pbkdf2.Key([]byte(passphrase), []byte(salt), 4096, 32, sha256.New)
	return key, nil
}

// Create AES-GCM Cipher
func CreateAesGcmCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesGcm, nil
}
