import { useAuthStore } from '~/stores/auth'

export const usePaquetesTuristicos = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  // Obtener lista de paquetes (endpoint público)
  const getPaquetes = async (
    params: {
      page?: number
      limit?: number
      search?: string
      frecuencia?: string
      nivel_dificultad?: string
      tipo_duracion?: string
      agencia_id?: number
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
    const url = qs ? `${apiBase}/public/paquetes?${qs}` : `${apiBase}/public/paquetes`
    return $fetch(url)
  }

  // Obtener detalle de paquete (endpoint público)
  const getPaquete = async (id: number) => {
    return $fetch(`${apiBase}/public/paquetes/${id}`)
  }

  // Obtener salidas de un paquete (endpoint público)
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
    const url = qs ? `${apiBase}/public/paquetes/${paqueteId}/salidas?${qs}` : `${apiBase}/public/paquetes/${paqueteId}/salidas`
    return $fetch(url)
  }

  return {
    getPaquetes,
    getPaquete,
    getSalidas
  }
}
