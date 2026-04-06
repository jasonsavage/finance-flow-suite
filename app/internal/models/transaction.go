package models

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Date            time.Time `json:"date"`
	Description     string    `json:"description"`
	Category        *string   `json:"category"`
	Deposit         float64   `json:"deposit"`
	Withdrawal      float64   `json:"withdrawal"`
	BankAccountName string    `json:"bank_account_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
