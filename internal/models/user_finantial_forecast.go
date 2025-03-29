package models

// UserFinancialForecast representa a estrutura da tabela
type UserFinancialForecast struct {
	UserFinancialForecastID        int     `json:"UserFinancialForecastID"`
	UserCategoryID                 *int    `json:"UserCategoryID"`
	FinancialUserItemID            int     `json:"FinancialUserItemID"`
	UserFinancialForecastAmount    float64 `json:"UserFinancialForecastAmount"`
	UserFinancialForecastBeginDate string  `json:"UserFinancialForecastBeginDate"`
	UserFinancialForecastEndDate   *string `json:"UserFinancialForecastEndDate,omitempty"`
	CurrencyID                     int     `json:"CurrencyID"`
	UserCategoryName               *string `json:"UserCategoryName"`
	FinancialUserItemName          string  `json:"FinancialUserItemName"`
	CurrencyName                   string  `json:"CurrencyName"`
	CreatedAt                      string  `json:"CreatedAt"`
}
