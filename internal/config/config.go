package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Redis     RedisConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string
	Environment string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	DefaultRequestsPerMinute int
	AuthRequestsPerMinute    int
	SearchRequestsPerMinute  int
	BookingRequestsPerMinute int
	ReviewRequestsPerMinute  int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret            string
	ExpirationHours   int
	RefreshExpiration int
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8081"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "airbnb_clone"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHours:   getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshExpiration: getEnvAsInt("JWT_REFRESH_EXPIRATION_HOURS", 168), // 7 days
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		RateLimit: RateLimitConfig{
			DefaultRequestsPerMinute: getEnvAsInt("RATE_LIMIT_DEFAULT", 100),
			AuthRequestsPerMinute:    getEnvAsInt("RATE_LIMIT_AUTH", 10),
			SearchRequestsPerMinute:  getEnvAsInt("RATE_LIMIT_SEARCH", 30),
			BookingRequestsPerMinute: getEnvAsInt("RATE_LIMIT_BOOKING", 5),
			ReviewRequestsPerMinute:  getEnvAsInt("RATE_LIMIT_REVIEW", 10),
		},
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as an integer with a fallback value
func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
