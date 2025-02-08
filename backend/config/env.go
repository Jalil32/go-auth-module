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
	Port  string
	DB    PostgresConfig
	OAuth OAuthConfig
	JWT   JWTConfig
	ClientLocal string
	ClientProxy string
	ClientFly string
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

func LoadConfig() (*Config, error) {

	cfg := &Config{
		Port: os.Getenv("PORT"),
		ClientLocal: os.Getenv("CLIENT_LOCAL"),
		ClientProxy: os.Getenv("CLIENT_PROXY"),
		ClientFly: os.Getenv("CLIENT_FLY"),
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
	}

	return cfg, nil
}
