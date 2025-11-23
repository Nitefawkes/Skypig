package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	CORSOrigins string
	JWTSecret   string

	// OAuth
	QRZClientID     string
	QRZClientSecret string
	QRZRedirectURL  string

	// Feature flags
	EnablePropagation bool
	EnableSDR         bool
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/hamradio?sslmode=disable"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),

		QRZClientID:     getEnv("QRZ_CLIENT_ID", ""),
		QRZClientSecret: getEnv("QRZ_CLIENT_SECRET", ""),
		QRZRedirectURL:  getEnv("QRZ_REDIRECT_URL", "http://localhost:8080/api/v1/auth/callback/qrz"),

		EnablePropagation: getEnv("ENABLE_PROPAGATION", "true") == "true",
		EnableSDR:         getEnv("ENABLE_SDR", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}
