package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	GRPC     GrpcConfig
}

type AppConfig struct {
	Port     string
	GRPCPort string
}

type GrpcConfig struct {
	UserService    string
	RequestService string
}
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Port:     getenv("PORT", getenv("SERVER_PORT", "8082")),
			GRPCPort: getenv("GRPC_PORT", "50052"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Expiry: getEnvDuration("JWT_EXPIRY_HOURS", 24) * time.Hour,
		},
		GRPC: GrpcConfig{
			UserService:    getenv("USER_SERVICE_GRPC_ADDRESS", "localhost:50051"),
			RequestService: getenv("REQUEST_SERVICE_GRPC_ADDRESS", "localhost:50053"),
		},
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

func getEnvDuration(key string, fallbackHours int) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return time.Duration(fallbackHours)
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return time.Duration(fallbackHours)
	}
	return time.Duration(n)
}

func (cfg *Config) validate() error {
	required := map[string]string{
		"DB_HOST":                      cfg.Database.Host,
		"DB_USER":                      cfg.Database.User,
		"DB_PASSWORD":                  cfg.Database.Password,
		"DB_NAME":                      cfg.Database.Name,
		"DB_PORT":                      cfg.Database.Port,
		"JWT_SECRET":                   cfg.JWT.Secret,
		"USER_SERVICE_GRPC_ADDRESS":    cfg.GRPC.UserService,
		"REQUEST_SERVICE_GRPC_ADDRESS": cfg.GRPC.RequestService,
	}
	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}
	return nil
}
