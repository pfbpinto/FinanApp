package tests

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

// Teste para a procedure CreateUserParentAsset
func TestCreateUserParentAsset_Success(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

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

	fmt.Printf("📨 Resposta da procedure: %s\n", response)
}

func TestCreateUserAssetParentIncome_Success(t *testing.T) {
	// Conecta ao banco
	db := setupDB(t)
	defer db.Close()

	var response string

	// Chamada da procedure
	err := db.QueryRow("CALL CreateUserAssetParentIncome($1, $2, $3, $4, $5, $6, $7, $8)",
		1,                // p_UserID
		1,                // p_UserAssetID
		"Aluguel Mensal", // p_FinancialUserItemName
		2,                // p_RecurrencyID
		11,               // p_FinancialUserEntityItemID
		2500.00,          // p_ParentIncomeAmount
		"2025-04-06",     // p_BeginDate
		&response,        // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	fmt.Printf("📨 Resposta da procedure: %s\n", response)
}

func TestCreateUserAssetChildIncomeTax_Success(t *testing.T) {
	// Conecta ao banco
	db := setupDB(t)
	defer db.Close()

	var response string

	// Chamada da procedure
	err := db.QueryRow("CALL CreateUserAssetChildIncomeTax($1, $2, $3, $4, $5, $6, $7)",
		1,                  // p_UserID
		1,                  // p_UserAssetID
		"IR sobre Aluguel", // p_FinancialUserItemName
		12,                 // p_FinancialUserEntityItemID
		4,                  // p_ParentFinancialUserItemID
		450.75,             // p_TaxIncomeAmount
		&response,          // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	fmt.Printf("📨 Resposta da procedure: %s\n", response)
}
