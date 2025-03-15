package factory

import (
	"finanapp/internal/models"
	"math/rand"
	"time"
)

// GenerateUser creates a dummy user for testing purposes
func GenerateUser() models.User {
	return models.User{
		FirstName:    randomString(8),
		EmailAddress: randomString(5) + "@example.com",
		Password:     "hashed_password", // Replace with a hashed password
		UserTypeID:   1,
		IsActive:     true,
		LastLogin:    nil,
	}
}

// randomString generates a random string of given length
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
