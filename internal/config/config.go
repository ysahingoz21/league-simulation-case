package config

import "os"

const (
	defaultPort   = "8080"
	defaultAppEnv = "development"
	defaultDBPath = "league.db"
)

type Config struct {
	Port   string
	AppEnv string
	DBPath string
}

func Load() Config {
	return Config{
		Port:   getEnv("PORT", defaultPort),
		AppEnv: getEnv("APP_ENV", defaultAppEnv),
		DBPath: getEnv("DB_PATH", defaultDBPath),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
