package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"
)

// GetTaxes API returns a list of Tax in the database.
func GetTaxes(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Create a variable to store Tax
	var taxType []models.TaxType
	var tax []models.Tax

	// Fetch all Tax types from the database
	if result := db.DB.Find(&taxType); result.Error != nil {
		http.Error(w, "Failed to retrieve asset types data", http.StatusInternalServerError)
		return
	}

	// Fetch all taxes from the database
	if result := db.DB.Where("user_id = ?", user.ID).Find(&tax); result.Error != nil {
		http.Error(w, "Failed to retrieve asset types data", http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"taxTypes": taxType,
		"taxes":    tax,
	}

	// Set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}
}

// CreateUserAsset creates a new Tax
func CreateTax(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		//log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the received body for debugging
	log.Printf("Received body: %s", string(body))

	// Define the structure to match the incoming JSON
	var taxData struct {
		UserID             uint   `json:"UserID"`
		TaxTypeID          string `json:"TaxTypeID"`
		TaxName            string `json:"TaxName"`
		TaxPercentage      string `json:"TaxPercentage"`
		TaxPercentageRange string `json:"TaxPercentageRange"`
		TaxApplicableCycle string `json:"TaxApplicableCycle"`
	}

	// Unmarshal the JSON into assetData
	if err := json.Unmarshal(body, &taxData); err != nil {
		//log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validation of required fields
	if taxData.TaxName == "" || taxData.TaxPercentage == "" || taxData.TaxPercentageRange == "" || taxData.TaxTypeID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	// Convert TaxPercentage from string to float64
	taxPercentage, err := strconv.ParseFloat(taxData.TaxPercentage, 64)
	if err != nil {
		log.Printf("Error converting TaxPercentage: %v", err)
		http.Error(w, "Error converting TaxPercentage", http.StatusBadRequest)
		return
	}

	// Convert TaxTypeID from string to uint
	taxTypeID, err := strconv.ParseUint(taxData.TaxTypeID, 10, 32)
	if err != nil {
		log.Printf("Error converting TaxTypeID: %v", err)
		http.Error(w, "Error converting TaxTypeID", http.StatusBadRequest)
		return
	}

	// Create the UserAsset model
	taxes := models.Tax{
		UserID:             uint(taxData.UserID),
		TaxName:            taxData.TaxName,
		TaxTypeID:          uint(taxTypeID),
		TaxPercentage:      taxPercentage,
		TaxPercentageRange: taxData.TaxPercentageRange,
		TaxApplicableCycle: taxData.TaxApplicableCycle,
	}

	// Save to the database using GORM
	if err := db.DB.Create(&taxes).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// Send response back with created Tax data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tax created successfully!",
		"taxes":   taxes,
	})

}

// DeleteTaxes deletes an existing Tax in the database.
func DeleteTaxes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	// Define the structure to match the incoming JSON
	var requestData struct {
		TaxID uint `json:"taxID"`
	}

	// Unmarshal the JSON into requestData
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validate if TaxID was provided
	if requestData.TaxID < 1 {
		http.Error(w, "Missing tax ID in payload", http.StatusBadRequest)
		return
	}

	// Delete from the database using GORM
	result := db.DB.Delete(&models.Tax{}, requestData.TaxID)
	if result.Error != nil {
		log.Printf("Error deleting tax: %v", result.Error)
		http.Error(w, "Error deleting tax", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tax deleted successfully!"})
}
