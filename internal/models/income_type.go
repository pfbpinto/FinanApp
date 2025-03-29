package models

type IncomeType struct {
	IncomeTypeID      int    `json:"income_type_id"`
	IncomeTypeName    string `json:"income_type_name"`
	IncomeDescription string `json:"income_description"`
	EntityID          int    `json:"entity_id"`
	CreatedAt         string `json:"created_at"`
}
