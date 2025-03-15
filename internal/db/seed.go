package db

import (
	"finanapp/internal/models"
	"finanapp/internal/utils"
	"log"

	"gorm.io/gorm"
)

// SeedDatabase populates the database with initial data
func SeedDatabase(db *gorm.DB) {
	// Seeding user types (admin, user)
	userTypes := []models.UserType{
		{
			Name: "admin",
		},
		{
			Name: "user",
		},
	}

	// Inserting userType if not exists
	for _, userType := range userTypes {
		var count int64
		db.Model(&models.UserType{}).Where("name = ?", userType.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&userType).Error; err != nil {
				log.Printf("Error creating user type %s: %v", userType.Name, err)
			}
			log.Println("User type seeding completed successfully.")
		} else {
			log.Printf("User type %s already exists.", userType.Name)
		}
	}

	// Seeding income types
	incomeTypes := []models.IncomeType{
		{
			IncomeTypeName: "Salary",
		},
		{
			IncomeTypeName: "Investiments",
		},
	}

	// Inserting incomeType if not exists
	for _, incomeType := range incomeTypes {
		var count int64
		db.Model(&models.IncomeType{}).Where("income_type_name = ?", incomeType.IncomeTypeName).Count(&count)
		if count == 0 {
			if err := db.Create(&incomeType).Error; err != nil {
				log.Printf("Error creating tax type %s: %v", incomeType.IncomeTypeName, err)
			}
			log.Println("Income type seeding completed successfully.")
		} else {
			log.Printf("Income type %s already exists.", incomeType.IncomeTypeName)
		}
	}

	// Seeding expenditure if not exists
	expenditureTypes := []models.ExpenditureType{
		{
			ExpenditureTypeName: "General",
		},
		{
			ExpenditureTypeName: "HomeCare",
		},
	}

	// Inserting expenditure if not exists
	for _, expenditureType := range expenditureTypes {
		var count int64
		db.Model(&models.ExpenditureType{}).Where("expenditure_type_name = ?", expenditureType.ExpenditureTypeName).Count(&count)
		if count == 0 {
			if err := db.Create(&expenditureType).Error; err != nil {
				log.Printf("Error creating expenditure type %s: %v", expenditureType.ExpenditureTypeName, err)
			}
			log.Println("Expenditure type seeding completed successfully.")
		} else {
			log.Printf("Expenditure type %s already exists.", expenditureType.ExpenditureTypeName)
		}
	}

	// Seeding tax types (assets, user)
	taxTypes := []models.TaxType{
		{
			TaxTypeName: "Asset",
		},
		{
			TaxTypeName: "Income",
		},
	}

	// Inserting taxType if not exists
	for _, taxType := range taxTypes {
		var count int64
		db.Model(&models.TaxType{}).Where("tax_type_name = ?", taxType.TaxTypeName).Count(&count)
		if count == 0 {
			if err := db.Create(&taxType).Error; err != nil {
				log.Printf("Error creating tax type %s: %v", taxType.TaxTypeName, err)
			}
			log.Println("Tax type seeding completed successfully.")
		} else {
			log.Printf("Tax type %s already exists.", taxType.TaxTypeName)
		}
	}

	// Sample admin user for seeding
	users := []models.User{
		{
			FirstName:    "admin",
			UserTypeID:   1, // Admin type
			EmailAddress: "admin@example.com",
			Password:     hashPasswordSafely("password123"),
			IsActive:     true,
		},
		{
			FirstName:    "user",
			UserTypeID:   2, // Regular user type
			EmailAddress: "user@example.com",
			Password:     hashPasswordSafely("password123"),
			IsActive:     true,
		},
	}

	// Checking if the user already exists before adding
	for _, user := range users {
		var count int64
		// Checking if the user already exists
		db.Model(&models.User{}).Where("email_address = ?", user.EmailAddress).Count(&count)
		if count == 0 {
			// If not, create the user
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Error creating user %s: %v", user.FirstName, err)
			}
			log.Println("User seeding completed successfully.")
		} else {
			log.Printf("User %s already exists.", user.FirstName)
		}
	}

	// Seeding taxes
	taxes := []models.Tax{
		{
			UserID:             1,
			TaxName:            "ICMS",
			TaxTypeID:          1,
			TaxPercentage:      18.00,
			TaxPercentageRange: "0-100%",
			TaxApplicableCycle: "Mensal",
		},
		{
			UserID:             1,
			TaxName:            "IPVA",
			TaxTypeID:          2,
			TaxPercentage:      4.00,
			TaxPercentageRange: "0-100%",
			TaxApplicableCycle: "Anual",
		},
	}

	// Inserting taxes if not exists
	for _, tax := range taxes {
		var count int64
		db.Model(&models.Tax{}).Where("tax_name = ?", tax.TaxName).Count(&count)
		if count == 0 {
			if err := db.Create(&tax).Error; err != nil {
				log.Printf("Error creating tax %s: %v", tax.TaxName, err)
			}
			log.Println("Tax seeding completed successfully.")
		} else {
			log.Printf("Tax %s already exists.", tax.TaxName)
		}
	}

	// Seeding asset types
	assetTypes := []models.AssetType{
		{
			AssetTypeName: "Apartament",
		},
		{
			AssetTypeName: "Car",
		},
	}

	// Inserting assetType if not exists
	for _, assetType := range assetTypes {
		var count int64
		db.Model(&models.AssetType{}).Where("asset_type_name = ?", assetType.AssetTypeName).Count(&count)
		if count == 0 {
			if err := db.Create(&assetType).Error; err != nil {
				log.Printf("Error creating asset type %s: %v", assetType.AssetTypeName, err)
			}
			log.Println("Asset seeding completed successfully.")
		} else {
			log.Printf("Asset type %s already exists.", assetType.AssetTypeName)
		}
	}

	// Seeding group types
	groupTypes := []models.GroupType{
		{
			GroupTypeName: "Familia",
		},
		{
			GroupTypeName: "Negocios",
		},
	}

	// Inserting groupType if not exists
	for _, groupType := range groupTypes {
		var count int64
		db.Model(&models.GroupType{}).Where("group_type_name = ?", groupType.GroupTypeName).Count(&count)
		if count == 0 {
			if err := db.Create(&groupType).Error; err != nil {
				log.Printf("Error creating Group type %s: %v", groupType.GroupTypeName, err)
			}
			log.Println("Group seeding completed successfully.")
		} else {
			log.Printf("Group type %s already exists.", groupType.GroupTypeName)
		}
	}

	// Seeding user roles types
	userRoles := []models.UserRole{
		{
			UserRoleName:   "Advance",
			ViewPermission: true,
		},
		{
			UserRoleName:   "Basic",
			ViewPermission: true,
		},
		{
			UserRoleName:   "Viewer",
			ViewPermission: false,
		},
	}

	// Inserting groupType if not exists
	for _, userRole := range userRoles {
		var count int64
		db.Model(&models.UserRole{}).Where("user_role_name = ?", userRole.UserRoleName).Count(&count)
		if count == 0 {
			if err := db.Create(&userRole).Error; err != nil {
				log.Printf("Error creating user roles %s: %v", userRole.UserRoleName, err)
			}
			log.Println("user roles seeding completed successfully.")
		} else {
			log.Printf("user roles %s already exists.", userRole.UserRoleName)
		}
	}

	// Seeding user group types
	userGroups := []models.UserGroup{
		{
			GroupTypeID: 1,
			UserID:      1,
			GroupName:   "Brum Pinto",
		},
	}

	// Inserting groupType if not exists
	for _, userGroup := range userGroups {
		var count int64
		db.Model(&models.UserGroup{}).Where("group_name = ?", userGroup.GroupName).Count(&count)
		if count == 0 {
			if err := db.Create(&userGroup).Error; err != nil {
				log.Printf("Error creating user group %s: %v", userGroup.GroupName, err)
			}
			log.Println("user group seeding completed successfully.")
		} else {
			log.Printf("user group %s already exists.", userGroup.GroupName)
		}
	}

	log.Println("Seeding completed successfully.")
}

// Helper function to generate the password hash and handle the error
func hashPasswordSafely(password string) string {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Error generating password hash: %v", err)
	}
	return hashedPassword
}
