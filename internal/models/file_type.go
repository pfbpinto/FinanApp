package models

type FileType struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	FileTypeName string `gorm:"size:100;not null;unique"`
}

func (FileType) TableName() string {
	return "file_types"
}
