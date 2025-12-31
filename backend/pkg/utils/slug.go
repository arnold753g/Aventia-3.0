package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GenerateSlug genera un slug URL-friendly desde un texto
func GenerateSlug(text string) string {
	// Convertir a minúsculas
	slug := strings.ToLower(text)

	// Remover acentos y caracteres especiales
	slug = removeAccents(slug)

	// Reemplazar espacios y caracteres no permitidos con guiones
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remover guiones al inicio y final
	slug = strings.Trim(slug, "-")

	// Limitar longitud
	if len(slug) > 200 {
		slug = slug[:200]
		// Asegurar que no termine con guión después del corte
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}

// removeAccents elimina acentos y diacríticos de un string
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// isMn verifica si el rune es un marcador no espaciador (acento)
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// GenerateUniqueSlug genera un slug único agregando un sufijo numérico si es necesario
func GenerateUniqueSlug(baseSlug string, existingSlugs []string) string {
	slug := baseSlug
	counter := 1

	// Crear un mapa para búsqueda rápida
	slugMap := make(map[string]bool)
	for _, s := range existingSlugs {
		slugMap[s] = true
	}

	// Si el slug base ya existe, agregar sufijo
	for slugMap[slug] {
		slug = baseSlug + "-" + string(rune(counter+'0'))
		counter++
	}

	return slug
}
