package models

import "time"

type GroupFiles struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	FileTypeID    uint   `gorm:"not null"`
	UserGroupID   uint   `gorm:"not null"`
	FileName      string `gorm:"size:150;not null;unique"`
	FileDate      *time.Time
	FileBeginDate *time.Time
	FileEndDate   *time.Time

	FileType  FileType  `gorm:"foreignKey:FileTypeID"`
	UserGroup UserGroup `gorm:"foreignKey:UserGroupID"`
}

func (GroupFiles) TableName() string {
	return "group_files"
}
