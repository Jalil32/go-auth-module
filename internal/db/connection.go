package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/jalil32/go-auth-module/config"
)

func InitDb(cfg *config.Config) (*sqlx.DB, error) {
	// Create connection string
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s", cfg.DB.User, cfg.DB.Name, cfg.DB.SslMode, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port)

	// Open database connection
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	return db, nil
}
