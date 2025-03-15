package db

import (
	"finanapp/internal/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Define the database connection string globally
var dsn = "postgres://postgres:Fpadminpostgre@localhost:5432/finanapp?sslmode=disable"

// Helper function to initialize database connection
func getDBInstance() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func TestUserTypeID(t *testing.T) {
	// Get the DB instance
	dbInstance, err := getDBInstance()
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}

	// Perform the migration
	if err := dbInstance.AutoMigrate(&models.UserType{}); err != nil {
		t.Fatalf("Error migrating the UserType table")
	}

	// Insert a UserType record
	userType := models.UserType{
		Name:     "user",
		IsActive: true,
	}
	if result := dbInstance.Create(&userType); result.Error != nil {
		t.Fatalf("Error inserting user: %v", result.Error)
	}

	// Validate the insertion
	var fetchedUserType models.UserType
	if err := dbInstance.First(&fetchedUserType, "name = ?", userType.Name).Error; err != nil {
		t.Fatalf("UserType not found after insertion: %v", err)
	}

	// Check if the data matches
	if fetchedUserType.Name != userType.Name {
		t.Errorf("Expected name %v, but got %v", userType.Name, fetchedUserType.Name)
	}

	t.Log("UserType inserted and validated successfully")
}

func TestInsertUser(t *testing.T) {
	// Get the DB instance
	dbInstance, err := getDBInstance()
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}

	// Perform the migration
	if err := dbInstance.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Error migrating the User table: %v", err)
	}

	// Create a user for testing
	user := models.User{
		FirstName:    "testuser",
		EmailAddress: "testuser@example.com",
		Password:     "password123",
		UserTypeID:   1,
		IsActive:     true,
	}

	// Insert the user
	if result := dbInstance.Create(&user); result.Error != nil {
		t.Fatalf("Error inserting user: %v", result.Error)
	}

	// Validate the insertion
	var fetchedUser models.User
	if err := dbInstance.First(&fetchedUser, "email_address = ?", user.EmailAddress).Error; err != nil {
		t.Fatalf("User not found after insertion: %v", err)
	}

	// Check if the data matches
	if fetchedUser.EmailAddress != user.EmailAddress {
		t.Errorf("Expected email %v, but got %v", user.EmailAddress, fetchedUser.EmailAddress)
	}

	t.Log("User inserted and validated successfully")
}

func TestUserLogin(t *testing.T) {
	// Get the DB instance
	dbInstance, err := getDBInstance()
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}

	// Ensure the table exists
	if err := dbInstance.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Error migrating the User table: %v", err)
	}

	// Check if the test user exists
	var fetchedUser models.User
	email := "testuser@example.com"
	if err := dbInstance.First(&fetchedUser, "email_address = ?", email).Error; err != nil {
		t.Fatalf("User not found for login test: %v", err)
	}

	// Simulate login by validating email and password
	if fetchedUser.EmailAddress == email && fetchedUser.Password == "password123" {
		t.Log("User login successful")
	} else {
		t.Errorf("User login failed for email: %v", email)
	}
}
