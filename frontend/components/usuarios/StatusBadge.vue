<template>
  <span
    v-if="displayLabel"
    :class="[
      'inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold',
      'transition-all duration-200',
      colorClasses
    ]"
  >
    <i v-if="icon" :class="['pi', icon]"></i>
    <slot>{{ displayLabel }}</slot>
  </span>
  <span v-else class="text-xs text-gray-400 italic">Sin asignar</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  type: 'rol' | 'status'
  value: string
  label?: string
  icon?: string
}>()

// Computed para obtener el label a mostrar
const displayLabel = computed(() => {
  if (props.label && props.label.trim() !== '') {
    return props.label
  }
  if (props.value && props.value.trim() !== '') {
    return props.value
  }
  return ''
})

const colorClasses = computed(() => {
  // Si no hay valor, devolver estilo por defecto
  if (!props.value || props.value.trim() === '') {
    return 'bg-gray-100 text-gray-500 hover:bg-gray-200'
  }

  if (props.type === 'rol') {
    const colors: Record<string, string> = {
      admin: 'bg-purple-100 text-purple-700 hover:bg-purple-200',
      turista: 'bg-blue-100 text-blue-700 hover:bg-blue-200',
      encargado_agencia: 'bg-green-100 text-green-700 hover:bg-green-200'
    }
    return colors[props.value] || 'bg-gray-100 text-gray-700 hover:bg-gray-200'
  }

  if (props.type === 'status') {
    const colors: Record<string, string> = {
      active: 'bg-green-100 text-green-700 hover:bg-green-200',
      inactive: 'bg-orange-100 text-orange-700 hover:bg-orange-200',
      suspended: 'bg-red-100 text-red-700 hover:bg-red-200'
    }
    return colors[props.value] || 'bg-gray-100 text-gray-700 hover:bg-gray-200'
  }

  return 'bg-gray-100 text-gray-700 hover:bg-gray-200'
})
</script>
