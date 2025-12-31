package models

import "time"

// AgenciaTurismo representa una agencia de turismo
type AgenciaTurismo struct {
	ID                   uint                  `gorm:"primaryKey" json:"id"`
	NombreComercial      string                `gorm:"size:255;not null;index" json:"nombre_comercial"`
	Slug                 string                `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Descripcion          string                `gorm:"type:text" json:"descripcion"`
	Direccion            string                `gorm:"size:500;not null" json:"direccion"`
	DepartamentoID       uint                  `gorm:"not null;index" json:"departamento_id"`
	Departamento         *Departamento         `gorm:"foreignKey:DepartamentoID" json:"departamento,omitempty"`
	Latitud              *float64              `gorm:"type:decimal(10,8)" json:"latitud"`
	Longitud             *float64              `gorm:"type:decimal(11,8)" json:"longitud"`
	Telefono             string                `gorm:"size:20;not null" json:"telefono"`
	Email                string                `gorm:"size:255;not null" json:"email"`
	SitioWeb             string                `gorm:"size:255" json:"sitio_web"`
	Facebook             string                `gorm:"size:255" json:"facebook"`
	Instagram            string                `gorm:"size:255" json:"instagram"`
	LicenciaTuristica    bool                  `gorm:"default:false" json:"licencia_turistica"`
	HorarioApertura      *string               `gorm:"type:time" json:"horario_apertura"`
	HorarioCierre        *string               `gorm:"type:time" json:"horario_cierre"`
	AceptaQR             bool                  `gorm:"default:true" json:"acepta_qr"`
	AceptaTransferencia  bool                  `gorm:"default:true" json:"acepta_transferencia"`
	AceptaEfectivo       bool                  `gorm:"default:true" json:"acepta_efectivo"`
	EncargadoPrincipalID *uint                 `gorm:"index" json:"encargado_principal_id"`
	EncargadoPrincipal   *Usuario              `gorm:"foreignKey:EncargadoPrincipalID" json:"encargado_principal,omitempty"`
	Status               string                `gorm:"size:20;default:'activa';index" json:"status"`
	VisiblePublico       bool                  `gorm:"default:true;index" json:"visible_publico"`
	Fotos                []AgenciaFoto         `gorm:"foreignKey:AgenciaID" json:"fotos,omitempty"`
	Especialidades       []AgenciaEspecialidad `gorm:"foreignKey:AgenciaID" json:"especialidades,omitempty"`
	// Many-to-many con tabla agencia_dias (columnas: agencia_id, dia_id)
	Dias      []Dia     `gorm:"many2many:agencia_dias;foreignKey:ID;joinForeignKey:AgenciaID;References:ID;joinReferences:DiaID" json:"dias,omitempty"`
	Politicas *PaquetePolitica `gorm:"-" json:"politicas,omitempty"`
	CreatedBy uint      `gorm:"not null" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AgenciaFoto representa una foto de una agencia
type AgenciaFoto struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AgenciaID   uint      `gorm:"not null;index" json:"agencia_id"`
	FotoURL     string    `gorm:"size:500;not null" json:"foto_url"`
	Titulo      string    `gorm:"size:255" json:"titulo"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	EsPrincipal bool      `gorm:"default:false" json:"es_principal"`
	Orden       int       `gorm:"default:0" json:"orden"`
	CreatedAt   time.Time `json:"created_at"`
}

// AgenciaEspecialidad representa la relación entre agencia y categoría
type AgenciaEspecialidad struct {
	ID          uint                `gorm:"primaryKey" json:"id"`
	AgenciaID   uint                `gorm:"not null;index" json:"agencia_id"`
	CategoriaID uint                `gorm:"not null;index" json:"categoria_id"`
	Categoria   *CategoriaAtraccion `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
	EsPrincipal bool                `gorm:"default:false" json:"es_principal"`
	CreatedAt   time.Time           `json:"created_at"`
}

func (AgenciaTurismo) TableName() string {
	return "agencias_turismo"
}

func (AgenciaFoto) TableName() string {
	return "agencia_fotos"
}

func (AgenciaEspecialidad) TableName() string {
	return "agencia_especialidades"
}
