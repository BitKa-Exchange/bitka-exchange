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
	// --- NEW FIELDS ---
	LogLevel    string // e.g., "debug", "info", "warn"
	ServiceName string // e.g., "auth-service"
	HTTPPort    string // e.g., "3000"
	InstanceID  string // The Container ID or Pod Name
}

// Load reads configuration.
// dbEnvKey: The specific env var name for this service's DB (e.g. "AUTH_DB_NAME")
func Load(dbEnvKey string) *Config {
	loadEnvFile()

	// Get Hostname (Container ID in Docker, Pod Name in K8s)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", "password"),
		// Critical Change: Look for the specific key first, then fallback to generic
		DBName: getEnv(dbEnvKey, getEnv("DB_NAME", "bitka_auth")),
		DBPort: getEnv("DB_PORT", "5432"),

		// Logging
		HTTPPort: getEnv("HTTP_PORT", "3000"),
		// Default to "info" if not set. Use "debug" in local .env
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// Default to "unknown" to alert if not set
		ServiceName: getEnv("SERVICE_NAME", "unknown-service"),
		InstanceID:  getEnv("INSTANCE_ID", hostname), // Fallback to env var if needed
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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
