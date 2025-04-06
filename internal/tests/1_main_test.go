package tests

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

type ProcedureResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Conecta com o banco apenas uma vez
func setupDB() *sql.DB {
	dsn := "postgres://postgres:Fpadminpostgre@localhost:5432/finanapp?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com o banco: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}
	log.Println("✅ Conexão com o banco de dados estabelecida.")
	return db
}

// Fornece a conexão já estabelecida aos testes
func getTestDB(t *testing.T) *sql.DB {
	if testDB == nil {
		t.Fatal("Banco de dados de teste não foi inicializado")
	}
	return testDB
}

// Executado uma vez antes de todos os testes
func TestMain(m *testing.M) {
	testDB = setupDB()
	defer testDB.Close()

	code := m.Run()
	os.Exit(code)
}

func LogProcedureResponse(t *testing.T, response string) {
	var parsed ProcedureResponse
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		t.Errorf("❌ Erro ao interpretar resposta JSON: %v", err)
		return
	}

	if parsed.Status == "success" {
		t.Logf("✅ Sucesso: %s", parsed.Message)
	} else {
		t.Errorf("❌ Falha: %s", parsed.Message)
	}
}
