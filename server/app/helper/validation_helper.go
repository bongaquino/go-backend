package helper

import (
	"fmt"
	"regexp"
)

// ValidatePassword checks if a given password meets specific criteria:
// - At least one uppercase letter (A–Z)
// - At least one lowercase letter (a–z)
// - At least one special character (e.g., ! @ # $ % ^ & *)
// - At least 8 characters long
func ValidatePassword(password string) (bool, error) {
	// Check length
	if len(password) < 8 {
		return false, fmt.Errorf("password must be at least 8 characters long")
	}

	// Define regex patterns for each requirement
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	// Check all conditions
	if hasUpperCase && hasLowerCase && hasSpecialChar {
		return true, nil
	}

	// Return error if any condition fails
	return false, fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one special character")
}
