package config

import (
	"os"
	"strings"
	"time"
)

type RateLimitConfig struct {
	MaxRequests     int
	RefillRate      float64
	CleanupInterval time.Duration
	IPTimeout       time.Duration
}


type Config struct {
	ServerHost string
	ServerPort string
	ServerAddr string

	PaymentServiceURL string

	CorsAllowedOrigins []string

	JWTSecret string

	Environment string

	RateLimit RateLimitConfig
}


func New() *Config{
	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", "8080")

	return &Config{
		ServerHost:        host,
		ServerPort:        port,
		ServerAddr:        host + ":" + port,
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8081"),
		CorsAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		RateLimit: RateLimitConfig{
			MaxRequests:     10,  // Maximum requests in bucket
			RefillRate:      0.1667, // 1 request every 6 seconds (10 request/minute)
			CleanupInterval: 5 * time.Minute, // Cleanup interval 5 minutes
			IPTimeout:       10 * time.Minute, // Delete IP after 10 minutes of inactivity
		},
	}
}

func getEnv(key, defaultValue string) string{
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}