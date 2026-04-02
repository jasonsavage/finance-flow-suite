package models

import "time"

type User struct {
	ID           int       `json:"id"`
	AccountID    string    `json:"account_id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // never serialize password hash to JSON
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
