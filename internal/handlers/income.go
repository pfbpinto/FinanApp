package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"log"
	"net/http"
)

type IncomeResponse struct {
	IncomeTypes        []models.IncomeType            `json:"income_types"`
	UserCategories     []models.UserCategory          `json:"user_categories"`
	FinancialForecasts []models.UserFinancialForecast `json:"user_finance_forecast"`
}

// UserIncomeForecast shows the logged-in user's Forecast Incomes
func UserIncomeForecast(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	// Retrieve the user from the request context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("UserDashboard: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	log.Println("User from Context : ", user)

	database := db.GetDB()

	// Starting collection for Income Forecast API
	// Starting to collect income forecast data
	incomeQuery := `
		SELECT IncomeTypeID, IncomeTypeName, IncomeDescription, CreatedAt
		FROM IncomeType
	`
	incomeRows, err := database.Query(incomeQuery)
	if err != nil {
		log.Println("Error fetching income types: ", err)
		http.Error(w, "Error fetching income types", http.StatusInternalServerError)
		return
	}
	defer incomeRows.Close()

	var incomeTypes []models.IncomeType
	for incomeRows.Next() {
		var income models.IncomeType
		if err := incomeRows.Scan(&income.IncomeTypeID, &income.IncomeTypeName, &income.IncomeDescription, &income.CreatedAt); err != nil {
			log.Println("Error processing income types: ", err)
			http.Error(w, "Error processing income types", http.StatusInternalServerError)
			return
		}
		incomeTypes = append(incomeTypes, income)
	}

	// Search for category with type "Income"
	// Querying for categories with type "Income"
	userCategoryQuery := `
    SELECT "usercategoryid", "usercategoryname", "itemtypename", "itemtypenameid", "createdat"
    FROM "usercategory"
    WHERE "itemtypename" = 'Income' AND "userprofileid" = $1`

	log.Println("Testing User ID: ", user.UserProfileID)
	userCategoryRows, err := database.Query(userCategoryQuery, user.UserProfileID)
	if err != nil {
		log.Println("Error fetching user categories: ", err)
		http.Error(w, "Error fetching user categories", http.StatusInternalServerError)
		return
	}
	defer userCategoryRows.Close()

	var userCategories []models.UserCategory
	for userCategoryRows.Next() {
		var category models.UserCategory
		if err := userCategoryRows.Scan(&category.UserCategoryID, &category.UserCategoryName, &category.ItemTypeName, &category.ItemTypeNameID, &category.CreatedAt); err != nil {
			log.Println("Error processing user categories: ", err)
			http.Error(w, "Error processing user categories", http.StatusInternalServerError)
			return
		}
		userCategories = append(userCategories, category)
	}

	// Fetch user financial forecast data
	// Query to fetch user financial forecast data
	userFinancialForecastQuery := `
    SELECT 
        "userfinancialforecastid", "userfinancialforecastname", "usercategoryid", 
        "entitytypename", "entitytypeid", "entityitemtypename", "entityitemtypeid", 
        "userfinancialforecastamount", "userfinancialforecastbegindate", 
        "userfinancialforecastenddate", "currencyid", "createdat"
    FROM "userfinancialforecast"
    WHERE "entitytypename" = 'User Income' AND "entitytypeid" = $1`

	log.Println("Executing query to fetch financial forecasts for user:", user.UserProfileID)
	rows, err := database.Query(userFinancialForecastQuery, user.UserProfileID)
	if err != nil {
		log.Println("Error fetching financial forecasts: ", err)
		http.Error(w, "Error fetching financial forecasts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Process results for user financial forecasts
	// Processing the results for the financial forecasts
	var financialForecasts []models.UserFinancialForecast
	for rows.Next() {
		var forecast models.UserFinancialForecast
		if err := rows.Scan(
			&forecast.UserFinancialForecastID, &forecast.UserFinancialForecastName, &forecast.UserCategoryID,
			&forecast.EntityTypeName, &forecast.EntityTypeID, &forecast.EntityItemTypeName, &forecast.EntityItemTypeID,
			&forecast.UserFinancialForecastAmount, &forecast.UserFinancialForecastBeginDate,
			&forecast.UserFinancialForecastEndDate, &forecast.CurrencyID, &forecast.CreatedAt,
		); err != nil {
			log.Println("Error processing financial forecast data: ", err)
			http.Error(w, "Error processing financial forecast data", http.StatusInternalServerError)
			return
		}
		financialForecasts = append(financialForecasts, forecast)
	}

	// Create the response object to send as JSON
	// Preparing the response data for JSON output
	response := IncomeResponse{
		IncomeTypes:        incomeTypes,
		UserCategories:     userCategories,
		FinancialForecasts: financialForecasts,
	}

	// Set Content-Type header and send the response as JSON
	// Sending the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
