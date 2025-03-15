package models

import (
	"time"
)

type UserSubscription struct {
	ID                    uint       `gorm:"primaryKey;autoIncrement"`
	UserID                uint       `gorm:"not null"`
	SubscriptionTypeID    uint       `gorm:"not null"`
	SubscriptionBeginDate time.Time  `gorm:"not null"`
	SubscriptionEndDate   time.Time  `gorm:"not null"`
	CancelledAt           *time.Time `gorm:"default:null"`
	IsActive              bool       `gorm:"default:true"`

	User         User         `gorm:"foreignKey:UserID"`
	Subscription Subscription `gorm:"foreignKey:SubscriptionTypeID"`
}

func (UserSubscription) TableName() string {
	return "user_subscriptions"
}
