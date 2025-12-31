<template>
  <Card class="h-full hover:shadow-lg transition-shadow cursor-pointer overflow-hidden" @click="goToDetail">
    <template #header>
      <div class="relative h-52 overflow-hidden">
        <img
          v-if="fotoPrincipal && !imageError"
          :src="fotoPrincipal"
          :alt="agencia.nombre_comercial"
          class="w-full h-full object-cover transition-transform duration-500 hover:scale-105"
          loading="lazy"
          @error="imageError = true"
        />
        <div v-else class="w-full h-full bg-gradient-to-br from-sky-50 to-emerald-100 flex items-center justify-center">
          <i class="pi pi-image text-6xl text-gray-400"></i>
        </div>

        <div class="absolute top-2 left-2 flex flex-wrap gap-2">
          <Tag v-if="agencia.licencia_turistica" value="Licencia verificada" severity="success" icon="pi pi-check-circle" />
        </div>
      </div>
    </template>

    <template #content>
      <div class="space-y-3">
        <div class="space-y-1">
          <h3 class="text-xl font-bold text-gray-900 line-clamp-2">
            {{ agencia.nombre_comercial }}
          </h3>
          <div class="flex items-center gap-2 text-sm text-gray-600">
            <i class="pi pi-map-marker"></i>
            <span>
              {{ agencia.departamento?.nombre || 'Sin departamento' }}
            </span>
          </div>
          <p class="text-sm text-gray-500 line-clamp-1">{{ agencia.direccion || 'Direccion por definir' }}</p>
        </div>

        <p class="text-sm text-gray-600 line-clamp-2">
          {{ agencia.descripcion || 'Agencia de turismo registrada en el sistema.' }}
        </p>

        <div v-if="especialidades.length" class="flex flex-wrap gap-1">
          <Chip v-for="(item, idx) in especialidades.slice(0, 3)" :key="idx" :label="item" class="text-xs" />
          <Chip v-if="especialidades.length > 3" :label="`+${especialidades.length - 3}`" class="text-xs" />
        </div>

        <div class="flex flex-wrap gap-3 pt-2 border-t text-xs text-gray-700">
          <span class="inline-flex items-center gap-1">
            <i class="pi pi-phone"></i>{{ agencia.telefono || 'N/D' }}
          </span>
          <span class="inline-flex items-center gap-1">
            <i class="pi pi-envelope"></i>{{ agencia.email || 'N/D' }}
          </span>
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = defineProps<{
  agencia: any
}>()

const config = useRuntimeConfig()
const assetsBase = String(config.public.apiBase || '').replace(/\/api\/v1\/?$/, '')

const imageError = ref(false)

const resolveFotoUrl = (path?: string) => {
  if (!path) return ''
  let normalized = String(path).replace(/\\\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const fotoPrincipal = computed(() => {
  const fotos = props.agencia?.fotos || []
  const principal = fotos.find((f: any) => f.es_principal) || fotos[0]
  const path = principal?.foto_url || principal?.foto
  return path ? resolveFotoUrl(path) : ''
})

const especialidades = computed(() => {
  return (props.agencia?.especialidades || [])
    .map((item: any) => item?.categoria?.nombre || item?.categoria?.nombre_categoria || '')
    .filter(Boolean)
})

const goToDetail = () => {
  if (!props.agencia?.id) return
  navigateTo(`/turista/agencias/${props.agencia.id}`)
}
</script>
