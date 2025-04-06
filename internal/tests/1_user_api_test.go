package tests

import (
	"fmt"
	"testing"
	"time"
)

// Testa a criação de um novo usuário via stored procedure
func TestCreateUser(t *testing.T) {
	db := getTestDB(t) // usa a conexão já aberta e reutilizável

	// Garante que o email seja único em cada execução
	email := fmt.Sprintf("francisco_%d@example.com", time.Now().UnixNano())

	var response string
	err := db.QueryRow("CALL CreateUser($1, $2, $3, $4, $5, $6)",
		"Francisco",
		"Pinto",
		email,
		"superSecure123!",
		"1982-04-05",
		&response).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure CreateUser: %v", err)
	}

	LogProcedureResponse(t, response)

	if response == "" {
		t.Errorf("A resposta da procedure está vazia")
	}
}
