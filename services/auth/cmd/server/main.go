package main

import (
	"log"

	"bitka/pkg/config"
	"bitka/pkg/logger"
	"bitka/services/auth/internal/app"
)

func main() {
	// TODO: Make all packages use the global logger
	// 1. Global Init
	logger.Init()
	cfg := config.Load("AUTH_DB_NAME")

	// 2. Create Server (Wiring)
	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// 3. Start Server
	log.Println("Starting Auth Service on :3000")
	if err := server.Listen(":3000"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
