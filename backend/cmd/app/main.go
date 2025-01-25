package main

import (
	"log"
	"wealthscope/backend/config"
	server "wealthscope/backend/internal/app"
	"wealthscope/backend/internal/app/db"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	db, err := db.InitDb(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	// Start the server
	if err := server.StartServer(cfg, db); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
