package tests

import (
	"bytes"
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestMain inicializa o banco de dados antes de rodar os testes
func TestMain(m *testing.M) {
	db.InitDB() // Inicializa a conexão com o banco de dados
	code := m.Run()
	os.Exit(code)
}

// inicializa o router com a rota da API
func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/register", handlers.RegisterReact).Methods("POST")
	return r
}

// payload enviado para a API
type RegisterPayload struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	UserPassword string `json:"user_password"`
	DateOfBirth  string `json:"date_of_birth"`
}

// resposta esperada da API
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
		EmailAddress: "jones@pipeline.com", // mesmo e-mail do teste anterior
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
