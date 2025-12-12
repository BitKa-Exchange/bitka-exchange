package main

import (
	"log"
	"os"

	"bitka/pkg/config"
	"bitka/pkg/logger"
	"bitka/services/account/internal/app"
)

func main() {
	logger.Init()
	cfg := config.Load("ACCOUNT_DB_NAME")

	// Override DB Name for Account Service if not set in env specific to service
	if os.Getenv("DB_NAME") == "" {
		cfg.DBName = "bitka_account"
	}

	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed init: %v", err)
	}

	port := config.GetEnv("HTTP_PORT" ,"3001")
	if port == "" {
		port = "3001"
	}

	log.Printf("Starting Account Service on :%s", port)
	server.Listen(":" + port)
}
