package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppEnv        string
	HTTPPort      string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration string
	WSAllowed     []string
}

func Load() Config {
	return Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", buildDatabaseURL()),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "1h"),
		WSAllowed:     parseCSV(getEnv("WS_ALLOWED_ORIGINS", "http://localhost:8080")),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func buildDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "bsmart")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
}

func parseCSV(value string) []string {
	var items []string
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items
}
