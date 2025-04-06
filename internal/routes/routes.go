package routes

import (
	"finanapp/internal/handlers"
	"finanapp/internal/middlewares"
	"net/http"

	"github.com/rs/cors"
)

func RegisterRoutes(mux *http.ServeMux, corsMiddleware *cors.Cors) {
	RegisterAuthRoutes(mux, corsMiddleware)
	RegisterUserRoutes(mux, corsMiddleware)
	RegisterIncomeRoutes(mux, corsMiddleware)
	RegisterAssetRoutes(mux, corsMiddleware)
	RegisterExpenseRoutes(mux, corsMiddleware)

	// Static
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// Home and health
	mux.HandleFunc("/", middlewares.AuthMiddleware(handlers.Home))
	mux.HandleFunc("/health", handlers.Health)
}
