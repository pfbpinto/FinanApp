package models

type FinancialUserItem struct {
	FinancialUserItemID       int    `json:"financialUserItemId"`
	FinancialUserItemName     string `json:"financialUserItemName"`
	EntityID                  int    `json:"entityId"`
	UserEntityID              int    `json:"userEntityId"`
	RecurrencyID              string `json:"recurrencyId"`
	FinancialUserEntityItemID int    `json:"financialUserEntityItemId"`
	ParentFinancialUserItemID int    `json:"parentFinancialUserItemId"`
	IsActive                  bool   `json:"isActive"`
	CreatedAt                 string `json:"createdAt"`

	// Relation
	EntityType     string `json:"entityType"`
	RecurrencyName string `json:"recurrencyName"`
	IncomeTypeName string `json:"incomeTypeName"`

	// Aditional field for the Create function]
	Amount     string `json:"amount"`
	CurrencyID string `json:"currencyId"`
}
