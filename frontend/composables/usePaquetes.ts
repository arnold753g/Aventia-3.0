import { useAuthStore } from '~/stores/auth'

export const usePaquetes = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  const getPaquetes = async (
    agenciaId: number,
    params: {
      page?: number
      limit?: number
      search?: string
      status?: string
      frecuencia?: string
      visible_publico?: string
      include_eliminado?: string
      sort_by?: string
      sort_order?: string
    } = {}
  ) => {
    const query = new URLSearchParams()
    Object.entries(params || {}).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') query.append(key, String(value))
    })
    const qs = query.toString()
    const url = qs ? `${apiBase}/agencias/${agenciaId}/paquetes?${qs}` : `${apiBase}/agencias/${agenciaId}/paquetes`

    return $fetch(url, { headers: authHeader() })
  }

  const createPaquete = async (agenciaId: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  const getPaquete = async (agenciaId: number, paqueteId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}`, {
      headers: authHeader()
    })
  }

  const updatePaquete = async (agenciaId: number, paqueteId: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}`, {
      method: 'PUT',
      headers: authHeader(),
      body: data
    })
  }

  const deletePaquete = async (agenciaId: number, paqueteId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Fotos
  const uploadPaqueteFoto = async (agenciaId: number, paqueteId: number, formData: FormData) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/fotos/upload`, {
      method: 'POST',
      headers: authHeader(),
      body: formData
    })
  }

  const removePaqueteFoto = async (agenciaId: number, paqueteId: number, fotoId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/fotos/${fotoId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Itinerario
  const getItinerario = async (agenciaId: number, paqueteId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/itinerario`, {
      headers: authHeader()
    })
  }

  const createItinerarioItem = async (agenciaId: number, paqueteId: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/itinerario`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  const updateItinerarioItem = async (agenciaId: number, paqueteId: number, itemId: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/itinerario/${itemId}`, {
      method: 'PUT',
      headers: authHeader(),
      body: data
    })
  }

  const deleteItinerarioItem = async (agenciaId: number, paqueteId: number, itemId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/itinerario/${itemId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Atracciones
  const getPaqueteAtracciones = async (agenciaId: number, paqueteId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/atracciones`, {
      headers: authHeader()
    })
  }

  const addPaqueteAtraccion = async (agenciaId: number, paqueteId: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/atracciones`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  const updatePaqueteAtraccion = async (agenciaId: number, paqueteId: number, itemId: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/atracciones/${itemId}`, {
      method: 'PUT',
      headers: authHeader(),
      body: data
    })
  }

  const removePaqueteAtraccion = async (agenciaId: number, paqueteId: number, itemId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/atracciones/${itemId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Salidas
  const getPaqueteSalidas = async (agenciaId: number, paqueteId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/salidas`, {
      headers: authHeader()
    })
  }

  const createPaqueteSalida = async (
    agenciaId: number,
    paqueteId: number,
    data: {
      fecha_salida: string
      tipo_salida?: string
      estado?: string
    }
  ) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/salidas`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  const updatePaqueteSalida = async (agenciaId: number, paqueteId: number, salidaId: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/paquetes/${paqueteId}/salidas/${salidaId}`, {
      method: 'PUT',
      headers: authHeader(),
      body: data
    })
  }

  return {
    getPaquetes,
    createPaquete,
    getPaquete,
    updatePaquete,
    deletePaquete,
    uploadPaqueteFoto,
    removePaqueteFoto,
    getItinerario,
    createItinerarioItem,
    updateItinerarioItem,
    deleteItinerarioItem,
    getPaqueteAtracciones,
    addPaqueteAtraccion,
    updatePaqueteAtraccion,
    removePaqueteAtraccion,
    getPaqueteSalidas,
    createPaqueteSalida,
    updatePaqueteSalida
  }
}
