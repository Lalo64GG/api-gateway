// ===== internal/middleware/rate_limit.go =====
package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Lalo64GG/api-gateway/internal/ratelimiter"
)

// RateLimit retorna un middleware que aplica rate limiting por IP
func RateLimit(ipRateLimiter *ratelimiter.IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractClientIP(r)
			
			if ip == "" {
				log.Printf("Could not extract IP from request: %s", r.RemoteAddr)
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			
			limiter := ipRateLimiter.GetLimiter(ip)
			
			if !limiter.Allow() {
				log.Printf("Rate limit exceeded for IP: %s", ip)
				
				// Añadir headers informativos para el cliente
				w.Header().Set("X-RateLimit-Limit", "10")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")
				
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}
			
			// Request permitido, continuar al siguiente handler
			next.ServeHTTP(w, r)
		})
	}
}

// extractClientIP extrae la IP real del cliente considerando proxies y load balancers
func extractClientIP(r *http.Request) string {
	// Prioridad 1: X-Forwarded-For (proxy/load balancer)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For puede tener múltiples IPs: "client, proxy1, proxy2"
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			// Tomar la primera IP (cliente original)
			clientIP := strings.TrimSpace(ips[0])
			if clientIP != "" && isValidIP(clientIP) {
				return clientIP
			}
		}
	}
	
	// Prioridad 2: X-Real-IP (algunos proxies usan este header)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if isValidIP(xri) {
			return xri
		}
	}
	
	// Prioridad 3: X-Forwarded (menos común pero existe)
	if xf := r.Header.Get("X-Forwarded"); xf != "" {
		if strings.Contains(xf, "for=") {
			// Formato: for=192.168.1.1
			parts := strings.Split(xf, "for=")
			if len(parts) > 1 {
				ip := strings.TrimSpace(strings.Split(parts[1], ",")[0])
				ip = strings.Trim(ip, "\"")
				if isValidIP(ip) {
					return ip
				}
			}
		}
	}
	
	// Último recurso: RemoteAddr directamente
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// Si RemoteAddr no tiene puerto, puede ser solo la IP
		if isValidIP(r.RemoteAddr) {
			return r.RemoteAddr
		}
		return ""
	}
	
	return ip
}

// isValidIP verifica si una string es una IP válida (IPv4 o IPv6)
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}