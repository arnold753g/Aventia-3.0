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

export interface CompraDetalle {
  id: number
  fecha_compra: string
  fecha_seleccionada: string
  tipo_compra: 'compartido' | 'privado'
  total_participantes: number
  precio_total: number
  status: 'pendiente_confirmacion' | 'confirmada' | 'rechazada' | 'cancelada' | 'completada'
  paquete: {
    id: number
    nombre: string
    frecuencia: string
    duracion_dias: number | null
    horario: string | null
  }
  salida: {
    id: number
    fecha_salida: string
    tipo_salida: string
    estado: string
    cupos_disponibles: number
  } | null
  ultimo_pago: {
    id: number
    metodo_pago: string
    monto: number
    estado: string
    fecha_confirmacion: string | null
  } | null
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

