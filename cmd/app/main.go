package main

import (
	"finanapp/config"
	"finanapp/internal/db"
	"finanapp/internal/messaging"
	"finanapp/internal/routes"

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
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Use ServeMux instead of default
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, corsMiddleware)

	// Start server
	log.Printf("Server running on port %s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, mux)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
