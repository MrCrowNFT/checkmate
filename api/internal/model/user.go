package model

import (
	"time"
)

type User struct {
	ID          string    `json:"id`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
}
