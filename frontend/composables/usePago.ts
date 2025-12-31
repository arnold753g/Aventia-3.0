import { useAuthStore } from '~/stores/auth'
import type { CrearPagoRequest } from '~/types/pago'

export const usePago = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  const crearPago = async (data: CrearPagoRequest) => {
    const formData = new FormData()
    formData.append('compra_id', String(data.compra_id))
    formData.append('metodo_pago', data.metodo_pago)
    formData.append('monto', String(data.monto))
    if (data.comprobante) formData.append('comprobante', data.comprobante)

    return $fetch(`${apiBase}/pagos`, {
      method: 'POST',
      headers: authHeader(),
      body: formData
    })
  }

  const confirmarPago = async (pagoId: number, notas_encargado?: string) => {
    return $fetch(`${apiBase}/pagos/${pagoId}/confirmar`, {
      method: 'PUT',
      headers: authHeader(),
      body: { notas_encargado }
    })
  }

  const rechazarPago = async (pagoId: number, razon_rechazo: string) => {
    return $fetch(`${apiBase}/pagos/${pagoId}/rechazar`, {
      method: 'PUT',
      headers: authHeader(),
      body: { razon_rechazo }
    })
  }

  return {
    crearPago,
    confirmarPago,
    rechazarPago
  }
}

