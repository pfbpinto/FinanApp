package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("yoursupersecurekey")

// CreateJWT creates a token JWT for an email
func CreateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expiração de 24 horas
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a token JWT and return it claims
func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the signing method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signature method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Ensure that the token is valid and contains the correct claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Here we could add any additional checks, e.g. expiration, as needed:
	if exp, ok := claims["exp"].(float64); ok {
		// Convert exp to time
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return nil, errors.New("token has expired")
		}
	}

	return claims, nil
}
