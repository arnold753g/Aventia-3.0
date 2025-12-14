package models

import "time"

// Departamento representa un departamento de Bolivia
type Departamento struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"size:100;uniqueIndex;not null" json:"nombre"`
	CreatedAt time.Time `json:"created_at"`
}

// Provincia representa una provincia de Bolivia
type Provincia struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	DepartamentoID  uint         `gorm:"not null;index;uniqueIndex:uniq_dep_nombre" json:"departamento_id"`
	Departamento    Departamento `gorm:"foreignKey:DepartamentoID;constraint:OnDelete:RESTRICT;" json:"departamento,omitempty"`
	Nombre          string       `gorm:"size:100;not null;uniqueIndex:uniq_dep_nombre" json:"nombre"`
	CreatedAt       time.Time    `json:"created_at"`
}

// Dia representa un dia de la semana
type Dia struct {
	ID     uint   `gorm:"primaryKey;column:id_dia" json:"id"`
	Nombre string `gorm:"size:20;uniqueIndex;not null" json:"nombre"`
}

// Mes representa un mes del año
type Mes struct {
	ID     uint   `gorm:"primaryKey;column:id_mes" json:"id"`
	Nombre string `gorm:"size:15;uniqueIndex;not null" json:"nombre"`
}

// CategoriaAtraccion representa una categoria principal
type CategoriaAtraccion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Nombre      string    `gorm:"size:100;uniqueIndex;not null" json:"nombre"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	Icono       string    `gorm:"size:50" json:"icono"`
	Orden       int       `gorm:"default:0" json:"orden"`
	CreatedAt   time.Time `json:"created_at"`
}

// SubcategoriaAtraccion representa una subcategoria
type SubcategoriaAtraccion struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	CategoriaID uint               `gorm:"not null;index" json:"categoria_id"`
	Categoria   CategoriaAtraccion `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
	Nombre      string             `gorm:"size:100;not null" json:"nombre"`
	Descripcion string             `gorm:"type:text" json:"descripcion"`
	Icono       string             `gorm:"size:50" json:"icono"`
	Orden       int                `gorm:"default:0" json:"orden"`
	CreatedAt   time.Time          `json:"created_at"`
}

// AtraccionTuristica representa una atraccion turistica
type AtraccionTuristica struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Nombre      string  `gorm:"size:255;not null" json:"nombre"`
	Descripcion string  `gorm:"type:text" json:"descripcion"`

	// Ubicacion
	ProvinciaID uint      `gorm:"not null;index" json:"provincia_id"`
	Provincia   Provincia `gorm:"foreignKey:ProvinciaID" json:"provincia,omitempty"`
	Direccion   string    `gorm:"type:text" json:"direccion"`
	Latitud     *float64  `gorm:"type:decimal(10,8);index:idx_coordenadas,composite:coordenadas" json:"latitud"`
	Longitud    *float64  `gorm:"type:decimal(11,8);index:idx_coordenadas,composite:coordenadas" json:"longitud"`

	// Horarios y precios
	HorarioApertura *string  `gorm:"type:time" json:"horario_apertura"`
	HorarioCierre   *string  `gorm:"type:time" json:"horario_cierre"`
	PrecioEntrada   float64  `gorm:"type:decimal(10,2);default:0" json:"precio_entrada"`

	// Caracteristicas
	NivelDificultad   string `gorm:"size:20;check:nivel_dificultad IN ('facil','medio','dificil','extremo');index" json:"nivel_dificultad"`
	RequiereAgencia   bool   `gorm:"default:false;index" json:"requiere_agencia"`
	AccesoParticular  bool   `gorm:"default:true" json:"acceso_particular"`

	// Mejor epoca
	MesInicioID *uint `gorm:"column:mes_inicio" json:"mes_inicio_id"`
	MesFinID    *uint `gorm:"column:mes_fin" json:"mes_fin_id"`
	MesInicio   *Mes  `gorm:"foreignKey:MesInicioID" json:"mes_inicio,omitempty"`
	MesFin      *Mes  `gorm:"foreignKey:MesFinID" json:"mes_fin,omitempty"`

	// Estado y gestion
	Status         string `gorm:"size:30;default:'activa';check:status IN ('activa','inactiva','mantenimiento','fuera_temporada');index" json:"status"`
	VisiblePublico bool   `gorm:"default:true;index" json:"visible_publico"`

	// Contacto
	Telefono  string `gorm:"size:20" json:"telefono"`
	Email     string `gorm:"size:100;index" json:"email"`
	SitioWeb  string `gorm:"size:200" json:"sitio_web"`
	Facebook  string `gorm:"size:200" json:"facebook"`
	Instagram string `gorm:"size:200" json:"instagram"`

	// Relaciones
	Subcategorias []AtraccionSubcategoria `gorm:"foreignKey:AtraccionID" json:"subcategorias,omitempty"`
	Fotos         []AtraccionFoto         `gorm:"foreignKey:AtraccionID" json:"fotos,omitempty"`
	Dias          []Dia                   `gorm:"many2many:atraccion_dias;foreignKey:ID;joinForeignKey:AtraccionID;References:ID;joinReferences:DiaID" json:"dias,omitempty"`

	// Auditoria
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AtraccionSubcategoria relacion atraccion-subcategoria
type AtraccionSubcategoria struct {
	ID             uint                  `gorm:"primaryKey" json:"id"`
	AtraccionID    uint                  `gorm:"not null;uniqueIndex:idx_atraccion_subcategoria,composite:atraccion_subcategoria" json:"atraccion_id"`
	SubcategoriaID uint                  `gorm:"not null;uniqueIndex:idx_atraccion_subcategoria,composite:atraccion_subcategoria" json:"subcategoria_id"`
	Subcategoria   SubcategoriaAtraccion `gorm:"foreignKey:SubcategoriaID" json:"subcategoria,omitempty"`
	EsPrincipal    bool                  `gorm:"default:false" json:"es_principal"`
	CreatedAt      time.Time             `json:"created_at"`
}

// AtraccionFoto representa una foto de atraccion
type AtraccionFoto struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AtraccionID uint      `gorm:"not null;index" json:"atraccion_id"`
	Foto        string    `gorm:"type:text;not null" json:"foto"`
	EsPrincipal bool      `gorm:"default:false" json:"es_principal"`
	Orden       int       `gorm:"default:0;index:idx_atraccion_orden,composite:atraccion_orden" json:"orden"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName overrides
func (Departamento) TableName() string {
	return "departamentos"
}

func (Provincia) TableName() string {
	return "provincias"
}

func (Dia) TableName() string {
	return "dias"
}

func (Mes) TableName() string {
	return "meses"
}

func (CategoriaAtraccion) TableName() string {
	return "categorias_atracciones"
}

func (SubcategoriaAtraccion) TableName() string {
	return "subcategorias_atracciones"
}

func (AtraccionTuristica) TableName() string {
	return "atracciones_turisticas"
}

func (AtraccionSubcategoria) TableName() string {
	return "atraccion_subcategorias"
}

func (AtraccionFoto) TableName() string {
	return "atraccion_fotos"
}
