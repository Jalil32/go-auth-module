package models

import "time"

type Transaction struct {
	TransactionID int       // Primary key: Auto-incremented in the database
	UserID        int       // Foreign key to the user who made the transaction
	Date          time.Time // Date of transaction
	AmountCents   int       // Transaction amount in cents
	Description   string    // Transaction description: Default is "No description"
}
