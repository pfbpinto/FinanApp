package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/smtp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a bcrypt hash from a plain password.
func HashPassword(password string) (string, error) {
	// Generate the password hash with a cost of 12 (which is the recommended default)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating password hash:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plain password with a hashed password.
func CheckPasswordHash(password, hashedPassword string) bool {
	// Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Password doesn't match
		return false
	}
	// Password matches
	return true
}

func SendPasswordResetEmail(email, resetToken string) {
	from := "your-email@example.com"
	password := "your-email-password"

	to := []string{email}
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	message := []byte("Subject: Password Reset Request\n" +
		"\n" +
		"Click the following link to reset your password:\n" +
		"http://example.com/reset-password?token=" + resetToken)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateResetToken generates a random token for password reset
func GenerateResetToken() string {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		log.Fatal("Error generating reset token:", err)
	}
	return hex.EncodeToString(token)
}
