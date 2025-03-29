package models

type Recurrency struct {
	RecurrencyID     int    `json:"recurrency_id"`
	RecurrencyName   string `json:"recurrency_name"`
	RecurrencyPeriod string `json:"recurrency_period"`
	CreatedAt        string `json:"created_at"`
}
