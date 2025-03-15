package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CreateUserExpense creates a new UserExpense SharedExpenditure
func CreateUserExpense(w http.ResponseWriter, r *http.Request) {
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
	var expenseData struct {
		ExpenditureName       string `json:"ExpenditureName"`
		ExpenditureValue      string `json:"ExpenditureValue"`
		ExpenditureType       string `json:"Expenditure_ID"`
		ExpenditureStartDate  string `json:"ExpenditureStartDate"`
		ExpenditureRecurrence string `json:"ExpenditureRecurrence"`
		Shared                *bool  `json:"SharedExpenditure"`
		UserID                uint   `json:"userID"`
	}

	// Unmarshal the JSON into expenseData
	if err := json.Unmarshal(body, &expenseData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validation of required fields
	if expenseData.ExpenditureName == "" || expenseData.ExpenditureValue == "" || expenseData.ExpenditureType == "" || expenseData.UserID < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	if expenseData.ExpenditureStartDate == "" {
		http.Error(w, "Invalid Date format", http.StatusBadRequest)
		return
	} else {
		// Parse the date
		_, err := time.Parse("2006-01-02", expenseData.ExpenditureStartDate)
		if err != nil {
			http.Error(w, "Invalid Date format", http.StatusBadRequest)
			return
		}
	}

	// Convert ExpenditureValue from string to float64
	expenseValue, err := strconv.ParseFloat(expenseData.ExpenditureValue, 64)
	if err != nil {
		log.Printf("Error converting ExpenseValue: %v", err)
		http.Error(w, "Error converting ExpenseValue", http.StatusBadRequest)
		return
	}

	// Convert ExpenseType from string to uint
	expenseTypeID, err := strconv.ParseUint(expenseData.ExpenditureType, 10, 32)
	if err != nil {
		log.Printf("Error converting ExpenseType: %v", err)
		http.Error(w, "Error converting ExpenseType", http.StatusBadRequest)
		return
	}

	// Parse the date
	parsedDate, err := time.Parse("2006-01-02", expenseData.ExpenditureStartDate)
	if err != nil {
		http.Error(w, "Error parsing Expenditure StartDate", http.StatusInternalServerError)
		return
	}

	// Create the UserExpense model
	userExpense := models.UserExpenditure{
		UserID:                expenseData.UserID,
		ExpenditureID:         uint(expenseTypeID), // Convert uint64 to uint
		ExpenditureName:       expenseData.ExpenditureName,
		ExpenditureStartDate:  &parsedDate,
		ExpenditureValue:      expenseValue,
		SharedExpenditure:     expenseData.Shared != nil && *expenseData.Shared || false,
		ExpenditureRecurrence: expenseData.ExpenditureRecurrence,
	}

	// Save to the database using GORM
	if err := db.DB.Create(&userExpense).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// Send response back with created Expense data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Expense created successfully!",
		"userExpense": userExpense,
	})
}

// GetAssetType returns a list of UserExpense in the database.
func GetExpenseType(w http.ResponseWriter, r *http.Request) {

	// Create a variable to store Expense types
	var expenseTypes []models.ExpenditureType

	// Fetch all Expense types from the database
	if result := db.DB.Find(&expenseTypes); result.Error != nil {
		http.Error(w, "Failed to retrieve expense types data", http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"expenseTypes": expenseTypes, // Returning the list of expenseTypes
	}

	// Set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}
}

// GetAsset returns an existing group of UserExpense in the database.
func GetExpense(w http.ResponseWriter, r *http.Request) {

	var expense []models.UserExpenditure
	result := db.DB.Preload("User").Preload("Expenditure").Find(&expense)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expense)
}

// UpdateUserAsset updates an existing UserExpense in the database.
func UpdateUserExpense(w http.ResponseWriter, r *http.Request) {
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
	var expenseData struct {
		ExpenditureName       string  `json:"ExpenditureName"`
		ExpenditureValue      float64 `json:"ExpenditureValue"`
		ExpenditureType       uint    `json:"ExpenditureID"`
		ExpenditureStartDate  string  `json:"ExpenditureStartDate"`
		ExpenditureRecurrence string  `json:"ExpenditureRecurrence"`
		Shared                *bool   `json:"SharedExpenditure"`
		UserID                uint    `json:"userID"`
		ID                    uint    `json:"ID"`
	}

	// Unmarshal the JSON data into the expenseData structure
	if err := json.Unmarshal(body, &expenseData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validate that required fields are filled
	if expenseData.ExpenditureName == "" || expenseData.ExpenditureValue < 1 || expenseData.ExpenditureType < 1 || expenseData.UserID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	// Validate the expense acquisition date format
	if expenseData.ExpenditureStartDate == "" {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Try parsing the first format (with timezone)
	parsedDate, err := time.Parse("2006-01-02T15:04:05-07:00", expenseData.ExpenditureStartDate)
	if err != nil {
		// Se o primeiro formato falhar, tente analisar o formato de data simples
		parsedDate, err = time.Parse("2006-01-02", expenseData.ExpenditureStartDate)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
	}

	// Find the existing UserExpense by its ID in the database
	var userExpense models.UserExpenditure
	if err := db.DB.First(&userExpense, expenseData.ID).Error; err != nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	// Update the existing UserExpense with the new data
	userExpense.ExpenditureName = expenseData.ExpenditureName
	userExpense.ExpenditureValue = expenseData.ExpenditureValue
	userExpense.ExpenditureID = expenseData.ExpenditureType
	userExpense.ExpenditureStartDate = &parsedDate
	userExpense.ExpenditureRecurrence = expenseData.ExpenditureRecurrence
	userExpense.SharedExpenditure = expenseData.Shared != nil && *expenseData.Shared || false

	log.Printf("Error converting ExpenseType: %v", userExpense)

	// Save the updated UserAsset back to the database
	if err := db.DB.Save(&userExpense).Error; err != nil {
		log.Printf("Error updating expense in database: %v", err)
		http.Error(w, "Error updating expense", http.StatusInternalServerError)
		return
	}

	// Respond with a success message and the updated expense data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Expense updated successfully!",
		"userExpense": userExpense,
	})
}

// DeleteAsset deletes an existing UserExpense in the database.
func DeleteExpense(w http.ResponseWriter, r *http.Request) {

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
	var expenseData struct {
		ExpenseId uint `json:"itemId"`
	}

	// Unmarshal the JSON data into the expenseData structure
	if err := json.Unmarshal(body, &expenseData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Delete the UserExpense
	var expense models.UserExpenditure
	result := db.DB.Delete(&expense, expenseData.ExpenseId)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Expense deleted successfully!",
	})
}
