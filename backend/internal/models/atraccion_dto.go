package models

import "time"

// CreateAtraccionRequest para crear atraccion
type CreateAtraccionRequest struct {
	Nombre           string   `json:"nombre" validate:"required,min=3,max=255"`
	Descripcion      string   `json:"descripcion"`
	ProvinciaID      uint     `json:"provincia_id" validate:"required"`
	Direccion        string   `json:"direccion"`
	Latitud          *float64 `json:"latitud" validate:"omitempty,min=-22.90,max=-9.67"`
	Longitud         *float64 `json:"longitud" validate:"omitempty,min=-69.65,max=-57.45"`
	HorarioApertura  string   `json:"horario_apertura" validate:"omitempty"`
	HorarioCierre    string   `json:"horario_cierre" validate:"omitempty"`
	PrecioEntrada    float64  `json:"precio_entrada" validate:"omitempty,min=0"`
	NivelDificultad  string   `json:"nivel_dificultad" validate:"omitempty,oneof=facil medio dificil extremo"`
	RequiereAgencia  bool     `json:"requiere_agencia"`
	AccesoParticular bool     `json:"acceso_particular"`
	MesInicioID      *uint    `json:"mes_inicio_id" validate:"omitempty,min=1,max=12"`
	MesFinID         *uint    `json:"mes_fin_id" validate:"omitempty,min=1,max=12"`
	Status           string   `json:"status" validate:"omitempty,oneof=activa inactiva mantenimiento fuera_temporada"`
	VisiblePublico   bool     `json:"visible_publico"`
	Telefono         string   `json:"telefono" validate:"omitempty"`
	Email            string   `json:"email" validate:"omitempty,email"`
	SitioWeb         string   `json:"sitio_web" validate:"omitempty,url"`
	Facebook         string   `json:"facebook" validate:"omitempty"`
	Instagram        string   `json:"instagram" validate:"omitempty"`
	SubcategoriasIDs []uint   `json:"subcategorias_ids" validate:"omitempty,min=1,max=4,dive,min=1"`
	DiasIDs          []uint   `json:"dias_ids" validate:"omitempty,min=1,max=7,dive,min=1,max=7"`
	Fotos            []string `json:"fotos" validate:"omitempty,max=10,dive,url"`
}

// UpdateAtraccionRequest para actualizar atraccion
type UpdateAtraccionRequest struct {
	Nombre           string   `json:"nombre" validate:"omitempty,min=3,max=255"`
	Descripcion      string   `json:"descripcion"`
	ProvinciaID      uint     `json:"provincia_id" validate:"omitempty"`
	Direccion        string   `json:"direccion"`
	Latitud          *float64 `json:"latitud" validate:"omitempty,min=-22.90,max=-9.67"`
	Longitud         *float64 `json:"longitud" validate:"omitempty,min=-69.65,max=-57.45"`
	HorarioApertura  string   `json:"horario_apertura"`
	HorarioCierre    string   `json:"horario_cierre"`
	PrecioEntrada    *float64 `json:"precio_entrada" validate:"omitempty,min=0"`
	NivelDificultad  string   `json:"nivel_dificultad" validate:"omitempty,oneof=facil medio dificil extremo"`
	RequiereAgencia  *bool    `json:"requiere_agencia"`
	AccesoParticular *bool    `json:"acceso_particular"`
	MesInicioID      *uint    `json:"mes_inicio_id" validate:"omitempty,min=1,max=12"`
	MesFinID         *uint    `json:"mes_fin_id" validate:"omitempty,min=1,max=12"`
	Status           string   `json:"status" validate:"omitempty,oneof=activa inactiva mantenimiento fuera_temporada"`
	VisiblePublico   *bool    `json:"visible_publico"`
	Telefono         string   `json:"telefono"`
	Email            string   `json:"email" validate:"omitempty,email"`
	SitioWeb         string   `json:"sitio_web" validate:"omitempty,url"`
	Facebook         string   `json:"facebook"`
	Instagram        string   `json:"instagram"`
}

// AtraccionPublic representa los datos publicos
type AtraccionPublic struct {
	ID               uint                    `json:"id"`
	Nombre           string                  `json:"nombre"`
	Descripcion      string                  `json:"descripcion"`
	ProvinciaID      uint                    `json:"provincia_id"`
	Provincia        Provincia               `json:"provincia"`
	Direccion        string                  `json:"direccion"`
	Latitud          *float64                `json:"latitud"`
	Longitud         *float64                `json:"longitud"`
	HorarioApertura  *string                 `json:"horario_apertura"`
	HorarioCierre    *string                 `json:"horario_cierre"`
	PrecioEntrada    float64                 `json:"precio_entrada"`
	NivelDificultad  string                  `json:"nivel_dificultad"`
	RequiereAgencia  bool                    `json:"requiere_agencia"`
	AccesoParticular bool                    `json:"acceso_particular"`
	MesInicio        *Mes                    `json:"mes_inicio"`
	MesFin           *Mes                    `json:"mes_fin"`
	Status           string                  `json:"status"`
	VisiblePublico   bool                    `json:"visible_publico"`
	Telefono         string                  `json:"telefono"`
	Email            string                  `json:"email"`
	SitioWeb         string                  `json:"sitio_web"`
	Facebook         string                  `json:"facebook"`
	Instagram        string                  `json:"instagram"`
	Subcategorias    []SubcategoriaAtraccion `json:"subcategorias"`
	Fotos            []AtraccionFoto         `json:"fotos"`
	Dias             []Dia                   `json:"dias"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// AddSubcategoriaRequest para agregar subcategoria
type AddSubcategoriaRequest struct {
	SubcategoriaID uint `json:"subcategoria_id" validate:"required"`
	EsPrincipal    bool `json:"es_principal"`
}

// AddFotoRequest para agregar foto
type AddFotoRequest struct {
	Foto        string `json:"foto" validate:"required,url"`
	EsPrincipal bool   `json:"es_principal"`
	Orden       int    `json:"orden" validate:"min=0"`
}
