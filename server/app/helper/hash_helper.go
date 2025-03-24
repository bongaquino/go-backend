package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash hashes a data using bcrypt
func Hash(data string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckHash compares a hashed data with a plain text data
func CheckHash(data, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}
