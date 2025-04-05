package routes

import (
	"finanapp/internal/handlers"
	"finanapp/internal/middlewares"
	"net/http"

	"github.com/rs/cors"
)

func RegisterAuthRoutes(mux *http.ServeMux, corsMiddleware *cors.Cors) {

	mux.Handle("/api/auth-status", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.AuthStatus),
	)))
	mux.Handle("/api/login", corsMiddleware.Handler(http.HandlerFunc(handlers.LoginReact)))
	mux.Handle("/api/logout", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.LogoutReact),
	)))
	mux.Handle("/api/register", corsMiddleware.Handler(http.HandlerFunc(handlers.RegisterReact)))
}
