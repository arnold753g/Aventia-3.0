import { useAuthStore } from '~/stores/auth'
import type { CrearCompraRequest } from '~/types/compra'

export const useCompra = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  const crearCompra = async (data: CrearCompraRequest) => {
    return $fetch(`${apiBase}/compras`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  const obtenerDetalleCompra = async (compraId: number) => {
    return $fetch(`${apiBase}/compras/${compraId}`, {
      headers: authHeader()
    })
  }

  const listarMisCompras = async (params: { page?: number; page_size?: number } = {}) => {
    const query = new URLSearchParams()
    if (params.page) query.set('page', String(params.page))
    if (params.page_size) query.set('page_size', String(params.page_size))

    const qs = query.toString()
    const url = qs ? `${apiBase}/mis-compras?${qs}` : `${apiBase}/mis-compras`

    return $fetch(url, { headers: authHeader() })
  }

  const cancelarCompra = async (compraId: number, razon?: string) => {
    return $fetch(`${apiBase}/compras/${compraId}/cancelar`, {
      method: 'POST',
      headers: authHeader(),
      body: razon ? { razon } : undefined
    })
  }

  return {
    crearCompra,
    obtenerDetalleCompra,
    listarMisCompras,
    cancelarCompra
  }
}
