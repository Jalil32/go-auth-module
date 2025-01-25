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
	DbUser     string
	DbName     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbSslMode  string
}

func LoadConfig() (*Config, error) {

	cfg := &Config{
		Port: os.Getenv("PORT"),
		DB: PostgresConfig{
			DbUser:     os.Getenv("POSTGRES_USER"),
			DbName:     os.Getenv("POSTGRES_DB"),
			DbPassword: os.Getenv("POSTGRES_PASSWORD"),
			DbHost:     os.Getenv("POSTGRES_HOST"),
			DbPort:     os.Getenv("POSTGRES_PORT"),
			DbSslMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
	}

	return cfg, nil
}
