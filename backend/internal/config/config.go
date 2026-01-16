package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Diagnocat DiagnocatConfig
}

type ServerConfig struct {
	Port            string
	Environment     string
	AllowedOrigins  []string
	MaxUploadSizeMB int64
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type DiagnocatConfig struct {
	APIKey string
}

func Load() *Config {
	// Load .env file if exists (for local dev)
	godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxUploadMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "500"), 10, 64)

	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			Environment:     getEnv("ENVIRONMENT", "development"),
			AllowedOrigins:  []string{getEnv("FRONTEND_URL", "http://localhost:3000")},
			MaxUploadSizeMB: maxUploadMB,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "dental_marketplace"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-me-in-production"),
			AccessTokenTTL:  24 * time.Hour,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		},
		Diagnocat: DiagnocatConfig{
			APIKey: os.Getenv("DIAGNOCAT_API_KEY"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
