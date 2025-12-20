export type EstadoPagoVenta = 'pendiente' | 'confirmado' | 'rechazado'

export interface AgenciaVentaPago {
  pago_id: number
  compra_id: number
  metodo_pago: 'efectivo' | 'qr' | 'transferencia'
  monto: number
  estado: EstadoPagoVenta
  comprobante_foto: string | null

  confirmado_por: number | null
  fecha_confirmacion: string | null
  razon_rechazo: string | null
  notas_encargado: string | null
  fecha_pago: string

  compra_status: string
  fecha_seleccionada: string
  horario_seleccionado: string | null
  tipo_compra: string
  total_participantes: number
  precio_total: number

  paquete_id: number
  paquete_nombre: string
  paquete_frecuencia: string
  paquete_duracion_dias: number | null
  paquete_horario: string | null

  turista_id: number
  turista_nombre: string
  turista_apellido_paterno: string
  turista_apellido_materno: string
  turista_phone: string
  turista_email: string
}

export interface VentasPagosData {
  pagos: AgenciaVentaPago[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface AgenciaVentaSalida {
  salida_id: number

  paquete_id: number
  paquete_nombre: string
  paquete_frecuencia: string
  paquete_duracion_dias: number | null
  paquete_horario: string | null

  fecha_salida: string
  tipo_salida: string
  cupo_minimo: number
  cupo_maximo: number
  cupos_reservados: number
  cupos_confirmados: number
  estado: string
  updated_at: string
}

export interface VentasSalidasData {
  salidas: AgenciaVentaSalida[]
}

export interface VentaSalidaDetalle {
  id: number
  paquete_id: number
  fecha_salida: string
  tipo_salida: string
  cupo_minimo: number
  cupo_maximo: number
  cupos_reservados: number
  cupos_confirmados: number
  estado: string
  created_at: string
  updated_at: string
}

export interface AgenciaVentaSalidaCompra {
  compra_id: number
  compra_status: string
  fecha_compra: string
  fecha_seleccionada: string
  horario_seleccionado: string | null
  tipo_compra: string
  total_participantes: number
  precio_total: number

  turista_id: number
  turista_nombre: string
  turista_apellido_paterno: string
  turista_apellido_materno: string
  turista_phone: string
  turista_email: string

  pago_id: number | null
  metodo_pago: 'efectivo' | 'qr' | 'transferencia' | string | null
  monto: number | null
  estado: EstadoPagoVenta | null
  comprobante_foto: string | null
  fecha_pago: string | null

  confirmado_por: number | null
  fecha_confirmacion: string | null
  razon_rechazo: string | null
  notas_encargado: string | null
}

export interface VentasSalidaComprasData {
  salida: VentaSalidaDetalle
  compras: AgenciaVentaSalidaCompra[]
}
