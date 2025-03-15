package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// GetTaxes returns a list of Tax in the database.
func GetCategories(w http.ResponseWriter, r *http.Request) {

	// Create a variable to store Tax
	var taxType []models.TaxType
	var assetType []models.AssetType
	var incomeType []models.IncomeType
	var expenditureType []models.ExpenditureType
	var groupType []models.GroupType
	var fileType []models.FileType

	// Fetch all Tax types from the database
	if result := db.DB.Find(&taxType); result.Error != nil {
		http.Error(w, "Failed to retrieve tax types data", http.StatusInternalServerError)
		return
	}
	// Fetch all asset types from the database
	if result := db.DB.Find(&assetType); result.Error != nil {
		http.Error(w, "Failed to retrieve asset types data", http.StatusInternalServerError)
		return
	}
	// Fetch all Income types from the database
	if result := db.DB.Find(&incomeType); result.Error != nil {
		http.Error(w, "Failed to retrieve income types data", http.StatusInternalServerError)
		return
	}
	// Fetch all Expenditure types from the database
	if result := db.DB.Find(&expenditureType); result.Error != nil {
		http.Error(w, "Failed to retrieve expenditure types data", http.StatusInternalServerError)
		return
	}
	// Fetch all Group types from the database
	if result := db.DB.Find(&groupType); result.Error != nil {
		http.Error(w, "Failed to retrieve group types data", http.StatusInternalServerError)
		return
	}
	// Fetch all File types from the database
	if result := db.DB.Find(&fileType); result.Error != nil {
		http.Error(w, "Failed to retrieve File types data", http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"assetType":       assetType,
		"taxTypes":        taxType,
		"incomeType":      incomeType,
		"expenditureType": expenditureType,
		"groupType":       groupType,
		"fileType":        fileType,
	}

	// Set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}
}

// Category Models
var databaseModels = map[string]interface{}{
	"assetType":       &models.AssetType{},
	"taxTypes":        &models.TaxType{},
	"expenditureType": &models.ExpenditureType{},
	"fileType":        &models.FileType{},
	"groupType":       &models.GroupType{},
	"incomeType":      &models.IncomeType{},
}

func CreateCategories(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for create-category")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categoryData struct {
		Category string `json:"category"`
		Field    string `json:"field"`
		Model    string `json:"model"`
		Name     string `json:"name"`
	}

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&categoryData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received category data: Category = %s, Field = %s, Model = %s, Name = %s",
		categoryData.Category, categoryData.Field, categoryData.Model, categoryData.Name)

	// Check if the model exists in the databaseModels map
	model, exists := databaseModels[categoryData.Model]
	if !exists {
		http.Error(w, "Model not found", http.StatusBadRequest)
		return
	}

	// Create a new instance of the model dynamically
	record := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	// Set the field dynamically using reflection
	v := reflect.ValueOf(record).Elem()
	f := v.FieldByName(categoryData.Field)

	if !f.IsValid() {
		http.Error(w, "Invalid field name", http.StatusBadRequest)
		return
	}

	// Assign the name to the corresponding field
	f.SetString(categoryData.Name)

	// Save the record to the database
	if err := db.DB.Create(record).Error; err != nil {

		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "Error 1062") {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // HTTP 409 - Conflict
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Category name already exists",
				"message": "A Category with this name already exists. Please choose a different name.",
			})
			return
		}

		log.Printf("Error saving to database: %v", err)
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// Send response with the created category
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Category created successfully!",
		"category": record, // Return the newly inserted category
	})
}

func DeleteCategories(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for delete-category")

	// Ensure the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define the structure to parse the incoming JSON payload
	var categoryData struct {
		CategoryID    uint   `json:"categoryID"`
		CategoryModel string `json:"model"`
	}

	// Parse request body into categoryData
	err := json.NewDecoder(r.Body).Decode(&categoryData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the provided model exists in the databaseModels map
	model, exists := databaseModels[categoryData.CategoryModel]
	if !exists {
		http.Error(w, "Model not found", http.StatusBadRequest)
		return
	}

	// Create a new instance of the correct model dynamically
	record := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	// Attempt to delete the category by ID
	if err := db.DB.Where("id = ?", categoryData.CategoryID).Delete(record).Error; err != nil {
		log.Printf("Error deleting category: %v", err)
		http.Error(w, "Error deleting category", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Category deleted successfully!",
	})
}
