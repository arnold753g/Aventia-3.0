<template>
  <!-- Este componente no tiene UI, solo maneja la conexión WebSocket -->
  <div v-show="false"></div>
</template>

<script setup lang="ts">
const { conectarWebSocket, desconectarWebSocket, cargarNotificaciones, actualizarContadorNoLeidas, solicitarPermisoNotificaciones } = useNotifications()

let intervalId: NodeJS.Timeout | null = null

// Conectar WebSocket cuando el componente se monta y el usuario está autenticado
onMounted(async () => {
  // Verificar que estamos en el cliente y hay un usuario autenticado
  if (!process.client) return

  try {
    const authStore = useAuthStore()

    if (authStore?.isAuthenticated) {
      // Solicitar permisos de notificación del navegador
      solicitarPermisoNotificaciones()

      // Cargar notificaciones iniciales
      await cargarNotificaciones(1, 20)

      // Conectar WebSocket para recibir notificaciones en tiempo real
      conectarWebSocket()

      // Actualizar contador periódicamente como fallback
      intervalId = setInterval(() => {
        if (authStore.isAuthenticated) {
          actualizarContadorNoLeidas()
        }
      }, 60000) // Cada minuto
    }
  } catch (error) {
    console.error('Error inicializando notificaciones:', error)
  }
})

// Limpiar al desmontar
onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
  desconectarWebSocket()
})

// Reconectar si el usuario hace login
if (process.client) {
  watch(
    () => {
      try {
        const authStore = useAuthStore()
        return authStore?.isAuthenticated ?? false
      } catch {
        return false
      }
    },
    (isAuth) => {
      if (isAuth) {
        conectarWebSocket()
        cargarNotificaciones(1, 20)
      } else {
        desconectarWebSocket()
      }
    }
  )

  // Reconectar cuando la pestaña vuelve a estar visible
  document.addEventListener('visibilitychange', () => {
    try {
      const authStore = useAuthStore()
      if (!document.hidden && authStore?.isAuthenticated) {
        conectarWebSocket()
        actualizarContadorNoLeidas()
      }
    } catch {
      // Ignorar si authStore no está disponible
    }
  })
}
</script>
