package main

import (
	"os"

	"bitka/pkg/config"
	"bitka/pkg/database"
	"bitka/pkg/logger"
	"bitka/pkg/token"

	"bitka/services/account/internal/app"
	"bitka/services/account/internal/delivery/event"
	"bitka/services/account/internal/domain"
	"bitka/services/account/internal/repository"
	"bitka/services/account/internal/usecase"

	"github.com/rs/zerolog/log"
)

func main() {
	// 1. Load Configuration
	// We pass "ACCOUNT_DB_NAME" to look for that specific env var override
	cfg := config.Load("ACCOUNT_DB_NAME")

	// 2. Initialize Logger
	logger.Init(logger.Config{
		Environment: cfg.AppEnv,
		LogLevel:    cfg.LogLevel,
		ServiceName: cfg.ServiceName,
		InstanceID:  cfg.InstanceID,
	})

	log.Info().Msg("Starting Account Service...")

	// 3. Log Config (Safely redacted)
	logger.LogConfigSafe(*cfg)

	// 4. Connect Database
	// We use the values directly from cfg (defaults are handled in pkg/config)
	db, err := database.Connect(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// AutoMigrate is fine for dev, but consider sql-migrate for prod
	if err := db.AutoMigrate(&domain.Profile{}); err != nil {
		log.Fatal().Err(err).Msg("Database migration failed")
	}

	// 5. Dependency Injection
	repo := repository.NewAccountRepo(db)
	uc := usecase.NewAccountUsecase(repo)

	// 6. JWT Validator
	// This is specific to Account service, so we handle the fallback here
	jwksURL := os.Getenv("AUTH_JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://localhost:3000/.well-known/jwks.json"
	}
	log.Info().Str("url", jwksURL).Msg("Initializing JWT Validator")
	validator := token.NewValidator(jwksURL)

	// 7. Kafka Consumer Server
	kafkaHandler := event.NewHandler(uc)
	kafkaServer := event.NewServer(kafkaHandler)

	// Start Kafka in a goroutine
	go func() {
		log.Info().Msg("Starting Kafka Consumer...")
		kafkaServer.Start()

	}()

	// 8. HTTP Server
	httpServer := app.NewServer(uc, validator)

	addr := ":" + cfg.HTTPPort
	log.Info().Str("addr", addr).Msg("HTTP Server listening")

	if err := httpServer.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("HTTP Server failed")
	}
}
