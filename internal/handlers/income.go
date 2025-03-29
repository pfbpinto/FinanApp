package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

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

func IncomeItem(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("UserDashboard: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Decode request payload into struct
	var payload struct {
		ItemID int `json:"itemId"`
	}

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get Database
	database := db.GetDB()

	// Consulta UserFinancialForecasts
	userFinancialForecastQuery := `
	SELECT 
		uff.UserFinancialForecastID, 
		uff.UserCategoryID, 
		uff.FinancialUserItemID, 
		uff.UserFinancialForecastAmount, 
		uff.UserFinancialForecastBeginDate, 
		uff.UserFinancialForecastEndDate, 
		uff.CurrencyID,
		uc.UserCategoryName, -- Exemplo de relacionamento com a tabela de categoria
		fui.FinancialUserItemName, -- Exemplo de relacionamento com a tabela de item financeiro
		c.CurrencyName -- Exemplo de relacionamento com a tabela de moedas
	FROM 
		userfinancialforecast uff
	LEFT JOIN 
		usercategory uc ON uff.UserCategoryID = uc.UserCategoryID
	JOIN 
		financialuseritem fui ON uff.FinancialUserItemID = fui.FinancialUserItemID
	JOIN 
		currency c ON uff.CurrencyID = c.CurrencyID
	WHERE
		fui.userentityid = $1 and fui.EntityID = 5 and fui.FinancialUserItemID = $2
	ORDER BY 
		uff.UserFinancialForecastBeginDate;
	`
	userFinancialForecastRows, err := database.Query(userFinancialForecastQuery, user.UserProfileID, payload.ItemID)
	if err != nil {
		log.Println("Erro ao buscar UserFinancialForecasts:", err)
		http.Error(w, "Erro ao buscar UserFinancialForecasts", http.StatusInternalServerError)
		return
	}
	defer userFinancialForecastRows.Close()

	var userFinancialForecasts []models.UserFinancialForecast

	for userFinancialForecastRows.Next() {
		var uff models.UserFinancialForecast
		// Atualizando o Scan para incluir os novos campos
		if err := userFinancialForecastRows.Scan(
			&uff.UserFinancialForecastID,
			&uff.UserCategoryID,
			&uff.FinancialUserItemID,
			&uff.UserFinancialForecastAmount,
			&uff.UserFinancialForecastBeginDate,
			&uff.UserFinancialForecastEndDate,
			&uff.CurrencyID,
			&uff.UserCategoryName,
			&uff.FinancialUserItemName,
			&uff.CurrencyName,
		); err != nil {
			log.Println("Erro ao escanear UserFinancialForecast:", err)
			continue
		}
		userFinancialForecasts = append(userFinancialForecasts, uff)
	}

	// Consulta UserFinancialActuals
	userFinancialActualQuery := `
	SELECT 
		ufa.UserFinancialActualID, 
		ufa.UserCategoryID, 
		ufa.FinancialUserItemID, 
		ufa.UserFinancialActualAmount, 
		ufa.UserFinancialActualtBeginDate, 
		ufa.UserFinancialActualEndDate, 
		ufa.CurrencyID,
		uc.UserCategoryName, -- Relacionamento com tabela de categorias
		fui.FinancialUserItemName, -- Relacionamento com tabela de itens financeiros
		c.CurrencyName -- Relacionamento com tabela de moedas
	FROM 
		userfinancialactual ufa
	JOIN 
		usercategory uc ON ufa.UserCategoryID = uc.UserCategoryID
	JOIN 
		financialuseritem fui ON ufa.FinancialUserItemID = fui.FinancialUserItemID
	JOIN 
		currency c ON ufa.CurrencyID = c.CurrencyID
	WHERE
		fui.userentityid = $1  and fui.EntityID = 5 and fui.FinancialUserItemID = $2   
	ORDER BY 
		ufa.UserFinancialActualtBeginDate;`

	userFinancialActualRows, err := database.Query(userFinancialActualQuery, user.UserProfileID, payload.ItemID)
	if err != nil {
		log.Println("Erro ao buscar UserFinancialActuals:", err)
		http.Error(w, "Erro ao buscar UserFinancialActuals", http.StatusInternalServerError)
		return
	}
	defer userFinancialActualRows.Close()

	var userFinancialActuals []models.UserFinancialActual
	for userFinancialActualRows.Next() {
		var ufa models.UserFinancialActual
		if err := userFinancialActualRows.Scan(
			&ufa.UserFinancialActualID,
			&ufa.UserCategoryID,
			&ufa.FinancialUserItemID,
			&ufa.UserFinancialActualAmount,
			&ufa.UserFinancialActualtBeginDate,
			&ufa.UserFinancialActualEndDate,
			&ufa.CurrencyID,
			&ufa.UserCategoryName,
			&ufa.FinancialUserItemName,
			&ufa.CurrencyName,
		); err != nil {
			log.Println("Erro ao escanear UserFinancialActual:", err)
			continue
		}
		userFinancialActuals = append(userFinancialActuals, ufa)
	}

	// Criar resposta final
	response := struct {
		UserFinancialForecasts []models.UserFinancialForecast `json:"user_financial_forecasts"`
		UserFinancialActuals   []models.UserFinancialActual   `json:"user_financial_actuals"`
	}{

		UserFinancialForecasts: userFinancialForecasts,
		UserFinancialActuals:   userFinancialActuals,
	}

	// Configurar cabeçalhos da resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Converter para JSON e enviar
	json.NewEncoder(w).Encode(response)

}

func CreateIncome(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateIncome: Unauthorized access - no user found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get database connection
	database := db.GetDB()

	// Decode request payload into struct
	var payload models.FinancialUserItem
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(payload.Amount, 64)
	if err != nil {
		log.Println("Invalid amount format:", err)
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	recurrencyID, err := strconv.Atoi(payload.RecurrencyID)
	if err != nil {
		log.Println("Invalid recurrency ID format:", err)
		http.Error(w, "Invalid recurrency ID format", http.StatusBadRequest)
		return
	}

	// Validate and parse BeginDate (if applicable)
	beginDate := time.Now()

	// Call the stored procedure
	var message string
	err = database.QueryRow(
		"CALL CreateUserParentIncome($1, $2, $3, $4, $5, $6, $7)",
		user.UserProfileID,
		payload.FinancialUserItemName,
		recurrencyID,
		5, // Hardcoded FinancialUserEntityItemID
		amount,
		beginDate,
		message).Scan(&message)

	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	fmt.Println("Mensagem da stored procedure:", message)

	// Parse JSON response from stored procedure
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing JSON response:", err)
		http.Error(w, "Invalid response format from stored procedure", http.StatusInternalServerError)
		return
	}

	// Send response based on status
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else if resp.Status == "success" {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Unknown response status:", resp.Status)
		w.WriteHeader(http.StatusInternalServerError)
		resp = models.Response{Status: "error", Message: "Unknown response status"}
	}

	json.NewEncoder(w).Encode(resp)
}

func UpdateIncome(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request payload into struct
	var payload struct {
		ItemID                int    `json:"FinancialUserItemId"`   // ID of the item to be updated
		FinancialUserItemName string `json:"FinancialUserItemName"` // Name of the item to be updated
	}

	// Parse the JSON body into the payload struct
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the required parameters are present
	if payload.ItemID == 0 || payload.FinancialUserItemName == "" {
		http.Error(w, "ItemID and FinancialUserItemName are required", http.StatusBadRequest)
		return
	}

	// Get database connection
	database := db.GetDB()

	// Prepare the SQL update statement
	query := `UPDATE FinancialUserItem SET FinancialUserItemName = $1 WHERE FinancialUserItemId = $2`

	// Execute the update statement
	_, err = database.Exec(query, payload.FinancialUserItemName, payload.ItemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating income: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Income updated successfully"})
}

func DeleteIncome(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request payload into struct
	var payload struct {
		ItemID int `json:"itemId"`
	}

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get database connection
	database := db.GetDB()

	// Begin a transaction
	tx, err := database.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Ensure rollback on error

	// Delete related records from 'userfinancialactual' table
	_, err = tx.Exec("DELETE FROM userfinancialactual WHERE FinancialUserItemID = $1", payload.ItemID)
	if err != nil {
		log.Println("Error deleting from userfinancialactual:", err)
		http.Error(w, "Failed to delete related records from userfinancialactual", http.StatusInternalServerError)
		return
	}

	// Delete related records from 'userfinancialforecast' table
	_, err = tx.Exec("DELETE FROM userfinancialforecast WHERE FinancialUserItemID = $1", payload.ItemID)
	if err != nil {
		log.Println("Error deleting from userfinancialforecast:", err)
		http.Error(w, "Failed to delete related records from userfinancialforecast", http.StatusInternalServerError)
		return
	}

	// Finally, delete the record from 'financialuseritem' table
	_, err = tx.Exec("DELETE FROM financialuseritem WHERE FinancialUserItemID = $1", payload.ItemID)
	if err != nil {
		log.Println("Error deleting from financialuseritem:", err)
		http.Error(w, "Failed to delete record from financialuseritem", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Income deleted successfully"})
}

func CreateIncomeCategory(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateIncomeCategory: Unauthorized access - no user found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode request payload into struct
	var payload struct {
		UserCategoryName string `json:"user_category_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Define hardcoded values
	entityID := 5
	financialGroupEntityItemID := 1
	isActive := true

	// Get database connection
	database := db.GetDB()

	// Insert into the database
	query := `INSERT INTO UserCategory (UserCategoryName, UserProfileID, EntityID, FinancialGroupEntityItemID, IsActive)
			  VALUES ($1, $2, $3, $4, $5) RETURNING UserCategoryID`
	var userCategoryID int
	err := database.QueryRow(query, payload.UserCategoryName, user.UserProfileID, entityID, financialGroupEntityItemID, isActive).Scan(&userCategoryID)
	if err != nil {
		log.Println("Error inserting new category:", err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	// Send response
	response := map[string]interface{}{
		"message":          "Category created successfully",
		"user_category_id": userCategoryID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func DeleteIncomeCategory(w http.ResponseWriter, r *http.Request) {
	// Checks if the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieves the user from the context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("DeleteIncomeCategory: Unauthorized access - no user found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decodes the JSON from the request body
	var payload struct {
		UserCategoryID int `json:"user_category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Checks if the ID was provided
	if payload.UserCategoryID == 0 {
		http.Error(w, "UserCategoryID is required", http.StatusBadRequest)
		return
	}

	// Connects to the database
	database := db.GetDB()

	// Executes the DELETE in the UserCategory table, ensuring it belongs to the authenticated user
	result, err := database.Exec(`
        DELETE FROM usercategory 
        WHERE UserCategoryID = $1 AND UserProfileID = $2`,
		payload.UserCategoryID, user.UserProfileID)

	if err != nil {
		log.Println("Error deleting category:", err)
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Category not found or unauthorized", http.StatusNotFound)
		return
	}

	// Responds with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Category deleted successfully"})
}
