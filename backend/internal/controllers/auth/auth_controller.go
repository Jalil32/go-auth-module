package auth

import (
	"context"
	"time"
	"wealthscope/backend/config"
	"wealthscope/backend/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type JWTGenerator interface {
	GenerateJWT(user *models.User) (string, error)
}

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

type UserRepository interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(ext sqlx.Ext, user *models.User) error
	UpdateUser(ext sqlx.Ext, user *models.User) error
	Beginx() (*sqlx.Tx, error)
}

type AuthController struct {
	UserDB          UserRepository
	RedisCache      RedisClient
	Logger          Logger
	JwtToken        string
	JwtExpiry       string
	Host            string
	Port            string
	Username        string
	Password        string
	FrontendAddress string
	JWTGenerator    JWTGenerator
}

// NewAuthController initializes a new AuthController
func NewAuthController(userRepo UserRepository, rdb RedisClient, logger Logger, jwtGenerator JWTGenerator, cfg *config.Config) (*AuthController, error) {

	return &AuthController{
		UserDB:          userRepo,
		RedisCache:      rdb,
		Logger:          logger,
		JwtToken:        cfg.JWT.Token,
		JwtExpiry:       cfg.JWT.Expiry,
		Host:            cfg.SMTP.Host,
		Port:            cfg.SMTP.Port,
		Username:        cfg.SMTP.Username,
		Password:        cfg.SMTP.Password,
		FrontendAddress: cfg.Frontend.Addr,
		JWTGenerator:    jwtGenerator,
	}, nil
}
