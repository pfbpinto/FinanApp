package handlers

import (
	"database/sql"
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO
func AssetItem(w http.ResponseWriter, r *http.Request) {

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

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateAsset: Unauthorized access - no user found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get database connection
	database := db.GetDB()

	// Decode request payload into struct
	var payload models.UserAsset
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(payload.UserAssetValueAmount) == "" {
		log.Println("Amount is empty")
		http.Error(w, "Amount cannot be empty", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(payload.UserAssetValueAmount, 64)
	if err != nil {
		log.Println("Invalid amount format:", err)
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	var beginDate, endDate sql.NullTime
	if payload.UserAssetAcquisitionBeginDate != "" {
		parsed, err := time.Parse("2006-01-02", payload.UserAssetAcquisitionBeginDate)
		if err != nil {
			http.Error(w, "Invalid acquisition begin date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		beginDate = sql.NullTime{Time: parsed, Valid: true}
	}

	if payload.UserAssetAcquisitionEndDate != "" {
		parsed, err := time.Parse("2006-01-02", payload.UserAssetAcquisitionEndDate)
		if err != nil {
			http.Error(w, "Invalid acquisition end date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		endDate = sql.NullTime{Time: parsed, Valid: true}
	}

	// Call stored procedure
	var message string
	err = database.QueryRow(
		"CALL CreateUserAsset($1, $2, $3, $4, $5, $6, $7)",
		payload.AssetTypeID,
		user.UserProfileID,
		payload.UserAssetName,
		amount,
		beginDate,
		endDate,
		message).Scan(&message)

	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	// Parse stored procedure response
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing stored procedure response:", err)
		http.Error(w, "Invalid response from stored procedure", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

// TODO
func UpdateAsset(w http.ResponseWriter, r *http.Request) {
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

func DeleteAsset(w http.ResponseWriter, r *http.Request) {
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

func CreateAssetParentIncome(w http.ResponseWriter, r *http.Request) {

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateAsset: Unauthorized access - no user found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload models.CreateUserAssetParentIncome

	// Get database connection
	database := db.GetDB()

	// Decodifica JSON
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if payload.UserID == 0 || payload.UserAssetID == 0 || payload.FinancialUserItemName == "" ||
		payload.RecurrencyID == 0 || payload.FinancialUserEntityItemID == 0 || payload.ParentIncomeAmount <= 0 || payload.BeginDate == "" {
		http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
		return
	}

	// Executa a procedure
	var message string
	err := database.QueryRow(`
		CALL CreateUserAssetParentIncome($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		user.UserProfileID,
		payload.UserAssetID,
		payload.FinancialUserItemName,
		payload.RecurrencyID,
		payload.FinancialUserEntityItemID,
		payload.ParentIncomeAmount,
		payload.BeginDate,
		message).Scan(&message)

	// Trata erro de execução
	if err != nil {
		log.Printf("Error calling procedure: %v", err)
		http.Error(w, "Error executing procedure", http.StatusInternalServerError)
		return
	}

	// Parse stored procedure response
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing stored procedure response:", err)
		http.Error(w, "Invalid response from stored procedure", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

func DeleteAssetParentIncome(w http.ResponseWriter, r *http.Request) {
	// Garante que o método seja DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Pega o usuário do contexto
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("DeleteAssetParentIncome: Unauthorized access - no user in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get database connection
	database := db.GetDB()

	// Lê o payload
	var payload struct {
		FinancialUserItemID int `json:"financial_user_item_id"`
		UserAssetID         int `json:"user_asset_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validação
	if payload.FinancialUserItemID == 0 || payload.UserAssetID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Executa a procedure
	var message string
	err := database.QueryRow(`
		CALL DeleteUserAssetParentIncome($1, $2, $3, $4)
	`,
		payload.FinancialUserItemID,
		user.UserProfileID,
		payload.UserAssetID,
		message).Scan(&message)

	if err != nil {
		log.Printf("Error executing stored procedure: %v", err)
		http.Error(w, "Error executing procedure", http.StatusInternalServerError)
		return
	}

	// Interpreta a resposta
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing stored procedure response:", err)
		http.Error(w, "Invalid response from stored procedure", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

func CreateAssetChildIncomeTax(w http.ResponseWriter, r *http.Request) {

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateAssetChildIncomeTax: Unauthorized access - no user found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload models.CreateUserAssetChildIncomeTax

	// Get database connection
	database := db.GetDB()

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Basic validations
	if payload.UserAssetID == 0 || payload.FinancialUserItemName == "" ||
		payload.FinancialUserEntityItemID == 0 || payload.ParentFinancialUserItemID == 0 || payload.TaxIncomeAmount <= 0 {
		http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
		return
	}

	// Execute stored procedure
	var message string
	err := database.QueryRow(`
		CALL CreateUserAssetChildIncomeTax($1, $2, $3, $4, $5, $6, $7)
	`,
		user.UserProfileID,
		payload.UserAssetID,
		payload.FinancialUserItemName,
		payload.FinancialUserEntityItemID,
		payload.ParentFinancialUserItemID,
		payload.TaxIncomeAmount,
		message).Scan(&message)

	if err != nil {
		log.Printf("Error calling procedure: %v", err)
		http.Error(w, "Error executing procedure", http.StatusInternalServerError)
		return
	}

	// Parse stored procedure response
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing stored procedure response:", err)
		http.Error(w, "Invalid response from stored procedure", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

func CreateAssetChildIncomeExpense(w http.ResponseWriter, r *http.Request) {

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateAssetChildIncomeExpense: Unauthorized access - no user found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload models.CreateUserAssetChildIncomeExpense

	// Get database connection
	database := db.GetDB()

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Basic validations
	if payload.UserAssetID == 0 || payload.FinancialUserItemName == "" ||
		payload.FinancialUserEntityItemID == 0 || payload.ParentFinancialUserItemID == 0 || payload.ExpenseAmount <= 0 {
		http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
		return
	}

	// Execute stored procedure
	var message string
	err := database.QueryRow(`
		CALL CreateUserAssetChildIncomeExpense($1, $2, $3, $4, $5, $6, $7)
	`,
		user.UserProfileID,
		payload.UserAssetID,
		payload.FinancialUserItemName,
		payload.FinancialUserEntityItemID,
		payload.ParentFinancialUserItemID,
		payload.ExpenseAmount,
		message).Scan(&message)

	if err != nil {
		log.Printf("Error calling procedure: %v", err)
		http.Error(w, "Error executing procedure", http.StatusInternalServerError)
		return
	}

	// Parse stored procedure response
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing stored procedure response:", err)
		http.Error(w, "Invalid response from stored procedure", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

func DeleteUserAssetChildIncomeExpense(w http.ResponseWriter, r *http.Request) {
	// Permite apenas método DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Pega o usuário do contexto
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("DeleteUserAssetChildIncomeExpense: Unauthorized - user not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get database connection
	database := db.GetDB()

	// Define struct temporária para receber o payload
	type Payload struct {
		FinancialUserItemID int `json:"financial_user_item_id"`
		UserAssetID         int `json:"user_asset_id"`
	}

	var input Payload

	// Faz o decode do JSON do body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("DeleteUserAssetChildIncomeExpense: Invalid request payload -", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Verifica campos obrigatórios
	if input.FinancialUserItemID == 0 || input.UserAssetID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Chama a stored procedure
	var message string

	err := database.QueryRow(`
		CALL DeleteUserAssetChildIncomeExpense($1, $2, $3, $4)
	`, input.FinancialUserItemID, user.UserProfileID, input.UserAssetID, message).Scan(&message)

	if err != nil {
		log.Printf("DeleteUserAssetChildIncomeExpense: Error calling procedure - %v", err)
		http.Error(w, "Error executing procedure", http.StatusInternalServerError)
		return
	}

	// Interpreta a resposta da stored procedure
	var response models.Response
	if err := json.Unmarshal([]byte(message), &response); err != nil {
		log.Println("DeleteUserAssetChildIncomeExpense: Error parsing procedure response -", err)
		http.Error(w, "Error parsing response", http.StatusInternalServerError)
		return
	}

	// Retorna resposta
	w.Header().Set("Content-Type", "application/json")
	if response.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteUserAssetChildIncomeTax(w http.ResponseWriter, r *http.Request) {
	// Permite apenas método DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Pega o usuário do contexto
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("DeleteUserAssetChildIncomeTax: Unauthorized - user not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get database connection
	database := db.GetDB()

	// Define struct temporária para receber o payload
	type Payload struct {
		FinancialUserItemID int `json:"financial_user_item_id"`
		UserAssetID         int `json:"user_asset_id"`
	}

	var input Payload

	// Faz o decode do JSON do body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("DeleteUserAssetChildIncomeTax: Invalid request payload -", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Verifica campos obrigatórios
	if input.FinancialUserItemID == 0 || input.UserAssetID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Chama a stored procedure
	var message string

	err := database.QueryRow(`
		CALL DeleteUserAssetChildIncomeTax($1, $2, $3, $4)
	`, input.FinancialUserItemID, user.UserProfileID, input.UserAssetID, message).Scan(&message)

	if err != nil {
		log.Printf("DeleteUserAssetChildIncomeTax: Error calling procedure - %v", err)
		http.Error(w, "Error executing procedure", http.StatusInternalServerError)
		return
	}

	// Interpreta a resposta da stored procedure
	var response models.Response
	if err := json.Unmarshal([]byte(message), &response); err != nil {
		log.Println("DeleteUserAssetChildIncomeTax: Error parsing procedure response -", err)
		http.Error(w, "Error parsing response", http.StatusInternalServerError)
		return
	}

	// Retorna resposta
	w.Header().Set("Content-Type", "application/json")
	if response.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

// TODO
func CreateAssetCategory(w http.ResponseWriter, r *http.Request) {
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

// TODO
func DeleteAssetCategory(w http.ResponseWriter, r *http.Request) {
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
