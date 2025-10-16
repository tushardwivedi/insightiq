package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength is the minimum required password length
	MinPasswordLength = 8
	// BcryptCost is the cost factor for bcrypt hashing
	BcryptCost = 12
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a plain text password with a hashed password
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength validates password strength requirements
func ValidatePasswordStrength(password string) bool {
	return len(password) >= MinPasswordLength
}
