package models

type GroupInvite struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	GroupID     uint   `gorm:"not null;index:idx_group_invite,unique"`
	InviteEmail string `gorm:"size:255;not null;index:idx_group_invite,unique"`

	UserGroup UserGroup `gorm:"foreignKey:GroupID;references:ID"`
}

func (GroupInvite) TableName() string {
	return "group_invite"
}
