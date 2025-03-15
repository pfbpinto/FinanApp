package models

type Tax struct {
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	UserID             uint    `gorm:"not null;uniqueIndex:idx_user_tax_name_type"`
	TaxName            string  `gorm:"size:150;not null;uniqueIndex:idx_user_tax_name_type"`
	TaxTypeID          uint    `gorm:"not null;uniqueIndex:idx_user_tax_name_type"`
	TaxPercentage      float64 `gorm:"type:decimal(5,2);not null"`
	TaxPercentageRange string  `gorm:"size:50"`
	TaxApplicableCycle string  `gorm:"size:50;not null"`

	User    User    `gorm:"foreignKey:UserID;constraint:onDelete:CASCADE"`
	TaxType TaxType `gorm:"foreignKey:TaxTypeID;references:ID"`
}

func (Tax) TableName() string {
	return "taxes"
}
