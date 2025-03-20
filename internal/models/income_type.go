package models

type IncomeType struct {
	IncomeTypeID      int    `json:"income_type_id"`
	IncomeTypeName    string `json:"income_type_name"`
	IncomeDescription string `json:"income_description"`
	CreatedAt         string `json:"created_at"`
}
