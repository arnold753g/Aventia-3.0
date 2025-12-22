<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Salidas confirmadas</h1>
          <p class="muted mt-1">
            Salidas habilitadas de paquetes turÍsticos que ya tienen cupos confirmados.
          </p>
        </div>
        <div class="flex gap-2">
          <Button
            label="Actualizar"
            icon="pi pi-refresh"
            severity="secondary"
            outlined
            :loading="loading"
            @click="reload"
          />
          <Button
            label="Ver paquetes"
            icon="pi pi-briefcase"
            severity="secondary"
            outlined
            @click="navigateTo('/turista/paquetes')"
          />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Card class="surface-card">
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-6 gap-4 items-end">
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Buscar paquete</label>
              <span class="p-input-icon-left w-full">
                <i class="pi pi-search" />
                <InputText
                  v-model="filters.search"
                  placeholder="Nombre del paquete"
                  class="w-full"
                  @input="debouncedReload"
                />
              </span>
            </div>

            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Desde</label>
              <DatePicker
                v-model="filters.desde"
                class="w-full"
                :minDate="minDesdeDate"
                dateFormat="dd/mm/yy"
                showIcon
                @update:modelValue="reload"
              />
            </div>

            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Hasta</label>
              <DatePicker
                v-model="filters.hasta"
                class="w-full"
                :minDate="filters.desde || minDesdeDate"
                dateFormat="dd/mm/yy"
                showIcon
                @update:modelValue="reload"
              />
            </div>

            <div class="md:col-span-6 flex justify-end">
              <Button label="Limpiar" icon="pi pi-filter-slash" severity="secondary" text @click="clearFilters" />
            </div>
          </div>
        </template>
      </Card>

      <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

      <div v-if="loading && salidas.length === 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card v-for="n in 6" :key="n" class="surface-card">
          <template #content>
            <Skeleton height="1.25rem" width="70%" class="mb-2" />
            <Skeleton height="1rem" width="45%" class="mb-4" />
            <Skeleton height="5rem" class="mb-3" />
            <Skeleton height="2.5rem" />
          </template>
        </Card>
      </div>

      <Card v-else-if="!loading && salidas.length === 0" class="surface-card">
        <template #content>
          <div class="text-center py-12">
            <i class="pi pi-calendar-times text-6xl text-gray-300 mb-4 block"></i>
            <p class="text-xl font-semibold text-gray-800">No hay salidas confirmadas</p>
            <p class="text-sm text-gray-500 mt-2">
              Ajusta los filtros o vuelve mÁs tarde.
            </p>
          </div>
        </template>
      </Card>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card
          v-for="salida in salidas"
          :key="salida.salida_id"
          class="h-full hover:shadow-lg transition-shadow overflow-hidden surface-card"
        >
          <template #header>
            <div class="relative h-52 overflow-hidden">
              <img
                v-if="portadaUrl(salida) && !isImageError(salida.salida_id)"
                :src="portadaUrl(salida)"
                :alt="salida.paquete_nombre"
                class="w-full h-full object-cover transition-transform duration-500 hover:scale-105"
                loading="lazy"
                @error="markImageError(salida.salida_id)"
              />
              <div v-else class="w-full h-full bg-gradient-to-br from-emerald-50 to-blue-100 flex items-center justify-center">
                <i class="pi pi-image text-5xl text-gray-400"></i>
              </div>

              <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/15 to-transparent" />
              <div class="absolute top-2 left-2 flex flex-wrap gap-2">
                <Tag :value="tipoSalidaLabel(salida.tipo_salida)" severity="info" />
                <Tag value="Salida confirmada" severity="success" icon="pi pi-check" />
                <Tag v-if="isMultiDay(salida)" value="Varios días" severity="secondary" icon="pi pi-sun" />
              </div>

              <div class="absolute bottom-2 left-2 bg-white/90 backdrop-blur px-2.5 py-1 rounded text-xs border border-gray-200">
                <i class="pi pi-calendar mr-1"></i>{{ formatFecha(salida.fecha_salida) }}
              </div>

              <div class="absolute bottom-2 right-2 text-xs text-gray-600 bg-white/90 backdrop-blur px-2.5 py-1 rounded border border-gray-200">
                <i class="pi pi-building mr-1"></i>{{ salida.agencia_nombre }}
              </div>
            </div>
          </template>

          <template #content>
            <div>
              <p class="text-base font-bold text-gray-900 line-clamp-2">
                {{ salida.paquete_nombre }}
              </p>
              <p class="text-xs text-gray-500 mt-1">
                {{ frecuenciaLabel(salida.paquete_frecuencia) }}
                <span v-if="isMultiDay(salida)">· {{ salida.paquete_duracion_dias }} días</span>
                <span v-else-if="horaSalidaText(salida)">· {{ horaSalidaText(salida) }}</span>
              </p>
            </div>

            <div class="mt-3 text-sm text-gray-600">
              <i class="pi pi-info-circle mr-2"></i>
              Esta fecha ya tiene compras confirmadas.
            </div>

            <div class="mt-4 flex justify-end gap-2">
              <Button
                label="Ver paquete"
                icon="pi pi-eye"
                severity="secondary"
                outlined
                size="small"
                @click="navigateTo(`/turista/paquetes/${salida.paquete_id}`)"
              />
            </div>
          </template>
        </Card>
      </div>

      <div class="flex justify-center">
        <Button
          v-if="canLoadMore"
          label="Cargar más"
          icon="pi pi-plus"
          :loading="loadingMore"
          @click="loadMore"
        />
      </div>

      <Toast />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import type { SalidaConfirmada } from '~/types/salida'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const { getSalidasConfirmadas } = useSalidasPublicas()
