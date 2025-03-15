package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	UserTypeID   uint           `gorm:"not null;index"`
	Password     string         `gorm:"size:255;not null"`
	FirstName    string         `gorm:"size:255;not null"`
	LastName     string         `gorm:"size:255;not null"`
	EmailAddress string         `gorm:"size:255;not null;unique"`
	DataOfBirth  *time.Time     `gorm:"type:date"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	IsActive     bool           `gorm:"default:true"`
	LastLogin    *time.Time     `gorm:"type:timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	UserType UserType `gorm:"foreignKey:UserTypeID"`
}

func (User) TableName() string {
	return "user"
}
