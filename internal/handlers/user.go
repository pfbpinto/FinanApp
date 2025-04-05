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

// UserIncome shows the logged-in user's Incomes
func UserIncome(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("UserDashboard: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Get Database
	database := db.GetDB()

	// Query currency
	currencyQuery := `SELECT currencyid, currencyname, currencyabreviation, currencysymbol, createdat FROM currency`
	currencyRows, err := database.Query(currencyQuery)
	if err != nil {
		log.Println("Error fetching Currency Type:", err)
		http.Error(w, "Error fetching Currency Type", http.StatusInternalServerError)
		return
	}
	defer currencyRows.Close()

	var currency []models.Currency
	for currencyRows.Next() {
		var cc models.Currency
		if err := currencyRows.Scan(&cc.CurrencyID, &cc.CurrencyName, &cc.CurrencyAbreviation, &cc.CurrencySymbol, &cc.CreatedAt); err != nil {
			log.Println("Error scanning Currency Type:", err)
			continue
		}
		currency = append(currency, cc)
	}

	// Query Recurrency
	recurrencyQuery := `SELECT recurrencyid, recurrencyname, recurrencyperiod, createdat FROM recurrency`
	recurrencyRows, err := database.Query(recurrencyQuery)
	if err != nil {
		log.Println("Error fetching Recurrency:", err)
		http.Error(w, "Error fetching Recurrency", http.StatusInternalServerError)
		return
	}
	defer recurrencyRows.Close()

	var recurrencies []models.Recurrency
	for recurrencyRows.Next() {
		var rr models.Recurrency
		if err := recurrencyRows.Scan(&rr.RecurrencyID, &rr.RecurrencyName, &rr.RecurrencyPeriod, &rr.CreatedAt); err != nil {
			log.Println("Error scanning Recurrency:", err)
			continue
		}
		recurrencies = append(recurrencies, rr)
	}

	// Query Income Type
	incomeTypeQuery := `SELECT incometypeid, incometypename, incomedescription, entityid, createdat FROM incometype i`
	incomeTypeRows, err := database.Query(incomeTypeQuery)
	if err != nil {
		log.Println("Error fetching Income Type:", err)
		http.Error(w, "Error fetching Income Type", http.StatusInternalServerError)
		return
	}
	defer incomeTypeRows.Close()

	var incomeType []models.IncomeType
	for incomeTypeRows.Next() {
		var it models.IncomeType
		if err := incomeTypeRows.Scan(&it.IncomeTypeID, &it.IncomeTypeName, &it.IncomeDescription, &it.EntityID, &it.CreatedAt); err != nil {
			log.Println("Error scanning Income Type:", err)
			continue
		}
		incomeType = append(incomeType, it)
	}

	// Query User Categories
	userCategoryQuery := `
		SELECT UserCategoryID, UserCategoryName, UserProfileID, EntityID, IsActive
		FROM usercategory WHERE EntityID = 5 AND userprofileid = $1
	`
	userCategoryRows, err := database.Query(userCategoryQuery, user.UserProfileID)
	if err != nil {
		log.Println("Error fetching User Categories:", err)
		http.Error(w, "Error fetching User Categories", http.StatusInternalServerError)
		return
	}
	defer userCategoryRows.Close()

	var userCategories []models.UserCategory
	for userCategoryRows.Next() {
		var uc models.UserCategory
		if err := userCategoryRows.Scan(&uc.UserCategoryID, &uc.UserCategoryName, &uc.UserProfileID, &uc.EntityID, &uc.IsActive); err != nil {
			log.Println("Error scanning User Category:", err)
			continue
		}
		userCategories = append(userCategories, uc)
	}

	// Query Financial User Items
	financialUserItemQuery := `
		SELECT 
			f.FinancialUserItemID, 
			f.FinancialUserItemName, 
			f.EntityID, 
			f.UserEntityID, 
			f.RecurrencyID, 
			f.FinancialUserEntityItemID,
			f.IsActive, 
			f.CreatedAt,
			e.EntityType,
			r.RecurrencyName,
			it.IncomeTypeName
		FROM 
			financialuseritem f
		JOIN 
			Entity e ON f.EntityID = e.EntityID
		JOIN 
			Recurrency r ON f.RecurrencyID = r.RecurrencyID 
		LEFT JOIN 
			incometype it ON f.UserEntityID = it.incometypeid
		WHERE 
			f.UserEntityID = $1 AND e.entityname = 'User' AND f.EntityID = 5
	`

	// Passing the logged-in user's ID to filter the SQL query
	financialUserItemRows, err := database.Query(financialUserItemQuery, user.UserProfileID)
	if err != nil {
		log.Println("Error fetching Financial User Items:", err)
		http.Error(w, "Error fetching Financial User Items", http.StatusInternalServerError)
		return
	}
	defer financialUserItemRows.Close()

	var financialUserItems []models.FinancialUserItem
	for financialUserItemRows.Next() {
		var fui models.FinancialUserItem
		if err := financialUserItemRows.Scan(
			&fui.FinancialUserItemID,
			&fui.FinancialUserItemName,
			&fui.EntityID,
			&fui.UserEntityID,
			&fui.RecurrencyID,
			&fui.FinancialUserEntityItemID,
			&fui.IsActive,
			&fui.CreatedAt,
			&fui.EntityType,
			&fui.RecurrencyName,
			&fui.IncomeTypeName,
		); err != nil {
			log.Println("Error scanning Financial User Item:", err)
			continue
		}

		financialUserItems = append(financialUserItems, fui)
	}

	// Create final response
	response := struct {
		Currency           []models.Currency          `json:"currency"`
		Recurrency         []models.Recurrency        `json:"recurrency"`
		IncomeType         []models.IncomeType        `json:"income_type"`
		UserCategories     []models.UserCategory      `json:"user_categories"`
		FinancialUserItems []models.FinancialUserItem `json:"financial_user_items"`
	}{
		Currency:           currency,
		Recurrency:         recurrencies,
		IncomeType:         incomeType,
		UserCategories:     userCategories,
		FinancialUserItems: financialUserItems,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert to JSON and send response
	json.NewEncoder(w).Encode(response)
}

// UserAsset returns the logged-in user's Assets and Income Types
func UserAsset(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("UserAsset: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Get Database
	database := db.GetDB()

	// Query Income Type
	incomeTypeQuery := `SELECT incometypeid, incometypename, incomedescription, entityid, createdat FROM incometype`
	incomeTypeRows, err := database.Query(incomeTypeQuery)
	if err != nil {
		log.Println("UserAsset: Error fetching Income Type:", err)
		http.Error(w, "Error fetching Income Type", http.StatusInternalServerError)
		return
	}
	defer incomeTypeRows.Close()

	var incomeType []models.IncomeType
	for incomeTypeRows.Next() {
		var it models.IncomeType
		if err := incomeTypeRows.Scan(&it.IncomeTypeID, &it.IncomeTypeName, &it.IncomeDescription, &it.EntityID, &it.CreatedAt); err != nil {
			log.Println("UserAsset: Error scanning Income Type:", err)
			continue
		}
		incomeType = append(incomeType, it)
	}

	// Query User Assets
	userAssetQuery := `
		SELECT 
			ua.UserAssetID,
			ua.UserAssetName,
			ua.UserAssetValueAmount,
			ua.UserAssetAcquisitionBeginDate,
			ua.UserAssetAcquisitionEndDate,
			ua.IsActive,
			ua.CreatedAt,
			at.AssetTypeName
		FROM 
			UserAsset ua
		JOIN 
			AssetType at ON ua.AssetTypeID = at.AssetTypeID
		WHERE 
			ua.UserProfileID = $1
	`

	userAssetRows, err := database.Query(userAssetQuery, user.UserProfileID)
	if err != nil {
		log.Println("UserAsset: Error fetching User Assets:", err)
		http.Error(w, "Error fetching User Assets", http.StatusInternalServerError)
		return
	}
	defer userAssetRows.Close()

	var userAssets []models.UserAsset
	for userAssetRows.Next() {
		var ua models.UserAsset
		if err := userAssetRows.Scan(
			&ua.UserAssetID,
			&ua.UserAssetName,
			&ua.UserAssetValueAmount,
			&ua.UserAssetAcquisitionBeginDate,
			&ua.UserAssetAcquisitionEndDate,
			&ua.IsActive,
			&ua.CreatedAt,
			&ua.AssetTypeName,
			&ua.UserProfileName,
		); err != nil {
			log.Println("UserAsset: Error scanning User Asset:", err)
			continue
		}
		userAssets = append(userAssets, ua)
	}

	// Create final response
	response := struct {
		IncomeType []models.IncomeType `json:"income_type"`
		UserAssets []models.UserAsset  `json:"user_assets"`
	}{
		IncomeType: incomeType,
		UserAssets: userAssets,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("UserAsset: Error encoding response:", err)
	}
}
