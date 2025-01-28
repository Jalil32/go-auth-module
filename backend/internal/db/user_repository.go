package db

import (
	"database/sql"
	"fmt"
	"wealthscope/backend/internal/models"

	"github.com/jmoiron/sqlx"
)

func FindUserByEmail(db *sqlx.DB, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email=$1`

	var user models.User
	err := db.Get(&user, query, email)

	if err != nil {
		// handle case where there are no rows
		if err == sql.ErrNoRows {
			return nil, nil
		}
		// handle case where another error occurs
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	return &user, nil
}

func CreateUser(db *sqlx.DB, user *models.User) error {
	// Validate email/password users
	if user.Provider == nil && user.PasswordHash == nil {
		return fmt.Errorf("password_hash is required for email/password users")
	}

	query := `INSERT INTO users (email, first_name, last_name, provider, password_hash)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Provider,
		user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}
