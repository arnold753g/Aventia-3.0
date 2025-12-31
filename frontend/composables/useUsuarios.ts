import { useAuthStore } from '~/stores/auth'

export const useUsuarios = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => ({
    Authorization: `Bearer ${authStore.token}`
  })

  const toFormData = (payload: any) => {
    if (payload instanceof FormData) return payload

    const fd = new FormData()

    Object.entries(payload || {}).forEach(([key, value]) => {
      if (value === undefined || value === null) {
        return // Skip undefined/null values
      }

      // Handle File objects properly
      if (value instanceof File) {
        fd.append(key, value, value.name)
      }
      // Handle Blob objects
      else if (value instanceof Blob) {
        fd.append(key, value)
      }
      // Handle other values as strings
      else {
        fd.append(key, String(value))
      }
    })

    return fd
  }

  const getUsuarios = async (params: {
    page?: number
    limit?: number
    search?: string
    rol?: string
    status?: string
    sort_by?: string
    sort_order?: string
  }) => {
    const queryParams = new URLSearchParams()

    if (params.page) queryParams.append('page', params.page.toString())
    if (params.limit) queryParams.append('limit', params.limit.toString())
    if (params.search) queryParams.append('search', params.search)
    if (params.rol) queryParams.append('rol', params.rol)
    if (params.status) queryParams.append('status', params.status)
    if (params.sort_by) queryParams.append('sort_by', params.sort_by)
    if (params.sort_order) queryParams.append('sort_order', params.sort_order)

    return $fetch(`${apiBase}/admin/usuarios?${queryParams.toString()}`, {
      headers: authHeader()
    })
  }

  const getUsuario = async (id: number) => {
    return $fetch(`${apiBase}/usuarios/${id}`, {
      headers: authHeader()
    })
  }

  const createUsuario = async (data: any) => {
    return $fetch(`${apiBase}/admin/usuarios`, {
      method: 'POST',
      headers: authHeader(),
      body: toFormData(data)
    })
  }

  const updateUsuario = async (id: number, data: any) => {
    return $fetch(`${apiBase}/usuarios/${id}`, {
      method: 'PUT',
      headers: authHeader(),
      body: toFormData(data)
    })
  }

  const updateRol = async (id: number, rol: string) => {
    return $fetch(`${apiBase}/admin/usuarios/${id}/rol`, {
      method: 'PATCH',
      headers: authHeader(),
      body: { rol }
    })
  }

  const updateStatus = async (id: number, status: string) => {
    return $fetch(`${apiBase}/admin/usuarios/${id}/status`, {
      method: 'PATCH',
      headers: authHeader(),
      body: { status }
    })
  }

  const deactivateUsuario = async (id: number) => {
    return $fetch(`${apiBase}/admin/usuarios/${id}/deactivate`, {
      method: 'POST',
      headers: authHeader()
    })
  }

  const getStats = async () => {
    return $fetch(`${apiBase}/admin/usuarios/stats`, {
      headers: authHeader()
    })
  }

  return {
    getUsuarios,
    getUsuario,
    createUsuario,
    updateUsuario,
    updateRol,
    updateStatus,
    deactivateUsuario,
    getStats
  }
}
