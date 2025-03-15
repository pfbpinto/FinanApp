package models

import (
	"time"
)

type UserAsset struct {
	ID                   uint       `gorm:"primaryKey;autoIncrement"`
	UserID               uint       `gorm:"not null;uniqueIndex:idx_user_asset"`
	AssetTypeID          uint       `gorm:"not null;uniqueIndex:idx_user_asset"`
	AssetName            string     `gorm:"size:150;not null;uniqueIndex:idx_user_asset"`
	AssetAquisitionDate  *time.Time `gorm:"column:AssetAquisitionDate"`
	AssetDispositionDate *time.Time `gorm:"column:AssetDispositionDate"`
	AssetValue           float64    `gorm:"type:decimal(10,2);not null"`
	SharedAsset          bool       `gorm:"default:false"`

	// Relationships
	User           User           `gorm:"foreignKey:UserID;constraint:onDelete:CASCADE"`
	AssetType      AssetType      `gorm:"foreignKey:AssetTypeID;constraint:onDelete:RESTRICT"`
	UserAssetTaxes []UserAssetTax `gorm:"foreignKey:UserAssetID;constraint:onDelete:CASCADE"`
}

func (UserAsset) TableName() string {
	return "user_assets"
}
