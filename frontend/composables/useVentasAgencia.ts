import { useAuthStore } from '~/stores/auth'

export const useVentasAgencia = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  const getVentasPagos = async (
    agenciaId: number,
    params: { estado?: string; page?: number; page_size?: number } = {}
  ) => {
    const query = new URLSearchParams()
    if (params.estado) query.set('estado', String(params.estado))
    if (params.page) query.set('page', String(params.page))
    if (params.page_size) query.set('page_size', String(params.page_size))

    const qs = query.toString()
    const url = qs ? `${apiBase}/agencias/${agenciaId}/ventas/pagos?${qs}` : `${apiBase}/agencias/${agenciaId}/ventas/pagos`

    return $fetch(url, { headers: authHeader() })
  }

  const getVentasSalidas = async (
    agenciaId: number,
    params: { desde?: string; hasta?: string; estado?: string } = {}
  ) => {
    const query = new URLSearchParams()
    if (params.desde) query.set('desde', String(params.desde))
    if (params.hasta) query.set('hasta', String(params.hasta))
    if (params.estado) query.set('estado', String(params.estado))

    const qs = query.toString()
    const url = qs ? `${apiBase}/agencias/${agenciaId}/ventas/salidas?${qs}` : `${apiBase}/agencias/${agenciaId}/ventas/salidas`

    return $fetch(url, { headers: authHeader() })
  }

  const getVentasSalidaCompras = async (agenciaId: number, salidaId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/ventas/salidas/${salidaId}/compras`, { headers: authHeader() })
  }

  return {
    getVentasPagos,
    getVentasSalidas,
    getVentasSalidaCompras
  }
}
