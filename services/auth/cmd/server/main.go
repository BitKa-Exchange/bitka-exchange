package main

import (
	"bitka/pkg/config"
	"bitka/pkg/logger"
	"bitka/services/auth/internal/app"

	"github.com/rs/zerolog/log"
)

func main() {
	// 1. Global Init
	cfg := config.Load("AUTH_DB_NAME")

	logger.Init(logger.Config{
		Environment: cfg.AppEnv,
		LogLevel:    cfg.LogLevel,
		ServiceName: cfg.ServiceName,
		InstanceID:  cfg.InstanceID,
	})

	// Now you can use the global logger immediately
	log.Info().Msg("Application starting...")

	logger.LogConfigSafe(*cfg)

	// 2. Create Server (Wiring)
	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize server")
	}

	// 3. Start Server
	addr := ":" + cfg.HTTPPort

	log.Info().
		Str("port", cfg.HTTPPort).
		Msg("Starting Auth Service")

	if err := server.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
