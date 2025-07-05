package model

import "time"

type User struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't expose password hash in JSON
	CreatedAt  time.Time `json:"created_at"`
}
