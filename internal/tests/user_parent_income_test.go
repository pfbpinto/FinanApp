package tests

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestCreateUserParentIncome_Success(t *testing.T) {
	// Configuração da conexão
	dsn := "postgres://postgres:Fpadminpostgre@localhost:5432/finanapp?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir conexão com o banco: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}
	log.Println("✅ Conexão com o banco de dados estabelecida.")

	var response string

	// Chamada da procedure
	err = db.QueryRow("CALL CreateUserParentIncome($1, $2, $3, $4, $5, $6, $7)",
		1,                // p_UserID
		"Salário Mensal", // p_FinancialUserItemName
		1,                // p_RecurrencyID (ex: mensal)
		1,                // p_FinancialUserEntityItemID (ex: empresa)
		15000.00,         // p_ParentIncomeAmount
		"2025-04-05",     // p_BeginDate
		&response,        // p_Message (OUT)
	).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure: %v", err)
	}

	fmt.Printf("📨 Resposta da procedure: %s\n", response)
}
