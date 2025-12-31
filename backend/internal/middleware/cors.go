package middleware

import (
	"os"
	"strings"

	"github.com/rs/cors"
)

func SetupCORS() *cors.Cors {
	// Orígenes permitidos desde variable de entorno o defaults
	allowedOrigins := getOrigins()

	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
			"X-Total-Count",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-Cache",
			"Cache-Control",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutos de caché para preflight requests
	})
}

func getOrigins() []string {
	// Leer desde variable de entorno ALLOWED_ORIGINS
	// Formato: "http://localhost:3000,http://localhost:3001,https://midominio.com"
	originsEnv := os.Getenv("ALLOWED_ORIGINS")

	if originsEnv != "" {
		origins := strings.Split(originsEnv, ",")
		trimmed := make([]string, 0, len(origins))
		for _, origin := range origins {
			if cleaned := strings.TrimSpace(origin); cleaned != "" {
				trimmed = append(trimmed, cleaned)
			}
		}
		if len(trimmed) > 0 {
			return trimmed
		}
	}

	// Defaults para desarrollo
	return []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:5173", // Vite
		"http://localhost:8080",
	}
}