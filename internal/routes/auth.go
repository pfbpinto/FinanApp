package routes

import (
	"encoding/json"
	"net/http"

	"finanapp/internal/auth"
)

// Estrutura para enviar resposta ao frontend
type UserResponse struct {
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

// Verificar autenticação
func VerifyAuth(w http.ResponseWriter, r *http.Request) {
	// Obter o token JWT do cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	// Validar o token usando a função existente
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Obter o email dos claims
	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Token não contém email válido", http.StatusUnauthorized)
		return
	}

	// Simulação de um avatar
	response := UserResponse{
		Email:  email,
		Avatar: "https://via.placeholder.com/150", // Substitua por um valor real, se disponível
	}

	// Retornar os dados do usuário
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
