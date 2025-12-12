package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort string
}

// Load reads configuration.
// dbEnvKey: The specific env var name for this service's DB (e.g. "AUTH_DB_NAME")
func Load(dbEnvKey string) *Config {
	loadEnvFile()

	return &Config{
		AppEnv: GetEnv("APP_ENV", "development"),
		DBHost: GetEnv("DB_HOST", "localhost"),
		DBUser: GetEnv("DB_USER", "postgres"),
		DBPass: GetEnv("DB_PASS", "password"),
		// Critical Change: Look for the specific key first, then fallback to generic
		DBName: GetEnv(dbEnvKey, GetEnv("DB_NAME", "bitka_auth")),
		DBPort: GetEnv("DB_PORT", "5432"),
	}
}

func loadEnvFile() {
	if os.Getenv("APP_ENV") == "production" {
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err == nil {
				log.Printf("Loaded environment from %s", envPath)
			}
			return
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}
}

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
