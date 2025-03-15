package handlers

import (
	"encoding/json"
	"finanapp/internal/db"
	"finanapp/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// GetGroups API returns a list of Groups in the database.
func GetGroups(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Create a variable to store Groups
	var groupType []models.GroupType
	var userGroup []models.UserGroup
	var userIncome []models.UserIncome
	var userExpense []models.UserExpenditure

	// Fetch all group types from the database
	if result := db.DB.Find(&groupType); result.Error != nil {
		http.Error(w, "Failed to retrieve group types data", http.StatusInternalServerError)
		return
	}

	// Fetch user group from the database
	if result := db.DB.
		Preload("GroupType").
		Preload("User").
		Preload("UserGroupIncomes.UserIncome").
		Preload("UserGroupExpenditures.UserExpenditure").
		Preload("GroupMembers.User").
		Preload("GroupMembers.UserRole").
		Preload("GroupInvites").
		Where("user_id = ?", user.ID).
		Find(&userGroup); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve user group data", http.StatusInternalServerError)
		return
	}

	// Fetch all user shared income
	if result := db.DB.
		Where("shared_income = ? AND user_id = ?", true, user.ID).
		Find(&userIncome); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve shared income data", http.StatusInternalServerError)
		return
	}

	// Fetch all user shared expense
	if result := db.DB.
		Where("shared_expenditure = ? AND user_id = ?", true, user.ID).
		Find(&userExpense); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to retrieve shared expense data", http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"userGroup":   userGroup,
		"userIncome":  userIncome,
		"userExpense": userExpense,
		"groupType":   groupType,
	}

	// Set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the data as JSON
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}
}

// CreateGroup creates a new User Group
func CreateGroup(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		//log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the received body for debugging
	log.Printf("Received body: %s", string(body))

	// Define the structure to match the incoming JSON
	var groupData struct {
		GroupName   string `json:"GroupName"`
		GroupTypeID string `json:"GroupTypeID"`
		UserID      uint   `json:"UserID"`
	}

	// Unmarshal the JSON into assetData
	if err := json.Unmarshal(body, &groupData); err != nil {
		//log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validation of required fields
	if groupData.GroupName == "" || groupData.GroupTypeID == "" || groupData.UserID < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	// Convert GroupTypeID from string to uint
	groupType, err := strconv.ParseUint(groupData.GroupTypeID, 10, 32)
	if err != nil {
		log.Printf("Error converting GroupTypeID: %v", err)
		http.Error(w, "Error converting GroupTypeID", http.StatusBadRequest)
		return
	}

	// Create the UserGroup model
	group := models.UserGroup{
		GroupName:   groupData.GroupName,
		GroupTypeID: uint(groupType),
		UserID:      groupData.UserID,
	}

	// Save to the database using GORM
	if err := db.DB.Create(&group).Error; err != nil {
		log.Printf("Error saving to database: %v", err)

		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "Error 1062") {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // HTTP 409 - Conflict
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Group name already exists",
				"message": "A group with this name already exists. Please choose a different name.",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Database error",
			"message": "Error saving to database",
		})
		return
	}

	// load GroupType to feed the html table with a new row
	if err := db.DB.Preload("GroupType").First(&group, group.ID).Error; err != nil {
		log.Printf("Error fetching group with GroupType: %v", err)
		http.Error(w, "Error fetching group", http.StatusInternalServerError)
		return
	}

	// Send response back with created Group data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Group created successfully!",
		"group":   group,
	})

}

