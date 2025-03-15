package models

type GroupType struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	GroupTypeName string `gorm:"size:100;not null;unique"`
}

func (GroupType) TableName() string {
	return "group_types"
}
