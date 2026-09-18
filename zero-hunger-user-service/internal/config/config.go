package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort  string
	GRPCPort    string
	DatabaseURL string
	JWTSecret   string
	AccessTTL   time.Duration
	RefreshTTL  time.Duration
}

func Load() Config {
	return Config{
		ServerPort:  getenv("SERVER_PORT", "8081"),
		GRPCPort:    getenv("GRPC_PORT", "50051"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://user:password@localhost:5433/user_db?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET", "change-me-in-development"),
		AccessTTL:   durationEnv("JWT_ACCESS_TTL", time.Hour),
		RefreshTTL:  durationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	if parsed, err := time.ParseDuration(value); err == nil && parsed > 0 {
		return parsed
	}
	return fallback
}
