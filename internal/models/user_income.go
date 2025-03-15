package models

import (
	"time"
)

type UserIncome struct {
	ID               uint       `gorm:"primaryKey;autoIncrement"`
	UserID           uint       `gorm:"not null;uniqueIndex:idx_user_income_name"`
	IncomeTypeID     uint       `gorm:"not null;uniqueIndex:idx_user_income_name"`
	IncomeName       string     `gorm:"size:150;not null;uniqueIndex:idx_user_income_name"`
	IncomeValue      float64    `gorm:"type:decimal(10,2);not null"`
	IncomeRecurrence string     `gorm:"size:100;not null"`
	IncomeStartDate  *time.Time `gorm:"not null"`
	IncomeEndDate    *time.Time
	SharedIncome     bool    `gorm:"default:false"`
	OwningPercentage float64 `gorm:"type:decimal(5,2);not null;default:100.00"`

	User       User       `gorm:"foreignKey:UserID;constraint:onDelete:CASCADE"`
	IncomeType IncomeType `gorm:"foreignKey:IncomeTypeID"`
	UserTaxes  []UserTax  `gorm:"foreignKey:UserIncomeID;constraint:onDelete:CASCADE"`
}

func (UserIncome) TableName() string {
	return "user_incomes"
}
