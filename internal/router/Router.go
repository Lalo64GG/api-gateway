package router

import (
	"net/http"

	"github.com/Lalo64GG/api-gateway/internal/config"
	mw "github.com/Lalo64GG/api-gateway/internal/middleware"

	"github.com/Lalo64GG/api-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(cfg *config.Config, serviceRegistry *services.Registry) http.Handler {

	r := chi.NewRouter()


	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(mw.RequestLogger)
	r.Use(mw.ConfigureCors(cfg).Handler)
	r.Use(mw.SecurityHeaders(cfg))


	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "OK", "version": "1.0.0"}`))
	})

	r.Route("/v1", func(r chi.Router) {

		r.Mount("/payment", serviceRegistry.Payment.Routes())
	})

	return r
}