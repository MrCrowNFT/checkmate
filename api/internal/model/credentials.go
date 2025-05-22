package model

import (
	"time"
)

// stored in sql db
type PlatformCredential struct {
	ID        int       `json:"id"`
	UserID    string    `json:"userId"`
	Platform  string    `json:"platform"`
	APIKey    string    `json:"apiKey"` // will be encrypted in storage
	CreatedAt time.Time `json:"createdAt"`
}

// credentials without sensitive info -> meant for the return values to the client
type SafeCredential struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Platform  string    `json:"platform"`          
	CreatedAt time.Time `json:"created_at"`
}

// user input
type PlatformCredentialInput struct {
	Platform string `json:"platform"`
	APIKey   string `json:"apiKey"`
}
