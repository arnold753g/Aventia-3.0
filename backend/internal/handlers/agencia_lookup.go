package handlers

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"gorm.io/gorm"
)

func resolveAgenciaByIDOrSlug(db *gorm.DB, idOrSlug string) (models.AgenciaTurismo, error) {
	var agencia models.AgenciaTurismo

	if id, err := strconv.ParseUint(idOrSlug, 10, 32); err == nil {
		return agencia, db.First(&agencia, id).Error
	}

	if err := db.Where("slug = ?", idOrSlug).First(&agencia).Error; err == nil {
		return agencia, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return agencia, err
	}

	lowerSlug := strings.ToLower(idOrSlug)
	if err := db.Where("LOWER(slug) = ?", lowerSlug).First(&agencia).Error; err == nil {
		return agencia, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return agencia, err
	}

	prefixID, err := resolveAgenciaIDBySlugPrefix(db, idOrSlug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return agencia, err
	}
	if prefixID > 0 {
		return agencia, db.First(&agencia, prefixID).Error
	}

	fallbackID, err := resolveAgenciaIDBySlugFallback(idOrSlug)
	if err != nil {
		return agencia, err
	}

	return agencia, db.First(&agencia, fallbackID).Error
}

func resolveAgenciaIDBySlugFallback(slug string) (uint, error) {
	normalized := strings.TrimSpace(slug)
	if normalized == "" {
		return 0, gorm.ErrRecordNotFound
	}

	tokens := strings.FieldsFunc(normalized, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	if len(tokens) == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	pattern := "%" + strings.Join(tokens, "%") + "%"

	var candidates []models.AgenciaTurismo
	if err := database.GetDB().
		Select("id", "nombre_comercial", "slug").
		Where("nombre_comercial ILIKE ?", pattern).
		Find(&candidates).Error; err != nil {
		return 0, err
	}

	slugCompact := strings.ReplaceAll(normalized, "-", "")
	for _, candidate := range candidates {
		generated := utils.GenerateSlug(candidate.NombreComercial)
		if generated == slug || strings.ReplaceAll(generated, "-", "") == slugCompact {
			if strings.TrimSpace(candidate.Slug) == "" {
				newSlug := generated
				_ = database.GetDB().
					Model(&models.AgenciaTurismo{}).
					Where("id = ?", candidate.ID).
					Update("slug", newSlug).Error
			}
			return candidate.ID, nil
		}
	}

	return 0, gorm.ErrRecordNotFound
}

func resolveAgenciaIDBySlugPrefix(db *gorm.DB, slug string) (uint, error) {
	pattern := strings.TrimSpace(slug) + "-%"
	if pattern == "-%" {
		return 0, gorm.ErrRecordNotFound
	}

	var candidates []models.AgenciaTurismo
	if err := db.Select("id", "slug").Where("slug ILIKE ?", pattern).Find(&candidates).Error; err != nil {
		return 0, err
	}

	if len(candidates) == 1 {
		return candidates[0].ID, nil
	}

	var numericMatches []models.AgenciaTurismo
	for _, candidate := range candidates {
		if hasNumericSuffix(candidate.Slug, slug) {
			numericMatches = append(numericMatches, candidate)
		}
	}

	if len(numericMatches) == 1 {
		return numericMatches[0].ID, nil
	}

	return 0, gorm.ErrRecordNotFound
}

func hasNumericSuffix(candidateSlug, base string) bool {
	prefix := strings.ToLower(strings.TrimSpace(base)) + "-"
	slug := strings.ToLower(strings.TrimSpace(candidateSlug))
	if !strings.HasPrefix(slug, prefix) {
		return false
	}
	suffix := strings.TrimPrefix(slug, prefix)
	if suffix == "" {
		return false
	}
	for _, r := range suffix {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
