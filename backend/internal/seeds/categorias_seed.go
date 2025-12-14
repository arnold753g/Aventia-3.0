package seeds

import (
	"fmt"
	"log"

	"andaria-backend/internal/models"
	"gorm.io/gorm"
)

// SeedCategoria estructura para organizar categorías y subcategorías
type SeedCategoria struct {
	Nombre        string
	Descripcion   string
	Icono         string
	Orden         int
	Subcategorias []SeedSubcategoria
}

type SeedSubcategoria struct {
	Nombre      string
	Descripcion string
	Icono       string
	Orden       int
}

// CategoriasBolivia contiene categorías y subcategorías de atracciones turísticas
// Ajustado al nuevo set:
// 1) Enoturismo
// 2) Cultural
// 3) Natural
// 4) Deportivo
// CategoriasBolivia contiene categorías y subcategorías de atracciones turísticas
// Ajustado al nuevo set con descripciones e íconos
var CategoriasBolivia = []SeedCategoria{
	{
		Nombre:      "Enoturismo",
		Descripcion: "Experiencias relacionadas con la producción y degustación de vinos",
		Icono:       "pi-star",
		Orden:       1,
		Subcategorias: []SeedSubcategoria{
			{
				Nombre:      "Bodegas",
				Descripcion: "Visitas a bodegas, recorridos y experiencias de vinificación",
				Icono:       "pi-building",
				Orden:       1,
			},
			{
				Nombre:      "Viñedos",
				Descripcion: "Recorridos por viñedos y experiencias agrícolas del vino",
				Icono:       "pi-leaf",
				Orden:       2,
			},
			{
				Nombre:      "Catas",
				Descripcion: "Degustaciones guiadas de vinos y maridajes",
				Icono:       "pi-star-fill",
				Orden:       3,
			},
			{
				Nombre:      "Vendimia",
				Descripcion: "Fiestas y actividades relacionadas con la cosecha de uva",
				Icono:       "pi-calendar",
				Orden:       4,
			},
			{
				Nombre:      "Vinoteca",
				Descripcion: "Tiendas especializadas y espacios de compra y degustación",
				Icono:       "pi-shopping-bag",
				Orden:       5,
			},
		},
	},
	{
		Nombre:      "Cultural",
		Descripcion: "Atracciones de valor histórico, arqueológico y cultural",
		Icono:       "pi-building",
		Orden:       2,
		Subcategorias: []SeedSubcategoria{
			{
				Nombre:      "Arqueología",
				Descripcion: "Ruinas, sitios arqueológicos y patrimonio prehispánico",
				Icono:       "pi-compass",
				Orden:       1,
			},
			{
				Nombre:      "Monumentos",
				Descripcion: "Monumentos, esculturas y espacios conmemorativos",
				Icono:       "pi-flag",
				Orden:       2,
			},
			{
				Nombre:      "Mercados",
				Descripcion: "Mercados tradicionales con valor cultural y turístico",
				Icono:       "pi-shopping-cart",
				Orden:       3,
			},
			{
				Nombre:      "Museos",
				Descripcion: "Museos, galerías y centros de interpretación cultural",
				Icono:       "pi-book",
				Orden:       4,
			},
			{
				Nombre:      "Pueblos",
				Descripcion: "Pueblos con identidad cultural, arquitectura y tradiciones",
				Icono:       "pi-map-marker",
				Orden:       5,
			},
			{
				Nombre:      "Festivales",
				Descripcion: "Festividades locales, celebraciones y ferias culturales",
				Icono:       "pi-calendar",
				Orden:       6,
			},
			{
				Nombre:      "Espectáculos",
				Descripcion: "Presentaciones artísticas, música, danza y shows",
				Icono:       "pi-ticket",
				Orden:       7,
			},
			{
				Nombre:      "Histórico",
				Descripcion: "Sitios y rutas con relevancia histórica nacional o regional",
				Icono:       "pi-clock",
				Orden:       8,
			},
		},
	},
	{
		Nombre:      "Natural",
		Descripcion: "Bellezas naturales y paisajes",
		Icono:       "pi-sun",
		Orden:       3,
		Subcategorias: []SeedSubcategoria{
			{
				Nombre:      "Montañas",
				Descripcion: "Cordilleras, cerros, nevados y paisajes de altura",
				Icono:       "pi-arrow-up",
				Orden:       1,
			},
			{
				Nombre:      "Lagos",
				Descripcion: "Lagos, lagunas y cuerpos de agua naturales",
				Icono:       "pi-circle",
				Orden:       2,
			},
			{
				Nombre:      "Ríos",
				Descripcion: "Ríos, afluentes y recorridos fluviales",
				Icono:       "pi-cloud",
				Orden:       3,
			},
			{
				Nombre:      "Cascadas",
				Descripcion: "Cascadas y saltos de agua accesibles para turismo",
				Icono:       "pi-cloud-download",
				Orden:       4,
			},
			{
				Nombre:      "Salares",
				Descripcion: "Salares y paisajes salinos de alto valor escénico",
				Icono:       "pi-circle-fill",
				Orden:       5,
			},
			{
				Nombre:      "Desiertos",
				Descripcion: "Zonas áridas, formaciones de arena y paisajes extremos",
				Icono:       "pi-stop",
				Orden:       6,
			},
			{
				Nombre:      "Playas",
				Descripcion: "Playas lacustres y fluviales, zonas de recreación",
				Icono:       "pi-sun",
				Orden:       7,
			},
			{
				Nombre:      "Fauna",
				Descripcion: "Observación de fauna silvestre en su hábitat natural",
				Icono:       "pi-eye",
				Orden:       8,
			},
			{
				Nombre:      "Flora",
				Descripcion: "Áreas de riqueza botánica, bosques y especies nativas",
				Icono:       "pi-leaf",
				Orden:       9,
			},
			{
				Nombre:      "Miradores",
				Descripcion: "Puntos panorámicos y vistas naturales destacadas",
				Icono:       "pi-search",
				Orden:       10,
			},
			{
				Nombre:      "Reservas protegidas",
				Descripcion: "Áreas protegidas, parques y zonas de conservación",
				Icono:       "pi-shield",
				Orden:       11,
			},
			{
				Nombre:      "Dunas",
				Descripcion: "Campos de dunas y paisajes eólicos",
				Icono:       "pi-flag-fill",
				Orden:       12,
			},
			{
				Nombre:      "Formaciones geológicas",
				Descripcion: "Cañones, rocas singulares y estructuras naturales únicas",
				Icono:       "pi-globe",
				Orden:       13,
			},
		},
	},
	{
		Nombre:      "Deportivo",
		Descripcion: "Actividades deportivas y de aventura",
		Icono:       "pi-bolt",
		Orden:       4,
		Subcategorias: []SeedSubcategoria{
			{
				Nombre:      "Senderismo",
				Descripcion: "Rutas de caminata, trekking y circuitos guiados",
				Icono:       "pi-directions",
				Orden:       1,
			},
			{
				Nombre:      "Escalada",
				Descripcion: "Escalada en roca, hielo y montañismo técnico",
				Icono:       "pi-angle-double-up",
				Orden:       2,
			},
			{
				Nombre:      "Ciclismo",
				Descripcion: "Rutas de ciclismo urbano, ruta y montaña",
				Icono:       "pi-spin",
				Orden:       3,
			},
			{
				Nombre:      "Acuáticos",
				Descripcion: "Rafting, kayak, paddle y deportes en agua",
				Icono:       "pi-heart",
				Orden:       4,
			},
			{
				Nombre:      "Aéreos",
				Descripcion: "Parapente, ala delta y otras experiencias de vuelo",
				Icono:       "pi-send",
				Orden:       5,
			},
			{
				Nombre:      "Nieve",
				Descripcion: "Actividades en nieve: esquí, snowboard y caminatas",
				Icono:       "pi-cloud",
				Orden:       6,
			},
		},
	},
}

