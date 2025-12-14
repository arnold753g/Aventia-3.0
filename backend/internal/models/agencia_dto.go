package models

import "time"

// CreateAgenciaRapidaRequest DTO para creación rápida de agencia
// Solo campos esenciales: nombre, departamento, teléfono, encargado
type CreateAgenciaRapidaRequest struct {
	NombreComercial      string `json:"nombre_comercial" validate:"required,min=2,max=255"`
	DepartamentoID       uint   `json:"departamento_id" validate:"required"`
	Telefono             string `json:"telefono" validate:"required"`
	EncargadoPrincipalID *uint  `json:"encargado_principal_id" validate:"required"`
}

// CreateAgenciaCompletaRequest DTO para creación completa
type CreateAgenciaCompletaRequest struct {
	NombreComercial      string   `json:"nombre_comercial" validate:"required,min=2,max=255"`
	Descripcion          string   `json:"descripcion"`
	Direccion            string   `json:"direccion" validate:"required"`
	DepartamentoID       uint     `json:"departamento_id" validate:"required"`
	Latitud              *float64 `json:"latitud"`
	Longitud             *float64 `json:"longitud"`
	Telefono             string   `json:"telefono" validate:"required"`
	Email                string   `json:"email" validate:"required,email"`
	SitioWeb             string   `json:"sitio_web"`
	Facebook             string   `json:"facebook"`
	Instagram            string   `json:"instagram"`
	LicenciaTuristica    bool     `json:"licencia_turistica"`
	HorarioApertura      string   `json:"horario_apertura"`
	HorarioCierre        string   `json:"horario_cierre"`
	AceptaQR             bool     `json:"acepta_qr"`
	AceptaTransferencia  bool     `json:"acepta_transferencia"`
	AceptaEfectivo       bool     `json:"acepta_efectivo"`
	EncargadoPrincipalID *uint    `json:"encargado_principal_id" validate:"required"`
	Status               string   `json:"status"`
	VisiblePublico       bool     `json:"visible_publico"`
	DiasIDs              []uint   `json:"dias_ids"`
}

// UpdateAgenciaRequest DTO para actualización
type UpdateAgenciaRequest struct {
	NombreComercial      *string  `json:"nombre_comercial"`
	Descripcion          *string  `json:"descripcion"`
	Direccion            *string  `json:"direccion"`
	DepartamentoID       *uint    `json:"departamento_id"`
	Latitud              *float64 `json:"latitud"`
	Longitud             *float64 `json:"longitud"`
	Telefono             *string  `json:"telefono"`
	Email                *string  `json:"email"`
	SitioWeb             *string  `json:"sitio_web"`
	Facebook             *string  `json:"facebook"`
	Instagram            *string  `json:"instagram"`
	LicenciaTuristica    *bool    `json:"licencia_turistica"`
	HorarioApertura      *string  `json:"horario_apertura"`
	HorarioCierre        *string  `json:"horario_cierre"`
	AceptaQR             *bool    `json:"acepta_qr"`
	AceptaTransferencia  *bool    `json:"acepta_transferencia"`
	AceptaEfectivo       *bool    `json:"acepta_efectivo"`
	EncargadoPrincipalID *uint    `json:"encargado_principal_id"`
	Status               *string  `json:"status"`
	VisiblePublico       *bool    `json:"visible_publico"`
	DiasIDs              []uint   `json:"dias_ids"`
}

// UpdateAgenciaStatusRequest DTO para cambio de status
type UpdateAgenciaStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=activa inactiva suspendida en_revision"`
}

// AgenciaPublic DTO para vista pública
type AgenciaPublic struct {
	ID                uint                  `json:"id"`
	NombreComercial   string                `json:"nombre_comercial"`
	Descripcion       string                `json:"descripcion"`
	Direccion         string                `json:"direccion"`
	Departamento      *Departamento         `json:"departamento"`
	Latitud           *float64              `json:"latitud"`
	Longitud          *float64              `json:"longitud"`
	Telefono          string                `json:"telefono"`
	Email             string                `json:"email"`
	SitioWeb          string                `json:"sitio_web"`
	LicenciaTuristica bool                  `json:"licencia_turistica"`
	HorarioApertura   string                `json:"horario_apertura"`
	HorarioCierre     string                `json:"horario_cierre"`
	Fotos             []AgenciaFoto         `json:"fotos"`
	Especialidades    []AgenciaEspecialidad `json:"especialidades"`
	Dias              []Dia                 `json:"dias"`
}

// AgenciaDetalle DTO para vista detallada
type AgenciaDetalle struct {
	AgenciaPublic
	AceptaQR            bool      `json:"acepta_qr"`
	AceptaTransferencia bool      `json:"acepta_transferencia"`
	AceptaEfectivo      bool      `json:"acepta_efectivo"`
	EncargadoPrincipal  *Usuario  `json:"encargado_principal"`
	Status              string    `json:"status"`
	VisiblePublico      bool      `json:"visible_publico"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
