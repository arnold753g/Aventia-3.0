# Seeds - Datos Base de Bolivia

Este paquete contiene los seeds (datos iniciales) para poblar la base de datos con información geográfica y temporal de Bolivia.

## 📦 Contenido

### 1. Departamentos y Provincias (`bolivia_data.go`)

Datos completos de los **9 departamentos** y **111 provincias** de Bolivia:

| Departamento | Provincias |
|--------------|-----------|
| Chuquisaca   | 10        |
| La Paz       | 20        |
| Cochabamba   | 16        |
| Oruro        | 15        |
| Potosí       | 16        |
| Tarija       | 6         |
| Santa Cruz   | 15        |
| Beni         | 8         |
| Pando        | 5         |

### 2. Días y Meses (`departamentos_seed.go`)

- **7 días de la semana** (Lunes - Domingo)
- **12 meses del año** (Enero - Diciembre)

## 🚀 Uso

Los seeds se ejecutan automáticamente al iniciar la aplicación en `database.Connect()`.

### Características

- ✅ **Idempotente**: Puede ejecutarse múltiples veces sin duplicar datos
- ✅ **Transaccional**: Si falla, revierte todos los cambios
- ✅ **Validado**: Verifica que los totales sean correctos
- ✅ **Logging detallado**: Muestra el progreso y resultados

### Ejemplo de salida

```
🌎 Iniciando seed de departamentos y provincias...
  ✓ Departamento creado: Chuquisaca
  ✓ Provincias de Chuquisaca: 10 registradas
  ✓ Departamento creado: La Paz
  ✓ Provincias de La Paz: 20 registradas
  ...

📊 Resumen del seed:
  • Departamentos nuevos: 9
  • Provincias nuevas: 112
  • Total departamentos en BD: 9 (esperado: 9)
  • Total provincias en BD: 111 (esperado: 111)

✅ Seed de departamentos y provincias completado exitosamente

📅 Iniciando seed de días y meses...
✅ Seed de días y meses completado

✅ Todos los seeds completados exitosamente
```

## 🔧 Ejecución Manual

Si necesitas ejecutar los seeds manualmente:

```go
import (
    "andaria-backend/internal/database"
    "andaria-backend/internal/seeds"
)

func main() {
    db := database.GetDB()

    // Ejecutar todos los seeds
    if err := seeds.RunAllSeeds(db); err != nil {
        log.Fatal(err)
    }

    // O ejecutar seeds individuales
    seeds.SeedDepartamentosProvincias(db)
    seeds.SeedDiasYMeses(db)
}
```

## 📝 Estructura de Modelos

### Departamento

```go
type Departamento struct {
    ID        uint
    Nombre    string
    CreatedAt time.Time
}
```

### Provincia

```go
type Provincia struct {
    ID              uint
    DepartamentoID  uint
    Departamento    Departamento
    Nombre          string
    CreatedAt       time.Time
}
```

**Índices:**
- Unique: `(departamento_id, nombre)` - Previene provincias duplicadas
- Foreign Key: `departamento_id` con `ON DELETE RESTRICT`

## 🛡️ Seguridad

- El índice único compuesto `(departamento_id, nombre)` previene:
  - Provincias duplicadas en el mismo departamento
  - Inconsistencias en los datos

- La restricción `ON DELETE RESTRICT` previene:
  - Borrado accidental de departamentos con provincias asociadas
  - Pérdida de integridad referencial

## 📚 Fuente de Datos

Los datos corresponden a la división político-administrativa oficial de Bolivia:
- 9 Departamentos
- 111 Provincias

Última actualización: Diciembre 2024