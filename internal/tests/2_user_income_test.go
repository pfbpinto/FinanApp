package tests

import (
	"testing"

	_ "github.com/lib/pq"
)

// Teste para a procedure CreateUserParentIncome
func TestCreateUserParentIncome(t *testing.T) {
	// Chama a função de conexão
	db := getTestDB(t)

	var response string

	// Chamada da procedure
	err := db.QueryRow("CALL CreateUserParentIncome($1, $2, $3, $4, $5, $6, $7)",
		1,                // p_UserID
		"Salário Mensal", // p_FinancialUserItemName
		1,                // p_RecurrencyID (ex: mensal)
		5,                // p_FinancialUserEntityItemID (ex: Income)
		15000.00,         // p_ParentIncomeAmount
		"2025-04-05",     // p_BeginDate
		&response,        // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}

// Teste para a procedure UpdateUserParentIncome
func TestUpdateUserParentIncome(t *testing.T) {
	// Chama a função de conexão
	db := getTestDB(t)

	var response string

	// Chamada da procedure UpdateUserParentIncome
	err := db.QueryRow("CALL UpdateUserParentIncome($1, $2, $3, $4, $5, $6, $7)",
		1,                    // p_FinancialUserItemID
		1,                    // p_UserID (ID do usuário criado no teste anterior)
		"Salário Atualizado", // p_NewFinancialUserItemName
		20000.00,             // p_NewParentIncomeAmount
		"2025-05-01",         // p_NewBeginDate
		true,                 // p_IsActive (usuário ativo)
		&response,            // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}

// Teste para a procedure CreateUserChildIncomeTax
func TestCreateUserChildIncomeTax(t *testing.T) {
	// Chama a função de conexão
	db := getTestDB(t)

	var response string

	// Chamada da procedure CreateUserChildIncomeTax
	err := db.QueryRow("CALL CreateUserChildIncomeTax($1, $2, $3, $4, $5, $6)",
		1,                  // p_UserID
		"Imposto de Renda", // p_FinancialUserItemName
		1,                  // p_RecurrencyID (ex: mensal)
		7,                  // p_FinancialUserEntityItemID (ex: Income Tax)
		1,                  // p_ParentFinancialUserItemID (ID do item do Income parent)
		&response,          // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}

// Teste para a procedure CreateUserChildIncomeExpense
func TestCreateUserChildIncomeExpense(t *testing.T) {
	// Chama a função de conexão
	db := getTestDB(t)

	var response string

	// Chamada da procedure CreateUserChildIncomeExpense
	err := db.QueryRow("CALL CreateUserChildIncomeExpense($1, $2, $3, $4, $5, $6)",
		1,                // p_UserID
		"Despesa Mensal", // p_FinancialUserItemName
		1,                // p_RecurrencyID (ex: mensal)
		8,                // p_FinancialUserEntityItemID (ex: Income Expense)
		1,                // p_ParentFinancialUserItemID (ID do item do Income parent)
		&response,        // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}

// Teste para a procedure DeleteUserParentIncome
func TestDeleteUserParentIncome(t *testing.T) {
	// Chama a função de conexão
	db := getTestDB(t)

	var response string

	// Chamada da procedure DeleteUserParentIncome
	err := db.QueryRow("CALL DeleteUserParentIncome($1, $2, $3)",
		1,         // p_FinancialUserItemID
		1,         // p_UserID (ID do Income a ser deletado)
		&response, // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}
