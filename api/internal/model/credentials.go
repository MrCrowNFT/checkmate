package model

import (
	"time"
)

//stored in sql db
type PlatformCredential struct {
	ID        int       `json:"id"`
	UserID    string    `json:"userId"`
	Platform  string    `json:"platform"`
	Name      string    `json:"name"`
	APIKey    string    `json:"apiKey"` // will be encrypted in storage
	CreatedAt time.Time `json:"createdAt"`
}

//user input
type PlatformCredentialInput struct {
	Platform string `json:"platform"`
	Name     string `json:"name"`
	APIKey   string `json:"apiKey"`
}
