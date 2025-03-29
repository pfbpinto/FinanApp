package models

type UserFinancialActual struct {
	UserFinancialActualID         int     `json:"UserFinancialActualID"`
	UserCategoryID                int     `json:"UserCategoryID"`
	FinancialUserItemID           int     `json:"FinancialUserItemID"`
	UserFinancialActualtBeginDate string  `json:"UserFinancialActualtBeginDate"`
	UserFinancialActualEndDate    *string `json:"UserFinancialActualEndDate,omitempty"`
	UserFinancialActualAmount     float64 `json:"UserFinancialActualAmount"`
	CurrencyID                    int     `json:"CurrencyID"`
	UserCategoryName              string  `json:"UserCategoryName"`
	FinancialUserItemName         string  `json:"FinancialUserItemName"`
	CurrencyName                  string  `json:"CurrencyName"`
	Note                          *string `json:"Note,omitempty"`
	CreatedAt                     string  `json:"CreatedAt"`
}
