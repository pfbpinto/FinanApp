package utils

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ValidateUsername checks if the username meets the required format
func ValidateUsername(username string) bool {
	// Username must be between 3 and 50 characters and can only contain letters, numbers, and underscores
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	return re.MatchString(username)
}

// ValidateEmail checks the format of the email
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ValidatePassword checks if the password meets the criteria
func ValidatePassword(password string) bool {
	// Min and max password size
	if len(password) < 8 || len(password) > 128 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

	// Returns true only if all criteria are met
	return hasLower && hasUpper && hasDigit && hasSpecial
}

// hashPasswordSafely hashes the user's password using bcrypt with a default cost.
func HashPasswordSafely(password string) (string, error) {
	// Generate the bcrypt hash from the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(hashedPassword), nil
}

func Capitalize(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.ToUpper(string(name[0])) + strings.ToLower(name[1:])
}
