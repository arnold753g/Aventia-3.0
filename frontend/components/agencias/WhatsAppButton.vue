<template>
  <Transition name="fade-scale">
    <button
      v-if="telefono"
      type="button"
      class="fixed bottom-6 right-6 z-50 h-14 w-14 rounded-full bg-green-500 hover:bg-green-600 shadow-lg hover:shadow-xl transition-all flex items-center justify-center group focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-green-300"
      :aria-label="`Contactar por WhatsApp a ${nombreAgencia}`"
      @click="openWhatsApp"
    >
      <i class="pi pi-whatsapp text-2xl text-white"></i>

      <!-- Tooltip -->
      <div class="absolute right-full mr-3 px-3 py-2 bg-gray-900 text-white text-sm rounded-lg whitespace-nowrap opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none">
        Contactar por WhatsApp
        <div class="absolute top-1/2 -right-1 -translate-y-1/2 w-2 h-2 bg-gray-900 rotate-45"></div>
      </div>
    </button>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  telefono?: string
  nombreAgencia?: string
  mensajeInicial?: string
}>()

const openWhatsApp = () => {
  if (!props.telefono) return

  // Limpiar el teléfono (remover espacios, guiones, etc.)
  const cleanPhone = props.telefono.replace(/\D/g, '')

  // Construir mensaje inicial
  const mensaje = props.mensajeInicial ||
    `Hola, estoy interesado en los servicios de ${props.nombreAgencia || 'su agencia'}. ¿Podrían darme más información?`

  // Construir URL de WhatsApp
  const whatsappUrl = `https://wa.me/591${cleanPhone}?text=${encodeURIComponent(mensaje)}`

  // Abrir en nueva ventana
  if (process.client) {
    window.open(whatsappUrl, '_blank', 'noopener,noreferrer')
  }
}
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.8);
}
</style>
