package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Função de teste para a inicialização do banco de dados
func TestInitDB(t *testing.T) {
	// Configuração do PostgreSQL direto no teste
	databaseURL := "postgres://postgres:Fpadminpostgre@localhost:5432/finanapp?sslmode=disable"

	// Tenta estabelecer a conexão com o banco de dados
	DB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		t.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}

	// Verifica se a conexão foi bem-sucedida
	sqlDB, err := DB.DB()
	if err != nil {
		t.Fatalf("Erro ao obter informações sobre DB: %v", err)
	}

	// Verifica se a conexão com o banco de dados foi bem-sucedida
	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Erro ao pingar o banco de dados: %v", err)
	}

	// Caso o teste passe, imprime uma mensagem de sucesso
	t.Log("Conexão com o PostgreSQL foi bem-sucedida")

	// Usando o assert para confirmar que a conexão foi bem-sucedida
	assert.NoError(t, err)
}
