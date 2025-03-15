package models

type UserGroupExpenditure struct {
	ID                uint `gorm:"primaryKey;autoIncrement"`
	UserExpenditureID uint `gorm:"not null;uniqueIndex:idx_user_expenditure_group"`
	GroupID           uint `gorm:"not null;uniqueIndex:idx_user_expenditure_group"`

	UserExpenditure UserExpenditure `gorm:"foreignKey:UserExpenditureID"`
	UserGroup       UserGroup       `gorm:"foreignKey:GroupID"`
}

func (UserGroupExpenditure) TableName() string {
	return "user_group_expenditure"
}
