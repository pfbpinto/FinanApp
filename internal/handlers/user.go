package handlers

import (
	"context"
	"encoding/json"
	"finanapp/internal/auth"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"finanapp/internal/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// UserDashboard shows the logged-in user's dashboard
func UserDashboard(w http.ResponseWriter, r *http.Request) {
	// Get the user session cookie value (JWT token)
	cookie, err := r.Cookie("user_session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate the JWT token
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to validate token", http.StatusUnauthorized)
		return
	}

	// Extract email from claims
	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Invalid token data", http.StatusUnauthorized)
		return
	}

	// Check if user data is cached in Redis
	ctx := context.Background()
	cachedUser, err := db.RDB.Get(ctx, email).Result()
	if err == nil && cachedUser != "" {
		// User is found in Redis; render using cached data
		data := map[string]interface{}{
			"User": email,
			"Name": cachedUser,
		}
		RenderTemplate(w, r, "user.html", data)
		return
	}

	// If not in Redis, load the user from the database
	var user models.User
	result := db.DB.Where("email_address = ?", email).First(&user)
	if result.Error != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Cache the user data in Redis for future requests
	err = db.RDB.Set(ctx, email, user.FirstName, 15*time.Minute).Err()
	if err != nil {
		// Log or handle the Redis cache error (optional)
	}

	// Render the template
	data := map[string]interface{}{
		"User": email,
		"Name": user.FirstName,
	}
	RenderTemplate(w, r, "user.html", data)
}

// UserDashboardReact shows the logged-in user's dashboard
func UserDashboardReact(w http.ResponseWriter, r *http.Request) {
	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	var userAsset []models.UserAsset
	var assetTypes []models.AssetType
	var userIncome []models.UserIncome
	var userExpense []models.UserExpenditure
	var taxes []models.Tax
	var userGroups []models.UserGroup

	// Query for user assets data with AssetType and UserAssetTax preload
	if result := db.DB.Preload("AssetType").
		Preload("UserAssetTaxes.Tax"). // Preload Tax relationship inside UserAssetTax
		Where("user_id = ?", user.ID).
		Find(&userAsset); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve user assets data", http.StatusInternalServerError)
		return
	}

	// Fetch all asset types from the database
	if result := db.DB.Find(&assetTypes); result.Error != nil {
		http.Error(w, "Failed to retrieve asset types data", http.StatusInternalServerError)
		return
	}

	// Query for user income data with IncomeType preload
	if result := db.DB.
		Preload("IncomeType").
		Preload("UserTaxes"). // Primeiro preenche UserTaxes
		Where("user_id = ?", user.ID).
		Find(&userIncome); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve user income data", http.StatusInternalServerError)
		return
	}

	// Agora carrega os impostos de cada UserTax individualmente
	for i := range userIncome {
		db.DB.Preload("Tax").Find(&userIncome[i].UserTaxes)
	}

	// Query for user expenditure data with Expenditure preload
	if result := db.DB.Preload("Expenditure").Where("user_id = ?", user.ID).Find(&userExpense); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve user expenditure data", http.StatusInternalServerError)
		return
	}

	// Query for Taxes data with TaxType preload
	if result := db.DB.Preload("TaxType").Find(&taxes); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve Taxes data", http.StatusInternalServerError)
		return
	}

	// Fetch User Group from the database
	if result := db.DB.Find(&userGroups); result.Error != nil {
		http.Error(w, "Failed to retrieve Group data", http.StatusInternalServerError)
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"user":        user,
		"userAsset":   userAsset,
		"assetTypes":  assetTypes,
		"userIncome":  userIncome,
		"userExpense": userExpense,
		"taxes":       taxes,
		"userGroups":  userGroups,
	}

	// Set the response header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Return the response data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}
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

	// Fetch the user to be updated
	var user models.User
	if result := db.DB.First(&user, requestData.UserId); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse the date of birth
	parsedDob, err := time.Parse("2006-01-02", requestData.DateOfBirth)
	if err != nil {
		http.Error(w, "Error parsing date of birth", http.StatusInternalServerError)
		return
	}

	// Prepare update data
	updateData := map[string]interface{}{
		"FirstName":   utils.Capitalize(requestData.FirstName),
		"LastName":    utils.Capitalize(requestData.LastName),
		"DataOfBirth": &parsedDob,
	}

	// Update the user
	if result := db.DB.Model(&user).Updates(updateData); result.Error != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	// For now, let's simulate a successful update
	//fmt.Printf("Received data: %+v\n", requestData)

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
