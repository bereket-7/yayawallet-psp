package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	YayaBaseURL  string
	YayaClientID string
	YayaClientSecret string
	Env          string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		YayaBaseURL:      getEnv("YAYA_BASE_URL", "https://pay.yayawallet.com"),
		YayaClientID:     os.Getenv("YAYA_CLIENT_ID"),
		YayaClientSecret: os.Getenv("YAYA_CLIENT_SECRET"),
		Env:              getEnv("ENV", "local"),
	}

	// Required checks
	if cfg.YayaClientID == "" {
		log.Fatal("Missing YAYA_CLIENT_ID")
	}
	if cfg.YayaClientSecret == "" {
		log.Fatal("Missing YAYA_CLIENT_SECRET")
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
