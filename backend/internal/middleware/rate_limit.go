package middleware

import (
	"net/http"
	"sync"
	"time"

	"andaria-backend/pkg/utils"
)

type visitor struct {
	lastSeen  time.Time
	requests  int
	resetTime time.Time
}

var (
	visitors   = make(map[string]*visitor)
	visitorsMu sync.RWMutex
)

// Limpiar visitantes inactivos cada 5 minutos
func init() {
	go cleanupVisitors()
}

func cleanupVisitors() {
	for {
		time.Sleep(5 * time.Minute)
		visitorsMu.Lock()
		now := time.Now()
		for ip, v := range visitors {
			if now.Sub(v.lastSeen) > 10*time.Minute {
				delete(visitors, ip)
			}
		}
		visitorsMu.Unlock()
	}
}

func getVisitor(ip string) *visitor {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		v = &visitor{
			lastSeen:  time.Now(),
			requests:  0,
			resetTime: time.Now().Add(1 * time.Minute),
		}
		visitors[ip] = v
	}

	return v
}

// RateLimitMiddleware limita el número de requests por minuto desde una IP
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener IP del cliente
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			v := getVisitor(ip)
			now := time.Now()

			visitorsMu.Lock()
			// Resetear contador si pasó el minuto
			if now.After(v.resetTime) {
				v.requests = 0
				v.resetTime = now.Add(1 * time.Minute)
			}

			v.lastSeen = now
			v.requests++
			currentRequests := v.requests
			visitorsMu.Unlock()

			// Verificar límite
			if currentRequests > requestsPerMinute {
				w.Header().Set("X-RateLimit-Limit", string(rune(requestsPerMinute)))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")
				utils.ErrorResponse(w, "RATE_LIMIT_EXCEEDED",
					"Demasiadas solicitudes. Por favor, intente más tarde.",
					nil, http.StatusTooManyRequests)
				return
			}

			// Agregar headers de rate limit
			remaining := requestsPerMinute - currentRequests
			if remaining < 0 {
				remaining = 0
			}
			w.Header().Set("X-RateLimit-Limit", string(rune(requestsPerMinute)))
			w.Header().Set("X-RateLimit-Remaining", string(rune(remaining)))

			next.ServeHTTP(w, r)
		})
	}
}
