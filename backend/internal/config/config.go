package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Stripe   StripeConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port        string
	Environment string
	CORSOrigins string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string
	QRZAPIKey     string
	SessionSecret string
}

// StripeConfig holds Stripe payment configuration
type StripeConfig struct {
	SecretKey      string
	WebhookSecret  string
	FreeTierPrice  string
	OperatorPrice  string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
			CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "hamradio_cloud"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", ""),
			QRZAPIKey:     getEnv("QRZ_API_KEY", ""),
			SessionSecret: getEnv("SESSION_SECRET", ""),
		},
		Stripe: StripeConfig{
			SecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
			FreeTierPrice: getEnv("STRIPE_FREE_TIER_PRICE", ""),
			OperatorPrice: getEnv("STRIPE_OPERATOR_PRICE", ""),
		},
	}

	// Validate critical configuration
	if cfg.Server.Environment == "production" {
		if cfg.Auth.JWTSecret == "" {
			return nil, fmt.Errorf("JWT_SECRET must be set in production")
		}
		if cfg.Database.Password == "" {
			return nil, fmt.Errorf("DB_PASSWORD must be set in production")
		}
	}

	return cfg, nil
}

// ConnectionString returns the PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
