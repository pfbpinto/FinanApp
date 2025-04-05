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
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("UserDashboard: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Get database connection
	database := db.GetDB()

	// Define a struct to match the incoming JSON payload
	var payload struct {
		FinancialUserItemId   int     `json:"FinancialUserItemId"`
		FinancialUserItemName string  `json:"FinancialUserItemName"`
		RecurrencyID          string  `json:"RecurrencyID"`
		CurrencyID            string  `json:"CurrencyID"`
		Amount                string  `json:"amount"`      // amount is string in payload
		IncomeValue           *string `json:"IncomeValue"` // nullable
	}

	// Parse the JSON payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if payload.FinancialUserItemId == 0 || payload.FinancialUserItemName == "" || payload.Amount == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Convert amount to float64
	amountFloat, err := strconv.ParseFloat(payload.Amount, 64)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Prepare values for stored procedure
	beginDate := time.Now().Format("2006-01-02")
	isActive := true

	// Call the stored procedure
	var message string
	err = database.QueryRow(
		"CALL UpdateUserParentIncome($1, $2, $3, $4, $5, $6, $7)",
		payload.FinancialUserItemId,
		user.UserProfileID,
		payload.FinancialUserItemName,
		amountFloat,
		beginDate,
		isActive,
		message).Scan(&message)

	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	// Return success message from procedure
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
	})
}

func DeleteIncome(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Get database connection
	database := db.GetDB()

	// Call the stored procedure
	var message string
	err := database.QueryRow(
		"CALL DeleteUserParentIncome($1, $2, $3)",
		payload.ItemID,
		user.UserProfileID,
		message).Scan(&message)

	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	// Return success message from procedure
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
	})
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

func CreateIncomeTax(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateIncomeTax: Request received")

	// Verifica se o método é POST
	if r.Method != http.MethodPost {
		log.Println("CreateIncomeTax: Invalid HTTP method:", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Recupera o usuário do contexto
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateIncomeTax: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	log.Println("CreateIncomeTax: User found, ID =", user.UserProfileID)

	// Define estrutura do payload esperado
	var payload struct {
		FinancialUserItemName     string `json:"financialUserItemName"`
		RecurrencyID              string `json:"recurrencyId"`
		FinancialUserEntityItemID string `json:"financialUserEntityItemId"`
		ParentFinancialUserItemID string `json:"parentFinancialUserItemId"`
	}

	// Decodifica o JSON da requisição
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("CreateIncomeTax: Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Printf("CreateIncomeTax: Payload decoded: %+v\n", payload)

	// Converte campos numéricos de string para int
	recurrencyID, err := strconv.Atoi(payload.RecurrencyID)
	if err != nil {
		log.Println("CreateIncomeTax: Invalid recurrencyId:", payload.RecurrencyID)
		http.Error(w, "Invalid recurrencyId", http.StatusBadRequest)
		return
	}

	financialUserEntityItemID, err := strconv.Atoi(payload.FinancialUserEntityItemID)
	if err != nil {
		log.Println("CreateIncomeTax: Invalid financialUserEntityItemId:", payload.FinancialUserEntityItemID)
		http.Error(w, "Invalid financialUserEntityItemId", http.StatusBadRequest)
		return
	}

	parentFinancialUserItemID, err := strconv.Atoi(payload.ParentFinancialUserItemID)
	if err != nil {
		log.Println("CreateIncomeTax: Invalid parentFinancialUserItemId:", payload.ParentFinancialUserItemID)
		http.Error(w, "Invalid parentFinancialUserItemId", http.StatusBadRequest)
		return
	}

	// Conexão com o banco de dados
	database := db.GetDB()

	// Chama a stored procedure
	var message string
	log.Println("CreateIncomeTax: Calling stored procedure CreateUserChildIncomeTax")
	err = database.QueryRow(
		"CALL CreateUserChildIncomeTax($1, $2, $3, $4, $5, $6)",
		user.UserProfileID,
		payload.FinancialUserItemName,
		recurrencyID,
		financialUserEntityItemID,
		parentFinancialUserItemID,
		message).Scan(&message)

	if err != nil {
		log.Println("CreateIncomeTax: Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	log.Println("CreateIncomeTax: Stored procedure executed successfully. Message:", message)

	// Retorna resposta de sucesso
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
	})
}

// CreateIncomeExpense lida com a criação de uma nova despesa de renda
func CreateIncomeExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verifica se o método é POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	// Recupera o usuário do contexto
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateIncomeExpense: Acesso não autorizado - nenhum usuário encontrado no contexto.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Estrutura para decodificar o payload JSON
	var payload struct {
		FinancialUserItemName     string `json:"financialUserItemName"`
		RecurrencyID              string `json:"recurrencyId"`
		FinancialUserEntityItemID string `json:"financialUserEntityItemId"`
		ParentFinancialUserItemID string `json:"parentFinancialUserItemId"`
	}

	// Decodifica o corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("CreateIncomeExpense: Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Converte os IDs de string para int
	recurrencyID, err := strconv.Atoi(payload.RecurrencyID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid recurrencyId"})
		return
	}

	financialUserEntityItemID, err := strconv.Atoi(payload.FinancialUserEntityItemID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid financialUserEntityItemId"})
		return
	}

	parentFinancialUserItemID, err := strconv.Atoi(payload.ParentFinancialUserItemID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid parentFinancialUserItemId"})
		return
	}

	// Conexão com o banco de dados
	database := db.GetDB()

	// Chama a procedure armazenada
	var message string
	err = database.QueryRow(
		"CALL CreateUserChildIncomeExpense($1, $2, $3, $4, $5, $6)",
		user.UserProfileID,
		payload.FinancialUserItemName,
		recurrencyID,
		financialUserEntityItemID,
		parentFinancialUserItemID,
		&message,
	).Scan(&message)

	if err != nil {
		log.Println("CreateIncomeExpense: Erro no banco de dados:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to execute stored procedure"})
		return
	}

	// Responde com sucesso
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
	})
}
