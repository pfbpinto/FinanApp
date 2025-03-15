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

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateUserAsset creates a new UserAsset
func CreateUserAsset(w http.ResponseWriter, r *http.Request) {
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
	var assetData struct {
		AssetName           string `json:"AssetName"`
		AssetValue          string `json:"AssetValue"`
		AssetType           string `json:"AssetTypeID"`
		AssetAquisitionDate string `json:"AssetAquisitionDate"`
		Shared              *bool  `json:"SharedAsset"`
		UserID              uint   `json:"userID"`
		UserAssetTaxes      []struct {
			ID            uint    `json:"TaxID"`
			TaxName       string  `json:"TaxName"`
			TaxPercentage float64 `json:"TaxPercentage"`
		} `json:"UserAssetTaxes"`
	}

	// Unmarshal the JSON into assetData
	if err := json.Unmarshal(body, &assetData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Tax Captured: %v", assetData.UserAssetTaxes)

	// Validation of required fields
	if assetData.AssetName == "" || assetData.AssetValue == "" || assetData.AssetType == "" || assetData.UserID < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	if assetData.AssetAquisitionDate == "" {
		http.Error(w, "Invalid Date format", http.StatusBadRequest)
		return
	} else {
		// Parse the date
		_, err := time.Parse("2006-01-02", assetData.AssetAquisitionDate)
		if err != nil {
			http.Error(w, "Invalid Date format", http.StatusBadRequest)
			return
		}
	}

	// Convert AssetValue from string to float64
	assetValue, err := strconv.ParseFloat(assetData.AssetValue, 64)
	if err != nil {
		log.Printf("Error converting AssetValue: %v", err)
		http.Error(w, "Error converting AssetValue", http.StatusBadRequest)
		return
	}

	// Convert AssetType from string to uint
	assetTypeID, err := strconv.ParseUint(assetData.AssetType, 10, 32)
	if err != nil {
		log.Printf("Error converting AssetType: %v", err)
		http.Error(w, "Error converting AssetType", http.StatusBadRequest)
		return
	}

	// Parse the date
	parsedDate, err := time.Parse("2006-01-02", assetData.AssetAquisitionDate)
	if err != nil {
		http.Error(w, "Error parsing Aquisition date", http.StatusInternalServerError)
		return
	}

	// Create the UserAsset model
	userAsset := models.UserAsset{
		UserID:              assetData.UserID,
		AssetTypeID:         uint(assetTypeID), // Convert uint64 to uint
		AssetName:           assetData.AssetName,
		AssetAquisitionDate: &parsedDate,
		AssetValue:          assetValue,
		SharedAsset:         assetData.Shared != nil && *assetData.Shared || false,
	}

	// Save to the database using GORM
	if err := db.DB.Create(&userAsset).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// Check if there are any taxes to process
	if len(assetData.UserAssetTaxes) > 0 {
		log.Printf("Processing taxes for UserAsset ID %d", userAsset.ID)

		// Iterate through the UserAssetTaxes array which contains full tax details
		for _, taxData := range assetData.UserAssetTaxes {
			log.Printf("Applying tax: %s (ID: %d, Percentage: %.2f%%)", taxData.TaxName, taxData.ID, taxData.TaxPercentage)

			// Calculate the tax value based on the AssetValue
			taxValue := assetValue * (taxData.TaxPercentage / 100)
			log.Printf("Calculated TaxValue: %.2f based on AssetValue: %.2f", taxValue, assetValue)

			// Create the record in the UserAssetTax table
			userAssetTax := models.UserAssetTax{
				UserAssetID: userAsset.ID, // ID of the newly created UserAsset
				TaxID:       taxData.ID,   // ID of the Tax
				TaxValue:    taxValue,
			}

			log.Printf("Creating UserAssetTax record: %+v", userAssetTax)

			// Save the UserAssetTax record in the database
			if err := db.DB.Create(&userAssetTax).Error; err != nil {
				log.Printf("Error saving UserAssetTax record: %v", err)
				continue
			}

			log.Printf("Successfully saved UserAssetTax record for TaxID %d", taxData.ID)
		}
	} else {
		log.Println("No taxes to process for UserAsset")
	}
	// Log successful creation
	log.Printf("UserAsset created: %+v", userAsset)

	// Send response back with created asset data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Asset created successfully!",
		"userAsset": userAsset,
	})
}

// GetAssetType returns a list of UserAsset in the database.
func GetAssetType(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Create a variable to store asset types
	var assetTypes []models.AssetType
	var taxes []models.Tax

	// Fetch all asset types from the database
	if result := db.DB.Find(&assetTypes); result.Error != nil {
		http.Error(w, "Failed to retrieve asset types data", http.StatusInternalServerError)
		return
	}

	// Fetch all tax from type 1: asset from the database
	if result := db.DB.Where("tax_type_id = ? AND user_id = ?", 1, user.ID).Find(&taxes); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve Tax data", http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"assetTypes": assetTypes, // Returning the list of assetTypes
		"taxes":      taxes,      // Returning the list of taxex
	}

	// Set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}
}