// CreateGroupItem creates a new Group Item
func CreateGroupItem(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Define the structure to match the incoming JSON
	var groupData struct {
		GroupItemSelected string `json:"GroupItemSelected"`
		GroupID           string `json:"GroupID"`
	}

	// Unmarshal the JSON into groupData
	if err := json.Unmarshal(body, &groupData); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validation of required fields
	if groupData.GroupItemSelected == "" || groupData.GroupID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	// Check if GroupItemSelected has the correct prefix and extract the ID
	var groupItemType string
	var itemID string

	if strings.HasPrefix(groupData.GroupItemSelected, "income_") {
		groupItemType = "income"
		itemID = strings.TrimPrefix(groupData.GroupItemSelected, "income_")
	} else if strings.HasPrefix(groupData.GroupItemSelected, "expense_") {
		groupItemType = "expense"
		itemID = strings.TrimPrefix(groupData.GroupItemSelected, "expense_")
	} else {
		http.Error(w, "Invalid GroupItemSelected type", http.StatusBadRequest)
		return
	}

	// Convert GroupID to uint
	groupID, err := strconv.ParseUint(groupData.GroupID, 10, 32)
	if err != nil {
		http.Error(w, "Error converting GroupID", http.StatusBadRequest)
		return
	}
	// Convert Item to uint
	groupItemID, err := strconv.ParseUint(itemID, 10, 32)
	if err != nil {
		http.Error(w, "Error converting ItemID", http.StatusBadRequest)
		return
	}

	// Based on the groupItemType, decide which model to create
	if groupItemType == "income" {
		// Create UserGroupIncome
		groupItem := models.UserGroupIncome{
			UserIncomeID: uint(groupItemID),
			GroupID:      uint(groupID),
		}

		// Save to the database using GORM
		if err := db.DB.Create(&groupItem).Error; err != nil {
			log.Printf("Error saving income to database: %v", err)

			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "UNIQUE constraint failed") ||
				strings.Contains(err.Error(), "Error 1062") {

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict) // HTTP 409 - Conflict
				json.NewEncoder(w).Encode(map[string]string{
					"error":   "Duplicate entry",
					"message": "This income item already exists in the group.",
				})
				return
			}

			// Retorna erro genérico se não for duplicidade
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Database error",
				"message": "Error saving income to database",
			})
			return
		}

		// Send response back with the created Group data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Group income created successfully!",
			"group":   groupItem,
		})

	} else if groupItemType == "expense" {
		// Create UserGroupExpenditure
		groupItem := models.UserGroupExpenditure{
			UserExpenditureID: uint(groupItemID),
			GroupID:           uint(groupID),
		}

		// Save to the database using GORM
		if err := db.DB.Create(&groupItem).Error; err != nil {
			log.Printf("Error saving expenditure to database: %v", err)

			// Verifica se o erro é de violação de unicidade
			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "UNIQUE constraint failed") ||
				strings.Contains(err.Error(), "Error 1062") {

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict) // HTTP 409 - Conflict
				json.NewEncoder(w).Encode(map[string]string{
					"error":   "Duplicate entry",
					"message": "This expenditure item already exists in the group.",
				})
				return
			}

			// Retorna erro genérico se não for duplicidade
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Database error",
				"message": "Error saving expenditure to database",
			})
			return
		}

		// Send response back with the created Group data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Group expenditure created successfully!",
			"group":   groupItem,
		})
	}

}

// CreateGroupInvite creates a new Group Invite
func CreateGroupInvite(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Define the structure to match the incoming JSON
	var inviteData struct {
		GroupID     string `json:"GroupID"`
		InviteEmail string `json:"InviteEmail"`
	}

	// Unmarshal the JSON into inviteData
	if err := json.Unmarshal(body, &inviteData); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validation of required fields
	if inviteData.InviteEmail == "" || inviteData.GroupID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Required fields not filled"})
		return
	}

	// Convert GroupID to uint
	groupID, err := strconv.ParseUint(inviteData.GroupID, 10, 32)
	if err != nil {
		http.Error(w, "Error converting GroupID", http.StatusBadRequest)
		return
	}

	// Find user by email
	var user models.User
	if err := db.DB.Where("email_address = ?", inviteData.InviteEmail).First(&user).Error; err != nil {
		// If user not found, create a new GroupInvite
		groupInvite := models.GroupInvite{
			GroupID:     uint(groupID),
			InviteEmail: inviteData.InviteEmail,
		}

		if err := db.DB.Create(&groupInvite).Error; err != nil {
			http.Error(w, "Error creating group invite", http.StatusInternalServerError)
			return
		}

		// Return response if user not found
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User not found, We've sent an Invite by Email",
		})
		return
	}

	// If user found, create GroupMember record
	groupMember := models.GroupMember{
		GroupID:    uint(groupID),
		UserID:     user.ID,
		UserRoleID: 3,     // Default UserRoleID
		Active:     false, // Default is false
	}

	if err := db.DB.Create(&groupMember).Error; err != nil {
		http.Error(w, "Error adding user to group", http.StatusInternalServerError)
		return
	}

	// Return response if user added as group member
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User Invite Sent",
	})
}

// DeleteGroup deletes an existing Group in the database.
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Define the structure to match the incoming JSON
	var requestData struct {
		GroupID uint `json:"groupID"`
	}

	// Unmarshal the JSON into requestData
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validate if GroupID was provided
	if requestData.GroupID < 1 {
		http.Error(w, "Missing group ID in payload", http.StatusBadRequest)
		return
	}

	// Delete from the database using GORM
	result := db.DB.Delete(&models.UserGroup{}, requestData.GroupID)
	if result.Error != nil {

		if strings.Contains(result.Error.Error(), "violates foreign") {

			http.Error(w, "You must delete the items before deleting group", http.StatusInternalServerError)
		}

		log.Printf("Error deleting group: %v", result.Error)
		http.Error(w, "Error deleting group", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group deleted successfully!"})
}
