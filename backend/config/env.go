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
	DB   PostgresConfig
	Port string
}

type PostgresConfig struct {
	User     string
	Name     string
	Password string
	Host     string
	Port     string
	SslMode  string
}

func LoadConfig() (*Config, error) {

	cfg := &Config{
		Port: os.Getenv("PORT"),
		DB: PostgresConfig{
			User:     os.Getenv("POSTGRES_USER"),
			Name:     os.Getenv("POSTGRES_NAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			SslMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
	}

	return cfg, nil
}
