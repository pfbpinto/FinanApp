package models

import "time"

// UserFinancialForecast representa a estrutura da tabela
type UserFinancialForecast struct {
	UserFinancialForecastID        int        `json:"user_financial_forecast_id"`
	UserFinancialForecastName      string     `json:"user_financial_forecast_name"`
	UserCategoryID                 int        `json:"user_category_id"`
	EntityTypeName                 string     `json:"entity_type_name"`
	EntityTypeID                   int        `json:"entity_type_id"`
	EntityItemTypeName             string     `json:"entity_item_type_name"`
	EntityItemTypeID               int        `json:"entity_item_type_id"`
	UserFinancialForecastAmount    float64    `json:"user_financial_forecast_amount"`
	UserFinancialForecastBeginDate time.Time  `json:"user_financial_forecast_begin_date"`
	UserFinancialForecastEndDate   *time.Time `json:"user_financial_forecast_end_date,omitempty"`
	CurrencyID                     int        `json:"currency_id"`
	CreatedAt                      time.Time  `json:"created_at"`
}
