package models

type UserGroupAsset struct {
	ID               uint    `gorm:"primaryKey;autoIncrement"`
	UserAssetID      uint    `gorm:"not null;uniqueIndex:idx_user_asset_group"`
	GroupID          uint    `gorm:"not null;uniqueIndex:idx_user_asset_group"`
	Owner            bool    `gorm:"default:false"`
	OwningPercentage float64 `gorm:"type:decimal(5,2)"`

	// Relacionamentos
	UserAsset UserAsset `gorm:"foreignKey:UserAssetID"`
	UserGroup UserGroup `gorm:"foreignKey:GroupID"`
}

func (UserGroupAsset) TableName() string {
	return "user_group_assets"
}
