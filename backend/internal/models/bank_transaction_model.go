package models

import "time"

type Transaction struct {
	TransactionId int       // Primary key: Auto-incremented in the database
	UserId        int       // Foreign key to the user who made the transaction
	Date          time.Time // Date of transaction
	AmountCents   int       // Transaction amount in cents
	Description   string    // Transaction description: Default is "No description"
}
