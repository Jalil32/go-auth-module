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
	Username string
	Password string
	URL      string
	Port     string
}

func LoadConfig() (*Config, error) {

	cfg := &Config{
		Port: os.Getenv("PORT"),

		// todo - we will need this configuration at some stage...
		// DB: PostgresConfig{
		// 	Username: os.Getenv("POSTGRES_USER"),
		// 	Password: os.Getenv("POSTGRES_PWD"),
		// 	URL:      os.Getenv("POSTGRES_URL"),
		// 	Port:     os.Getenv("POSTGRES_PORT"),
		// },
	}

	return cfg, nil
}
