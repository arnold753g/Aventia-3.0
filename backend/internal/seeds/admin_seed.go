package seeds

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"gorm.io/gorm"
)

func SeedAdminFromEnv(db *gorm.DB) error {
	email := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
	password := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
	if email == "" || password == "" {
		log.Println("Admin seed skipped: ADMIN_EMAIL or ADMIN_PASSWORD not set")
		return nil
	}

	nombreRaw := strings.TrimSpace(os.Getenv("ADMIN_NOMBRE"))
	apellidoPaternoRaw := strings.TrimSpace(os.Getenv("ADMIN_APELLIDO_PATERNO"))
	apellidoMaternoRaw := strings.TrimSpace(os.Getenv("ADMIN_APELLIDO_MATERNO"))
	ciRaw := strings.TrimSpace(os.Getenv("ADMIN_CI"))
	expedidoRaw := strings.TrimSpace(os.Getenv("ADMIN_EXPEDIDO"))
	phoneRaw := strings.TrimSpace(os.Getenv("ADMIN_PHONE"))
	ciudadRaw := strings.TrimSpace(os.Getenv("ADMIN_CIUDAD"))
	fechaNacimientoRaw := strings.TrimSpace(os.Getenv("ADMIN_FECHA_NACIMIENTO"))

	fechaNacimiento, err := parseAdminDate(fechaNacimientoRaw)
	if err != nil {
		return err
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing admin password: %w", err)
	}

	var existing models.Usuario
	err = db.Where("email = ?", email).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking admin user: %w", err)
	}

	acceptedAt := time.Now()

	if err == nil {
		updates := map[string]interface{}{
			"password_hash":    passwordHash,
			"rol":              "admin",
			"status":           "active",
			"email_verified":   true,
			"terms_accepted":   true,
			"terms_accepted_at": &acceptedAt,
		}

		if nombreRaw != "" {
			updates["nombre"] = nombreRaw
		}
		if apellidoPaternoRaw != "" {
			updates["apellido_paterno"] = apellidoPaternoRaw
		}
		if apellidoMaternoRaw != "" {
			updates["apellido_materno"] = apellidoMaternoRaw
		}
		if ciRaw != "" {
			updates["ci"] = strings.ToUpper(ciRaw)
		}
		if expedidoRaw != "" {
			updates["expedido"] = expedidoRaw
		}
		if phoneRaw != "" {
			updates["phone"] = phoneRaw
		}
		if ciudadRaw != "" {
			updates["ciudad"] = ciudadRaw
		}
		if fechaNacimientoRaw != "" {
			updates["fecha_nacimiento"] = fechaNacimiento
		}

		if err := db.Model(&existing).Updates(updates).Error; err != nil {
			return fmt.Errorf("error updating admin user %s: %w", email, err)
		}
		log.Printf("Admin user updated: %s", email)
		return nil
	}

	nombre := nombreRaw
	if nombre == "" {
		nombre = "Admin"
	}
	apellidoPaterno := apellidoPaternoRaw
	if apellidoPaterno == "" {
		apellidoPaterno = "Andaria"
	}
	apellidoMaterno := apellidoMaternoRaw
	if apellidoMaterno == "" {
		apellidoMaterno = "Admin"
	}
	ci := ciRaw
	if ci == "" {
		ci = "99999999"
	}
	expedido := expedidoRaw
	if expedido == "" {
		expedido = "LP"
	}
	phone := phoneRaw
	if phone == "" {
		phone = "70000000"
	}
	ciudad := ciudadRaw
	if ciudad == "" {
		ciudad = "La Paz"
	}

	usuario := models.Usuario{
		Nombre:          nombre,
		ApellidoPaterno: apellidoPaterno,
		ApellidoMaterno: apellidoMaterno,
		Email:           email,
		PasswordHash:    passwordHash,
		CI:              strings.ToUpper(ci),
		Expedido:        expedido,
		Phone:           phone,
		FechaNacimiento: fechaNacimiento,
		Ciudad:          ciudad,
		Rol:             "admin",
		Status:          "active",
		Nationality:     "Bolivia",
		EmailVerified:   true,
		TermsAccepted:   true,
		TermsAcceptedAt: &acceptedAt,
	}

	if err := db.Create(&usuario).Error; err != nil {
		return fmt.Errorf("error creating admin user %s: %w", email, err)
	}

	log.Printf("Admin user created: %s", email)
	return nil
}

func parseAdminDate(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), nil
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid ADMIN_FECHA_NACIMIENTO (expected YYYY-MM-DD)")
	}

	return parsed, nil
}
