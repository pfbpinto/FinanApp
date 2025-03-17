package models

import "time"

// UserProfile represents the user data structure used throughout the app
type UserProfile struct {
	UserProfileID    int       `json:"userProfileID"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	DateOfBirth      string    `json:"dob"`
	UserPassword     string    `json:"userPassword"`
	EmailAddress     string    `json:"emailAddress"`
	UserSubscription bool      `json:"userSubscription"`
	CreatedAt        time.Time `json:"createdAt"`
}
