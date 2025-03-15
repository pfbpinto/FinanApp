package models

type UserGroupIncome struct {
	ID           uint `gorm:"primaryKey;autoIncrement"`
	UserIncomeID uint `gorm:"not null;uniqueIndex:idx_user_income_group"`
	GroupID      uint `gorm:"not null;uniqueIndex:idx_user_income_group"`

	UserIncome UserIncome `gorm:"foreignKey:UserIncomeID"`
	UserGroup  UserGroup  `gorm:"foreignKey:GroupID"`
}

func (UserGroupIncome) TableName() string {
	return "user_group_incomes"
}
