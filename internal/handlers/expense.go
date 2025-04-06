package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	// Check if it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("CreateUserParentExpense: Unauthorized - user not found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get DB connection
	database := db.GetDB()

	// Decode JSON payload
	var payload models.UserParentExpense
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding payload:", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate and convert fields
	if strings.TrimSpace(payload.ParentExpenseAmount) == "" {
		log.Println("ParentExpenseAmount is empty")
		http.Error(w, "Amount cannot be empty", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(payload.ParentExpenseAmount, 64)
	if err != nil {
		log.Println("Invalid amount format:", err)
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	if payload.BeginDate == "" {
		log.Println("BeginDate is empty")
		http.Error(w, "BeginDate cannot be empty", http.StatusBadRequest)
		return
	}

	beginDate, err := time.Parse("2006-01-02", payload.BeginDate)
	if err != nil {
		log.Println("Invalid BeginDate format:", err)
		http.Error(w, "Invalid BeginDate format. Expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Call the stored procedure
	var message string
	err = database.QueryRow(
		"CALL CreateUserParentExpense($1, $2, $3, $4, $5, $6, $7)",
		user.UserProfileID,
		payload.FinancialUserItemName,
		payload.RecurrencyID,
		payload.FinancialUserEntityItemID,
		amount,
		beginDate,
		message).Scan(&message)

	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	// Unmarshal the response from stored procedure
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing stored procedure response:", err)
		http.Error(w, "Invalid response from stored procedure", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	// Check if it's a PUT request
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("UpdateUserParentExpense: Unauthorized - user not found in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get DB connection
	database := db.GetDB()

	// Decode JSON payload
	var payload models.UserParentExpenseUpdate
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding payload:", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if payload.FinancialUserItemID == 0 {
		http.Error(w, "Missing FinancialUserItemID", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(payload.NewFinancialUserItemName) == "" {
		http.Error(w, "Missing FinancialUserItemName", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(payload.NewParentExpenseAmount) == "" {
		http.Error(w, "Missing ParentExpenseAmount", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(payload.NewBeginDate) == "" {
		http.Error(w, "Missing BeginDate", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(payload.NewParentExpenseAmount, 64)
	if err != nil {
		log.Println("Invalid amount format:", err)
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	beginDate, err := time.Parse("2006-01-02", payload.NewBeginDate)
	if err != nil {
		log.Println("Invalid BeginDate format:", err)
		http.Error(w, "Invalid BeginDate format. Expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Call the stored procedure
	var message string
	err = database.QueryRow(
		"CALL UpdateUserParentExpense($1, $2, $3, $4, $5, $6, $7)",
		payload.FinancialUserItemID,
		user.UserProfileID,
		payload.NewFinancialUserItemName,
		amount,
		beginDate,
		payload.IsActive,
		message).Scan(&message)

	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	// Unmarshal stored procedure JSON response
	var resp models.Response
	if err := json.Unmarshal([]byte(message), &resp); err != nil {
		log.Println("Error parsing stored procedure response:", err)
		http.Error(w, "Invalid response from stored procedure", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	// Garante que o método HTTP seja DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Recupera o usuário do contexto
	user, ok := r.Context().Value("user").(models.UserProfile)
	if !ok {
		log.Println("DeleteExpense: Unauthorized access - no user found in context.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Estrutura para capturar o payload da requisição
	var payload struct {
		ItemID int `json:"itemId"`
	}

	// Decodifica o JSON do corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("DeleteExpense: Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Conecta ao banco de dados
	database := db.GetDB()

	// Variável para capturar a resposta da stored procedure
	var message string

	// Executa a stored procedure
	err := database.QueryRow(
		"CALL DeleteUserParentExpense($1, $2, $3)",
		payload.ItemID,
		user.UserProfileID,
		message).Scan(&message)

	if err != nil {
		log.Println("DeleteExpense: Database error:", err)
		http.Error(w, "Failed to execute stored procedure", http.StatusInternalServerError)
		return
	}

	// Retorna a resposta da stored procedure em formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": message,
	})
}
