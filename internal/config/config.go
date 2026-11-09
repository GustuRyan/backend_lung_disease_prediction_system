package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment     string
	Address         string
	DatabaseDSN     string
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Environment:     env("APP_ENV", "development"),
		Address:         env("HTTP_ADDRESS", ":8080"),
		DatabaseDSN:     strings.TrimSpace(os.Getenv("DATABASE_DSN")),
		AllowedOrigins:  splitCSV(env("CORS_ALLOWED_ORIGINS", "http://localhost:8501")),
		ShutdownTimeout: durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
	}

	if cfg.DatabaseDSN == "" {
		return Config{}, errors.New("DATABASE_DSN is required")
	}
	if len(cfg.AllowedOrigins) == 0 {
		return Config{}, errors.New("CORS_ALLOWED_ORIGINS must contain at least one origin")
	}
	if cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			result = append(result, value)
		}
	}
	return result
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	if parsed, err := time.ParseDuration(value); err == nil {
		return parsed
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return 0
}
