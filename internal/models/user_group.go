package models

type UserGroup struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	GroupTypeID uint   `gorm:"not null;uniqueIndex:idx_user_group_name"`
	UserID      uint   `gorm:"not null;uniqueIndex:idx_user_group_name"`
	GroupName   string `gorm:"size:150;not null;uniqueIndex:idx_user_group_name"`

	GroupType             GroupType              `gorm:"foreignKey:GroupTypeID"`
	User                  User                   `gorm:"foreignKey:UserID"`
	UserGroupIncomes      []UserGroupIncome      `gorm:"foreignKey:GroupID"`
	UserGroupExpenditures []UserGroupExpenditure `gorm:"foreignKey:GroupID"`
	GroupMembers          []GroupMember          `gorm:"foreignKey:GroupID"`
	GroupInvites          []GroupInvite          `gorm:"foreignKey:GroupID"`
}

func (UserGroup) TableName() string {
	return "user_groups"
}
