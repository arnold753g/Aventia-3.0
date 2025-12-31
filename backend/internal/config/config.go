package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DBHost               string
    DBPort               string
    DBUser               string
    DBPassword           string
    DBName               string
    DBSSLMode            string
    ServerPort           string
    ServerHost           string
    JWTSecret            string
    JWTExpiration        string
    JWTRefreshExpiration string
    AppEnv               string
    // Configuración de expiración de compras (en minutos)
    CompraExpiracionMinutos string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }

    return &Config{
        DBHost:               getEnv("DB_HOST", "localhost"),
        DBPort:               getEnv("DB_PORT", "5432"),
        DBUser:               getEnv("DB_USER", "postgres"),
        DBPassword:           getEnv("DB_PASSWORD", "1234"),
        DBName:               getEnv("DB_NAME", "Andaria_01"),
        DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
        ServerPort:           getEnv("SERVER_PORT", "5750"),
        ServerHost:           getEnv("SERVER_HOST", "localhost"),
        JWTSecret:            getEnv("JWT_SECRET", "secret"),
        JWTExpiration:        getEnv("JWT_EXPIRATION", "24h"),
        JWTRefreshExpiration: getEnv("JWT_REFRESH_EXPIRATION", "168h"),
        AppEnv:               getEnv("APP_ENV", "development"),
        CompraExpiracionMinutos: getEnv("COMPRA_EXPIRACION_MINUTOS", "30"),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}