export const useNotifications = () => {
  const store = useNotificacionesStore()

  return {
    // State
    notificaciones: computed(() => store.notificaciones),
    noLeidas: computed(() => store.noLeidas),
    notificacionesNoLeidas: computed(() => store.notificacionesNoLeidas),
    notificacionesRecientes: computed(() => store.notificacionesRecientes),
    loading: computed(() => store.loading),
    error: computed(() => store.error),
    isConnected: computed(() => store.isConnected),

    // Actions
    cargarNotificaciones: store.cargarNotificaciones,
    actualizarContadorNoLeidas: store.actualizarContadorNoLeidas,
    marcarComoLeida: store.marcarComoLeida,
    marcarTodasLeidas: store.marcarTodasLeidas,
    eliminarNotificacion: store.eliminarNotificacion,
    conectarWebSocket: store.conectarWebSocket,
    desconectarWebSocket: store.desconectarWebSocket,
    solicitarPermisoNotificaciones: store.solicitarPermisoNotificaciones,
  }
}
