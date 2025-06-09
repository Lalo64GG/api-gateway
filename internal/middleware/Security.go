package middleware

import (
	"net/http"

	"github.com/Lalo64GG/api-gateway/internal/config"
)

func SecurityHeaders(cfg *config.Config) func(nex http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w  http.ResponseWriter, r *http.Request) {

			// Evit clickjacking
			w.Header().Set("X-Frame-Options", "DENY")

			//Polict for content security
			w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")

			// Evit XSS 
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			//HSTS (required for HTTPS)
			if cfg.Environment == "production" {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}

			// Policy of referrer
			w.Header().Set("Referrer-Policy", "strict-origin")


			// Evit MIME sniffing

			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Policy of permissions
			w.Header().Set("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")

			next.ServeHTTP(w, r)
		})
	}
}