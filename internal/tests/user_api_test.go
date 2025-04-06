package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"finanapp/internal/handlers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	// Conexão direta para testes
	connStr := "host=localhost port=5432 user=postgres password=Fpadminpostgre dbname=finanapp sslmode=disable"
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Erro ao conectar no banco: %v\n", err)
		os.Exit(1)
	}
	defer testDB.Close()

	// Verifica a conexão
	err = testDB.Ping()
	if err != nil {
		fmt.Printf("Erro ao testar conexão com o banco: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/register", handlers.RegisterReact).Methods("POST")
	return r
}

type RegisterPayload struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	UserPassword string `json:"user_password"`
	DateOfBirth  string `json:"date_of_birth"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestRegisterUser_Success(t *testing.T) {
	router := setupRouter()

	payload := RegisterPayload{
		FirstName:    "Pipeline",
		LastName:     "Jones",
		EmailAddress: "jones@pipeline.com",
		UserPassword: "password123",
		DateOfBirth:  "2025-07-30",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response.Status)
	assert.Contains(t, response.Message, "Usuário criado")
}

func TestRegisterUser_EmailDuplicado(t *testing.T) {
	router := setupRouter()

	payload := RegisterPayload{
		FirstName:    "Pipeline",
		LastName:     "Jones",
		EmailAddress: "jones@pipeline.com",
		UserPassword: "password123",
		DateOfBirth:  "2025-07-30",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "fail", response.Status)
	assert.Contains(t, response.Message, "já está cadastrado")
}
