package main

import (
	"bitka/pkg/config"
	"bitka/services/account/internal/app"

	"github.com/rs/zerolog/log"
)

func main() {
	// 1. Load Configuration
	// We pass "ACCOUNT_DB_NAME" to look for that specific env var override
	cfg := config.Load("ACCOUNT_DB_NAME")

	// Override DB Name for Account Service if not set in env specific to service
	cfg.DBName = config.GetEnv("DB_NAME", "bitka_account")
	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize server")
	}

	log.Printf("Starting Account Service on :%s", cfg.HTTPPort)
	server.Listen(":" + cfg.HTTPPort)
}
