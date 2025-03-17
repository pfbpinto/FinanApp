package middlewares

import (
	"context"
	"database/sql"
	"finanapp/internal/auth"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"log"
	"net/http"
)

// AuthMiddleware checks if the user is authenticated and passes the data to the context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Check if the session cookie exists
		cookie, err := r.Cookie("user_session")
		if err != nil || cookie.Value == "" {
			log.Println("AUTH-MID: Unauthorized access - missing or invalid session cookie")
			unauthorized(w, r) // Respond with unauthorized if cookie is not found
			return
		}

		// 2. Validate the JWT token using the 'ValidateJWT' function from 'auth'
		claims, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			log.Println("AUTH-MID: Unauthorized access - invalid JWT token")
			unauthorized(w, r) // Respond with unauthorized if JWT is invalid
			return
		}

		// 3. Extract the email from the JWT claims
		email, ok := claims["email"].(string)
		if !ok {
			log.Println("AUTH-MID: Unauthorized access - email not found in JWT claims")
			unauthorized(w, r) // Respond with unauthorized if email is not found in claims
			return
		}

		// 4. Load the user from the database using the extracted email
		log.Printf("AUTH-MID: Looking up user in the database for email: %s\n", email)
		var user models.UserProfile

		// Correct query with properly named columns
		query := `SELECT userprofileid, firstname, lastname, emailaddress, dateofbirth, createdat 
          FROM userprofile WHERE emailaddress = $1`

		log.Printf("AUTH-MID: Executing query: %s with email: %s\n", query, email)

		err = db.GetDB().QueryRow(query, email).Scan(
			&user.UserProfileID, // userprofileid
			&user.FirstName,     // firstname
			&user.LastName,      // lastname
			&user.EmailAddress,  // emailaddress
			&user.DateOfBirth,   // dateofbirth
			&user.CreatedAt,     // createdat
		)

		// 5. Check if the user was found
		if err == sql.ErrNoRows {
			log.Println("AUTH-MID: Unauthorized access - user not found in the database")
			unauthorized(w, r) // Respond with unauthorized if user is not found in the database
			return
		} else if err != nil {
			log.Printf("AUTH-MID: Database error - %v\n", err)
			http.Error(w, "Database error", http.StatusInternalServerError) // Handle database errors
			return
		}
		// 6. Add the user to the request context
		log.Printf("AUTH-MID: User authenticated - %s %s (ID: %d)\n", user.FirstName, user.LastName, user.UserProfileID)
		ctx := context.WithValue(r.Context(), "authenticated", true)
		ctx = context.WithValue(ctx, "user", user)

		// 7. Call the next handler, passing the updated context
		log.Println("AUTH-MID: Passing control to the next handler")
		next(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
