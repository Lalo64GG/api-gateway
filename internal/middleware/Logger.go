package middleware

import (
	"log"
	"net/http"
	"time"
)

// Wraper for http.ResponseWriter to capture the status code
type responseWriter struct {
	w      http.ResponseWriter
	status int
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{w: w, status: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		log.Printf("%s - %s %s - %d - %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			ww.status,
			duration,
		)
	})
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.w.WriteHeader(statusCode)
}