package tests

import (
	"database/sql"
	"log"
	"testing"
)

// Função para estabelecer a conexão com o banco
func setupDB(t *testing.T) *sql.DB {
	// Configuração da conexão
	dsn := "postgres://postgres:Fpadminpostgre@localhost:5432/finanapp?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir conexão com o banco: %v", err)
	}
	if err = db.Ping(); err != nil {
		t.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}
	log.Println("✅ Conexão com o banco de dados estabelecida.")
	return db
}
