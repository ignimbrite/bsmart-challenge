package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv        string
	HTTPPort      string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration string
	WSAllowed     []string
	SeedOnStart   bool
}

func Load() Config {
	return Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", buildDatabaseURL()),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "1h"),
		WSAllowed:     parseCSV(getEnv("WS_ALLOWED_ORIGINS", "http://localhost,http://127.0.0.1,http://localhost:8080,http://127.0.0.1:8080,https://ignimbrite.github.io,https://8113c6fc74a6.ngrok-free.app")),
		SeedOnStart:   getEnvAsBool("SEED_ON_START", false),
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

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}
