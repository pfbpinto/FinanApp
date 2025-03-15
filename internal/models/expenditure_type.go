package models

type ExpenditureType struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement"`
	ExpenditureTypeName string `gorm:"size:100;not null;unique"`
}

func (ExpenditureType) TableName() string {
	return "expenditures"
}
