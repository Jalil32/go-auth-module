package auth

import (
	"fmt"
	"log/slog"
	"wealthscope/backend/config"

	"github.com/jmoiron/sqlx"
)

type AuthController struct {
	DB        *sqlx.DB
	Logger    *slog.Logger
	JwtToken  string
	JwtExpiry string
}

// NewAuthController initializes a new AuthController
func NewAuthController(db *sqlx.DB, logger *slog.Logger) (*AuthController, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return nil, fmt.Errorf("failed to load config: %w", err) // Error message clarified
	}

	return &AuthController{
		DB:        db,
		Logger:    logger,
		JwtToken:  cfg.JWT.Token,
		JwtExpiry: cfg.JWT.Expiry,
	}, nil
}
