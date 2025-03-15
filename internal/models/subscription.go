package models

type Subscription struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	Name         string  `gorm:"size:100;not null"`
	Description  string  `gorm:"size:255"`
	Duration     int     `gorm:"not null"`
	DurationUnit string  `gorm:"size:20;not null"` // "months", "days", etc.
	Price        float64 `gorm:"type:decimal(10,2);not null"`
}
