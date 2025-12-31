<template>
  <button
    @click="togglePanel"
    class="relative p-2 text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
    title="Notificaciones"
  >
    <i class="pi pi-bell text-xl"></i>

    <!-- Badge con contador -->
    <span
      v-if="noLeidas > 0"
      class="absolute -top-1 -right-1 inline-flex items-center justify-center min-w-[20px] h-5 px-1 text-xs font-bold text-white bg-red-500 rounded-full animate-pulse"
    >
      {{ noLeidas > 99 ? '99+' : noLeidas }}
    </span>

    <!-- Indicador de conexión WebSocket -->
    <span
      v-if="!isConnected"
      class="absolute bottom-1 right-1 w-2 h-2 bg-yellow-500 rounded-full"
      title="Reconectando..."
    ></span>
  </button>
</template>

<script setup lang="ts">
const { noLeidas, isConnected } = useNotifications()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()

const togglePanel = () => {
  emit('toggle')
}
</script>
