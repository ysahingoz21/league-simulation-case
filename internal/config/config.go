package config

import "os"

const (
	defaultPort   = "8080"
	defaultAppEnv = "development"
)

type Config struct {
	Port   string
	AppEnv string
}

func Load() Config {
	return Config{
		Port:   getEnv("PORT", defaultPort),
		AppEnv: getEnv("APP_ENV", defaultAppEnv),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
