package models

type UserCategory struct {
	UserCategoryID   int    `json:"user_category_id"`
	UserCategoryName string `json:"user_category_name"`
	ItemTypeName     string `json:"item_type_name"`
	ItemTypeNameID   int    `json:"item_type_name_id"`
	CreatedAt        string `json:"created_at"`
}
