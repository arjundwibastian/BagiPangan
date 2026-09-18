package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	GRPCPort        string
	DatabaseURL     string
	JWTSecret       string
	UserServiceGRPC string
	RequestTTL      time.Duration
}

func Load() Config {
	// Load local .env when present. Existing process environment variables win.
	_ = godotenv.Load()
	return Config{
		ServerPort:      getenv("SERVER_PORT", "8083"),
		GRPCPort:        getenv("GRPC_PORT", "50053"),
		DatabaseURL:     getenv("DATABASE_URL", "postgres://user:password@localhost:5435/request_db?sslmode=disable"),
		JWTSecret:       getenv("JWT_SECRET", "change-me-in-development"),
		UserServiceGRPC: getenv("USER_SERVICE_GRPC_ADDRESS", "localhost:50051"),
		RequestTTL:      durationEnv("REQUEST_TTL", 24*time.Hour),
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
