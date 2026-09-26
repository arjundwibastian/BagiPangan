package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort      string
	GRPCPort        string
	DatabaseURL     string
	JWTSecret       string
	UserServiceGRPC string
	RequestTTL      time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:      getenv("SERVER_PORT", "8083"),
		GRPCPort:        getenv("GRPC_PORT", "50053"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		UserServiceGRPC: getenv("USER_SERVICE_GRPC_ADDRESS", "localhost:50051"),
		RequestTTL:      durationEnv("REQUEST_TTL", 24*time.Hour),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
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

func (cfg *Config) validate() error {
	required := map[string]string{
		"DATABASE_URL": cfg.DatabaseURL,
		"JWT_SECRET":   cfg.JWTSecret,
	}
	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}
	return nil
}
