export const useAgenciaVisitas = () => {
  const config = useRuntimeConfig()

  /**
   * Registra una visita a la página pública de una agencia
   * No lanza errores para no interrumpir la experiencia del usuario
   */
  const registrarVisita = async (agenciaIdOrSlug: number | string) => {
    // Prevenir visitas duplicadas en la misma sesión
    const visitKey = `visited_agencia_${agenciaIdOrSlug}`

    if (process.client && sessionStorage.getItem(visitKey)) {
      return // Ya visitó en esta sesión
    }

    try {
      const response = await $fetch(`${config.public.apiBase}/public/agencias/${agenciaIdOrSlug}/visitas`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      })

      // Marcar como visitado en esta sesión
      if (process.client) {
        sessionStorage.setItem(visitKey, 'true')
      }

      return response
    } catch (error) {
      // Silenciosamente fallar - no queremos interrumpir la experiencia del usuario
      console.debug('No se pudo registrar la visita:', error)
      return null
    }
  }

  /**
   * Obtiene las estadísticas de visitas de una agencia (requiere autenticación)
   */
  const getEstadisticas = async (agenciaId: number) => {
    try {
      const { getAuthHeaders } = useAuth()
      const headers = getAuthHeaders()

      const response: any = await $fetch(`${config.public.apiBase}/agencias/${agenciaId}/estadisticas-visitas`, {
        method: 'GET',
        headers
      })

      if (response.success) {
        return response.data
      }

      throw new Error(response.error?.message || 'Error al obtener estadísticas')
    } catch (error: any) {
      throw error
    }
  }

  /**
   * Obtiene el detalle de visitas con paginación (requiere autenticación)
   */
  const getVisitasDetalle = async (agenciaId: number, params?: { page?: number; limit?: number }) => {
    try {
      const { getAuthHeaders } = useAuth()
      const headers = getAuthHeaders()

      const queryParams = new URLSearchParams()
      if (params?.page) queryParams.append('page', params.page.toString())
      if (params?.limit) queryParams.append('limit', params.limit.toString())

      const url = `${config.public.apiBase}/agencias/${agenciaId}/visitas-detalle?${queryParams.toString()}`

      const response: any = await $fetch(url, {
        method: 'GET',
        headers
      })

      if (response.success) {
        return response.data
      }

      throw new Error(response.error?.message || 'Error al obtener visitas')
    } catch (error: any) {
      throw error
    }
  }

  return {
    registrarVisita,
    getEstadisticas,
    getVisitasDetalle
  }
}
