package seeds

import (
	"fmt"
	"log"

	"andaria-backend/internal/models"
	"gorm.io/gorm"
)

// SeedDepartamentosProvincias precarga departamentos y provincias de Bolivia
// Esta función es idempotente: puede ejecutarse múltiples veces sin duplicar datos
func SeedDepartamentosProvincias(db *gorm.DB) error {
	log.Println("🌎 Iniciando seed de departamentos y provincias...")

	return db.Transaction(func(tx *gorm.DB) error {
		departamentosCreados := 0
		provinciasCreadas := 0

		for _, d := range Bolivia {
			dep := models.Departamento{Nombre: d.Nombre}

			// FirstOrCreate: inserta solo si no existe
			result := tx.Where("nombre = ?", d.Nombre).FirstOrCreate(&dep)
			if result.Error != nil {
				return fmt.Errorf("error al crear departamento %s: %w", d.Nombre, result.Error)
			}

			if result.RowsAffected > 0 {
				departamentosCreados++
				log.Printf("  ✓ Departamento creado: %s", d.Nombre)
			}

			// Crear provincias del departamento
			for _, p := range d.Provincias {
				prov := models.Provincia{
					DepartamentoID: dep.ID,
					Nombre:         p,
				}

				result := tx.Where("departamento_id = ? AND nombre = ?", dep.ID, p).FirstOrCreate(&prov)
				if result.Error != nil {
					return fmt.Errorf("error al crear provincia %s (%s): %w", p, d.Nombre, result.Error)
				}

				if result.RowsAffected > 0 {
					provinciasCreadas++
				}
			}

			log.Printf("  ✓ Provincias de %s: %d registradas", d.Nombre, len(d.Provincias))
		}

		// Validar totales
		var countDep, countProv int64
		tx.Model(&models.Departamento{}).Count(&countDep)
		tx.Model(&models.Provincia{}).Count(&countProv)

		log.Printf("\n📊 Resumen del seed:")
		log.Printf("  • Departamentos nuevos: %d", departamentosCreados)
		log.Printf("  • Provincias nuevas: %d", provinciasCreadas)
		log.Printf("  • Total departamentos en BD: %d (esperado: %d)", countDep, TotalDepartamentos)
		log.Printf("  • Total provincias en BD: %d (esperado: %d)", countProv, TotalProvincias)

		if countDep != TotalDepartamentos {
			return fmt.Errorf("total de departamentos incorrecto: %d, esperado: %d", countDep, TotalDepartamentos)
		}

		if countProv != TotalProvincias {
			return fmt.Errorf("total de provincias incorrecto: %d, esperado: %d", countProv, TotalProvincias)
		}

		log.Println("✅ Seed de departamentos y provincias completado exitosamente")
		return nil
	})
}

// SeedDiasYMeses precarga los días de la semana y meses del año
func SeedDiasYMeses(db *gorm.DB) error {
	log.Println("📅 Iniciando seed de días y meses...")

	return db.Transaction(func(tx *gorm.DB) error {
		// Días de la semana
		dias := []models.Dia{
			{ID: 1, Nombre: "Lunes"},
			{ID: 2, Nombre: "Martes"},
			{ID: 3, Nombre: "Miércoles"},
			{ID: 4, Nombre: "Jueves"},
			{ID: 5, Nombre: "Viernes"},
			{ID: 6, Nombre: "Sábado"},
			{ID: 7, Nombre: "Domingo"},
		}

		for _, dia := range dias {
			// Usar id_dia en lugar de id (nombre de columna personalizado)
			result := tx.Where("id_dia = ?", dia.ID).FirstOrCreate(&dia)
			if result.Error != nil {
				return fmt.Errorf("error al crear día %s: %w", dia.Nombre, result.Error)
			}
		}

		// Meses del año
		meses := []models.Mes{
			{ID: 1, Nombre: "Enero"},
			{ID: 2, Nombre: "Febrero"},
			{ID: 3, Nombre: "Marzo"},
			{ID: 4, Nombre: "Abril"},
			{ID: 5, Nombre: "Mayo"},
			{ID: 6, Nombre: "Junio"},
			{ID: 7, Nombre: "Julio"},
			{ID: 8, Nombre: "Agosto"},
			{ID: 9, Nombre: "Septiembre"},
			{ID: 10, Nombre: "Octubre"},
			{ID: 11, Nombre: "Noviembre"},
			{ID: 12, Nombre: "Diciembre"},
		}

		for _, mes := range meses {
			// Usar id_mes en lugar de id (nombre de columna personalizado)
			result := tx.Where("id_mes = ?", mes.ID).FirstOrCreate(&mes)
			if result.Error != nil {
				return fmt.Errorf("error al crear mes %s: %w", mes.Nombre, result.Error)
			}
		}

		log.Println("✅ Seed de días y meses completado")
		return nil
	})
}

// RunAllSeeds ejecuta todos los seeds en el orden correcto
func RunAllSeeds(db *gorm.DB) error {
	log.Println("\n🚀 Ejecutando todos los seeds...")

	// Orden importante: primero datos base, luego categorías
	if err := SeedDepartamentosProvincias(db); err != nil {
		return fmt.Errorf("error en seed de departamentos/provincias: %w", err)
	}

	if err := SeedDiasYMeses(db); err != nil {
		return fmt.Errorf("error en seed de días/meses: %w", err)
	}

	if err := SeedCategoriasSubcategorias(db); err != nil {
		return fmt.Errorf("error en seed de categorías/subcategorías: %w", err)
	}

	log.Println("✅ Todos los seeds completados exitosamente")
	return nil
}
