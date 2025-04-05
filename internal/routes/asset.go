package routes

import (
	"finanapp/internal/handlers"
	"finanapp/internal/middlewares"
	"net/http"

	"github.com/rs/cors"
)

func RegisterAssetRoutes(mux *http.ServeMux, corsMiddleware *cors.Cors) {

	mux.Handle("/api/asset", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateAsset),
	)))
	mux.Handle("/api/asset-update/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UpdateAsset),
	)))
	mux.Handle("/api/delete-asset/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteAsset),
	)))
	mux.Handle("/api/asset-item/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.AssetItem),
	)))
	mux.Handle("/api/asset-income", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateAssetParentIncome),
	)))
	mux.Handle("/api/delete-asset-income", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteAssetParentIncome),
	)))

	mux.Handle("/api/asset-income-tax", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateAssetChildIncomeTax),
	)))
	mux.Handle("/api/asset-income-expense", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateAssetChildIncomeExpense),
	)))

	mux.Handle("/api/delete-asset-income-tax", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteUserAssetChildIncomeTax),
	)))

	mux.Handle("/api/delete-asset-income-expense", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteUserAssetChildIncomeExpense),
	)))

	mux.Handle("/api/asset-category", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateAssetCategory),
	)))
	mux.Handle("/api/delete-asset-category", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteAssetCategory),
	)))
}
