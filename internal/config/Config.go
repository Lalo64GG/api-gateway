package config

import (
	"os"
	"strings"
)

type Config struct {
	ServerHost string
	ServerPort string
	ServerAddr string

	PaymentServiceURL string

	CorsAllowedOrigins []string

	JWTSecret string

	Environment string
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
	}
}

func getEnv(key, defaultValue string) string{
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}