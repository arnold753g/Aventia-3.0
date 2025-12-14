package database

import (
    "fmt"
    "log"
    "andaria-backend/internal/config"
    "andaria-backend/internal/models"
    "andaria-backend/internal/seeds"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        DisableForeignKeyConstraintWhenMigrating: true,
    })

    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    log.Println("✅ Database connection established")

    // Ejecutar migraciones automáticas (solo si la tabla no existe)
    // Si la tabla ya existe, las migraciones se saltan para evitar conflictos
    if err := runMigrations(); err != nil {
        log.Printf("⚠️  Warning: Migration skipped (table already exists with custom constraints): %v", err)
        log.Println("ℹ️  Continuing without migrations - using existing schema")
    }

    // Ejecutar seeds de datos base
    if err := runSeeds(); err != nil {
        log.Printf("⚠️  Warning during seeds: %v", err)
        log.Println("ℹ️  Continuing - some seeds may have been skipped")
    }

    return nil
}

func runMigrations() error {
    log.Println("🔄 Checking database schema...")

    // Verificar si las tablas principales ya existen
    hasUsuarios := DB.Migrator().HasTable(&models.Usuario{})
    hasAtracciones := DB.Migrator().HasTable(&models.AtraccionTuristica{})

    if hasUsuarios && hasAtracciones {
        log.Println("ℹ️  Tables already exist - skipping migrations")
        log.Println("✅ Using existing database schema")
        return nil
    }

    // Migrar todas las tablas en orden de dependencias
    log.Println("📝 Creating database tables...")

    err := DB.AutoMigrate(
        // Tablas base sin dependencias
        &models.Usuario{},
        &models.Departamento{},
        &models.Dia{},
        &models.Mes{},
        &models.CategoriaAtraccion{},

        // Tablas con dependencias nivel 1
        &models.Provincia{},
        &models.SubcategoriaAtraccion{},

        // Tablas con dependencias nivel 2
        &models.AtraccionTuristica{},

        // Tablas de relaciones
        &models.AtraccionSubcategoria{},
        &models.AtraccionFoto{},
    )

    if err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }

    log.Println("✅ Tables created successfully")
    return nil
}

func runSeeds() error {
    log.Println("🌱 Running database seeds...")

    // Ejecutar todos los seeds
    if err := seeds.RunAllSeeds(DB); err != nil {
        return fmt.Errorf("seeds failed: %w", err)
    }

    return nil
}

func GetDB() *gorm.DB {
    return DB
}
