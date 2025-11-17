package config

import (
	"time"

	"bitka/config"
)

type Config struct {
	ListenAddr  string
	Issuer      string
	Audience    string
	AccessTTL   time.Duration
	RefreshTTL  time.Duration
	KeyID       string
	PrivateKeyPath string
	PublicKeyPath  string
}

func Load() *Config {
	return &Config{
		ListenAddr:     config.GetEnv("AUTH_LISTEN", ":8080"),
		Issuer:         config.GetEnv("AUTH_ISSUER", "http://localhost:8080"),
		Audience:       config.GetEnv("AUTH_AUDIENCE", "bitka-api"),
		AccessTTL:      config.GetEnvDuration("AUTH_ACCESS_TTL", 15*time.Minute),
		RefreshTTL:     config.GetEnvDuration("AUTH_REFRESH_TTL", 7*24*time.Hour),
		KeyID:          config.GetEnv("AUTH_KEY_ID", "dev-key-1"),
		PrivateKeyPath: config.GetEnv("AUTH_PRIV_KEY", "keys/rsa_private.pem"),
		PublicKeyPath:  config.GetEnv("AUTH_PUB_KEY",  "keys/rsa_public.pem"),
	}
}
