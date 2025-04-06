package tests

import (
	"testing"

	_ "github.com/lib/pq"
)

// Teste para a procedure CreateUserParentAsset
func TestCreateUserParentAsset(t *testing.T) {
	db := getTestDB(t)

	var response string

	err := db.QueryRow("CALL CreateUserAsset($1, $2, $3, $4, $5, $6, $7)",
		1,                  // p_AssetTypeID
		1,                  // p_UserProfileID
		"Novo Apartamento", // p_UserAssetName
		150000.00,          // p_UserAssetValueAmount
		"2025-04-01",       // p_UserAssetAcquisitionBeginDate
		"2025-12-01",       // p_UserAssetAcquisitionEndDate
		&response,          // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}

func TestCreateUserAssetParentIncome(t *testing.T) {
	// Conecta ao banco
	db := getTestDB(t)

	var response string

	// Chamada da procedure
	err := db.QueryRow("CALL CreateUserAssetParentIncome($1, $2, $3, $4, $5, $6, $7, $8)",
		1,                     // p_UserID
		1,                     // p_UserAssetID
		"Aluguel Apartamento", // p_FinancialUserItemName
		2,                     // p_RecurrencyID
		11,                    // p_FinancialUserEntityItemID
		2500.00,               // p_ParentIncomeAmount
		"2025-04-06",          // p_BeginDate
		&response,             // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}

func TestCreateUserAssetChildIncomeTax(t *testing.T) {
	// Conecta ao banco
	db := getTestDB(t)

	var response string

	// Chamada da procedure
	err := db.QueryRow("CALL CreateUserAssetChildIncomeTax($1, $2, $3, $4, $5, $6, $7)",
		1,                            // p_UserID
		1,                            // p_UserAssetID
		"Tax IR Aluguel Apartamento", // p_FinancialUserItemName
		12,                           // p_FinancialUserEntityItemID
		4,                            // p_ParentFinancialUserItemID
		450.75,                       // p_TaxIncomeAmount
		&response,                    // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)
}

// Testa a criação de uma despesa associada a um ativo do usuário
func TestCreateUserAssetChildIncomeExpense(t *testing.T) {
	db := getTestDB(t)

	var response string

	err := db.QueryRow("CALL CreateUserAssetChildIncomeExpense($1, $2, $3, $4, $5, $6, $7)",
		1,                   // p_UserID
		1,                   // p_UserAssetID
		"Condomínio Mensal", // p_FinancialUserItemName
		13,                  // p_FinancialUserEntityItemID
		4,                   // p_ParentFinancialUserItemID
		250.00,              // p_ExpenseAmount
		&response,           // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)

	if response == "" {
		t.Errorf("A resposta da procedure está vazia")
	}
}

// Testa a exclusão de um imposto associado a um ativo do usuário
func TestDeleteUserAssetChildIncomeTax(t *testing.T) {
	db := getTestDB(t)

	var response string

	err := db.QueryRow("CALL DeleteUserAssetChildIncomeTax($1, $2, $3, $4)",
		5,         // p_FinancialUserItemID
		1,         // p_UserID
		1,         // p_UserAssetID
		&response, // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)

	if response == "" {
		t.Errorf("A resposta da procedure está vazia")
	}
}

// Testa a exclusão de uma despesa associada a um ativo do usuário
func TestDeleteUserAssetChildIncomeExpense(t *testing.T) {
	db := getTestDB(t)

	var response string

	err := db.QueryRow("CALL DeleteUserAssetChildIncomeExpense($1, $2, $3, $4)",
		6,         // p_FinancialUserItemID
		1,         // p_UserID
		1,         // p_UserAssetID
		&response, // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	LogProcedureResponse(t, response)

	if response == "" {
		t.Errorf("A resposta da procedure está vazia")
	}
}
