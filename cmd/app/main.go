package main

import (
	"context"
	"github.com/jalil32/go-auth-module/config"
	server "github.com/jalil32/go-auth-module/internal/app"
	"github.com/jalil32/go-auth-module/internal/db"
	"github.com/lmittmann/tint"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	// Initialise structures logger
	logger := slog.New(tint.NewHandler(os.Stdout, nil))

	// Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
	}

	// Load redis database and convert to int
	database, err := strconv.Atoi(cfg.Redis.Database)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       database,
	})

	// Ensure the connection is properly closed gracefully
	defer rdb.Close()

	ctx := context.Background()

	// Test the connection
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		logger.Error("Redis connection was refused", "error", err)
	} else {
		logger.Info("Redis successfully connected")
	}

	// Connect to the database
	db, err := db.InitDb(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to database", "error", err)
	} else {
		logger.Info(("Successfully connected to postgres database"))
	}

	defer db.Close()

	// Start the server
	if err := server.StartServer(cfg, db, rdb, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
