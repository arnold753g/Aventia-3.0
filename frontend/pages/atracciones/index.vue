<template>
  <div class="page-shell">
    <!-- Hero + filtros -->
    <section class="relative overflow-hidden bg-white border-b border-gray-200">
      <div class="absolute inset-0 bg-gradient-to-br from-blue-50 via-white to-white" />
      <div class="relative max-w-7xl mx-auto px-4 py-10">
        <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-8">
          <div>
            <p class="text-xs uppercase tracking-[0.25em] text-blue-700/70">Descubre Bolivia</p>
            <h1 class="text-4xl font-bold text-gray-900 mt-2">Atracciones turísticas</h1>
            <p class="text-gray-600 mt-2 max-w-2xl">
              Explora destinos increíbles, guarda tus favoritos y visita su ubicación en el mapa.
            </p>
          </div>

          <ClientOnly>
            <ScrollParallax :speed="0.08">
              <div class="w-28 h-28 rounded-3xl bg-white/80 border border-gray-200 shadow-sm flex items-center justify-center">
                <i class="pi pi-map-marker text-5xl text-blue-600"></i>
              </div>
            </ScrollParallax>
          </ClientOnly>
        </div>

        <div class="mt-8 grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
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
            <label class="block text-sm font-medium text-gray-700 mb-2">Departamento</label>
            <Dropdown
              v-model="filters.departamento_id"
              :options="departamentos"
              optionLabel="nombre"
              optionValue="id"
              placeholder="Todos"
              class="w-full"
              showClear
              @change="onDepartamentoChange"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Provincia</label>
            <Dropdown
              v-model="filters.provincia_id"
              :options="provincias"
              optionLabel="nombre"
              optionValue="id"
              placeholder="Todas"
              class="w-full"
              showClear
              :disabled="!filters.departamento_id"
              @change="reload"
            />
          </div>

          <div class="md:col-span-2">
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

          <div class="md:col-span-2 flex justify-end gap-2">
            <Button
              label="Limpiar filtros"
              icon="pi pi-filter-slash"
              severity="secondary"
              text
              @click="clearFilters"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- Listado -->
    <div class="max-w-7xl mx-auto px-4 py-8">
      <div v-if="loading && atracciones.length === 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
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
        <div ref="cardsGrid" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <AtraccionesAtraccionCard
            v-for="atraccion in atracciones"
            :key="atraccion.id"
            :atraccion="atraccion"
          />
        </div>

        <div v-if="pagination.total > 0" class="mt-6 text-sm text-gray-500">
          Mostrando {{ atracciones.length }} de {{ pagination.total }} atracciones
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
  layout: 'home'
})

const toast = useToast()
const { getAtracciones, getDepartamentos, getProvincias } = useAtracciones()

const atracciones = ref<any[]>([])
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
  departamento_id: null as number | null,
  provincia_id: null as number | null,
  nivel_dificultad: null as string | null
})

const nivelesOptions = [
  { label: 'Fácil', value: 'facil' },
  { label: 'Medio', value: 'medio' },
  { label: 'Difícil', value: 'dificil' },
  { label: 'Extremo', value: 'extremo' }
]

const departamentos = ref<any[]>([])
const provincias = ref<any[]>([])

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
    visible_publico: 'true',
    status: 'activa',
    sort_by: 'created_at',
    sort_order: 'desc'
  }
  if (filters.value.search) params.search = filters.value.search
  if (filters.value.provincia_id) params.provincia_id = String(filters.value.provincia_id)
  if (filters.value.departamento_id) params.departamento_id = String(filters.value.departamento_id)
  if (filters.value.nivel_dificultad) params.nivel_dificultad = String(filters.value.nivel_dificultad)
  return params
}

const reload = async () => {
  loading.value = true
  error.value = null
  pagination.value.page = 1
  try {
    const response: any = await getAtracciones(buildParams() as any)
    if (response.success) {
      atracciones.value = response.data?.atracciones || []
      pagination.value = { ...pagination.value, ...(response.data?.pagination || {}) }
      await animateCards()
      return
    }
    error.value = response?.error?.message || 'No se pudieron cargar las atracciones'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar las atracciones'
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
    const response: any = await getAtracciones(buildParams() as any)
    if (response.success) {
      const more = response.data?.atracciones || []
      atracciones.value = [...atracciones.value, ...more]
      pagination.value = { ...pagination.value, ...(response.data?.pagination || {}) }
      await animateCards()
      return
    }
    error.value = response?.error?.message || 'No se pudieron cargar más atracciones'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar más atracciones'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loadingMore.value = false
  }
}

const loadDepartamentos = async () => {
  try {
    const response: any = await getDepartamentos()
    if (response.success) departamentos.value = response.data || []
  } catch {
    // no-op
  }
}

const loadProvincias = async (departamentoId: number) => {
  try {
    const response: any = await getProvincias(departamentoId)
    if (response.success) provincias.value = response.data || []
  } catch {
    provincias.value = []
  }
}

const onDepartamentoChange = async () => {
  filters.value.provincia_id = null
  if (filters.value.departamento_id) {
    await loadProvincias(filters.value.departamento_id)
  } else {
    provincias.value = []
  }
  reload()
}

const clearFilters = () => {
  filters.value = {
    search: '',
    departamento_id: null,
    provincia_id: null,
    nivel_dificultad: null
  }
  provincias.value = []
  reload()
}

onMounted(async () => {
  await loadDepartamentos()
  await reload()
})
</script>
