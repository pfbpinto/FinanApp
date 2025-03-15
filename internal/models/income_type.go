package models

type IncomeType struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	IncomeTypeName string `gorm:"size:100;not null;unique"` // Example: "Salary", "Rental", etc.
}

func (IncomeType) TableName() string {
	return "income_types"
}
