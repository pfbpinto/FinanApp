package models

type UserCategory struct {
	UserCategoryID             int    `json:"user_category_id"`
	UserCategoryName           string `json:"user_category_name"`
	UserProfileID              int    `json:"user_profile_id"`
	EntityID                   int    `json:"entity_id"`
	FinancialGroupEntityItemID int    `json:"financial_group_entity_item_id"`
	IsActive                   bool   `json:"is_active"`
	CreatedAt                  string `json:"created_at"`
}
