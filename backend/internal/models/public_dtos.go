package models

// DTOs para endpoints públicos - sin información sensible

// AgenciaPublicaDTO representa una agencia para visualización pública
type AgenciaPublicaDTO struct {
	ID              uint     `json:"id"`
	NombreComercial string   `json:"nombre_comercial"`
	Descripcion     *string  `json:"descripcion"`
	Direccion       string   `json:"direccion"`
	DepartamentoID  uint     `json:"departamento_id"`
	Latitud         *float64 `json:"latitud"`
	Longitud        *float64 `json:"longitud"`
	// Contacto público (no personal)
	Telefono  *string `json:"telefono"`
	Email     *string `json:"email"`
	SitioWeb  *string `json:"sitio_web"`
	Facebook  *string `json:"facebook"`
	Instagram *string `json:"instagram"`
	// Horarios
	HorarioApertura *string `json:"horario_apertura"`
	HorarioCierre   *string `json:"horario_cierre"`
	// Fotos y especialidades (relaciones)
	Fotos           []AgenciaFoto         `json:"fotos,omitempty"`
	Especialidades  []AgenciaEspecialidad `json:"especialidades,omitempty"`
	DiasOperacion   []Dia                 `json:"dias_operacion,omitempty"`
}

// PaquetePublicoDTO representa un paquete para visualización pública
type PaquetePublicoDTO struct {
	ID                         uint                   `json:"id"`
	AgenciaID                  uint                   `json:"agencia_id"`
	Nombre                     string                 `json:"nombre"`
	Descripcion                *string                `json:"descripcion"`
	Frecuencia                 string                 `json:"frecuencia"`
	DuracionDias               *int                   `json:"duracion_dias"`
	DuracionNoches             *int                   `json:"duracion_noches"`
	FechaSalidaFija            *string                `json:"fecha_salida_fija"`
	Horario                    *string                `json:"horario"`
	HoraSalida                 *string                `json:"hora_salida"`
	NivelDificultad            *string                `json:"nivel_dificultad"`
	CupoMinimo                 int                    `json:"cupo_minimo"`
	CupoMaximo                 int                    `json:"cupo_maximo"`
	PermitePrivado             bool                   `json:"permite_privado"`
	PrecioBaseNacionales       float64                `json:"precio_base_nacionales"`
	PrecioAdicionalExtranjeros float64                `json:"precio_adicional_extranjeros"`
	Incluye                    StringArray            `json:"incluye"`
	NoIncluye                  StringArray            `json:"no_incluye"`
	QueLlevar                  StringArray            `json:"que_llevar"`
	// Relaciones públicas
	Fotos                      []PaqueteFoto          `json:"fotos,omitempty"`
	Itinerario                 []PaqueteItinerario    `json:"itinerario,omitempty"`
	Atracciones                []PaqueteAtraccion     `json:"atracciones,omitempty"`
	Agencia                    *AgenciaPublicaDTO     `json:"agencia,omitempty"`
	Politicas                  *PaquetePoliticaDTO    `json:"politicas,omitempty"`
	AgenciaDatosPago           *AgenciaDatosPagoDTO   `json:"agencia_datos_pago,omitempty"`
}

// PaquetePoliticaDTO sin información sensible
type PaquetePoliticaDTO struct {
	EdadMinimaPago           int     `json:"edad_minima_pago"`
	RecargoPrivadoPorcentaje float64 `json:"recargo_privado_porcentaje"`
	PoliticaCancelacion      *string `json:"politica_cancelacion"`
}

// AgenciaDatosPagoDTO solo información pública de pago
type AgenciaDatosPagoDTO struct {
	QRFoto *string `json:"qr_foto,omitempty"`
	// NO incluir: numero_cuenta, nombre_titular, banco
}

// SalidaPublicaDTO representa una salida para visualización pública
type SalidaPublicaDTO struct {
	ID                    uint     `json:"id"`
	PaqueteID             uint     `json:"paquete_id"`
	FechaSalida           string   `json:"fecha_salida"`
	TipoSalida            string   `json:"tipo_salida"`
	CuposDisponibles      int      `json:"cupos_disponibles"`
	CupoMinimo            int      `json:"cupo_minimo"`
	CupoMaximo            int      `json:"cupo_maximo"`
	Estado                string   `json:"estado"`
	PuntoEncuentro        *string  `json:"punto_encuentro,omitempty"`
	HoraEncuentro         *string  `json:"hora_encuentro,omitempty"`
	InstruccionesTuristas *string  `json:"instrucciones_turistas,omitempty"`
	// NO incluir: guia_nombre, guia_telefono, notas_logistica
}

// AtraccionPublicaDTO representa una atracción para visualización pública
type AtraccionPublicaDTO struct {
	ID                uint                   `json:"id"`
	Nombre            string                 `json:"nombre"`
	Descripcion       *string                `json:"descripcion"`
	ProvinciaID       uint                   `json:"provincia_id"`
	Direccion         *string                `json:"direccion"`
	Latitud           *float64               `json:"latitud"`
	Longitud          *float64               `json:"longitud"`
	HorarioApertura   *string                `json:"horario_apertura"`
	HorarioCierre     *string                `json:"horario_cierre"`
	PrecioEntrada     *float64               `json:"precio_entrada"`
	NivelDificultad   *string                `json:"nivel_dificultad"`
	RequiereAgencia   bool                   `json:"requiere_agencia"`
	AccesoParticular  bool                   `json:"acceso_particular"`
	MesInicio         *int                   `json:"mes_inicio"`
	MesFin            *int                   `json:"mes_fin"`
	Fotos             []AtraccionFoto        `json:"fotos,omitempty"`
	Subcategorias     []AtraccionSubcategoria `json:"subcategorias,omitempty"`
}

