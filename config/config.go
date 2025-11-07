package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    Port        string
    YayaBaseURL string
    YayaApiKey  string
}

func Load() *Config {
    _ = godotenv.Load()
    cfg := &Config{
        Port:        getEnv("PORT", "8080"),
        YayaBaseURL: getEnv("YAYA_BASE_URL", "https://api.yayawallet.com/v1"),
        YayaApiKey:  os.Getenv("YAYA_API_KEY"),
    }
    if cfg.YayaApiKey == "" {
        log.Fatal("Missing YAYA_API_KEY")
    }
    return cfg
}

func getEnv(k, d string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return d
}