// GetAsset returns an existing group of UserAsset in the database.
func GetAssets(w http.ResponseWriter, r *http.Request) {

	var assets []models.UserAsset
	result := db.DB.Preload("User").Preload("AssetType").Find(&assets)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assets)
}

// GetAsset returns an existing UserAsset in the database.
func GetAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	var asset models.UserAsset
	result := db.DB.Preload("User").Preload("AssetType").First(&asset, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(asset)
}

// UpdateUserAsset updates an existing UserAsset in the database.
func UpdateUserAsset(w http.ResponseWriter, r *http.Request) {
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
	var assetData struct {
		AssetName           string  `json:"AssetName"`
		AssetValue          float64 `json:"AssetValue"`
		AssetTypeID         uint    `json:"AssetTypeID"`
		AssetAquisitionDate string  `json:"AssetAquisitionDate"`
		Shared              *bool   `json:"SharedAsset"`
		UserID              uint    `json:"userID"`
		ID                  uint    `json:"ID"`
		UserAssetTaxes      []struct {
			ID       uint    `json:"ID"`
			TaxID    uint    `json:"TaxID"`
			TaxValue float64 `json:"TaxValue"`
		} `json:"UserAssetTaxes"`
	}

	// Unmarshal the JSON data into the assetData structure
	if err := json.Unmarshal(body, &assetData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validate that required fields are filled
	if assetData.AssetName == "" || assetData.AssetValue < 1 || assetData.AssetTypeID == 0 || assetData.UserID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	// Validate the asset acquisition date format
	if assetData.AssetAquisitionDate == "" {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Try parsing the first format (with timezone)
	parsedDate, err := time.Parse("2006-01-02T15:04:05-07:00", assetData.AssetAquisitionDate)
	if err != nil {
		// Se o primeiro formato falhar, tente analisar o formato de data simples
		parsedDate, err = time.Parse("2006-01-02", assetData.AssetAquisitionDate)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
	}

	// Find the existing UserAsset by its ID in the database
	var userAsset models.UserAsset
	if err := db.DB.First(&userAsset, assetData.ID).Error; err != nil {
		log.Printf("Error finding asset with ID %d: %v", assetData.ID, err)
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	// Update the existing UserAsset with the new data
	userAsset.AssetName = assetData.AssetName
	userAsset.AssetValue = assetData.AssetValue
	userAsset.AssetTypeID = assetData.AssetTypeID
	userAsset.AssetAquisitionDate = &parsedDate
	userAsset.SharedAsset = assetData.Shared != nil && *assetData.Shared || false

	// Save the updated UserAsset back to the database
	if err := db.DB.Save(&userAsset).Error; err != nil {
		log.Printf("Error updating asset in database: %v", err)
		http.Error(w, "Error updating asset", http.StatusInternalServerError)
		return
	}

	// Clear existing UserAssetTaxes for the asset (if any)
	if err := db.DB.Where("user_asset_id = ?", userAsset.ID).Delete(&models.UserAssetTax{}).Error; err != nil {
		log.Printf("Error clearing old taxes: %v", err)
		http.Error(w, "Error clearing old taxes", http.StatusInternalServerError)
		return
	}

	// Add new UserAssetTaxes
	for _, tax := range assetData.UserAssetTaxes {
		newTax := models.UserAssetTax{
			UserAssetID: userAsset.ID,
			TaxID:       tax.TaxID,
			TaxValue:    tax.TaxValue,
		}

		// Save the new UserAssetTax entry
		if err := db.DB.Create(&newTax).Error; err != nil {
			log.Printf("Error adding tax: %v", err)
			http.Error(w, "Error adding tax", http.StatusInternalServerError)
			return
		}
	}

	// Log successful update
	log.Printf("UserAsset updated: %+v", userAsset)

	// Respond with a success message and the updated asset data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Asset updated successfully!",
		"userAsset": userAsset,
	})
}

// DeleteAsset deletes an existing UserAsset in the database.
func DeleteAsset(w http.ResponseWriter, r *http.Request) {

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
	var assetData struct {
		AssetId uint `json:"itemId"`
	}

	// Unmarshal the JSON data into the assetData structure
	if err := json.Unmarshal(body, &assetData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Delete associated UserAssetTaxes first to avoid foreign key violations
	var userAssetTaxes []models.UserAssetTax
	if err := db.DB.Where("user_asset_id = ?", assetData.AssetId).Delete(&userAssetTaxes).Error; err != nil {
		log.Printf("Error deleting associated taxes: %v", err)
		http.Error(w, "Error deleting associated taxes", http.StatusInternalServerError)
		return
	}

	// Now delete the UserAsset
	var asset models.UserAsset
	result := db.DB.Delete(&asset, assetData.AssetId)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Asset and associated taxes deleted successfully!",
	})
}
