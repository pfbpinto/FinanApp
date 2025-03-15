package models

type TaxType struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	TaxTypeName string `gorm:"size:100;not null;unique"`
}

func (TaxType) TableName() string {
	return "taxes_types"
}
