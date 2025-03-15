package db

import (
	"crypto/md5"
	"encoding/hex"
	"finanapp/internal/models"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// CreateSchemaMigrationTable ensures the 'schema_migrations' table exists in the database.
func CreateSchemaMigrationTable(db *gorm.DB) {
	// Check if the 'schema_migrations' table exists
	if !db.Migrator().HasTable(&models.SchemaMigration{}) {
		// Create the table
		if err := db.Migrator().CreateTable(&models.SchemaMigration{}); err != nil {
			log.Fatalf("Error creating 'schema_migrations' table: %v", err)
		} else {
			log.Println("Table 'schema_migrations' created successfully.")
		}
	}
}

// GenerateVersion dynamically generates a version string based on the models provided.
// This ensures that any change to the models automatically updates the migration version.
func GenerateVersion(models []interface{}) string {
	data := ""
	for _, model := range models {
		data += fmt.Sprintf("%T", model)
	}
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])[:8] // Return the first 8 characters of the hash
}

// RunMigrations executes the migrations for the database, ensuring new models are migrated
// and the schema version is recorded in the 'schema_migrations' table.
func RunMigrations(db *gorm.DB) {
	// Ensure the schema_migrations table exists
	CreateSchemaMigrationTable(db)

	// Begin transaction
	tx := db.Begin()

	// Define all models that need migration
	modelsToMigrate := []interface{}{
		// Base types and foundational tables
		&models.UserType{},
		&models.IncomeType{},
		&models.Subscription{},
		&models.TaxType{},
		&models.Tax{},
		&models.AssetType{},
		&models.FileType{},
		&models.PersonalDocumentType{},

		// Primary entities
		&models.User{},
		&models.UserIncome{},
		&models.ExpenditureType{},

		// Relationships and dependent entities
		&models.UserAsset{},
		&models.UserIncome{},
		&models.UserExpenditure{},
		&models.UserTax{},
		&models.UserAssetTax{},
		&models.UserSubscription{},
		&models.UserDocument{},
		&models.UserFiles{},
		&models.UserRole{},
		&models.UserGroup{},
		&models.UserGroupIncome{},
		&models.UserGroupExpenditure{},
		&models.UserGroupAsset{},
		&models.GroupMember{},
		&models.GroupInvite{},
		&models.GroupFiles{},
	}

	// Generate the new version dynamically based on models
	newVersion := GenerateVersion(modelsToMigrate)

	// Get the last applied migration
	var lastMigration models.SchemaMigration
	err := tx.Order("applied_at desc").First(&lastMigration).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		log.Fatalf("Error fetching last migration: %v", err)
	}

	// Check if the new version matches the last applied migration
	if lastMigration.Version == newVersion {
		log.Println("No new migrations to apply. Database is up to date.")
		tx.Rollback()
		return
	}

	// Apply new migrations
	for _, model := range modelsToMigrate {
		err := tx.AutoMigrate(model)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Error migrating table %T: %v", model, err)
		}
	}

	if err := CreateStoredProcedures(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Error creating stored procedures: %v", err)
	}
	log.Println("Store Procedures successfully created.")

	// Record the new migration version in the database
	newMigration := models.SchemaMigration{
		Version:   newVersion,
		AppliedAt: time.Now(),
	}
	if err := tx.Create(&newMigration).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Error recording migration version: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error committing migrations: %v", err)
	}

	log.Printf("Migrations applied successfully. Current version: %s\n", newVersion)
}

func CreateStoredProcedures(db *gorm.DB) error {
	procedures := []string{
		`CREATE OR REPLACE FUNCTION CreateUser(
			FirstName VARCHAR(255),
			LastName VARCHAR(255),
			EmailAddress VARCHAR(255),
			Password VARCHAR(255),
			DataOfBirth VARCHAR(255)  -- Date as string
		)
		RETURNS JSONB AS $$  -- Return JSONB instead of a string
		BEGIN
			-- Validation
			IF FirstName IS NULL OR LENGTH(FirstName) = 0 THEN
				RETURN jsonb_build_object('status', 'error', 'message', 'First Name is mandatory.');
			END IF;

			IF LastName IS NULL OR LENGTH(LastName) = 0 THEN
				RETURN jsonb_build_object('status', 'error', 'message', 'Last Name is mandatory.');
			END IF;

			IF EmailAddress IS NULL OR LENGTH(EmailAddress) = 0 THEN
				RETURN jsonb_build_object('status', 'error', 'message', 'Email Address is mandatory.');
			END IF;

			IF POSITION('@' IN EmailAddress) = 0 THEN
				RETURN jsonb_build_object('status', 'error', 'message', 'Invalid email address.');
			END IF;

			-- Check if the email already exists in the system
			IF EXISTS (SELECT 1 FROM "user" WHERE email_address = EmailAddress) THEN
				RETURN jsonb_build_object('status', 'error', 'message', 'A user with this email already exists.');
			END IF;

			-- Insert the user into the database
			INSERT INTO "user" (first_name, last_name, user_type_id, email_address, password, data_of_birth, created_at, updated_at)
			VALUES (
				FirstName,
				LastName,
				2,  -- Assuming 2 is the default user_type_id
				EmailAddress,
				Password,
				TO_DATE(DataOfBirth, 'YYYY-MM-DD'),  -- Convert DataOfBirth from string to date
				CURRENT_TIMESTAMP,
				CURRENT_TIMESTAMP
			);

			-- Return success message after creating the user
			RETURN jsonb_build_object('status', 'success', 'message', 'User created successfully.');
		END;
		$$ LANGUAGE plpgsql;`,

		`CREATE OR REPLACE FUNCTION UpdateUser(
			UserID INT,
			FirstName VARCHAR(255),
			LastName VARCHAR(255),
			EmailAddress VARCHAR(255),
			Password VARCHAR(255),
			DataOfBirth VARCHAR(255)
		)
		RETURNS TEXT AS $$  -- Retorna um texto com a mensagem de erro ou sucesso
		BEGIN
			-- Validation
			IF NOT EXISTS (SELECT 1 FROM "user" WHERE id = UserID) THEN
				RETURN 'User not found.';
			END IF;

			IF FirstName IS NULL OR LENGTH(FirstName) = 0 THEN
				RETURN 'First Name is mandatory.';
			END IF;

			IF LastName IS NULL OR LENGTH(LastName) = 0 THEN
				RETURN 'Last Name is mandatory.';
			END IF;

			-- Check if email is different and exists
			IF EmailAddress IS NOT NULL AND EmailAddress <> (SELECT email_address FROM "user" WHERE id = UserID) THEN
				RETURN 'Email address cannot be updated.';
			END IF;

			IF EXISTS (SELECT 1 FROM "user" WHERE email_address = EmailAddress AND id <> UserID) THEN
				RETURN 'A user with this email already exists.';
			END IF;

			-- Update user information
			UPDATE "user"
			SET first_name = FirstName, last_name = LastName, 
				data_of_birth = DataOfBirth, updated_at = CURRENT_TIMESTAMP
			WHERE id = UserID;

			RETURN 'User successfully updated';  -- Success message
		END;
		$$ LANGUAGE plpgsql;`,
	}

	for _, procedure := range procedures {
		if err := db.Exec(procedure).Error; err != nil {
			return err
		}
	}
	return nil
}
