package middlewares

import (
	"context"
	"finanapp/internal/auth"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"net/http"
)

// AuthMiddleware checks if the user is authenticated and passes the data to the context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Checks if there is a session cookie
		cookie, err := r.Cookie("user_session")
		if err != nil || cookie.Value == "" {
			unauthorized(w, r)
			return
		}

		// Use the ValidateJWT function from 'auth' to validate the JWT token in the cookie
		claims, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			unauthorized(w, r)
			return
		}

		// Extract the email from the token claims
		email, ok := claims["email"].(string)
		if !ok {
			unauthorized(w, r)
			return
		}

		// Load user with UserType from database
		var user models.User
		result := db.DB.Preload("UserType").Where("email_address = ?", email).First(&user)
		if result.Error != nil {
			unauthorized(w, r)
			return
		}

		// Add User to context
		r = r.WithContext(context.WithValue(r.Context(), "authenticated", true))
		r = r.WithContext(context.WithValue(r.Context(), "user", user))

		// Call next handler
		next(w, r)
	}
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
