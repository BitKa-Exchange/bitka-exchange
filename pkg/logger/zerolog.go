package logger

import (
	"bitka/pkg/config"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds setup parameters
type Config struct {
	Environment string // "dev" or "prod"
	LogLevel    string // "debug", "info", "warn", "error"
	ServiceName string // e.g., "auth-service"
	InstanceID  string // e.g., "auth-pod-xyz" or container ID
}

// Init initializes the global logger
func Init(cfg Config) {
	// 1. Set Level
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// 2. Configure Output (Dev vs Prod)
	if cfg.Environment == "dev" || cfg.Environment == "development" {
		// Pretty printing for humans
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	} else {
		// JSON for machines (Default)
		// Unix timestamps are faster and easier for log aggregators
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	// 3. Add Service Context (Crucial for microservices!)
	// We add InstanceID so we know WHICH replica generated the log.
	log.Logger = log.Logger.With().
		Str("service", cfg.ServiceName).
		Str("instance_id", cfg.InstanceID).
		Logger()
}

// logConfigSafe prints the config without revealing secrets
func LogConfigSafe(cfg config.Config) {
	// Redact sensitive fields
	cfg.DBPass = "[REDACTED]"

	// You can add other redactions here if needed (e.g. API Keys)

	log.Info().
		Interface("config", cfg).
		Msg("Loaded Configuration")
}
