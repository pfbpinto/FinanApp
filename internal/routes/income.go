package routes

import (
	"finanapp/internal/handlers"
	"finanapp/internal/middlewares"
	"net/http"

	"github.com/rs/cors"
)

func RegisterIncomeRoutes(mux *http.ServeMux, corsMiddleware *cors.Cors) {

	mux.Handle("/api/income", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateIncome),
	)))
	mux.Handle("/api/income-update/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UpdateIncome),
	)))
	mux.Handle("/api/delete-income/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteIncome),
	)))
	mux.Handle("/api/income-item/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.IncomeItem),
	)))
	mux.Handle("/api/income-category", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateIncomeCategory),
	)))
	mux.Handle("/api/delete-income-category", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteIncomeCategory),
	)))
	mux.Handle("/api/create-income-tax", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateIncomeTax),
	)))
	mux.Handle("/api/create-income-expense", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateIncomeExpense),
	)))
}
