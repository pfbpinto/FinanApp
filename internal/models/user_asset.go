package models

type UserAsset struct {
	UserAssetID                   int    `json:"userAssetId"`
	UserAssetName                 string `json:"userAssetName"`
	UserAssetValueAmount          string `json:"userAssetValueAmount"`
	UserAssetAcquisitionBeginDate string `json:"acquisitionBeginDate"`
	UserAssetAcquisitionEndDate   string `json:"acquisitionEndDate"`
	IsActive                      bool   `json:"isActive"`
	CreatedAt                     string `json:"createdAt"`

	// Relation names
	AssetTypeName   string `json:"assetTypeName"`
	UserProfileName string `json:"userProfileName"`

	// Optional fields for Create/Update, se precisar
	AssetTypeID   int `json:"assetTypeId,omitempty"`
	UserProfileID int `json:"userProfileId,omitempty"`
}

type CreateUserAssetParentIncome struct {
	UserID                    int     `json:"user_id"`
	UserAssetID               int     `json:"user_asset_id"`
	FinancialUserItemName     string  `json:"financial_user_item_name"`
	RecurrencyID              int     `json:"recurrency_id"`
	FinancialUserEntityItemID int     `json:"financial_user_entity_item_id"`
	ParentIncomeAmount        float64 `json:"parent_income_amount"`
	BeginDate                 string  `json:"begin_date"`
}

type CreateUserAssetChildIncomeTax struct {
	UserAssetID               int     `json:"user_asset_id"`
	FinancialUserItemName     string  `json:"financial_user_item_name"`
	FinancialUserEntityItemID int     `json:"financial_user_entity_item_id"`
	ParentFinancialUserItemID int     `json:"parent_financial_user_item_id"`
	TaxIncomeAmount           float64 `json:"tax_income_amount"`
}

type CreateUserAssetChildIncomeExpense struct {
	UserID                    int     `json:"user_id"`
	UserAssetID               int     `json:"user_asset_id"`
	FinancialUserItemName     string  `json:"financial_user_item_name"`
	FinancialUserEntityItemID int     `json:"financial_user_entity_item_id"`
	ParentFinancialUserItemID int     `json:"parent_financial_user_item_id"`
	ExpenseAmount             float64 `json:"expense_amount"`
}
