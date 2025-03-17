package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"finanapp/internal/utils"
	"log"
	"net/http"
	"time"
)

// UserDashboardReact shows the logged-in user's dashboard
func UserDashboard(w http.ResponseWriter, r *http.Request) {
	log.Println("UserDashboard: Request received.")

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("UserDashboard: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	log.Printf("UserDashboard: User retrieved from context - ID: %d, Name: %s %s, Email: %s",
		user.UserProfileID, user.FirstName, user.LastName, user.EmailAddress)

	// Prepare response data
	responseData := map[string]interface{}{
		"user": user,
	}

	// Set the response header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Return the response data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Printf("UserDashboard: Error encoding response data: %v", err)
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	log.Println("UserDashboard: Response sent successfully.")
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Method not allowed",
		})
		return
	}

	// Parse JSON body
	var requestData struct {
		UserId      int    `json:"userId"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		DateOfBirth string `json:"dateOfBirth"`
	}

	// Decode the incoming JSON
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid input data: "+err.Error(), http.StatusBadRequest)
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

	if requestData.DateOfBirth == "" {
		validationErrors = append(validationErrors, "Date of birth is required")
	} else {
		// Parse the date
		_, err := time.Parse("2006-01-02", requestData.DateOfBirth)
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

	// Fetch the user profile to be updated
	var userProfile models.UserProfile
	err = db.DB.QueryRow(`
		SELECT user_profile_id, first_name, last_name, email_address, user_password, date_of_birth
		FROM user_profile WHERE user_profile_id = $1`, requestData.UserId).Scan(
		&userProfile.UserProfileID, &userProfile.FirstName, &userProfile.LastName,
		&userProfile.EmailAddress, &userProfile.UserPassword, &userProfile.DateOfBirth)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse the date of birth
	parsedDob, err := time.Parse("2006-01-02", requestData.DateOfBirth)
	if err != nil {
		http.Error(w, "Error parsing date of birth", http.StatusInternalServerError)
		return
	}

	// Update the user profile in the database
	_, err = db.DB.Exec(`
		UPDATE user_profile
		SET first_name = $1, last_name = $2, date_of_birth = $3, updated_at = CURRENT_TIMESTAMP
		WHERE user_profile_id = $4`,
		utils.Capitalize(requestData.FirstName),
		utils.Capitalize(requestData.LastName),
		parsedDob,
		requestData.UserId,
	)

	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"message":     "Profile updated successfully",
		"userId":      requestData.UserId,
		"firstName":   requestData.FirstName,
		"lastName":    requestData.LastName,
		"dateOfBirth": requestData.DateOfBirth,
	})
}
