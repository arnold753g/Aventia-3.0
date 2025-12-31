<template>
  <div
    v-if="isOpen"
    class="absolute right-0 top-full mt-2 w-96 bg-white dark:bg-gray-800 rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 z-50 max-h-[600px] flex flex-col"
  >
    <!-- Header -->
    <div class="p-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <i class="pi pi-bell text-blue-500"></i>
        <h3 class="font-semibold text-gray-900 dark:text-white">
          Notificaciones
        </h3>
        <span
          v-if="noLeidas > 0"
          class="px-2 py-0.5 text-xs font-bold text-white bg-red-500 rounded-full"
        >
          {{ noLeidas }}
        </span>
      </div>

      <button
        v-if="noLeidas > 0"
        @click="marcarTodasLeidas"
        class="text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
      >
        Marcar todas como leídas
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center p-8">
      <i class="pi pi-spin pi-spinner text-2xl text-blue-500"></i>
    </div>

    <!-- Lista de notificaciones -->
    <div v-else-if="notificaciones.length > 0" class="overflow-y-auto flex-1">
      <div
        v-for="notif in notificaciones"
        :key="notif.id"
        class="p-4 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer"
        :class="{ 'bg-blue-50 dark:bg-blue-900/20': !notif.leida }"
        @click="handleNotificationClick(notif)"
      >
        <div class="flex gap-3">
          <!-- Icono según tipo -->
          <div class="flex-shrink-0">
            <i
              :class="getIconClass(notif.tipo)"
              class="text-xl"
            ></i>
          </div>

          <!-- Contenido -->
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-2">
              <h4
                class="font-semibold text-sm"
                :class="notif.leida ? 'text-gray-700 dark:text-gray-300' : 'text-gray-900 dark:text-white'"
              >
                {{ notif.titulo }}
              </h4>
              <button
                @click.stop="eliminarNotificacion(notif.id)"
                class="text-gray-400 hover:text-red-500 dark:hover:text-red-400 flex-shrink-0"
              >
                <i class="pi pi-times text-sm"></i>
              </button>
            </div>

            <p
              class="text-sm mt-1 line-clamp-2"
              :class="notif.leida ? 'text-gray-600 dark:text-gray-400' : 'text-gray-700 dark:text-gray-300'"
            >
              {{ notif.mensaje }}
            </p>

            <!-- Datos adicionales -->
            <div v-if="notif.datos_json" class="mt-2 space-y-1">
              <div
                v-if="notif.datos_json.monto"
                class="text-xs text-gray-600 dark:text-gray-400"
              >
                <i class="pi pi-money-bill mr-1"></i>
                Monto: Bs {{ notif.datos_json.monto }}
              </div>
              <div
                v-if="notif.datos_json.paquete_nombre"
                class="text-xs text-gray-600 dark:text-gray-400"
              >
                <i class="pi pi-box mr-1"></i>
                {{ notif.datos_json.paquete_nombre }}
              </div>
            </div>

            <div class="flex items-center justify-between mt-2">
              <span class="text-xs text-gray-500 dark:text-gray-400">
                {{ formatearFecha(notif.created_at) }}
              </span>
              <span
                v-if="!notif.leida"
                class="w-2 h-2 bg-blue-500 rounded-full"
              ></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else class="flex flex-col items-center justify-center p-8 text-gray-500 dark:text-gray-400">
      <i class="pi pi-inbox text-4xl mb-2"></i>
      <p class="text-sm">No tienes notificaciones</p>
    </div>

    <!-- Footer con link a ver todas -->
    <div
      v-if="notificaciones.length > 0 && hasAllNotificationsRoute"
      class="p-3 border-t border-gray-200 dark:border-gray-700 text-center"
    >
      <NuxtLink
        to="/notificaciones"
        class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
        @click="$emit('close')"
      >
        Ver todas las notificaciones
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Notificacion } from '~/types/notificacion'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const {
  notificaciones,
  noLeidas,
  loading,
  marcarComoLeida,
  marcarTodasLeidas,
  eliminarNotificacion,
} = useNotifications()

const router = useRouter()
const hasAllNotificationsRoute = computed(() =>
  router.getRoutes().some((route) => route.path === '/notificaciones')
)

const handleNotificationClick = async (notif: Notificacion) => {
  if (!notif.leida) {
    await marcarComoLeida(notif.id)
  }

  // Redirigir según el tipo de notificación
  emit('close')

  switch (notif.tipo) {
    case 'nuevo_pago_pendiente':
      if (notif.datos_json.compra_id) {
        router.push(`/agencia/compras/${notif.datos_json.compra_id}`)
      }
      break
    case 'pago_confirmado':
    case 'pago_rechazado':
    case 'compra_expirada':
      if (notif.datos_json.compra_id) {
        router.push(`/turista/mis-compras/${notif.datos_json.compra_id}`)
      }
      break
  }
}

const getIconClass = (tipo: string) => {
  const iconMap: Record<string, string> = {
    nuevo_pago_pendiente: 'pi pi-dollar text-blue-500',
    pago_confirmado: 'pi pi-check-circle text-green-500',
    pago_rechazado: 'pi pi-times-circle text-red-500',
    compra_expirada: 'pi pi-clock text-orange-500',
  }
  return iconMap[tipo] || 'pi pi-info-circle text-gray-500'
}

const formatearFecha = (fecha: string) => {
  const date = new Date(fecha)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  const segundos = Math.floor(diff / 1000)
  const minutos = Math.floor(segundos / 60)
  const horas = Math.floor(minutos / 60)
  const dias = Math.floor(horas / 24)

  if (segundos < 60) return 'hace un momento'
  if (minutos < 60) return `hace ${minutos} ${minutos === 1 ? 'minuto' : 'minutos'}`
  if (horas < 24) return `hace ${horas} ${horas === 1 ? 'hora' : 'horas'}`
  if (dias < 7) return `hace ${dias} ${dias === 1 ? 'día' : 'días'}`

  return date.toLocaleDateString('es-ES', {
    day: '2-digit',
    month: 'short',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
  })
}

// Cerrar al presionar Escape
if (process.client) {
  const handleEscape = (e: KeyboardEvent) => {
    if (e.key === 'Escape' && props.isOpen) {
      emit('close')
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleEscape)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleEscape)
  })
}
</script>
