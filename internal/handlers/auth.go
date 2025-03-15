package handlers

import (
	"context"
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

// Return Messages
const (
	ErrInvalidUser           = "Invalid User"
	ErrInvalidCredentials    = "Invalid Credentials"
	ErrInternalServerError   = "Internal Server Error"
	ErrEmailRequired         = "Email and password are required"
	ErrEmailAlreadyInUse     = "Email already in use"
	ErrEmailInvalidFormat    = "Invalid email format"
	ErrWeakPassword          = "Weak password. Include uppercase, lowercase, numbers, and special characters."
	ErrUsernameAlreadyInUse  = "Username already in use"
	ErrInvalidUsernameFormat = "Username must be between 3 and 50 characters, and can only contain letters, numbers, and underscores."
	ErrRedisSession          = "Failed to create session in redis."
	ErrRedisSessionExpire    = "Failed to set session expiration."
	SuccessRegistration      = "Registration successful! Please log in."
)

// Helper function for HTTP redirection with error or success messages
func redirectWithMessage(w http.ResponseWriter, r *http.Request, path, messageType, message string) {
	http.Redirect(w, r, fmt.Sprintf("%s?%s=%s", path, messageType, message), http.StatusSeeOther)
}

// Struct for validating email input
type LoginRequest struct {
	EmailAddress string `validate:"required,email"`
	Password     string `validate:"required"`
}

type LoginRequestReact struct {
	EmailAddress string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("Attempt to login")
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate input fields
		if email == "" || password == "" {
			redirectWithMessage(w, r, "/login", "error", ErrEmailRequired)
			return
		}

		// Search for the user in the database by email
		var user models.User
		result := db.DB.Where("email_address = ?", email).First(&user)
		if result.Error != nil {
			redirectWithMessage(w, r, "/login", "error", ErrInvalidUser)
			return
		}

		// Verify password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			redirectWithMessage(w, r, "/login", "error", ErrInvalidCredentials)
			return
		}

		// Create token JWT
		token, err := auth.CreateJWT(email)
		if err != nil {
			redirectWithMessage(w, r, "/login", "error", ErrInternalServerError)
			return
		}

		// Store session data in Redis with the token as the key
		sessionData := map[string]string{
			"user_id": fmt.Sprintf("%d", user.ID),
			"email":   email,
			"role":    user.UserType.Name,
		}
		ctx := context.Background()
		err = rdb.HSet(ctx, token, sessionData).Err()
		if err != nil {
			redirectWithMessage(w, r, "/login", "error", ErrRedisSession)
			return
		}

		// Set TTL for the session
		err = rdb.Expire(ctx, token, 24*time.Hour).Err()
		if err != nil {
			redirectWithMessage(w, r, "/login", "error", ErrRedisSessionExpire)
			return
		}

		// Set a secure session cookie with the JWT
		http.SetCookie(w, &http.Cookie{
			Name:     "user_session",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true, // HTTPS only
			Path:     "/",
		})

		// Update user's last login
		now := time.Now()
		user.LastLogin = &now
		if err := db.DB.Save(&user).Error; err != nil {
			redirectWithMessage(w, r, "/login", "error", "Erro ao atualizar dados de login")
			return
		}

		// Redirect to the user dashboard
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	// Render the login page
	data := map[string]interface{}{
		"ErrorMessage":   r.URL.Query().Get("error"),
		"SuccessMessage": r.URL.Query().Get("success"),
	}
	RenderTemplate(w, r, "login.html", data)
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
		var user models.User
		result := db.DB.Preload("UserType").Where("email_address = ?", loginReq.EmailAddress).First(&user)
		if result.Error != nil {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
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
			"user_id":  fmt.Sprintf("%d", user.ID),
			"firsName": user.FirstName,
			"email":    loginReq.EmailAddress,
			"role":     user.UserType.Name,
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
				"role":      user.UserType.Name,
			},
			"token": token,
		}

		// Update user's last login
		now := time.Now()
		user.LastLogin = &now
		if err := db.DB.Save(&user).Error; err != nil {
			redirectWithMessage(w, r, "/login", "error", "Erro ao atualizar dados de login")
			return
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

// Logout handles the user logout request
func Logout(w http.ResponseWriter, r *http.Request) {
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
			} else {
				log.Printf("Session for token %s successfully removed from Redis", cookie.Value)
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

	// Redirect the user to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register handles the user registration request
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Collect form parameters
		email := r.FormValue("email")
		password := r.FormValue("password")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		dob := r.FormValue("dob")

		// Validate input data
		if firstName == "" || email == "" || password == "" || lastName == "" || dob == "" {
			redirectWithMessage(w, r, "/register", "error", "All fields are required")
			return
		}

		if !utils.ValidateUsername(firstName) {
			redirectWithMessage(w, r, "/register", "error", ErrInvalidUsernameFormat)
			return
		}

		if !utils.ValidateUsername(lastName) {
			redirectWithMessage(w, r, "/register", "error", ErrInvalidUsernameFormat)
			return
		}

		if !utils.ValidateEmail(email) {
			redirectWithMessage(w, r, "/register", "error", ErrEmailInvalidFormat)
			return
		}

		if !utils.ValidatePassword(password) {
			redirectWithMessage(w, r, "/register", "error", ErrWeakPassword)
			return
		}

		// Check if username or email already exists
		var user models.User
		if result := db.DB.Where("email_address = ?", email).First(&user); result.Error == nil {
			redirectWithMessage(w, r, "/register", "error", ErrEmailAlreadyInUse)
			return
		}

		hashedPassword, err := utils.HashPasswordSafely(password)
		if err != nil {
			redirectWithMessage(w, r, "/register", "error", ErrInternalServerError)
			return
		}

		const dateFormat = "2006-01-02"

		// Converta a string para time.Time
		parsedDob, err := time.Parse(dateFormat, dob)
		if err != nil {
			log.Println("Error parsing date:", err)
			redirectWithMessage(w, r, "/register", "error", "Invalid date format")
			return
		}

		user = models.User{
			UserTypeID:   2,
			Password:     hashedPassword,
			FirstName:    firstName,
			LastName:     lastName,
			EmailAddress: email,
			DataOfBirth:  &parsedDob,
		}

		// Create the user in the database
		if result := db.DB.Create(&user); result.Error != nil {
			redirectWithMessage(w, r, "/register", "error", ErrInternalServerError)
			return
		}

		/*
			// Prepare data to send to NSQ
			userData := map[string]interface{}{
				"username": username,
				"email_address": email,
				"password": password, // Send the unencrypted password for the Consumer to process
				"user_type_id": 1,   // Mudado para user_type_id
			}

			// Convert data to JSON
			message, err := json.Marshal(userData)
			if err != nil {
				log.Println("Error creating message for NSQ:", err)
				redirectWithMessage(w, r, "/register", "error", ErrInternalServerError)
				return
			}

			// Create the producer and publish the message
			producer, err := messaging.NewProducer()
			if err != nil {
				log.Printf("Error creating NSQ producer: %v", err)
				redirectWithMessage(w, r, "/register", "error", ErrInternalServerError)
				return
			}

			// Publish the message to NSQ
			err = producer.UserRegistration(message)
			if err != nil {
				log.Printf("Error publishing message to NSQ: %v", err)
				redirectWithMessage(w, r, "/register", "error", ErrInternalServerError)
				return
			}
		*/

		// Redirect the user to the login page with a success message
		redirectWithMessage(w, r, "/login", "success", SuccessRegistration)
		return
	}

	// Display the registration form
	data := map[string]interface{}{"ErrorMessage": r.URL.Query().Get("error")}
	RenderTemplate(w, r, "register.html", data)
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
		query := `SELECT CreateUser($1, $2, $3, $4, $5)` // Call to the function

		// Capture the response JSON returned by the stored procedure
		var responseMessage string
		err = db.DB.Raw(query, requestData.FirstName, requestData.LastName, requestData.Email, hashedPassword, requestData.Dob).Scan(&responseMessage).Error
		if err != nil {
			http.Error(w, "Error executing stored procedure", http.StatusInternalServerError)
			return
		}

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

// ForgotPassword handles the password reset request
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		var user models.User
		if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
			redirectWithMessage(w, r, "/forgot-password", "error", "User not found")
			return
		}

		resetToken := utils.GenerateResetToken()
		utils.SendPasswordResetEmail(user.EmailAddress, resetToken)
		redirectWithMessage(w, r, "/forgot-password", "success", "Check your email for the reset link")
		return
	}

	RenderTemplate(w, r, "forgot-password.html", nil)
}

