package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type cachedResponse struct {
	data      []byte
	expiresAt time.Time
	headers   http.Header
}

type responseWriter struct {
	http.ResponseWriter
	body    *bytes.Buffer
	status  int
	headers http.Header
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
		headers:        make(http.Header),
		status:         http.StatusOK,
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

// CacheMiddleware crea un middleware de caché simple para endpoints públicos
func CacheMiddleware(duration time.Duration) func(http.Handler) http.Handler {
	cache := make(map[string]cachedResponse)
	mu := sync.RWMutex{}

	// Limpieza periódica de caché expirado
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			now := time.Now()
			for key, cached := range cache {
				if now.After(cached.expiresAt) {
					delete(cache, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Solo cachear GET
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// Generar clave de caché (URL + query params)
			cacheKey := r.URL.String()

			// Verificar si existe en caché
			mu.RLock()
			cached, exists := cache[cacheKey]
			mu.RUnlock()

			if exists && time.Now().Before(cached.expiresAt) {
				// Cache hit - copiar headers pero evitar duplicar CORS
				for key, values := range cached.headers {
					// Saltar headers que CORS ya maneja globalmente
					if key == "Access-Control-Allow-Origin" ||
						key == "Access-Control-Allow-Credentials" ||
						key == "Access-Control-Allow-Methods" ||
						key == "Access-Control-Allow-Headers" ||
						key == "Access-Control-Expose-Headers" ||
						key == "Vary" {
						continue
					}
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				w.Header().Set("X-Cache", "HIT")
				w.Header().Set("Cache-Control", "public, max-age="+duration.String())
				w.WriteHeader(http.StatusOK)
				w.Write(cached.data)
				return
			}

			// Cache miss - capturar respuesta
			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			// Solo cachear respuestas exitosas
			if rw.status == http.StatusOK && rw.body.Len() > 0 {
				mu.Lock()
				cache[cacheKey] = cachedResponse{
					data:      rw.body.Bytes(),
					expiresAt: time.Now().Add(duration),
					headers:   rw.Header(),
				}
				mu.Unlock()

				w.Header().Set("X-Cache", "MISS")
			}
		})
	}
}
