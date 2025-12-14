<template>
  <Card class="h-full hover:shadow-lg transition-shadow cursor-pointer" @click="navigateTo(`/admin/atracciones/${atraccion.id}`)">
    <template #header>
      <div class="relative h-48 overflow-hidden">
        <img
          v-if="fotoPrincipal && !imageError"
          :src="fotoPrincipal"
          :alt="atraccion.nombre"
          class="w-full h-full object-cover"
          @error="handleImageError"
          @load="handleImageLoad"
        />
        <div v-else class="w-full h-full bg-gradient-to-br from-blue-100 to-green-100 flex items-center justify-center">
          <i class="pi pi-image text-6xl text-gray-400"></i>
        </div>

        <!-- Badges superpuestos -->
        <div class="absolute top-2 left-2 flex flex-wrap gap-2">
          <span
            v-if="atraccion.nivel_dificultad"
            :class="[
              'px-2 py-1 rounded text-xs font-semibold',
              `bg-${getNivelDificultadColor(atraccion.nivel_dificultad)}-100`,
              `text-${getNivelDificultadColor(atraccion.nivel_dificultad)}-700`
            ]"
          >
            {{ getNivelDificultadLabel(atraccion.nivel_dificultad) }}
          </span>
          <span
            v-if="atraccion.precio_entrada === 0"
            class="px-2 py-1 bg-green-500 text-white rounded text-xs font-semibold"
          >
            GRATIS
          </span>
        </div>

        <div class="absolute top-2 right-2">
          <span
            :class="[
              'px-2 py-1 rounded text-xs font-semibold',
              `bg-${getStatusAtraccionColor(atraccion.status)}-100`,
              `text-${getStatusAtraccionColor(atraccion.status)}-700`
            ]"
          >
            {{ getStatusAtraccionLabel(atraccion.status) }}
          </span>
        </div>
      </div>
    </template>

    <template #content>
      <div class="space-y-3">
        <!-- Título -->
        <h3 class="text-xl font-bold text-gray-900 line-clamp-2">
          {{ atraccion.nombre }}
        </h3>

        <!-- Ubicación -->
        <div class="flex items-center gap-2 text-sm text-gray-600">
          <i class="pi pi-map-marker"></i>
          <span>
            {{ atraccion.provincia?.nombre }}, {{ atraccion.provincia?.departamento?.nombre }}
          </span>
        </div>

        <!-- Descripción -->
        <p class="text-sm text-gray-600 line-clamp-2">
          {{ atraccion.descripcion }}
        </p>

        <!-- Subcategorías -->
        <div v-if="atraccion.subcategorias && atraccion.subcategorias.length > 0" class="flex flex-wrap gap-1">
          <Chip
            v-for="subcat in atraccion.subcategorias.slice(0, 3)"
            :key="subcat.id"
            :label="subcat.subcategoria?.nombre"
            class="text-xs"
          />
          <Chip
            v-if="atraccion.subcategorias.length > 3"
            :label="`+${atraccion.subcategorias.length - 3}`"
            class="text-xs"
          />
        </div>

        <!-- Información adicional -->
        <div class="flex items-center justify-between pt-2 border-t">
          <div class="flex items-center gap-4 text-sm">
            <div v-if="atraccion.requiere_agencia" class="flex items-center gap-1 text-orange-600">
              <i class="pi pi-building"></i>
              <span>Agencia</span>
            </div>
            <div v-if="atraccion.acceso_particular" class="flex items-center gap-1 text-green-600">
              <i class="pi pi-car"></i>
              <span>Particular</span>
            </div>
          </div>
       <!--   <div class="text-lg font-bold text-blue-600">
            {{ formatPrecioBoliviano(atraccion.precio_entrada) }}
          </div> -->
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  getNivelDificultadLabel,
  getNivelDificultadColor,
  getStatusAtraccionLabel,
  getStatusAtraccionColor,
  formatPrecioBoliviano
} from '~/utils/formatters-atraccion'

const props = defineProps<{
  atraccion: any
}>()

const config = useRuntimeConfig()
const assetsBase = config.public.apiBase.replace(/\/api\/v1\/?$/, '')

const imageError = ref(false)

const resolveFotoUrl = (path?: string) => {
  if (!path) return ''
  let normalized = path.replace(/\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const fotoPrincipal = computed(() => {
  const principal = props.atraccion.fotos?.find((f: any) => f.es_principal)
  const path = principal?.foto || props.atraccion.fotos?.[0]?.foto
  return resolveFotoUrl(path)
})

const handleImageError = (event: Event) => {
  console.error('Error cargando foto de atracción:', {
    src: fotoPrincipal.value,
    atraccion: props.atraccion.nombre
  })
  imageError.value = true
}

const handleImageLoad = (event: Event) => {
  console.log('Foto cargada exitosamente:', {
    src: fotoPrincipal.value,
    atraccion: props.atraccion.nombre
  })
}
</script>

<script setup lang="ts">
</script>