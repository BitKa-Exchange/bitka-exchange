package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env              string // "development" | "production"
	Port             string // ":8080" or "8080" (normalize later)
	DatabaseDSN      string
	JwtAccessSecret  string
	JwtRefreshSecret string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration
	LogLevel         string // debug, info, warn, error
}

func loadDotEnv() {
	// ignore error if file not present
	_ = godotenv.Load()
}

// getEnv with default
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// parseDurationFromMinutes tries to parse numeric minutes or fallback to default minutes
func parseMinutes(key string, defMinutes int) time.Duration {
	s := os.Getenv(key)
	if s == "" {
		return time.Duration(defMinutes) * time.Minute
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return time.Duration(defMinutes) * time.Minute
	}
	return time.Duration(n) * time.Minute
}

func New() (*Config, error) {
	// if .env exists, load it. For production, the app will use env vars set in system/containers.
	loadDotEnv()

	env := getEnv("APP_ENV", "development")
	port := getEnv("PORT", "8080")
	dsn := getEnv("DATABASE_DSN", "")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required")
	}

	accessTTL := parseMinutes("JWT_ACCESS_TTL_MIN", 15)
	refreshTTL := parseMinutes("JWT_REFRESH_TTL_MIN", 60*24*7) // 7 days default

	cfg := &Config{
		Env:              env,
		Port:             port,
		DatabaseDSN:      dsn,
		JwtAccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
		JwtRefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
		AccessTTL:        accessTTL,
		RefreshTTL:       refreshTTL,
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	// minimal validation
	if cfg.JwtAccessSecret == "" || cfg.JwtRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET required")
	}

	return cfg, nil
}
