export interface CrearCompraRequest {
  paquete_id: number
  fecha_seleccionada: string // YYYY-MM-DD
  tipo_compra: 'compartido' | 'privado'
  extranjero: boolean
  cantidad_adultos: number
  cantidad_ninos_pagan: number
  cantidad_ninos_gratis: number
  tiene_discapacidad: boolean
  descripcion_discapacidad?: string | null
  notas_turista?: string | null
}

export interface CrearCompraData {
  compra_id: number
  salida_id: number
  precio_total: number
}

export interface ParticipanteDetalle {
  nombre?: string
  edad?: number
  tipo: 'adulto' | 'nino_paga' | 'nino_gratis'
  notas?: string
}

export interface ParticipantesDetalle {
  adultos?: ParticipanteDetalle[]
  ninos_pagan?: ParticipanteDetalle[]
  ninos_gratis?: ParticipanteDetalle[]
}

export interface PaqueteFoto {
  id: number
  foto: string
  es_principal: boolean
  orden: number
}

export interface PaqueteItinerarioDetalle {
  id: number
  dia_numero: number
  titulo: string
  descripcion?: string | null
  actividades?: string[]
  hospedaje_info?: string | null
}

export interface PaquetePoliticaDetalle {
  edad_minima_pago: number
  recargo_privado_porcentaje: number
  politica_cancelacion?: string | null
}

export interface AgenciaDatosPagoDetalle {
  nombre_banco?: string | null
  numero_cuenta?: string | null
  nombre_titular?: string | null
  qr_pago_foto?: string | null
}

export interface AgenciaCompleta {
  id: number
  nombre_comercial: string
  direccion: string
  departamento?: { nombre: string }
  telefono?: string
  email?: string
  whatsapp?: string
}

export interface PaqueteCompleto {
  id: number
  nombre: string
  descripcion?: string | null
  frecuencia?: string
  duracion_dias?: number | null
  horario?: string | null
  hora_salida?: string | null
  nivel_dificultad?: string | null
  permite_privado?: boolean
  incluye?: string[]
  no_incluye?: string[]
  que_llevar?: string[]
  requisitos?: string[]
  itinerario?: PaqueteItinerarioDetalle[]
  fotos?: PaqueteFoto[]
  agencia?: AgenciaCompleta
  politicas?: PaquetePoliticaDetalle
  agencia_datos_pago?: AgenciaDatosPagoDetalle
  politica_cancelacion?: string | null
  contacto_emergencia?: string | null
  recomendaciones_climaticas?: string | null
}

export interface PagoDetalle {
  id: number
  metodo_pago: string
  monto: number
  estado: 'pendiente' | 'confirmado' | 'rechazado'
  comprobante_foto?: string | null
  fecha_confirmacion: string | null
  fecha_registro?: string
}

export interface CompraDetalle {
  id: number
  codigo_confirmacion?: string | null
  fecha_compra: string
  fecha_seleccionada: string
  fecha_confirmacion?: string | null
  tipo_compra: 'compartido' | 'privado'
  extranjero?: boolean
  cantidad_adultos?: number
  cantidad_ninos_pagan?: number
  cantidad_ninos_gratis?: number
  participantes?: ParticipantesDetalle
  total_participantes: number
  precio_total: number
  status: 'pendiente_confirmacion' | 'confirmada' | 'rechazada' | 'cancelada' | 'completada'
  tiene_discapacidad?: boolean
  descripcion_discapacidad?: string | null
  notas_turista?: string | null
  paquete: PaqueteCompleto
  salida: {
    id: number
    fecha_salida: string
    tipo_salida: string
    estado: string
    cupos_disponibles: number
    punto_encuentro?: string | null
    hora_encuentro?: string | null
    instrucciones_turistas?: string | null
    hora_salida?: string | null
  } | null
  ultimo_pago: PagoDetalle | null
}

export interface ComprasListData {
  compras: CompraDetalle[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}
