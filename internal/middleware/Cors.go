package middleware

import (
	"github.com/go-chi/cors"
	"github.com/Lalo64GG/api-gateway/internal/config"
)

func ConfigureCors(cfg *config.Config) *cors.Cors{
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.CorsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge: 300,
	})
}