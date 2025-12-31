import { defineStore } from 'pinia'
import type { Notificacion } from '~/types/notificacion'

interface NotificacionesState {
  notificaciones: Notificacion[]
  noLeidas: number
  loading: boolean
  error: string | null
  ws: WebSocket | null
  reconnectAttempts: number
  maxReconnectAttempts: number
  reconnectDelay: number
}

export const useNotificacionesStore = defineStore('notificaciones', {
  state: (): NotificacionesState => ({
    notificaciones: [],
    noLeidas: 0,
    loading: false,
    error: null,
    ws: null,
    reconnectAttempts: 0,
    maxReconnectAttempts: 5,
    reconnectDelay: 3000, // 3 segundos
  }),

  getters: {
    notificacionesNoLeidas: (state) =>
      state.notificaciones.filter(n => !n.leida),

    notificacionesRecientes: (state) =>
      state.notificaciones.slice(0, 5),

    isConnected: (state) =>
      state.ws !== null && state.ws.readyState === WebSocket.OPEN,
  },

  actions: {
    async cargarNotificaciones(page = 1, limit = 10) {
      this.loading = true
      this.error = null

      try {
        const authStore = useAuthStore()
        const config = useRuntimeConfig()

        const response = await $fetch<{
          success: boolean
          data: {
            notificaciones: Notificacion[]
            no_leidas: number
            pagination: {
              page: number
              limit: number
              total: number
              total_pages: number
            }
          }
        }>(`${config.public.apiBase}/notificaciones?page=${page}&limit=${limit}`, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })

        if (response?.success) {
          this.notificaciones = response.data.notificaciones
          this.noLeidas = response.data.no_leidas
        }
      } catch (err: any) {
        this.error = err.message || 'Error al cargar notificaciones'
        console.error('Error cargando notificaciones:', err)
      } finally {
        this.loading = false
      }
    },

    async actualizarContadorNoLeidas() {
      try {
        const authStore = useAuthStore()
        const config = useRuntimeConfig()

        const response = await $fetch<{
          success: boolean
          data: {
            no_leidas: number
          }
        }>(`${config.public.apiBase}/notificaciones/no-leidas/count`, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })

        if (response?.success) {
          this.noLeidas = response.data.no_leidas
        }
      } catch (err) {
        console.error('Error actualizando contador:', err)
      }
    },

    async marcarComoLeida(id: number) {
      try {
        const authStore = useAuthStore()
        const config = useRuntimeConfig()

        await $fetch(`${config.public.apiBase}/notificaciones/${id}/marcar-leida`, {
          method: 'PUT',
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })

        // Actualizar localmente
        const notif = this.notificaciones.find(n => n.id === id)
        if (notif && !notif.leida) {
          notif.leida = true
          notif.fecha_leida = new Date().toISOString()
          this.noLeidas = Math.max(0, this.noLeidas - 1)
        }
      } catch (err) {
        console.error('Error marcando notificación como leída:', err)
      }
    },

    async marcarTodasLeidas() {
      try {
        const authStore = useAuthStore()
        const config = useRuntimeConfig()

        await $fetch(`${config.public.apiBase}/notificaciones/marcar-todas-leidas`, {
          method: 'PUT',
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })

        // Actualizar localmente
        this.notificaciones.forEach(n => {
          if (!n.leida) {
            n.leida = true
            n.fecha_leida = new Date().toISOString()
          }
        })
        this.noLeidas = 0
      } catch (err) {
        console.error('Error marcando todas como leídas:', err)
      }
    },

    async eliminarNotificacion(id: number) {
      try {
        const authStore = useAuthStore()
        const config = useRuntimeConfig()

        await $fetch(`${config.public.apiBase}/notificaciones/${id}`, {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })

        const index = this.notificaciones.findIndex(n => n.id === id)
        if (index !== -1) {
          const wasUnread = !this.notificaciones[index].leida
          this.notificaciones.splice(index, 1)
          if (wasUnread) {
            this.noLeidas = Math.max(0, this.noLeidas - 1)
          }
        }
      } catch (err) {
        console.error('Error eliminando notificación:', err)
      }
    },

    agregarNotificacion(notificacion: Notificacion) {
      // Agregar al inicio de la lista
      this.notificaciones.unshift(notificacion)

      // Si es nueva y no leída, incrementar contador
      if (!notificacion.leida) {
        this.noLeidas++
      }

      // Limitar a las últimas 50 notificaciones en memoria
      if (this.notificaciones.length > 50) {
        this.notificaciones = this.notificaciones.slice(0, 50)
      }
    },

    // WebSocket connection management
    conectarWebSocket() {
      try {
        const authStore = useAuthStore()

        if (!authStore?.token) {
          console.warn('No hay token disponible para WebSocket')
          return
        }

        if (this.ws?.readyState === WebSocket.OPEN) {
          console.log('WebSocket ya está conectado')
          return
        }

        const config = useRuntimeConfig()
        const wsUrl = `${config.public.wsBaseUrl}/api/v1/ws?token=${authStore.token}`

        this.ws = new WebSocket(wsUrl)

        this.ws.onopen = () => {
          console.log('✅ WebSocket conectado')
          this.reconnectAttempts = 0
        }

        this.ws.onmessage = (event) => {
          try {
            const notificacion: Notificacion = JSON.parse(event.data)
            this.agregarNotificacion(notificacion)

            // Mostrar notificación del sistema
            if ('Notification' in window && Notification.permission === 'granted') {
              new Notification(notificacion.titulo, {
                body: notificacion.mensaje,
                icon: '/favicon.ico'
              })
            }
          } catch (err) {
            console.error('Error procesando mensaje WebSocket:', err)
          }
        }

        this.ws.onerror = (error) => {
          console.error('Error en WebSocket:', error)
        }

        this.ws.onclose = () => {
          console.log('WebSocket desconectado')
          this.ws = null
          this.intentarReconectar()
        }
      } catch (err) {
        console.error('Error conectando WebSocket:', err)
        this.intentarReconectar()
      }
    },

    intentarReconectar() {
      const authStore = useAuthStore()
      if (!authStore?.token || !authStore.isAuthenticated) {
        this.reconnectAttempts = 0
        return
      }
      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        console.error('Máximo de intentos de reconexión alcanzado')
        return
      }

      this.reconnectAttempts++
      const delay = this.reconnectDelay * this.reconnectAttempts

      console.log(`Reintentando conexión en ${delay}ms (intento ${this.reconnectAttempts}/${this.maxReconnectAttempts})`)

      setTimeout(() => {
        this.conectarWebSocket()
      }, delay)
    },

    desconectarWebSocket() {
      if (this.ws) {
        this.ws.close()
        this.ws = null
      }
      this.reconnectAttempts = 0
    },

    solicitarPermisoNotificaciones() {
      if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission().then(permission => {
          console.log('Permiso de notificaciones:', permission)
        })
      }
    }
  }
})
