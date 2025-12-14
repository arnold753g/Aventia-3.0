export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  const get = async (endpoint: string) => {
    try {
      const response = await $fetch(`${apiBase}${endpoint}`)
      return response
    } catch (error) {
      console.error('API Error:', error)
      throw error
    }
  }

  const post = async (endpoint: string, data: any) => {
    try {
      const response = await $fetch(`${apiBase}${endpoint}`, {
        method: 'POST',
        body: data
      })
      return response
    } catch (error) {
      console.error('API Error:', error)
      throw error
    }
  }

  const put = async (endpoint: string, data: any) => {
    try {
      const response = await $fetch(`${apiBase}${endpoint}`, {
        method: 'PUT',
        body: data
      })
      return response
    } catch (error) {
      console.error('API Error:', error)
      throw error
    }
  }

  const del = async (endpoint: string) => {
    try {
      const response = await $fetch(`${apiBase}${endpoint}`, {
        method: 'DELETE'
      })
      return response
    } catch (error) {
      console.error('API Error:', error)
      throw error
    }
  }

  return {
    get,
    post,
    put,
    delete: del
  }
}
