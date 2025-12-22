import { useAuthStore } from '~/stores/auth'

export const usePaquetesTuristicos = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  const getPaquetes = async (
    params: {
      page?: number
      limit?: number
      search?: string
      frecuencia?: string
      nivel_dificultad?: string
      tipo_duracion?: string
      precio_min?: string | number
      precio_max?: string | number
      sort_by?: string
      sort_order?: string
    } = {}
  ) => {
    const query = new URLSearchParams()
    Object.entries(params || {}).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') query.append(key, String(value))
    })
    const qs = query.toString()
    const url = qs ? `${apiBase}/paquetes?${qs}` : `${apiBase}/paquetes`
    return $fetch(url, { headers: authHeader() })
  }

  const getPaquete = async (id: number) => {
    return $fetch(`${apiBase}/paquetes/${id}`, { headers: authHeader() })
  }

  const getSalidas = async (
    paqueteId: number,
    params: {
      fecha?: string
      tipo?: 'compartido' | 'privado'
    } = {}
  ) => {
    const query = new URLSearchParams()
    if (params.fecha) query.append('fecha', params.fecha)
    if (params.tipo) query.append('tipo', params.tipo)
    const qs = query.toString()
    const url = qs ? `${apiBase}/paquetes/${paqueteId}/salidas?${qs}` : `${apiBase}/paquetes/${paqueteId}/salidas`
    return $fetch(url, { headers: authHeader() })
  }

  return {
    getPaquetes,
    getPaquete,
    getSalidas
  }
}
