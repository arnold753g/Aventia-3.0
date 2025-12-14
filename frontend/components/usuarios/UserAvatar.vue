<template>
  <div class="relative inline-block">
    <!-- Avatar con imagen -->
    <img
      v-if="src && !imageError"
      :src="src"
      :alt="alt"
      @error="handleImageError"
      @load="handleImageLoad"
      :class="[
        'rounded-full object-cover',
        sizeClasses,
        'ring-2 ring-white shadow-lg'
      ]"
    />

    <!-- Avatar con iniciales (fallback) -->
    <div
      v-else
      :class="[
        'rounded-full flex items-center justify-center font-semibold',
        sizeClasses,
        colorClasses,
        'ring-2 ring-white shadow-lg'
      ]"
    >
      {{ initials }}
    </div>

    <!-- Badge de status -->
    <div
      v-if="showStatus"
      :class="[
        'absolute bottom-0 right-0 rounded-full ring-2 ring-white',
        statusSizeClasses,
        statusColorClasses
      ]"
    ></div>

    <!-- Badge de rol -->
    <div
      v-if="showRolIcon"
      class="absolute -bottom-1 -right-1 bg-white rounded-full p-1 shadow-md"
    >
      <i :class="['text-xs', rolIconClasses]"></i>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { getInitials, getRolColor, getRolIcon } from '~/utils/formatters'

const props = defineProps<{
  src?: string
  nombre: string
  apellido: string
  status?: string
  rol?: string
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  showStatus?: boolean
  showRolIcon?: boolean
}>()

const imageError = ref(false)

const handleImageError = (event: Event) => {
  const target = event.target as HTMLImageElement
  console.error('❌ Error cargando imagen:', {
    src: props.src,
    usuario: `${props.nombre} ${props.apellido}`,
    errorEvent: event
  })
  imageError.value = true
}

const handleImageLoad = () => {
  //console.log('✅ Imagen cargada exitosamente:', {
  //   src: props.src,
  //   usuario: `${props.nombre} ${props.apellido}`
//  })
}

const alt = computed(() => `${props.nombre} ${props.apellido}`)

const initials = computed(() => getInitials(props.nombre, props.apellido))

const sizeClasses = computed(() => {
  const sizes = {
    xs: 'w-8 h-8 text-xs',
    sm: 'w-10 h-10 text-sm',
    md: 'w-12 h-12 text-base',
    lg: 'w-16 h-16 text-lg',
    xl: 'w-24 h-24 text-2xl'
  }
  return sizes[props.size || 'md']
})

const statusSizeClasses = computed(() => {
  const sizes = {
    xs: 'w-2 h-2',
    sm: 'w-2.5 h-2.5',
    md: 'w-3 h-3',
    lg: 'w-4 h-4',
    xl: 'w-5 h-5'
  }
  return sizes[props.size || 'md']
})

const colorClasses = computed(() => {
  if (!props.rol) return 'bg-gray-100 text-gray-600'

  const colors: Record<string, string> = {
    admin: 'bg-purple-100 text-purple-700',
    turista: 'bg-blue-100 text-blue-700',
    encargado_agencia: 'bg-green-100 text-green-700'
  }
  return colors[props.rol] || 'bg-gray-100 text-gray-600'
})

const statusColorClasses = computed(() => {
  if (!props.status) return ''

  const colors: Record<string, string> = {
    active: 'bg-green-500',
    inactive: 'bg-orange-500',
    suspended: 'bg-red-500'
  }
  return colors[props.status] || 'bg-gray-500'
})

const rolIconClasses = computed(() => {
  if (!props.rol) return ''

  const color = getRolColor(props.rol)
  const icon = getRolIcon(props.rol)

  return `${icon} text-${color}-600`
})
</script>