const config = useRuntimeConfig()
const assetsBase = String(config.public.apiBase || '').replace(/\/api\/v1\/?$/, '')

const salidas = ref<SalidaConfirmada[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const error = ref<string | null>(null)
const imageErrors = ref<Record<number, boolean>>({})

const resolveFotoUrl = (path?: string | null) => {
  if (!path) return ''
  let normalized = String(path).replace(/\\\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const portadaUrl = (salida: SalidaConfirmada) => {
  if (!salida?.paquete_foto) return ''
  return resolveFotoUrl(salida.paquete_foto)
}

const isImageError = (salidaId: number) => Boolean(imageErrors.value[salidaId])
const markImageError = (salidaId: number) => {
  imageErrors.value[salidaId] = true
}

const pagination = ref({
  page: 1,
  limit: 12,
  total: 0,
  total_pages: 0
})

const filters = ref({
  search: '',
  desde: getTomorrowDate(),
  hasta: null as Date | null
})

const canLoadMore = computed(() => pagination.value.page < pagination.value.total_pages)

function getTomorrowDate() {
  const d = new Date()
  d.setHours(0, 0, 0, 0)
  d.setDate(d.getDate() + 1)
  return d
}

const minDesdeDate = getTomorrowDate()

const pad2 = (n: number) => String(n).padStart(2, '0')
const toDateOnly = (d: Date) => `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())}`

let debounceTimeout: ReturnType<typeof setTimeout> | null = null
const debouncedReload = () => {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => reload(), 450)
}

watch(
  () => filters.value.desde,
  (value) => {
    if (!value) {
      filters.value.desde = getTomorrowDate()
      return
    }

    const normalized = new Date(value)
    normalized.setHours(0, 0, 0, 0)
    if (normalized < minDesdeDate) {
      filters.value.desde = getTomorrowDate()
    }
  }
)

watch(
  () => filters.value.hasta,
  (value) => {
    if (!value) return
    if (value < minDesdeDate) {
      filters.value.hasta = null
      return
    }
    if (filters.value.desde && value < filters.value.desde) {
      filters.value.hasta = null
    }
  }
)

const clearFilters = () => {
  filters.value.search = ''
  filters.value.desde = getTomorrowDate()
  filters.value.hasta = null
  reload()
}

const isMultiDay = (s: SalidaConfirmada) => Number(s.paquete_duracion_dias || 1) > 1

const tipoSalidaLabel = (tipo: string) => {
  const t = String(tipo || '').toLowerCase()
  if (t === 'privado') return 'Privado'
  if (t === 'compartido') return 'Compartido'
  return tipo || 'N/D'
}

const frecuenciaLabel = (frecuencia: string) => {
  const map: Record<string, string> = { salida_diaria: 'Salida diaria', salida_unica: 'Salida única' }
  return map[frecuencia] || (frecuencia || 'N/D')
}

const horaSalidaText = (s: SalidaConfirmada) => {
  const raw = s.paquete_hora_salida
  if (!raw) return ''
  const value = String(raw).trim()
  if (!value) return ''
  const match = value.match(/(\d{1,2}):(\d{2})/)
  if (match) return `${match[1].padStart(2, '0')}:${match[2]}`
  return value
}

const monthShortEs = ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic']
const fechaBadge = (value: string) => {
  const raw = String(value || '').split('T')[0].split(' ')[0]
  const match = raw.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (!match) return { day: '—', month: '' }
  const monthIndex = Number(match[2]) - 1
  return { day: match[3], month: monthShortEs[monthIndex] || match[2] }
}

const formatFecha = (value: string) => {
  const raw = String(value || '').split('T')[0].split(' ')[0]
  const match = raw.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`
  return raw || value
}

const load = async (opts: { append: boolean }) => {
  const targetPage = opts.append ? pagination.value.page + 1 : 1
  error.value = null

  const desde = filters.value.desde ? toDateOnly(filters.value.desde) : undefined
  const hasta = filters.value.hasta ? toDateOnly(filters.value.hasta) : undefined

  try {
    const response: any = await getSalidasConfirmadas({
      page: targetPage,
      limit: pagination.value.limit,
      search: filters.value.search?.trim() || undefined,
      desde,
      hasta
    })

    if (!response?.success) {
      throw new Error(response?.error?.message || 'No se pudieron cargar las salidas')
    }

    const data = response.data
    const items: SalidaConfirmada[] = data?.salidas || []
    if (opts.append) {
      salidas.value = salidas.value.concat(items)
    } else {
      salidas.value = items
    }

    pagination.value.page = data?.pagination?.page || targetPage
    pagination.value.limit = data?.pagination?.limit || pagination.value.limit
    pagination.value.total = data?.pagination?.total || 0
    pagination.value.total_pages = data?.pagination?.total_pages || 0
  } catch (err: any) {
    const msg = err?.data?.error?.message || err?.message || 'No se pudieron cargar las salidas'
    error.value = msg
    toast.add({ severity: 'error', summary: 'Error', detail: msg, life: 3000 })
  }
}

const reload = async () => {
  loading.value = true
  try {
    await load({ append: false })
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  if (!canLoadMore.value || loadingMore.value) return
  loadingMore.value = true
  try {
    await load({ append: true })
  } finally {
    loadingMore.value = false
  }
}

onMounted(() => {
  reload()
})
</script>
