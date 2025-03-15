package models

type UserDocument struct {
	ID                     uint   `gorm:"primaryKey;autoIncrement"`
	PersonalDocumentTypeID uint   `gorm:"not null"`
	UserID                 uint   `gorm:"not null"`
	PersonalIdentifier     string `gorm:"size:150;not null"`

	User                 User                 `gorm:"foreignKey:UserID"`
	PersonalDocumentType PersonalDocumentType `gorm:"foreignKey:PersonalDocumentTypeID"`
}
