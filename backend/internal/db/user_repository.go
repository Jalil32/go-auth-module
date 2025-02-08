package db

import (
	"database/sql"
	"fmt"
	"wealthscope/backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindUserByEmail(string) (*models.User, error)
	CreateUser(*models.User) error
	UpdateUser(user *models.User) error
}

type UserDB struct {
	*sqlx.DB
}

type UserTX struct {
	*sqlx.Tx
}

func (db *UserDB) FindUserByEmail(email string) (*models.User, error) {
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

func (db *UserDB) CreateUser(ext sqlx.Ext, user *models.User) error {
	if user.Provider == nil && user.PasswordHash == nil {
		return fmt.Errorf("password_hash is required for email/password users")
	}

	query := `INSERT INTO users (email, first_name, last_name, provider, password_hash)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := ext.Exec(
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

func (db *UserDB) UpdateUser(ext sqlx.Ext, user *models.User) error {
	query := `UPDATE users 
              SET email = $1, 
                  first_name = $2, 
                  last_name = $3, 
                  provider = $4, 
                  password_hash = $5, 
				  is_active = $6,
				  verified = $7
              WHERE id = $8`

	_, err := ext.Exec(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Provider,
		user.PasswordHash,
		user.IsActive,
		user.Verified,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
