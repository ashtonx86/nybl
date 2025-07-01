package schemas

import (
	"time"
)

// Represents the account of a user.
// `Platform` is the OAuth2 provider used during sign-up
type Account struct {
	ID string `json:"id"`
	Name string `json:"name"`
	EmailID string `json:"email_id"`
	Verified bool `json:"verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestAccountCreate struct {
	Name string `json:"name" validate:"required,min=1,max=42"`
	EmailID string `json:"email_id" validate:"required,min=1,email"`
}

type OTP struct {
	ID string 
	Code string
	EmailID string 
	AccountID string

	RequestedAt time.Time `json:"requested_at"`
}