package seeds

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"gorm.io/gorm"
)

type demoUsuarioSeed struct {
	Nombre          string
	ApellidoPaterno string
	ApellidoMaterno string
	Email           string
	Password        string
	CI              string
	Expedido        string
	Phone           string
	FechaNacimiento time.Time
	Ciudad          string
	Rol             string
}

// SeedUsuariosDemo crea usuarios de prueba si no existen.
func SeedUsuariosDemo(db *gorm.DB) error {
	log.Println("Seeding demo users...")

	usuarios := []demoUsuarioSeed{
		{
			Nombre:          "Admin",
			ApellidoPaterno: "Andaria",
			ApellidoMaterno: "Admin",
			Email:           "admin@andaria.bo",
			Password:        "admin123",
			CI:              "10000001",
			Expedido:        "LP",
			Phone:           "70000001",
			FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			Ciudad:          "La Paz",
			Rol:             "admin",
		},
		{
			Nombre:          "Juan",
			ApellidoPaterno: "Perez",
			ApellidoMaterno: "Lopez",
			Email:           "juan.perez@email.com",
			Password:        "turista123",
			CI:              "10000002",
			Expedido:        "CB",
			Phone:           "70000002",
			FechaNacimiento: time.Date(1995, time.May, 10, 0, 0, 0, 0, time.UTC),
			Ciudad:          "Cochabamba",
			Rol:             "turista",
		},
		{
			Nombre:          "Maria",
			ApellidoPaterno: "Lopez",
			ApellidoMaterno: "Rojas",
			Email:           "maria.lopez@agencia.com",
			Password:        "agencia123",
			CI:              "10000003",
			Expedido:        "SC",
			Phone:           "70000003",
			FechaNacimiento: time.Date(1992, time.August, 20, 0, 0, 0, 0, time.UTC),
			Ciudad:          "Santa Cruz",
			Rol:             "encargado_agencia",
		},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		created := 0

		for _, u := range usuarios {
			email := strings.ToLower(u.Email)
			ci := strings.ToUpper(u.CI)

			var existing models.Usuario
			err := tx.Where("email = ? OR ci = ?", email, ci).First(&existing).Error
			if err == nil {
				log.Printf("Demo user already exists: %s", email)
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("error checking demo user %s: %w", email, err)
			}

			passwordHash, err := utils.HashPassword(u.Password)
			if err != nil {
				return fmt.Errorf("error hashing password for %s: %w", email, err)
			}

			acceptedAt := time.Now()
			usuario := models.Usuario{
				Nombre:          u.Nombre,
				ApellidoPaterno: u.ApellidoPaterno,
				ApellidoMaterno: u.ApellidoMaterno,
				Email:           email,
				PasswordHash:    passwordHash,
				CI:              ci,
				Expedido:        u.Expedido,
				Phone:           u.Phone,
				FechaNacimiento: u.FechaNacimiento,
				Ciudad:          u.Ciudad,
				Rol:             u.Rol,
				Status:          "active",
				Nationality:     "Bolivia",
				EmailVerified:   true,
				TermsAccepted:   true,
				TermsAcceptedAt: &acceptedAt,
			}

			if err := tx.Create(&usuario).Error; err != nil {
				return fmt.Errorf("error creating demo user %s: %w", email, err)
			}

			created++
			log.Printf("Demo user created: %s", email)
		}

		log.Printf("Demo users created: %d", created)
		return nil
	})
}
