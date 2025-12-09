package main

import (
	"log"
	"os"

	"bitka/pkg/config"
	"bitka/pkg/logger"
	"bitka/services/account/internal/app"
	"bitka/services/account/internal/kafka"
)

func main() {
	logger.Init()
	cfg := config.Load("ACCOUNT_DB_NAME")

	// Override DB Name for Account Service if not set in env specific to service
	if os.Getenv("DB_NAME") == "" {
		cfg.DBName = "bitka_account"
	}

	server,uc, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed init: %v", err)
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "3001"
	}

	kafkaConsumer := kafka.NewConsumer(uc)
	kafkaConsumer.Start()

	log.Printf("Starting Account Service on :%s", port)
	server.Listen(":" + port)
}