// SeedCategoriasSubcategorias precarga categorías y subcategorías
func SeedCategoriasSubcategorias(db *gorm.DB) error {
	log.Println("🎨 Iniciando seed de categorías y subcategorías...")

	return db.Transaction(func(tx *gorm.DB) error {
		categoriasCreadas := 0
		subcategoriasCreadas := 0

		for _, c := range CategoriasBolivia {
			cat := models.CategoriaAtraccion{
				Nombre: c.Nombre,
			}

			// Assign ayuda a mantener actualizado descripcion/icono/orden si ya existe
			result := tx.
				Where("nombre = ?", c.Nombre).
				Assign(models.CategoriaAtraccion{
					Descripcion: c.Descripcion,
					Icono:       c.Icono,
					Orden:       c.Orden,
				}).
				FirstOrCreate(&cat)

			if result.Error != nil {
				return fmt.Errorf("error al crear categoría %s: %w", c.Nombre, result.Error)
			}

			if result.RowsAffected > 0 {
				categoriasCreadas++
				log.Printf("  ✓ Categoría creada: %s", c.Nombre)
			}

			// Crear subcategorías de la categoría
			for _, s := range c.Subcategorias {
				subcat := models.SubcategoriaAtraccion{
					CategoriaID: cat.ID,
					Nombre:      s.Nombre,
				}

				// Igual: asignamos campos opcionales + orden
				result := tx.
					Where("categoria_id = ? AND nombre = ?", cat.ID, s.Nombre).
					Assign(models.SubcategoriaAtraccion{
						Descripcion: s.Descripcion,
						Icono:       s.Icono,
						Orden:       s.Orden,
					}).
					FirstOrCreate(&subcat)

				if result.Error != nil {
					return fmt.Errorf("error al crear subcategoría %s (%s): %w", s.Nombre, c.Nombre, result.Error)
				}

				if result.RowsAffected > 0 {
					subcategoriasCreadas++
				}
			}

			log.Printf("  ✓ Subcategorías de %s: %d registradas", c.Nombre, len(c.Subcategorias))
		}

		// Validar totales
		var countCat, countSubcat int64
		tx.Model(&models.CategoriaAtraccion{}).Count(&countCat)
		tx.Model(&models.SubcategoriaAtraccion{}).Count(&countSubcat)

		log.Printf("\n📊 Resumen del seed de categorías:")
		log.Printf("  • Categorías nuevas: %d", categoriasCreadas)
		log.Printf("  • Subcategorías nuevas: %d", subcategoriasCreadas)
		log.Printf("  • Total categorías en BD: %d", countCat)
		log.Printf("  • Total subcategorías en BD: %d", countSubcat)

		log.Println("✅ Seed de categorías y subcategorías completado exitosamente")
		return nil
	})
}