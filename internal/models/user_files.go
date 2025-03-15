package models

import "time"

type UserFiles struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	FileTypeID    uint   `gorm:"not null"`
	UserID        uint   `gorm:"not null"`
	FileName      string `gorm:"size:150;not null;unique"`
	FileDate      *time.Time
	FileBeginDate *time.Time
	FileEndDate   *time.Time

	User     User     `gorm:"foreignKey:UserID"`
	FileType FileType `gorm:"foreignKey:FileTypeID"`
}

func (UserFiles) TableName() string {
	return "user_files"
}
