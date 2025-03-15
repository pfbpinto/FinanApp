package models

import "time"

// SchemaMigration represents a record of migrations applied.
type SchemaMigration struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Version   string    `gorm:"size:255;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// TableName defines the name of the table for schema migrations
func (SchemaMigration) TableName() string {
	return "schema_migrations"
}
