package auth

import (
	"fmt"
	"log/slog"
	"wealthscope/backend/config"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthController struct {
	DB              *sqlx.DB
	RedisCache      *redis.Client
	Logger          *slog.Logger
	JwtToken        string
	JwtExpiry       string
	Host            string
	Port            string
	Username        string
	Password        string
	FrontendAddress string
}

// NewAuthController initializes a new AuthController
func NewAuthController(db *sqlx.DB, rdb *redis.Client, logger *slog.Logger) (*AuthController, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return nil, fmt.Errorf("failed to load config: %w", err) // Error message clarified
	}

	return &AuthController{
		DB:              db,
		RedisCache:      rdb,
		Logger:          logger,
		JwtToken:        cfg.JWT.Token,
		JwtExpiry:       cfg.JWT.Expiry,
		Host:            cfg.SMTP.Host,
		Port:            cfg.SMTP.Port,
		Username:        cfg.SMTP.Username,
		Password:        cfg.SMTP.Password,
		FrontendAddress: cfg.Frontend.IP + ":" + cfg.Frontend.Port,
	}, nil
}
