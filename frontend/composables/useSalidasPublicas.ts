export const useSalidasPublicas = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  // Obtener salidas confirmadas (endpoint público)
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
    const url = qs ? `${apiBase}/public/salidas-confirmadas?${qs}` : `${apiBase}/public/salidas-confirmadas`
    return $fetch(url)
  }

  return {
    getSalidasConfirmadas
  }
}
