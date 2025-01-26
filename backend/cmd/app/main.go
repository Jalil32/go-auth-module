package main

import (
	"log/slog"
	"os"
	"wealthscope/backend/config"
	server "wealthscope/backend/internal/app"
	"wealthscope/backend/internal/app/db"

	"github.com/lmittmann/tint"
)

func main() {
	// Initialise structures logger
	logger := slog.New(tint.NewHandler(os.Stdout, nil))

	// Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
	}

	// Connect to the database
	db, err := db.InitDb(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	} else {
		logger.Info(("Successfully connected to postgres database"))
	}

	defer db.Close()

	// Start the server
	if err := server.StartServer(cfg, db, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
