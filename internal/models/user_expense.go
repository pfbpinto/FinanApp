package models

type UserParentExpense struct {
	FinancialUserItemName     string `json:"financialUserItemName"`
	RecurrencyID              int    `json:"recurrencyId"`
	FinancialUserEntityItemID int    `json:"financialUserEntityItemId"`
	ParentExpenseAmount       string `json:"parentExpenseAmount"` // como string vindo do front
	BeginDate                 string `json:"beginDate"`           // "YYYY-MM-DD"
}

type UserParentExpenseUpdate struct {
	FinancialUserItemID      int    `json:"financial_user_item_id"`
	NewFinancialUserItemName string `json:"financial_user_item_name"`
	NewParentExpenseAmount   string `json:"parent_expense_amount"`
	NewBeginDate             string `json:"begin_date"`
	IsActive                 bool   `json:"is_active"`
}
