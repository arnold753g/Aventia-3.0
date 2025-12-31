<template>
  <Card class="h-full hover:shadow-lg transition-shadow cursor-pointer overflow-hidden" @click="goToDetail">
    <template #header>
      <div class="relative h-52 overflow-hidden">
        <img
          v-if="fotoPrincipal && !imageError"
          :src="fotoPrincipal"
          :alt="paquete.nombre"
          class="w-full h-full object-cover transition-transform duration-500 hover:scale-105"
          loading="lazy"
          @error="imageError = true"
        />
        <div v-else class="w-full h-full bg-gradient-to-br from-emerald-50 to-blue-100 flex items-center justify-center">
          <i class="pi pi-image text-6xl text-gray-400"></i>
        </div>

        <div class="absolute top-2 left-2 flex flex-wrap gap-2">
          <Tag :value="frecuenciaLabel" severity="info" />
          <Tag v-if="duracionChip" :value="duracionChip" severity="secondary" />
          <Tag v-if="dificultadLabel" :value="dificultadLabel" :severity="dificultadSeverity" />
          <Tag v-if="paquete.permite_privado" value="Privado" severity="warning" icon="pi pi-lock" />
        </div>

        <div
          v-if="paquete.frecuencia === 'salida_unica' && paquete.fecha_salida_fija"
          class="absolute bottom-2 left-2 bg-white/90 backdrop-blur px-2.5 py-1 rounded text-xs border border-gray-200"
        >
          <i class="pi pi-calendar mr-1"></i>{{ formatFecha(paquete.fecha_salida_fija) }}
        </div>

        <div class="absolute bottom-2 right-2 bg-white/90 backdrop-blur px-2.5 py-1 rounded text-xs border border-gray-200">
          <span class="font-semibold text-emerald-700">Bs. {{ formatMoney(paquete.precio_base_nacionales) }}</span>
        </div>
      </div>
    </template>

    <template #content>
      <div class="space-y-3">
        <div class="space-y-1">
          <h3 class="text-xl font-bold text-gray-900 line-clamp-2">
            {{ paquete.nombre }}
          </h3>

          <div v-if="paquete.agencia?.nombre_comercial" class="flex items-center gap-2 text-sm text-gray-600">
            <i class="pi pi-building"></i>
            <span class="truncate">
              {{ paquete.agencia.nombre_comercial }}
              <span v-if="paquete.agencia?.departamento?.nombre">· {{ paquete.agencia.departamento.nombre }}</span>
            </span>
          </div>
        </div>

        <p class="text-sm text-gray-600 line-clamp-2">
          {{ paquete.descripcion || 'Descubre un paquete turístico diseñado por expertos.' }}
        </p>

        <div class="flex flex-wrap gap-3 pt-3 border-t text-xs text-gray-700">
          <span class="inline-flex items-center gap-1">
            <i class="pi pi-clock"></i>{{ duracionTexto }}
          </span>
          <span class="inline-flex items-center gap-1">
            <i class="pi pi-users"></i>{{ paquete.cupo_minimo }}-{{ paquete.cupo_maximo }}
          </span>
          <span class="inline-flex items-center gap-1">
            <i class="pi pi-calendar-plus"></i>{{ paquete.dias_previos_compra || 1 }} días ant.
          </span>
        </div>

        <div v-if="Number(paquete.precio_adicional_extranjeros || 0) > 0" class="text-xs text-gray-500">
          +Bs. {{ formatMoney(paquete.precio_adicional_extranjeros) }} extranjeros
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = defineProps<{
  paquete: any
}>()

const route = useRoute()
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
  const fotos = props.paquete?.fotos || []
  const principal = fotos.find((f: any) => f.es_principal) || fotos[0]
  return principal?.foto ? resolveFotoUrl(principal.foto) : ''
})

const frecuenciaLabel = computed(() => {
  const map: Record<string, string> = {
    salida_diaria: 'Salida diaria',
    salida_unica: 'Salida única'
  }
  return map[props.paquete?.frecuencia || ''] || (props.paquete?.frecuencia || 'N/D')
})

const duracionChip = computed(() => {
  const dias = Number(props.paquete?.duracion_dias || 1)
  if (!Number.isFinite(dias) || dias <= 1) return null
  const noches = Number(props.paquete?.duracion_noches || (dias - 1))
  return `${dias}D/${Math.max(0, noches)}N`
})

const duracionTexto = computed(() => {
  const dias = Number(props.paquete?.duracion_dias || 1)
  if (!Number.isFinite(dias) || dias <= 1) return '1 día'
  const noches = Number(props.paquete?.duracion_noches || (dias - 1))
  return `${dias} días / ${Math.max(0, noches)} noches`
})

const dificultadLabel = computed(() => {
  const nivel = props.paquete?.nivel_dificultad
  if (!nivel) return null
  const map: Record<string, string> = { facil: 'Fácil', medio: 'Medio', dificil: 'Difícil', extremo: 'Extremo' }
  return map[nivel] || String(nivel)
})

const dificultadSeverity = computed(() => {
  const nivel = props.paquete?.nivel_dificultad
  const map: Record<string, string> = { facil: 'success', medio: 'info', dificil: 'warning', extremo: 'danger' }
  return map[nivel] || 'secondary'
})

const formatMoney = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatFecha = (value?: any) => {
  if (!value) return ''
  const raw = String(value)
  const datePart = raw.split('T')[0].split(' ')[0]
  const match = datePart.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`
  return datePart || raw
}

const goToDetail = () => {
  // Si está en página de admin
  if (route.path.startsWith('/admin')) {
    navigateTo(`/admin/paquetes/${props.paquete.id}`)
    return
  }

  // Si está en página de turista
  if (route.path.startsWith('/turista')) {
    navigateTo(`/turista/paquetes/${props.paquete.id}`)
    return
  }

  // Si está en página de agencia
  if (route.path.startsWith('/agencia')) {
    navigateTo(`/turista/paquetes/${props.paquete.id}`)
    return
  }

  // Si está en página pública (home, agencias públicas, etc)
  // Intentar ir a la página de turista (que pedirá login si no está autenticado)
  navigateTo(`/turista/paquetes/${props.paquete.id}`)
}
</script>
