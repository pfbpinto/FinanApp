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
			DateOfBirth VARCHAR(255)  -- Date as string
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
			IF EXISTS (SELECT 1 FROM "userprofile" WHERE emailaddress = EmailAddress) THEN
				RETURN jsonb_build_object('status', 'error', 'message', 'A user with this email already exists.');
			END IF;

			-- Insert the user into the database
			INSERT INTO "userprofile" (
				firstname,
				lastname,
				emailaddress,
				userpassword,
				dateofbirth,
				usersubscription,
				createdat
			)
			VALUES (
				FirstName,
				LastName,
				EmailAddress,
				Password,  -- Consider hashing the password before inserting
				TO_DATE(DateOfBirth, 'YYYY-MM-DD'),  -- Convert DateOfBirth from string to date
				FALSE,  -- Assuming user_subscription defaults to FALSE
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
			DateOfBirth VARCHAR(255)
		)
		RETURNS TEXT AS $$  -- Retorna um texto com a mensagem de erro ou sucesso
		BEGIN
			-- Validation: Check if the user exists
			IF NOT EXISTS (SELECT 1 FROM "userprofile" WHERE userprofileid = UserID) THEN
				RETURN 'User not found.';
			END IF;

			-- Validation: Check mandatory fields
			IF FirstName IS NULL OR LENGTH(FirstName) = 0 THEN
				RETURN 'First Name is mandatory.';
			END IF;

			IF LastName IS NULL OR LENGTH(LastName) = 0 THEN
				RETURN 'Last Name is mandatory.';
			END IF;

			-- Check if email is different and exists
			IF EmailAddress IS NOT NULL AND EmailAddress <> (SELECT emailaddress FROM "userprofile" WHERE userprofileid = UserID) THEN
				-- Optional: if you want to allow email change, remove this check
				RETURN 'Email address cannot be updated.';
			END IF;

			IF EXISTS (SELECT 1 FROM "userprofile" WHERE emailaddress = EmailAddress AND userprofileid <> UserID) THEN
				RETURN 'A user with this email already exists.';
			END IF;

			-- If Password is provided, hash it (consider hashing before this step in actual code)
			IF Password IS NOT NULL AND LENGTH(Password) > 0 THEN
				-- Consider using a utility to hash the password securely
				Password := Password;  -- This should be hashed in practice, not just stored as-is
			END IF;

			-- Update user information in the database
			UPDATE "userprofile"
			SET 
				firstname = FirstName,
				lastname = LastName,
				emailaddress = EmailAddress,
				userpassword = Password,  -- Password should be hashed if necessary
				dateofbirth = TO_DATE(DateOfBirth, 'YYYY-MM-DD'),  -- Convert DateOfBirth from string to date
				updatedat = CURRENT_TIMESTAMP
			WHERE userprofileid = UserID;

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
