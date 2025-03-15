package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetAssetType returns a list of UserAsset in the database.
func GetIncomeType(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Create a variable to store Income types
	var incomeTypes []models.IncomeType
	var taxes []models.Tax

	// Fetch all Income types from the database
	if result := db.DB.Find(&incomeTypes); result.Error != nil {
		http.Error(w, "Failed to retrieve asset types data", http.StatusInternalServerError)
		return
	}

	// Fetch all tax from type 1: Income from the database
	if result := db.DB.Where("tax_type_id = ? AND user_id = ?", 2, user.ID).Find(&taxes); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve Tax data", http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"incomeTypes": incomeTypes, // Returning the list of IncomeTypes
		"taxes":       taxes,       // Returning the list of taxex
	}

	// Set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}
}

// CreateUserIncome creates a new UserIncome
func CreateUserIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the received body for debugging
	log.Printf("Received body: %s", string(body))

	// Define the structure to match the incoming JSON
	var incomeData struct {
		IncomeName       string `json:"IncomeName"`
		IncomeValue      string `json:"IncomeValue"`
		IncomeType       string `json:"IncomeTypeID"`
		OwningPercentage string `json:"OwningPercentage"`
		IncomeRecurrence string `json:"IncomeRecurrence"`
		IncomeStartDate  string `json:"IncomeStartDate"`
		Shared           *bool  `json:"SharedIncome"`
		UserID           uint   `json:"userID"`
		UserTaxes        []struct {
			ID            uint    `json:"TaxID"`
			TaxName       string  `json:"TaxName"`
			TaxPercentage float64 `json:"TaxPercentage"`
		} `json:"UserTaxes"`
	}

	// Unmarshal the JSON into incomeData
	if err := json.Unmarshal(body, &incomeData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Tax Captured: %v", incomeData.UserTaxes)

	// Validation of required fields
	if incomeData.IncomeName == "" || incomeData.IncomeValue == "" || incomeData.IncomeType == "" || incomeData.UserID < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	if incomeData.IncomeStartDate == "" {
		http.Error(w, "Invalid Date format", http.StatusBadRequest)
		return
	} else {
		// Parse the date
		_, err := time.Parse("2006-01-02", incomeData.IncomeStartDate)
		if err != nil {
			http.Error(w, "Invalid Date format", http.StatusBadRequest)
			return
		}
	}

	// Convert IncomeValue from string to float64
	IncomeValue, err := strconv.ParseFloat(incomeData.IncomeValue, 64)
	if err != nil {
		log.Printf("Error converting IncomeValue: %v", err)
		http.Error(w, "Error converting IncomeValue", http.StatusBadRequest)
		return
	}

	// Convert IncomeType from string to uint
	IncomeTypeID, err := strconv.ParseUint(incomeData.IncomeType, 10, 32)
	if err != nil {
		log.Printf("Error converting IncomeType: %v", err)
		http.Error(w, "Error converting IncomeType", http.StatusBadRequest)
		return
	}

	// Convert OwningPercentage from string to float64
	opValue, err := strconv.ParseFloat(incomeData.OwningPercentage, 64)
	if err != nil {
		log.Printf("Error converting Owning Percentage: %v", err)
		http.Error(w, "Error converting Owning Percentage", http.StatusBadRequest)
		return
	}

	// Parse the date
	parsedDate, err := time.Parse("2006-01-02", incomeData.IncomeStartDate)
	if err != nil {
		http.Error(w, "Error parsing Aquisition date", http.StatusInternalServerError)
		return
	}

	// Create the UserIncome model
	userIncome := models.UserIncome{
		UserID:           incomeData.UserID,
		IncomeTypeID:     uint(IncomeTypeID), // Convert uint64 to uint
		IncomeName:       incomeData.IncomeName,
		IncomeStartDate:  &parsedDate,
		IncomeValue:      IncomeValue,
		IncomeRecurrence: incomeData.IncomeRecurrence,
		SharedIncome:     incomeData.Shared != nil && *incomeData.Shared || false,
		OwningPercentage: opValue,
	}

	// Save to the database using GORM
	if err := db.DB.Create(&userIncome).Error; err != nil {

		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "Error 1062") {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // HTTP 409 - Conflict
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Income name already exists",
				"message": "A Income with this name already exists. Please choose a different name.",
			})
			return
		}

		log.Printf("Error saving to database: %v", err)
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// Check if there are any taxes to process
	if len(incomeData.UserTaxes) > 0 {
		log.Printf("Processing taxes for UserAsset ID %d", userIncome.ID)

		// Iterate through the UserAssetTaxes array which contains full tax details
		for _, taxData := range incomeData.UserTaxes {
			log.Printf("Applying tax: %s (ID: %d, Percentage: %.2f%%)", taxData.TaxName, taxData.ID, taxData.TaxPercentage)

			// Calculate the tax value based on the AssetValue
			taxValue := IncomeValue * (taxData.TaxPercentage / 100)
			log.Printf("Calculated TaxValue: %.2f based on AssetValue: %.2f", taxValue, IncomeValue)

			// Create the record in the UserAssetTax table
			userIncomeTax := models.UserTax{
				UserIncomeID: userIncome.ID, // ID of the newly created UserAsset
				TaxID:        taxData.ID,    // ID of the Tax
				TaxValue:     taxValue,
			}

			log.Printf("Creating UserTax record: %+v", userIncomeTax)

			// Save the UserAssetTax record in the database
			if err := db.DB.Create(&userIncomeTax).Error; err != nil {
				log.Printf("Error saving UserTax record: %v", err)
				continue
			}

			log.Printf("Successfully saved UserTax record for TaxID %d", taxData.ID)
		}
	} else {
		log.Println("No taxes to process for UserIncome")
	}
	// Log successful creation
	log.Printf("UserIncome created: %+v", userIncome)

	// Send response back with created asset data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Income created successfully!",
		"userIncome": userIncome,
	})
}

