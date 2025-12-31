package utils

import (
	"os"
	"strconv"
	"strings"
)

// GetEnvInt obtiene un entero desde variables de entorno con fallback.
func GetEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
