package models

import "time"

type UserExpenditure struct {
	ID                    uint    `gorm:"primaryKey;autoIncrement"`
	UserID                uint    `gorm:"not null;uniqueIndex:idx_user_expenditure"`
	ExpenditureID         uint    `gorm:"not null;uniqueIndex:idx_user_expenditure"`
	ExpenditureName       string  `gorm:"size:150;not null;uniqueIndex:idx_user_expenditure"`
	ExpenditureValue      float64 `gorm:"type:decimal(10,2);not null"`
	ExpenditureStartDate  *time.Time
	ExpenditureEndDate    *time.Time
	ExpenditureRecurrence string `gorm:"size:100;not null"`
	SharedExpenditure     bool   `gorm:"default:false"`

	// Relationships
	User        User            `gorm:"foreignKey:UserID;constraint:onDelete:CASCADE"`
	Expenditure ExpenditureType `gorm:"foreignKey:ExpenditureID;constraint:onDelete:CASCADE"`
}

func (UserExpenditure) TableName() string {
	return "user_expenditures"
}
