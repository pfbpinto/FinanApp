package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestRunDatabaseSQL(t *testing.T) {
	connStr := "host=localhost port=5432 user=postgres password=Fpadminpostgre dbname=finanapp sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	files := []string{
		"internal/db/database.sql",
		"internal/db/procedures.sql",
		"internal/db/seeder.sql",
	}

	for _, file := range files {
		t.Logf("Executando: %s", file)

		sqlContent, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Erro ao ler %s: %v", file, err)
		}

		if _, err := db.Exec(string(sqlContent)); err != nil {
			t.Fatalf("Erro ao executar %s: %v", file, err)
		}
	}

	log.Println("Banco populado com sucesso com todos os arquivos.")
}
