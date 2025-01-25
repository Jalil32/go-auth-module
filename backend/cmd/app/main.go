package main

import (
	"log"
	"wealthscope/backend/config"
	server "wealthscope/backend/internal/app"
)

// todo - should start the server from main
func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	server.StartServer(cfg)
}
