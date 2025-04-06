package routes

import (
	"finanapp/internal/handlers"
	"finanapp/internal/middlewares"
	"net/http"

	"github.com/rs/cors"
)

func RegisterExpenseRoutes(mux *http.ServeMux, corsMiddleware *cors.Cors) {

	mux.Handle("/api/expense", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateExpense),
	)))
	mux.Handle("/api/expense-update", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UpdateExpense),
	)))
	mux.Handle("/api/delete-expense", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteExpense),
	)))

}
