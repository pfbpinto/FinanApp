package models

type UserType struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:255;not null;unique"`
	IsActive bool   `gorm:"default:true"`
}

func (UserType) TableName() string {
	return "user_types"
}
