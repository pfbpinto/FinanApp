package models

type PersonalDocumentType struct {
	ID                   uint   `gorm:"primaryKey;autoIncrement"`
	PersonalDocumentName string `gorm:"size:150;not null;unique"`
	DocumentCountry      string `gorm:"size:100"`
}

func (PersonalDocumentType) TableName() string {
	return "personal_document_types"
}
