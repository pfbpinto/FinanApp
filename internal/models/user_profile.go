package models

// UserProfile represents the user data structure used throughout the app

type UserProfile struct {
	UserProfileID    int    `json:"user_profile_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	DateOfBirth      string `json:"date_of_birth"`
	UserPassword     string `json:"user_password"`
	EmailAddress     string `json:"email_address"`
	UserSubscription bool   `json:"user_subscription"`
	CreatedAt        string `json:"created_at"`
}
