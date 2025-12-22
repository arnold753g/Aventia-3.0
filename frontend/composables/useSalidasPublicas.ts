import { useAuthStore } from '~/stores/auth'

export const useSalidasPublicas = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  const getSalidasConfirmadas = async (
    params: {
      page?: number
      limit?: number
      search?: string
      desde?: string // YYYY-MM-DD
      hasta?: string // YYYY-MM-DD
    } = {}
  ) => {
    const query = new URLSearchParams()
    Object.entries(params || {}).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') query.append(key, String(value))
    })
    const qs = query.toString()
    const url = qs ? `${apiBase}/salidas-confirmadas?${qs}` : `${apiBase}/salidas-confirmadas`
    return $fetch(url, { headers: authHeader() })
  }

  return {
    getSalidasConfirmadas
  }
}
