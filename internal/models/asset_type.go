package models

type AssetType struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	AssetTypeName string `gorm:"size:100;not null;unique"`
}

func (AssetType) TableName() string {
	return "asset_types"
}
