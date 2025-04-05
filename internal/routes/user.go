package routes

import (
	"finanapp/internal/handlers"
	"finanapp/internal/middlewares"
	"net/http"

	"github.com/rs/cors"
)

func RegisterUserRoutes(mux *http.ServeMux, corsMiddleware *cors.Cors) {

	mux.Handle("/api/user", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UserDashboard),
	)))
	mux.Handle("/api/user-edit", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UserUpdate),
	)))
	mux.Handle("/api/user-income", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UserIncome),
	)))
	mux.Handle("/api/user-asset", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UserAsset),
	)))

}
