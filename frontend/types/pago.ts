export interface CrearPagoRequest {
  compra_id: number
  metodo_pago: 'efectivo' | 'qr' | 'transferencia'
  monto: number
  comprobante?: File | null
}

export interface PagoData {
  id: number
  compra_id: number
  metodo_pago: string
  monto: number
  estado: string
  comprobante_foto: string | null
  mensaje: string
}

