package db

import (
	"fmt"
	"wealthscope/backend/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDb(cfg *config.Config) (*sqlx.DB, error) {
	// Create connection string
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s", cfg.DB.DbUser, cfg.DB.DbName, cfg.DB.DbSslMode, cfg.DB.DbPassword, cfg.DB.DbHost, cfg.DB.DbPort)

	// Open database connection
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	return db, nil
}
