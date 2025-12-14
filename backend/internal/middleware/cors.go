package middleware

import (
    "github.com/rs/cors"
)

func SetupCORS() *cors.Cors {
    return cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:3000",
            "http://localhost:3001",
        },
        AllowedMethods: []string{
            "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
        },
        AllowedHeaders: []string{
            "Accept", "Authorization", "Content-Type", "X-CSRF-Token",
        },
        ExposedHeaders: []string{
            "Link",
        },
        AllowCredentials: true,
        MaxAge:           300,
    })
}