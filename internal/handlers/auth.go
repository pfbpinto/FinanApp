package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"finanapp/internal/auth"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"finanapp/internal/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var rdb *redis.Client // Redis client already iniciated on main.go

func init() {
	// Inicialize o cliente Redis (exemplo)
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	log.Printf("Redis Running")
}

// Struct for validating email input
type LoginRequestReact struct {
	EmailAddress string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
}

func LoginReact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var loginReq LoginRequestReact
		// Parse JSON
		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate fields
		if loginReq.EmailAddress == "" || loginReq.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Search DB for User
		var user models.UserProfile
		query := `SELECT userprofileid, firstname, lastname, userpassword, emailaddress FROM userprofile WHERE emailaddress = $1`
		err = db.GetDB().QueryRow(query, loginReq.EmailAddress).Scan(
			&user.UserProfileID, &user.FirstName, &user.LastName, &user.UserPassword, &user.EmailAddress,
		)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(loginReq.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Token JWT creation
		token, err := auth.CreateJWT(loginReq.EmailAddress)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Save data on Redis session
		sessionData := map[string]string{
			"user_id":  fmt.Sprintf("%d", user.UserProfileID),
			"firsName": user.FirstName,
			"email":    loginReq.EmailAddress,
		}
		ctx := context.Background()
		err = rdb.HSet(ctx, token, sessionData).Err()
		if err != nil {
			http.Error(w, "Redis session error", http.StatusInternalServerError)
			return
		}

		// Set Redis session expiration
		err = rdb.Expire(ctx, token, 24*time.Hour).Err()
		if err != nil {
			http.Error(w, "Redis session expiration error", http.StatusInternalServerError)
			return
		}

		// set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "user_session",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		response := map[string]interface{}{
			"status": "success",
			"user": map[string]interface{}{
				"email":     user.EmailAddress,
				"firstName": user.FirstName,
			},
			"token": token,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

// LogoutReact handles the user logout request for React
func LogoutReact(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session token from the cookie
	cookie, err := r.Cookie("user_session")
	if err == nil && cookie.Value != "" {

		// Use the token (JWT) as the key to check if the session exists in Redis
		ctx := context.Background()
		exists, err := rdb.Exists(ctx, cookie.Value).Result()
		if err != nil {
			log.Printf("Redis exists check error: %v", err)
		} else if exists == 0 {
			log.Printf("No session found in Redis for token: %s", cookie.Value)
		} else {
			// Delete the session data from Redis
			err := rdb.Del(ctx, cookie.Value).Err()
			if err != nil {
				log.Printf("Redis error during logout: %v", err)
			}
		}
	}

	// Invalidate the cookie on the client
	http.SetCookie(w, &http.Cookie{
		Name:     "user_session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expired in the past
		HttpOnly: true,
		Secure:   true, // HTTPS only
		Path:     "/",
	})

	// Respond with a JSON success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Successfully logged out"}`))
}

func RegisterReact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parsing JSON body
		var requestData struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Dob       string `json:"dob"`
		}

		// Decode the request body into the requestData struct
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPasswordSafely(requestData.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Preparing the call to the stored procedure in the database
		query := `SELECT CreateUser($1, $2, $3, $4, $5)`

		var responseMessage string

		err = db.GetDB().QueryRow(query, requestData.FirstName, requestData.LastName, requestData.Email, hashedPassword, requestData.Dob).Scan(&responseMessage)
		if err != nil {
			log.Printf("Error executing stored procedure: %v", err)
			http.Error(w, "Error executing stored procedure", http.StatusInternalServerError)
			return
		}

		log.Printf("Stored procedure executed successfully. Response: %s", responseMessage)

		// Parse the response message (which is in JSON format) into a map
		var response map[string]interface{}
		err = json.Unmarshal([]byte(responseMessage), &response)
		if err != nil {
			http.Error(w, "Error parsing response JSON", http.StatusInternalServerError)
			return
		}

		// Send appropriate HTTP status code based on the response status
		if response["status"] == "error" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest) // Send 400 Bad Request on error
			json.NewEncoder(w).Encode(response)  // Return the error response
		} else {
			// If successful, send 200 OK
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)        // Send 200 OK
			json.NewEncoder(w).Encode(response) // Return success message
		}
		return
	}

	// If the request method is not POST, return method not allowed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	response := map[string]interface{}{
		"status":  "error",
		"message": "Method not allowed",
	}
	json.NewEncoder(w).Encode(response)
}

