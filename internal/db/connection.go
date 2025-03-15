package db

import (
	"context"
	"finanapp/config"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB will be the global instance of the database connection
var DB *gorm.DB

var RDB *redis.Client

// InitDB initializes the database connection
func InitDB() {
	// Load configurations
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	// Connect to the database
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Connected to the PostgreSQL database")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	// Checks if the connection has been initialized; if not, calls InitDB
	if DB == nil {
		InitDB()
	}
	return DB
}

func InitRedis() {
	// Get environment variables
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost" // Default value for local environment
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379" // Default port for Redis
	}

	// Configure the Redis client
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Redis address
		Password: "",        // Password (if any)
		DB:       0,         // Default database
	})

	// Test the connection to Redis
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully")
}

// EnsureDatabaseExists ensures that migrations and seeds are applied.
func EnsureDatabaseExists(db *gorm.DB) {
	log.Println("Ensuring database is up-to-date...")

	// Run migrations to ensure the database schema is up-to-date
	RunMigrations(db)

	// Run seeds to populate the database with initial data
	SeedDatabase(db)

	log.Println("Database setup completed successfully.")
}
