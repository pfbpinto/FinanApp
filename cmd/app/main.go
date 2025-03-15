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
	// Check for open migrations
	db.EnsureDatabaseExists(db.DB)
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
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", middlewares.AuthMiddleware(handlers.Logout))
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/forgot-password", handlers.ForgotPassword)

	http.HandleFunc("/user", middlewares.AuthMiddleware(handlers.UserDashboard))
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
		middlewares.AuthMiddleware(handlers.UserDashboardReact),
	)))
	http.Handle("/api/user-edit", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UserUpdate))))

	// Assets API's
	http.Handle("/api/asset-type", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.GetAssetType),
	)))
	http.Handle("/api/assets", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateUserAsset))))

	http.Handle("/api/assets/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UpdateUserAsset))))

	http.Handle("/api/delete-assets/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteAsset))))

	// Tax API's
	http.Handle("/api/get-taxes", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.GetTaxes),
	)))
	http.Handle("/api/create-taxes", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateTax),
	)))
	http.Handle("/api/delete-tax/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteTaxes),
	)))

	// Category API
	http.Handle("/api/categories", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.GetCategories),
	)))
	http.Handle("/api/create-category", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateCategories),
	)))
	http.Handle("/api/delete-category/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteCategories),
	)))

	// Income API's
	http.Handle("/api/income-type", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.GetIncomeType),
	)))
	http.Handle("/api/income", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateUserIncome),
	)))
	http.Handle("/api/income/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UpdateUserIncome))))

	http.Handle("/api/delete-income/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteIncome))))

	// Expense API's
	http.Handle("/api/expense-type", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.GetExpenseType),
	)))
	http.Handle("/api/expense", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateUserExpense),
	)))
	http.Handle("/api/expense/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.UpdateUserExpense))))

	http.Handle("/api/delete-expense/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteExpense))))

	// User Group API's
	http.Handle("/api/user-group", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.GetGroups),
	)))
	http.Handle("/api/create-group", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateGroup),
	)))
	http.Handle("/api/create-group-item", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateGroupItem),
	)))

	http.Handle("/api/create-group-invite", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.CreateGroupInvite),
	)))

	http.Handle("/api/delete-group/{id}", corsMiddleware.Handler(http.HandlerFunc(
		middlewares.AuthMiddleware(handlers.DeleteGroup),
	)))

	// Start the server
	log.Printf("Server running on port %s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
