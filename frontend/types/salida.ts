export interface SalidaConfirmada {
  salida_id: number

  paquete_id: number
  paquete_nombre: string
  paquete_frecuencia: string
  paquete_duracion_dias: number | null
  paquete_horario: string | null
  paquete_hora_salida: string | null
  paquete_foto: string | null

  agencia_id: number
  agencia_nombre: string

  fecha_salida: string // YYYY-MM-DD
  tipo_salida: 'compartido' | 'privado' | string
  estado: string

  cupo_minimo: number
  cupo_maximo: number
  cupos_reservados: number
  cupos_confirmados: number
  cupos_disponibles: number
}

export interface SalidasConfirmadasListData {
  salidas: SalidaConfirmada[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}
