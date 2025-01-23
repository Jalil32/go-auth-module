package main

import (
	"log"
	"wealthscope/backend/config"
	"wealthscope/backend/server"
)

// todo - should start the server from main
func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
    
    server.StartServer(cfg)
}