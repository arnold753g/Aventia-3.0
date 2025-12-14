import { useAuthStore } from '~/stores/auth'

export const useAtracciones = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  // Obtener lista de atracciones
  const getAtracciones = async (params: {
    page?: number
    limit?: number
    search?: string
    provincia_id?: string
    departamento_id?: string
    categoria_id?: string
    subcategoria_id?: string
    nivel_dificultad?: string
    status?: string
    requiere_agencia?: string
    visible_publico?: string
    sort_by?: string
    sort_order?: string
  }) => {
    const queryParams = new URLSearchParams()
    Object.entries(params || {}).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        queryParams.append(key, value.toString())
      }
    })

    return $fetch(`${apiBase}/atracciones?${queryParams}`, {
      headers: authHeader()
    })
  }

  // Obtener una atracción por ID
  const getAtraccion = async (id: number) => {
    return $fetch(`${apiBase}/atracciones/${id}`, {
      headers: authHeader()
    })
  }

  // Crear atracción
  const createAtraccion = async (data: any) => {
    return $fetch(`${apiBase}/atracciones`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  // Actualizar atracción
  const updateAtraccion = async (id: number, data: any) => {
    return $fetch(`${apiBase}/atracciones/${id}`, {
      method: 'PUT',
      headers: authHeader(),
      body: data
    })
  }

  // Desactivar atracción
  const deleteAtraccion = async (id: number) => {
    return $fetch(`${apiBase}/admin/atracciones/${id}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Agregar subcategoría
  const addSubcategoria = async (atraccionId: number, subcategoriaId: number, esPrincipal: boolean = false) => {
    return $fetch(`${apiBase}/atracciones/${atraccionId}/subcategorias`, {
      method: 'POST',
      headers: authHeader(),
      body: {
        subcategoria_id: subcategoriaId,
        es_principal: esPrincipal
      }
    })
  }

  // Eliminar subcategoría
  const removeSubcategoria = async (atraccionId: number, subcategoriaId: number) => {
    return $fetch(`${apiBase}/atracciones/${atraccionId}/subcategorias/${subcategoriaId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Agregar foto
  const addFoto = async (atraccionId: number, foto: string, esPrincipal: boolean = false, orden: number = 0) => {
    return $fetch(`${apiBase}/atracciones/${atraccionId}/fotos`, {
      method: 'POST',
      headers: authHeader(),
      body: {
        foto,
        es_principal: esPrincipal,
        orden
      }
    })
  }

  // Eliminar foto
  const removeFoto = async (atraccionId: number, fotoId: number) => {
    return $fetch(`${apiBase}/atracciones/${atraccionId}/fotos/${fotoId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Obtener estadísticas
  const getStats = async () => {
    return $fetch(`${apiBase}/admin/atracciones/stats`, {
      headers: authHeader()
    })
  }

  // Datos auxiliares
  const getCategorias = async () => {
    return $fetch(`${apiBase}/categorias`)
  }

  const getSubcategorias = async (categoriaId?: number) => {
    const url = categoriaId
      ? `${apiBase}/subcategorias?categoria_id=${categoriaId}`
      : `${apiBase}/subcategorias`
    return $fetch(url)
  }

  const getDepartamentos = async () => {
    return $fetch(`${apiBase}/departamentos`)
  }

  const getProvincias = async (departamentoId?: number) => {
    const url = departamentoId
      ? `${apiBase}/provincias?departamento_id=${departamentoId}`
      : `${apiBase}/provincias`
    return $fetch(url)
  }

  const getDias = async () => {
    return $fetch(`${apiBase}/dias`)
  }

  const getMeses = async () => {
    return $fetch(`${apiBase}/meses`)
  }

  return {
    getAtracciones,
    getAtraccion,
    createAtraccion,
    updateAtraccion,
    deleteAtraccion,
    addSubcategoria,
    removeSubcategoria,
    addFoto,
    removeFoto,
    getStats,
    getCategorias,
    getSubcategorias,
    getDepartamentos,
    getProvincias,
    getDias,
    getMeses
  }
}
