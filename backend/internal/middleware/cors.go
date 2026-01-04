package middleware

import (
    "os"
    "strings"

    "github.com/rs/cors"
)

func SetupCORS() *cors.Cors {
    allowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))
    allowedSet := make(map[string]struct{}, len(allowedOrigins))
    for _, origin := range allowedOrigins {
        allowedSet[origin] = struct{}{}
    }

    return cors.New(cors.Options{
        AllowOriginFunc: func(origin string) bool {
            if origin == "" {
                return false
            }
            _, ok := allowedSet[origin]
            return ok
        },
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
        MaxAge:           300,
    })
}

func parseAllowedOrigins(raw string) []string {
    raw = strings.TrimSpace(raw)
    if raw == "" {
        return []string{
            "http://localhost:3000",
            "http://localhost:3001",
            "http://localhost:5173",
            "http://localhost:8080",
        }
    }

    parts := strings.Split(raw, ",")
    out := make([]string, 0, len(parts))
    for _, part := range parts {
        if origin := strings.TrimSpace(part); origin != "" {
            out = append(out, origin)
        }
    }

    if len(out) == 0 {
        return []string{
            "http://localhost:3000",
            "http://localhost:3001",
            "http://localhost:5173",
            "http://localhost:8080",
        }
    }

    return out
}
