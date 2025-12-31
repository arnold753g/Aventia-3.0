<template>
  <div class="page-shell">
    <!-- Hero + filtros -->
    <section class="relative overflow-hidden bg-white border-b border-gray-200">
      <div class="absolute inset-0 bg-gradient-to-br from-emerald-50 via-white to-white" />
      <div class="relative max-w-7xl mx-auto px-4 py-10">
        <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-8">
          <div>
            <p class="text-xs uppercase tracking-[0.25em] text-emerald-700/70">Explora experiencias</p>
            <h1 class="text-4xl font-bold text-gray-900 mt-2">Paquetes turísticos</h1>
            <p class="text-gray-600 mt-2 max-w-2xl">
              Encuentra tours diseñados por agencias locales: salidas diarias o fechas únicas, con itinerarios y atracciones incluidas.
            </p>
          </div>

          <ClientOnly>
            <ScrollParallax :speed="0.08">
              <div class="w-28 h-28 rounded-3xl bg-white/80 border border-gray-200 shadow-sm flex items-center justify-center">
                <i class="pi pi-briefcase text-5xl text-emerald-600"></i>
              </div>
            </ScrollParallax>
          </ClientOnly>
        </div>

        <div class="mt-8 grid grid-cols-1 md:grid-cols-6 gap-4 items-end">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">Buscar</label>
            <span class="p-input-icon-left w-full">
              <i class="pi pi-search" />
              <InputText
                v-model="filters.search"
                placeholder="Nombre o descripción…"
                class="w-full"
                @input="debouncedReload"
              />
            </span>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Frecuencia</label>
            <Dropdown
              v-model="filters.frecuencia"
              :options="frecuenciaOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Todas"
              class="w-full"
              showClear
              @change="reload"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Duración</label>
            <Dropdown
              v-model="filters.tipo_duracion"
              :options="duracionOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Todas"
              class="w-full"
              showClear
              @change="reload"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Dificultad</label>
            <Dropdown
              v-model="filters.nivel_dificultad"
              :options="nivelesOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Todas"
              class="w-full"
              showClear
              @change="reload"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Orden</label>
            <Dropdown
              v-model="filters.sort"
              :options="sortOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Más recientes"
              class="w-full"
              @change="reload"
            />
          </div>

          <div class="md:col-span-3">
            <label class="block text-sm font-medium text-gray-700 mb-2">Precio mínimo (Bs.)</label>
            <InputNumber
              v-model="filters.precio_min"
              class="w-full"
              :min="0"
              :useGrouping="true"
              :maxFractionDigits="2"
              placeholder="0"
              @update:modelValue="debouncedReload"
            />
          </div>

          <div class="md:col-span-3">
            <label class="block text-sm font-medium text-gray-700 mb-2">Precio máximo (Bs.)</label>
            <InputNumber
              v-model="filters.precio_max"
              class="w-full"
              :min="0"
              :useGrouping="true"
              :maxFractionDigits="2"
              placeholder="Sin límite"
              @update:modelValue="debouncedReload"
            />
          </div>

          <div class="md:col-span-6 flex justify-end gap-2">
            <Button label="Limpiar filtros" icon="pi pi-filter-slash" severity="secondary" text @click="clearFilters" />
          </div>
        </div>
      </div>
    </section>

    <!-- Listado -->
    <div class="max-w-7xl mx-auto px-4 py-8">
      <div v-if="loading && paquetes.length === 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card v-for="n in 6" :key="n" class="surface-card">
          <template #content>
            <Skeleton height="200px" class="mb-4" />
            <Skeleton width="70%" height="1.4rem" class="mb-2" />
            <Skeleton width="100%" height="1rem" class="mb-2" />
            <Skeleton width="60%" height="1rem" />
          </template>
        </Card>
      </div>

      <div v-else-if="error" class="space-y-3">
        <Message severity="warn" :closable="false">{{ error }}</Message>
        <Button label="Reintentar" icon="pi pi-refresh" outlined @click="reload" />
      </div>

      <div v-else>
        <div v-if="paquetes.length" ref="cardsGrid" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <PaquetesPaqueteCard v-for="p in paquetes" :key="p.id" :paquete="p" />
        </div>

        <div v-else class="text-center py-16">
          <i class="pi pi-briefcase text-6xl text-gray-300 mb-3"></i>
          <p class="text-gray-800 font-semibold">No se encontraron paquetes</p>
          <p class="text-sm text-gray-500">Prueba ajustando los filtros o busca otro destino.</p>
          <Button label="Limpiar filtros" icon="pi pi-filter-slash" class="mt-4" outlined @click="clearFilters" />
        </div>

        <div v-if="pagination.total > 0" class="mt-6 text-sm text-gray-500">
          Mostrando {{ paquetes.length }} de {{ pagination.total }} paquetes
        </div>

        <div class="flex justify-center mt-8">
          <Button
            v-if="canLoadMore"
            label="Cargar más"
            icon="pi pi-plus"
            :loading="loadingMore"
            @click="loadMore"
          />
        </div>
      </div>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import ScrollParallax from 'vue3-parallax/src/components/ScrollParallax.vue'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const { getPaquetes } = usePaquetesTuristicos()

