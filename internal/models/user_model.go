package models

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash *string   `db:"password_hash" json:"passwordHash"`
	FirstName    string    `db:"first_name" json:"firstName"`
	LastName     string    `db:"last_name" json:"lastName"`
	Provider     *string   `db:"provider" json:"provider"`
	IsActive     string    `db:"is_active" json:"isActive"`
	Verified     bool      `db:"verified" json:"verified"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
