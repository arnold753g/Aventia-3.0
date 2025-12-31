import { defineStore } from 'pinia'

interface User {
  id: number
  nombre: string
  apellido_paterno: string
  apellido_materno: string
  email: string
  rol: string
  status: string
  profile_photo: string | null
  ciudad: string | null
}

interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    refreshToken: null,
    isAuthenticated: false
  }),

  getters: {
    fullName: (state) => {
      if (!state.user) return ''
      return `${state.user.nombre} ${state.user.apellido_paterno}`
    },
    isAdmin: (state) => state.user?.rol === 'admin',
    isTurista: (state) => state.user?.rol === 'turista',
    isEncargado: (state) => state.user?.rol === 'encargado_agencia'
  },

  actions: {
    async login(email: string, password: string) {
      try {
        const config = useRuntimeConfig()
        const response: any = await $fetch(`${config.public.apiBase}/auth/login`, {
          method: 'POST',
          body: { email, password }
        })

        if (response.success) {
          this.setAuth(response.data)
          return { success: true, message: response.message }
        }

        return { success: false, error: response.error }
      } catch (error: any) {
        console.error('Login error:', error)
        return {
          success: false,
          error: error.data?.error?.message || 'Error al iniciar sesión'
        }
      }
    },

    async register(userData: any) {
      try {
        const config = useRuntimeConfig()

        // Convertir a FormData si hay un archivo
        let body: any = userData
        const hasFile = userData.profile_photo instanceof File

        if (hasFile) {
          const formData = new FormData()
          Object.entries(userData).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
              if (value instanceof File) {
                formData.append(key, value, value.name)
              } else {
                formData.append(key, String(value))
              }
            }
          })
          body = formData
          console.log('Registro con FormData (incluye foto)')
        } else {
          console.log('Registro sin foto (JSON)')
        }

        const response: any = await $fetch(`${config.public.apiBase}/auth/register`, {
          method: 'POST',
          body
        })

        if (response.success) {
          return { success: true, message: response.message }
        }

        return { success: false, error: response.error }
      } catch (error: any) {
        console.error('Register error:', error)
        return {
          success: false,
          error: error.data?.error?.message || 'Error al registrar usuario'
        }
      }
    },

    async refreshAccessToken() {
      if (!this.refreshToken) return false

      try {
        const config = useRuntimeConfig()
        const response: any = await $fetch(`${config.public.apiBase}/auth/refresh`, {
          method: 'POST',
          body: { refresh_token: this.refreshToken }
        })

        if (response.success) {
          this.token = response.data.token
          this.refreshToken = response.data.refresh_token
          return true
        }

        return false
      } catch (error) {
        console.error('Token refresh error:', error)
        this.logout()
        return false
      }
    },

    async getProfile() {
      if (!this.token) return false

      try {
        const config = useRuntimeConfig()
        const response: any = await $fetch(`${config.public.apiBase}/profile`, {
          headers: {
            Authorization: `Bearer ${this.token}`
          }
        })

        if (response.success) {
          this.user = response.data
          return true
        }

        return false
      } catch (error) {
        console.error('Get profile error:', error)
        return false
      }
    },

    setAuth(data: any) {
      this.user = data.user
      this.token = data.token
      this.refreshToken = data.refresh_token
      this.isAuthenticated = true
    },

    logout() {
      this.user = null
      this.token = null
      this.refreshToken = null
      this.isAuthenticated = false
      navigateTo('/login')
    }
  },

  persist: true
})
