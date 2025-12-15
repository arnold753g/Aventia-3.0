import { useAuthStore } from '~/stores/auth'

export const useAgencias = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  // Obtener lista de agencias
  const getAgencias = async (params: {
    page?: number
    limit?: number
    search?: string
    departamento_id?: string
    status?: string
    licencia_turistica?: string
    especialidad_id?: string
    encargado_id?: string
    visible_publico?: string
    sort_by?: string
    sort_order?: string
  } = {}) => {
    const query = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        query.append(key, String(value))
      }
    })

    const qs = query.toString()
    const url = qs ? `${apiBase}/agencias?${qs}` : `${apiBase}/agencias`

    return $fetch(url, {
      headers: authHeader()
    })
  }

  // Obtener una agencia
  const getAgencia = async (id: number) => {
    return $fetch(`${apiBase}/agencias/${id}`, {
      headers: authHeader()
    })
  }

  // Obtener mi agencia (encargado)
  const getMiAgencia = async () => {
    return $fetch(`${apiBase}/agencias/me`, {
      headers: authHeader()
    })
  }

  // Crear agencia rápida
  const createAgenciaRapida = async (data: {
    nombre_comercial: string
    departamento_id: number
    telefono: string
    encargado_principal_id: number
  }) => {
    return $fetch(`${apiBase}/agencias/rapida`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  // Crear agencia completa
  const createAgenciaCompleta = async (data: any) => {
    return $fetch(`${apiBase}/agencias/completa`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  // Actualizar agencia
  const updateAgencia = async (id: number, data: any) => {
    return $fetch(`${apiBase}/agencias/${id}`, {
      method: 'PUT',
      headers: authHeader(),
      body: data
    })
  }

  // Eliminar agencia
  const deleteAgencia = async (id: number) => {
    return $fetch(`${apiBase}/admin/agencias/${id}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Actualizar status
  const updateStatus = async (id: number, status: string) => {
    return $fetch(`${apiBase}/admin/agencias/${id}/status`, {
      method: 'PATCH',
      headers: authHeader(),
      body: { status }
    })
  }

  // Agregar foto (con archivo)
  const uploadFoto = async (agenciaId: number, formData: FormData) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/fotos/upload`, {
      method: 'POST',
      headers: authHeader(),
      body: formData
    })
  }

  // Eliminar foto
  const removeFoto = async (agenciaId: number, fotoId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/fotos/${fotoId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Agregar especialidad
  const addEspecialidad = async (agenciaId: number, data: {
    categoria_id: number
    es_principal?: boolean
  }) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/especialidades`, {
      method: 'POST',
      headers: authHeader(),
      body: data
    })
  }

  // Eliminar especialidad
  const removeEspecialidad = async (agenciaId: number, especialidadId: number) => {
    return $fetch(`${apiBase}/agencias/${agenciaId}/especialidades/${especialidadId}`, {
      method: 'DELETE',
      headers: authHeader()
    })
  }

  // Obtener estadísticas
  const getStats = async () => {
    return $fetch(`${apiBase}/admin/agencias/stats`, {
      headers: authHeader()
    })
  }

  // Obtener departamentos
  const getDepartamentos = async () => {
    return $fetch(`${apiBase}/agencias/data/departamentos`)
  }

  // Obtener categorías
  const getCategorias = async () => {
    return $fetch(`${apiBase}/agencias/data/categorias`)
  }

  // Obtener días
  const getDias = async () => {
    return $fetch(`${apiBase}/agencias/data/dias`)
  }

  // Obtener encargados
  const getEncargados = async (params: {
    only_unassigned?: boolean
    agencia_id?: number
  } = {}) => {
    const query = new URLSearchParams()
    if (params.only_unassigned !== undefined) query.set('only_unassigned', String(params.only_unassigned))
    if (params.agencia_id !== undefined) query.set('agencia_id', String(params.agencia_id))

    const qs = query.toString()
    const url = qs ? `${apiBase}/agencias/data/encargados?${qs}` : `${apiBase}/agencias/data/encargados`

    return $fetch(url)
  }

  return {
    getAgencias,
    getAgencia,
    getMiAgencia,
    createAgenciaRapida,
    createAgenciaCompleta,
    updateAgencia,
    deleteAgencia,
    updateStatus,
    uploadFoto,
    removeFoto,
    addEspecialidad,
    removeEspecialidad,
    getStats,
    getDepartamentos,
    getCategorias,
    getDias,
    getEncargados
  }
}