// AuthStatus checks if the user is authenticated by validating the JWT in the cookie
func AuthStatus(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	userId := uint(0)
	firstName := ""
	lastName := ""
	emailAddress := ""
	dataOfBirth := ""
	createdAt := ""

	// 1. Check the session cookie
	cookie, err := r.Cookie("user_session")
	if err != nil || cookie.Value == "" {
		log.Println("AUTH-STATS: Auth failed - Missing or invalid session cookie")
		respondUnauthorized(w)
		return
	}

	// 2. Validate the JWT token
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		log.Println("AUTH-STATS: Auth failed - Invalid JWT token")
		respondUnauthorized(w)
		return
	}

	// 3. Get the email from the JWT token
	email, ok := claims["email"].(string)
	if !ok {
		log.Println("AUTH-STATS: Auth failed - Email not found in JWT claims")
		respondUnauthorized(w)
		return
	}

	// 4. Retrieve the user data from the database
	// Query the database for the user's information using the email
	query := `SELECT userprofileid, firstname, lastname, emailaddress, dateofbirth, createdat FROM userprofile WHERE emailaddress = $1`

	log.Printf("AUTH-STATS: Executing query: %s with email: %s\n", query, email)

	// Execute the query and store the result in 'row'
	row := db.GetDB().QueryRow(query, email)

	// Scan the result into variables
	var dob sql.NullString

	err = row.Scan(
		&userId, &firstName, &lastName, &emailAddress,
		&dob, &createdAt,
	)
	// Check if no user was found
	if err == sql.ErrNoRows {
		log.Printf("AUTH-STATS: Auth failed - User not found for email: %s\n", email)
		respondUnauthorized(w)
		return
	} else if err != nil {
		log.Printf("AUTH-STATS: Database error - %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 5. If the user is found, format the data
	// If 'dob' is not NULL, format it as a string
	if dob.Valid {
		dataOfBirth = dob.String
	}

	authenticated = true

	// 6. Respond with the user data in JSON format
	// Prepare the response data
	response := map[string]interface{}{
		"authenticated": authenticated,
		"UserId":        userId,
		"firstName":     firstName,
		"lastName":      lastName,
		"emailAddress":  emailAddress,
		"dataOfBirth":   dataOfBirth,
		"createdAt":     createdAt,
	}

	// Log success message for debugging purposes
	log.Printf("AUTH-STATS: Auth successful for user: %s\n", email)

	// Set the response header and encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to respond with a 401 Unauthorized status
func respondUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
}

/*
func RegisterReact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parsing JSON body
		var requestData struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Dob       string `json:"dob"`
		}

		// Decode the request body into the requestData struct
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}

		// Validate input data
		var validationErrors []string

		if requestData.FirstName == "" {
			validationErrors = append(validationErrors, "First name is required")
		} else if !utils.ValidateUsername(requestData.FirstName) {
			validationErrors = append(validationErrors, "Invalid First name format")
		}

		if requestData.LastName == "" {
			validationErrors = append(validationErrors, "Last name is required")
		} else if !utils.ValidateUsername(requestData.LastName) {
			validationErrors = append(validationErrors, "Invalid Last name format")
		}

		if requestData.Email == "" {
			validationErrors = append(validationErrors, "Email is required")
		} else if !utils.ValidateEmail(requestData.Email) {
			validationErrors = append(validationErrors, "Invalid email format")
		}

		if requestData.Password == "" {
			validationErrors = append(validationErrors, "Password is required")
		} else if !utils.ValidatePassword(requestData.Password) {
			validationErrors = append(validationErrors, "Password is too weak")
		}

		if requestData.Dob == "" {
			validationErrors = append(validationErrors, "Date of birth is required")
		} else {
			// Parse the date
			_, err := time.Parse("2006-01-02", requestData.Dob)
			if err != nil {
				validationErrors = append(validationErrors, "Invalid date format")
			}
		}

		// If there are validation errors, return them to the frontend
		if len(validationErrors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			response := map[string]interface{}{
				"status":  "error",
				"message": validationErrors,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if email already exists
		var existingUser models.User
		if result := db.DB.Where("email_address = ?", requestData.Email).First(&existingUser); result.Error == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			response := map[string]interface{}{
				"status":  "error",
				"message": []string{"Email is already in use"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPasswordSafely(requestData.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Parse the date of birth
		parsedDob, err := time.Parse("2006-01-02", requestData.Dob)
		if err != nil {
			http.Error(w, "Error parsing date of birth", http.StatusInternalServerError)
			return
		}

		// Create the user
		user := models.User{
			UserTypeID:   2,
			Password:     hashedPassword,
			FirstName:    utils.Capitalize(requestData.FirstName),
			LastName:     utils.Capitalize(requestData.LastName),
			EmailAddress: requestData.Email,
			DataOfBirth:  &parsedDob,
		}

		if result := db.DB.Create(&user); result.Error != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Send a success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"status":  "success",
			"message": "User successfully registered",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// If the request method is not POST, return method not allowed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	response := map[string]interface{}{
		"status":  "error",
		"message": "Method not allowed",
	}
	json.NewEncoder(w).Encode(response)
}

*/