// AuthStatus checks if the user is authenticated by validating the JWT in the cookie
func AuthStatus(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	userId := uint(0)
	firstName := ""
	lastName := ""
	emailAddress := ""
	userTypeID := uint(0)
	role := ""
	dataOfBirth := ""
	createdAt := ""
	lastLogin := ""

	// 1. Verificar o cookie da sessão
	cookie, err := r.Cookie("user_session")
	if err != nil || cookie.Value == "" {
		// Se o cookie não existir ou estiver inválido, a autenticação falha
		respondUnauthorized(w)
		return
	}

	// 2. Validar o token JWT
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		// Se o token for inválido, a autenticação falha
		respondUnauthorized(w)
		return
	}

	// 3. Obter o e-mail do token JWT
	email, ok := claims["email"].(string)
	if !ok {
		// Se o e-mail não estiver no token, a autenticação falha
		respondUnauthorized(w)
		return
	}

	// 4. Buscar os dados do usuário no banco de dados
	var user models.User
	result := db.DB.Preload("UserType").Where("email_address = ?", email).First(&user)
	if result.Error != nil {
		// Se não encontrar o usuário, a autenticação falha
		respondUnauthorized(w)
		return
	}

	// 5. Se tudo estiver correto, definir os dados do usuário
	authenticated = true
	userId = user.ID
	firstName = user.FirstName
	lastName = user.LastName
	emailAddress = user.EmailAddress
	userTypeID = user.UserTypeID
	role = user.UserType.Name
	if user.DataOfBirth != nil {
		dataOfBirth = user.DataOfBirth.Format("2006-01-02")
	}
	createdAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	if user.LastLogin != nil {
		lastLogin = user.LastLogin.Format("2006-01-02 15:04:05")
	}

	// 6. Responder com os dados do usuário
	response := map[string]interface{}{
		"authenticated": authenticated,
		"UserId":        userId,
		"firstName":     firstName,
		"lastName":      lastName,
		"emailAddress":  emailAddress,
		"userTypeID":    userTypeID,
		"role":          role,
		"dataOfBirth":   dataOfBirth,
		"createdAt":     createdAt,
		"lastLogin":     lastLogin,
	}

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
