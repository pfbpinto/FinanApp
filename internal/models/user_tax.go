package models

type UserTax struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	UserIncomeID uint    `gorm:"not null"`
	TaxID        uint    `gorm:"not null"`
	TaxValue     float64 `gorm:"type:decimal(10,2);not null"`
	Payed        bool    `gorm:"default:false"`

	UserIncome UserIncome `gorm:"foreignKey:UserIncomeID;references:ID;constraint:onDelete:CASCADE"`
	Tax        Tax        `gorm:"foreignKey:TaxID;references:ID;constraint:onDelete:CASCADE"`
}

func (UserTax) TableName() string {
	return "user_taxes"
}
