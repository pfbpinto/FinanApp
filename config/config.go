package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config structure to store configuration values
type Config struct {
	Port        string
	DatabaseURL string
}

// Load loads the environment variables
func Load() (*Config, error) {
	// Attempts to get the directory where the code is being executed
	_, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error retrieving working directory: %v", err)
	}

	// Project root directory where the .env file should be located
	envFilePath := filepath.Join(".", ".env")

	// Load environment variables from the .env file
	if err := godotenv.Load(envFilePath); err != nil {
		log.Printf("Warning: .env not found at path: %v", envFilePath)
	}

	return &Config{
		Port:        os.Getenv("PORT"), // Agora falha se PORT não estiver definido
    	DatabaseURL: os.Getenv("DATABASE_URL"),
	}, nil
}
