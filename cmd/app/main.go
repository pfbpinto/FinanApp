package main

import (
	"finanapp/config"
	"finanapp/internal/db"
	"finanapp/internal/handlers"
	"finanapp/internal/messaging"
	"finanapp/internal/middlewares"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Change the working directory to the project root
	err := os.Chdir("../../") // Change the working directory to the project root
	if err != nil {
		log.Fatalf("Error changing working directory: %v", err)
	}

	// Load configurations
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	// Tente carregar o arquivo .env
	err2 := godotenv.Load()
	if err2 != nil {
		log.Printf("Warning: .env not found at path: %s", ".env")
	} else {
		log.Println(".env loaded successfully")
	}

	// Initialize database connection
	db.InitDB()
	db.InitRedis()
	// Routing for static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// Initialize NSQ Consumer
	consumer, err := messaging.NewConsumer()
	if err != nil {
		log.Fatalf("Error creating NSQ consumer: %v", err)
	}
	log.Printf("NSQ Consumer running")
	go consumer.StartConsumer()

	// CORS middleware configuration
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},  // Permitir apenas o frontend React
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Métodos permitidos
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Permitir envio de cookies/sessões
	})

	// Configure routes
	http.HandleFunc("/", middlewares.AuthMiddleware(handlers.Home))
	http.HandleFunc("/health", handlers.Health)

	// React Frontend (com suporte ao middleware de CORS)
	http.Handle("/api/auth-status", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.AuthStatus),
	)))
	http.Handle("/api/login", corsMiddleware.Handler(http.HandlerFunc(handlers.LoginReact)))
	http.Handle("/api/logout", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.LogoutReact),
	)))
	http.Handle("/api/register", corsMiddleware.Handler(http.HandlerFunc(handlers.RegisterReact)))
	http.Handle("/api/user", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UserDashboard),
	)))
	http.Handle("/api/user-edit", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UserUpdate))))

	// Start the server
	log.Printf("Server running on port %s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