const paquetes = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const error = ref<string | null>(null)
const cardsGrid = ref<HTMLElement | null>(null)

const pagination = ref({
  page: 1,
  limit: 12,
  total: 0,
  total_pages: 0
})

const filters = ref({
  search: '',
  frecuencia: null as string | null,
  tipo_duracion: null as string | null,
  nivel_dificultad: null as string | null,
  precio_min: null as number | null,
  precio_max: null as number | null,
  sort: 'created_desc'
})

const frecuenciaOptions = [
  { label: 'Salida diaria', value: 'salida_diaria' },
  { label: 'Salida única', value: 'salida_unica' }
]

const duracionOptions = [
  { label: 'Un día', value: 'un_dia' },
  { label: 'Varios días', value: 'varios_dias' }
]

const nivelesOptions = [
  { label: 'Fácil', value: 'facil' },
  { label: 'Medio', value: 'medio' },
  { label: 'Difícil', value: 'dificil' },
  { label: 'Extremo', value: 'extremo' }
]

const sortOptions = [
  { label: 'Más recientes', value: 'created_desc', sort_by: 'created_at', sort_order: 'desc' },
  { label: 'Más antiguos', value: 'created_asc', sort_by: 'created_at', sort_order: 'asc' },
  { label: 'Precio: menor a mayor', value: 'precio_asc', sort_by: 'precio', sort_order: 'asc' },
  { label: 'Precio: mayor a menor', value: 'precio_desc', sort_by: 'precio', sort_order: 'desc' },
  { label: 'Nombre A-Z', value: 'nombre_asc', sort_by: 'nombre', sort_order: 'asc' },
  { label: 'Nombre Z-A', value: 'nombre_desc', sort_by: 'nombre', sort_order: 'desc' }
]

const selectedSort = computed(() => {
  return sortOptions.find(s => s.value === filters.value.sort) || sortOptions[0]
})

const canLoadMore = computed(() => {
  if (!pagination.value.total_pages) return false
  return pagination.value.page < pagination.value.total_pages
})

let debounceTimeout: ReturnType<typeof setTimeout> | null = null
const debouncedReload = () => {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => reload(), 450)
}

const animateCards = async () => {
  if (!cardsGrid.value) return
  await nextTick()
  try {
    const { gsap } = await import('gsap')
    gsap.from(cardsGrid.value.children, {
      opacity: 0,
      y: 16,
      duration: 0.45,
      stagger: 0.03,
      ease: 'power2.out'
    })
  } catch {
    // no-op
  }
}

const buildParams = () => {
  const params: Record<string, string> = {
    page: String(pagination.value.page),
    limit: String(pagination.value.limit),
    sort_by: String(selectedSort.value.sort_by),
    sort_order: String(selectedSort.value.sort_order)
  }
  if (filters.value.search) params.search = filters.value.search
  if (filters.value.frecuencia) params.frecuencia = String(filters.value.frecuencia)
  if (filters.value.tipo_duracion) params.tipo_duracion = String(filters.value.tipo_duracion)
  if (filters.value.nivel_dificultad) params.nivel_dificultad = String(filters.value.nivel_dificultad)
  if (filters.value.precio_min !== null && filters.value.precio_min !== undefined) params.precio_min = String(filters.value.precio_min)
  if (filters.value.precio_max !== null && filters.value.precio_max !== undefined) params.precio_max = String(filters.value.precio_max)
  return params
}

const reload = async () => {
  loading.value = true
  error.value = null
  pagination.value.page = 1
  try {
    const response: any = await getPaquetes(buildParams() as any)
    if (response.success) {
      paquetes.value = response.data?.paquetes || []
      pagination.value = { ...pagination.value, ...(response.data?.pagination || {}) }
      await animateCards()
      return
    }
    error.value = response?.error?.message || 'No se pudieron cargar los paquetes'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar los paquetes'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  if (!canLoadMore.value) return
  loadingMore.value = true
  error.value = null
  try {
    pagination.value.page += 1
    const response: any = await getPaquetes(buildParams() as any)
    if (response.success) {
      const more = response.data?.paquetes || []
      paquetes.value = [...paquetes.value, ...more]
      pagination.value = { ...pagination.value, ...(response.data?.pagination || {}) }
      await animateCards()
      return
    }
    error.value = response?.error?.message || 'No se pudieron cargar más paquetes'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar más paquetes'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loadingMore.value = false
  }
}

const clearFilters = () => {
  filters.value = {
    search: '',
    frecuencia: null,
    tipo_duracion: null,
    nivel_dificultad: null,
    precio_min: null,
    precio_max: null,
    sort: 'created_desc'
  }
  reload()
}

onMounted(async () => {
  await reload()
})
</script>

