package models

type Currency struct {
	CurrencyID          int    `json:"currency_id"`
	CurrencyName        string `json:"currency_name"`
	CurrencyAbreviation string `json:"currency_abreviation"`
	CurrencySymbol      string `json:"currency_symbol"`
	CreatedAt           string `json:"created_at"`
}
