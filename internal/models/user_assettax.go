package models

type UserAssetTax struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	UserAssetID uint    `gorm:"not null;index:idx_user_asset_tax,unique"`
	TaxID       uint    `gorm:"not null;index:idx_user_asset_tax,unique"`
	TaxValue    float64 `gorm:"type:decimal(10,2);not null"`
	Payed       bool    `gorm:"default:false"`

	// Relationships
	UserAsset UserAsset `gorm:"foreignKey:UserAssetID;constraint:onDelete:CASCADE"`
	Tax       Tax       `gorm:"foreignKey:TaxID;constraint:onDelete:RESTRICT"`
}

func (UserAssetTax) TableName() string {
	return "user_asset_taxes"
}