// Métodos de conversión de modelos completos a DTOs públicos

func (a *AgenciaTurismo) ToPublicDTO() *AgenciaPublicaDTO {
	desc := &a.Descripcion
	if a.Descripcion == "" {
		desc = nil
	}
	tel := &a.Telefono
	if a.Telefono == "" {
		tel = nil
	}
	email := &a.Email
	if a.Email == "" {
		email = nil
	}
	web := &a.SitioWeb
	if a.SitioWeb == "" {
		web = nil
	}
	fb := &a.Facebook
	if a.Facebook == "" {
		fb = nil
	}
	ig := &a.Instagram
	if a.Instagram == "" {
		ig = nil
	}

	return &AgenciaPublicaDTO{
		ID:              a.ID,
		NombreComercial: a.NombreComercial,
		Descripcion:     desc,
		Direccion:       a.Direccion,
		DepartamentoID:  a.DepartamentoID,
		Latitud:         a.Latitud,
		Longitud:        a.Longitud,
		Telefono:        tel,
		Email:           email,
		SitioWeb:        web,
		Facebook:        fb,
		Instagram:       ig,
		HorarioApertura: a.HorarioApertura,
		HorarioCierre:   a.HorarioCierre,
		Fotos:           a.Fotos,
		Especialidades:  a.Especialidades,
		DiasOperacion:   a.Dias, // Corregido
	}
}

func (p *PaqueteTuristico) ToPublicDTO() *PaquetePublicoDTO {
	dto := &PaquetePublicoDTO{
		ID:                         p.ID,
		AgenciaID:                  p.AgenciaID,
		Nombre:                     p.Nombre,
		Descripcion:                p.Descripcion,
		Frecuencia:                 p.Frecuencia,
		DuracionDias:               p.DuracionDias,
		DuracionNoches:             p.DuracionNoches,
		FechaSalidaFija:            p.FechaSalidaFija,
		Horario:                    p.Horario,
		HoraSalida:                 p.HoraSalida,
		NivelDificultad:            p.NivelDificultad,
		CupoMinimo:                 p.CupoMinimo,
		CupoMaximo:                 p.CupoMaximo,
		PermitePrivado:             p.PermitePrivado,
		PrecioBaseNacionales:       p.PrecioBaseNacionales,
		PrecioAdicionalExtranjeros: p.PrecioAdicionalExtranjeros,
		Incluye:                    p.Incluye,
		NoIncluye:                  p.NoIncluye,
		QueLlevar:                  p.QueLlevar,
		Fotos:                      p.Fotos,
		Itinerario:                 p.Itinerario,
		Atracciones:                p.Atracciones,
	}

	if p.Agencia != nil {
		dto.Agencia = p.Agencia.ToPublicDTO()
	}

	if p.Politicas != nil {
		dto.Politicas = &PaquetePoliticaDTO{
			EdadMinimaPago:           p.Politicas.EdadMinimaPago,
			RecargoPrivadoPorcentaje: p.Politicas.RecargoPrivadoPorcentaje,
			PoliticaCancelacion:      p.Politicas.PoliticaCancelacion,
		}
	}

	if p.AgenciaDatosPago != nil {
		dto.AgenciaDatosPago = &AgenciaDatosPagoDTO{
			QRFoto: p.AgenciaDatosPago.QrPagoFoto,
		}
	}

	return dto
}

func (s *PaqueteSalidaHabilitada) ToPublicDTO() *SalidaPublicaDTO {
	cuposDisponibles := s.CupoMaximo - s.CuposReservados - s.CuposConfirmados
	if cuposDisponibles < 0 {
		cuposDisponibles = 0
	}

	return &SalidaPublicaDTO{
		ID:                    s.ID,
		PaqueteID:             s.PaqueteID,
		FechaSalida:           s.FechaSalida,
		TipoSalida:            s.TipoSalida,
		CuposDisponibles:      cuposDisponibles,
		CupoMinimo:            s.CupoMinimo,
		CupoMaximo:            s.CupoMaximo,
		Estado:                s.Estado,
		PuntoEncuentro:        s.PuntoEncuentro,
		HoraEncuentro:         s.HoraEncuentro,
		InstruccionesTuristas: s.InstruccionesTuristas,
	}
}

func (a *AtraccionTuristica) ToPublicDTO() *AtraccionPublicaDTO {
	desc := &a.Descripcion
	if a.Descripcion == "" {
		desc = nil
	}
	dir := &a.Direccion
	if a.Direccion == "" {
		dir = nil
	}
	precio := &a.PrecioEntrada
	dif := &a.NivelDificultad
	if a.NivelDificultad == "" {
		dif = nil
	}

	// Convertir *uint a *int para mes_inicio y mes_fin
	var mesInicio *int
	if a.MesInicioID != nil {
		temp := int(*a.MesInicioID)
		mesInicio = &temp
	}
	var mesFin *int
	if a.MesFinID != nil {
		temp := int(*a.MesFinID)
		mesFin = &temp
	}

	return &AtraccionPublicaDTO{
		ID:               a.ID,
		Nombre:           a.Nombre,
		Descripcion:      desc,
		ProvinciaID:      a.ProvinciaID,
		Direccion:        dir,
		Latitud:          a.Latitud,
		Longitud:         a.Longitud,
		HorarioApertura:  a.HorarioApertura,
		HorarioCierre:    a.HorarioCierre,
		PrecioEntrada:    precio,
		NivelDificultad:  dif,
		RequiereAgencia:  a.RequiereAgencia,
		AccesoParticular: a.AccesoParticular,
		MesInicio:        mesInicio,
		MesFin:           mesFin,
		Fotos:            a.Fotos,
		Subcategorias:    a.Subcategorias,
	}
}
