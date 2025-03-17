package db

import (
	"context"
	"database/sql"
	"finanapp/config"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// DB will be the global instance of the database connection
var DB *sql.DB

var RDB *redis.Client

// InitDB initializes the database connection
func InitDB() {
	// Load configurations
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	// Connect to the database
	DB, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error testing the database connection: %v", err)
	}

	fmt.Println("Connected to the PostgreSQL database successfully")
}

// GetDB returns the database instance
func GetDB() *sql.DB {
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
