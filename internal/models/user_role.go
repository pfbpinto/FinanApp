package models

type UserRole struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	UserRoleName   string `gorm:"size:100;not null;unique"`
	ViewPermission bool   `gorm:"default:false"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
