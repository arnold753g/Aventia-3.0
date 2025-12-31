export interface Notificacion {
  id: number
  tipo: 'nuevo_pago_pendiente' | 'pago_confirmado' | 'pago_rechazado' | 'compra_expirada'
  titulo: string
  mensaje: string
  datos_json: NotificacionDatos
  leida: boolean
  fecha_leida: string | null
  created_at: string
}

export interface NotificacionDatos {
  pago_id?: number
  compra_id?: number
  paquete_id?: number
  paquete_nombre?: string
  turista_nombre?: string
  confirmado_por?: string
  fecha_salida?: string
  monto?: number
  metodo_pago?: string
  comprobante_foto?: string
  razon_rechazo?: string
  puede_reintentar?: boolean
}

export interface Pagination {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface NotificacionesResponse {
  notificaciones: Notificacion[]
  no_leidas: number
  pagination?: Pagination
}
