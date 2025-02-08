/*
todo - add to this comment and config as required
this is for loading configuration from environment variables
usage: create a .env and populate with the following fields
PORT:
*/
package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Backend  BackendConfig
	Frontend FrontendConfig
	SMTP     SMTPConfig
	DB       PostgresConfig
	OAuth    OAuthConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

type BackendConfig struct {
	Port string
	IP   string
}

type FrontendConfig struct {
	Port string
	IP   string
}

type PostgresConfig struct {
	User     string
	Name     string
	Password string
	Host     string
	Port     string
	SslMode  string
}

type OAuthConfig struct {
	ClientID          string
	ClientSecret      string
	ClientCallbackURL string
}

type JWTConfig struct {
	Token  string
	Expiry string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type RedisConfig struct {
	Address  string
	Database string
	Password string
}

func LoadConfig() (*Config, error) {

	cfg := &Config{
		Frontend: FrontendConfig{
			IP:   os.Getenv("FRONTEND_IP"),
			Port: os.Getenv("FRONTEND_PORT"),
		},
		Backend: BackendConfig{
			IP:   os.Getenv("BACKEND_IP"),
			Port: os.Getenv("BACKEND_PORT"),
		},
		DB: PostgresConfig{
			User:     os.Getenv("POSTGRES_USER"),
			Name:     os.Getenv("POSTGRES_NAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			SslMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
		OAuth: OAuthConfig{
			ClientID:          os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret:      os.Getenv("GOOGLE_CLIENT_SECRET"),
			ClientCallbackURL: os.Getenv("GOOGLE_CLIENT_CALLBACK_URL"),
		},
		JWT: JWTConfig{
			Token:  os.Getenv("JWT_TOKEN"),
			Expiry: os.Getenv("JWT_EXPIRY"),
		},
		SMTP: SMTPConfig{
			Host:     os.Getenv("EMAIL_HOST"),
			Port:     os.Getenv("EMAIL_PORT"),
			Username: os.Getenv("EMAIL_USERNAME"),
			Password: os.Getenv("EMAIL_PASSWORD"),
		},
		Redis: RedisConfig{
			Address:  os.Getenv("REDIS_ADDRESS"),
			Database: os.Getenv("REDIS_DATABASE"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}

	return cfg, nil
}
