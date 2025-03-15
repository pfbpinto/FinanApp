package models

type GroupMember struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	GroupID    uint `gorm:"not null;index:idx_group_user,unique"`
	UserID     uint `gorm:"not null;index:idx_group_user,unique"`
	UserRoleID uint `gorm:"not null"`
	Active     bool `gorm:"default:false"`

	// Relationships
	User      User      `gorm:"foreignKey:UserID"`
	UserRole  UserRole  `gorm:"foreignKey:UserRoleID"`
	UserGroup UserGroup `gorm:"foreignKey:GroupID"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
