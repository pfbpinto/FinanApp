package tests

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var database *sql.DB

// Inicializa a conexão com o banco de dados
func init() {
	dsn := "postgres://postgres:Fpadminpostgre@localhost:5432/finanapp?sslmode=disable"

	var err error
	database, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com o banco: %v", err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco: %v", err)
	}

	fmt.Println("✅ Conexão com o banco de dados estabelecida.")
}

// Testa a criação de um novo usuário via stored procedure
func TestCreateUser_Success(t *testing.T) {
	// Garante que o email seja único em cada execução
	email := fmt.Sprintf("francisco_%d@example.com", time.Now().UnixNano())

	var response string
	err := database.QueryRow("CALL CreateUser($1, $2, $3, $4, $5, $6)",
		"Francisco",
		"Pinto",
		email,
		"superSecure123!",
		"1982-04-05",
		&response).Scan(&response)

	if err != nil {
		t.Fatalf("Erro ao executar a procedure CreateUser: %v", err)
	}

	fmt.Printf("📨 Resposta da procedure: %s\n", response)

	if response == "" {
		t.Errorf("A resposta da procedure está vazia")
	}
}