// UpdateUserIncome updates an existing UserIncome in the database.
func UpdateUserIncome(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the received body for debugging purposes
	log.Printf("Received body: %s", string(body))

	// Define a structure to match the incoming JSON
	var incomeData struct {
		IncomeName       string  `json:"IncomeName"`
		IncomeRecurrence string  `json:"IncomeRecurrence"`
		IncomeValue      float64 `json:"IncomeValue"`
		IncomeTypeID     uint    `json:"IncomeTypeID"`
		OwningPercentage float64 `json:"OwningPercentage"`
		IncomeStartDate  string  `json:"IncomeStartDate"`
		Shared           *bool   `json:"SharedIncome"`
		UserID           uint    `json:"userID"`
		ID               uint    `json:"ID"`
		UserTaxes        []struct {
			ID            uint    `json:"TaxID"`
			TaxName       string  `json:"TaxName"`
			TaxPercentage float64 `json:"TaxPercentage"`
		} `json:"UserTaxes"`
	}

	// Unmarshal the JSON data into the incomeData structure
	if err := json.Unmarshal(body, &incomeData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if incomeData.IncomeName == "" || incomeData.IncomeValue < 1 || incomeData.IncomeTypeID < 1 || incomeData.UserID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	// Validate the income start date format
	if incomeData.IncomeStartDate == "" {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Try parsing the date in different formats
	parsedDate, err := time.Parse("2006-01-02T15:04:05-07:00", incomeData.IncomeStartDate)
	if err != nil {
		parsedDate, err = time.Parse("2006-01-02", incomeData.IncomeStartDate)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
	}

	// Find the existing UserIncome by its ID in the database
	var userIncome models.UserIncome
	if err := db.DB.First(&userIncome, incomeData.ID).Error; err != nil {
		log.Printf("Error finding income with ID %d: %v", incomeData.ID, err)
		http.Error(w, "Income not found", http.StatusNotFound)
		return
	}

	// Update the existing UserIncome with the new data
	userIncome.IncomeName = incomeData.IncomeName
	userIncome.IncomeRecurrence = incomeData.IncomeRecurrence
	userIncome.IncomeValue = incomeData.IncomeValue
	userIncome.IncomeTypeID = incomeData.IncomeTypeID
	userIncome.IncomeStartDate = &parsedDate
	userIncome.OwningPercentage = incomeData.OwningPercentage
	userIncome.SharedIncome = incomeData.Shared != nil && *incomeData.Shared || false

	// Save the updated UserIncome back to the database
	if err := db.DB.Save(&userIncome).Error; err != nil {

		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "Error 1062") {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // HTTP 409 - Conflict
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Income name already exists",
				"message": "A Income with this name already exists. Please choose a different name.",
			})
			return
		}

		log.Printf("Error updating income in database: %v", err)
		http.Error(w, "Error updating income", http.StatusInternalServerError)
		return
	}

	// Clear existing UserTaxes for the income (if any)
	if err := db.DB.Where("user_income_id = ?", userIncome.ID).Delete(&models.UserTax{}).Error; err != nil {
		log.Printf("Error clearing old taxes: %v", err)
		http.Error(w, "Error clearing old taxes", http.StatusInternalServerError)
		return
	}

	// Add new UserTaxes
	for _, tax := range incomeData.UserTaxes {
		newTax := models.UserTax{
			UserIncomeID: userIncome.ID,
			TaxID:        tax.ID,
			TaxValue:     tax.TaxPercentage,
		}

		// Save the new UserTax entry
		if err := db.DB.Create(&newTax).Error; err != nil {
			log.Printf("Error adding tax: %v", err)
			http.Error(w, "Error adding tax", http.StatusInternalServerError)
			return
		}
	}

	// Log successful update
	log.Printf("UserIncome updated: %+v", userIncome)

	// Respond with a success message and the updated income data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Income updated successfully!",
		"userIncome": userIncome,
	})
}

// DeleteIncome deletes an existing UserIncome in the database.
func DeleteIncome(w http.ResponseWriter, r *http.Request) {

	// Ensure the method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the received body for debugging purposes
	log.Printf("Received body: %s", string(body))

	// Define a structure to match the incoming JSON
	var incomeData struct {
		IncomeId uint `json:"itemId"`
	}

	// Unmarshal the JSON data into the assetData structure
	if err := json.Unmarshal(body, &incomeData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Delete associated UserIncomeTaxes first to avoid foreign key violations
	var userIncomeTaxes []models.UserTax
	if err := db.DB.Where("user_income_id = ?", incomeData.IncomeId).Delete(&userIncomeTaxes).Error; err != nil {
		log.Printf("Error deleting associated taxes: %v", err)
		http.Error(w, "Error deleting associated taxes", http.StatusInternalServerError)
		return
	}

	// Now delete the UserIncome
	var income models.UserIncome
	result := db.DB.Delete(&income, incomeData.IncomeId)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Income and associated taxes deleted successfully!",
	})
}
