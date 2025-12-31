<template>
  <div class="page-shell">
    <!-- Hero + filtros -->
    <section class="relative overflow-hidden bg-white border-b border-gray-200">
      <div class="absolute inset-0 bg-gradient-to-br from-sky-50 via-white to-white" />
      <div class="relative max-w-7xl mx-auto px-4 py-10">
        <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-8">
          <div>
            <p class="text-xs uppercase tracking-[0.25em] text-sky-700/70">Explora agencias</p>
            <h1 class="text-4xl font-bold text-gray-900 mt-2">Agencias de turismo</h1>
            <p class="text-gray-600 mt-2 max-w-2xl">
              Conoce agencias registradas, revisa su informacion y descubre sus paquetes disponibles.
            </p>
          </div>

          <ClientOnly>
            <ScrollParallax :speed="0.08">
              <div class="w-28 h-28 rounded-3xl bg-white/80 border border-gray-200 shadow-sm flex items-center justify-center">
                <i class="pi pi-building text-5xl text-sky-600"></i>
              </div>
            </ScrollParallax>
          </ClientOnly>
        </div>

        <div class="mt-8 grid grid-cols-1 md:grid-cols-5 gap-4 items-end">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">Buscar</label>
            <span class="p-input-icon-left w-full">
              <i class="pi pi-search" />
              <InputText
                v-model="filters.search"
                placeholder="Nombre o descripcion"
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
              @change="reload"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Especialidad</label>
            <Dropdown
              v-model="filters.especialidad_id"
              :options="especialidades"
              optionLabel="nombre"
              optionValue="id"
              placeholder="Todas"
              class="w-full"
              showClear
              @change="reload"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Licencia</label>
            <Dropdown
              v-model="filters.licencia_turistica"
              :options="licenciaOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Todas"
              class="w-full"
              showClear
              @change="reload"
            />
          </div>

          <div class="md:col-span-5 flex justify-end gap-2">
            <Button label="Limpiar filtros" icon="pi pi-filter-slash" severity="secondary" text @click="clearFilters" />
          </div>
        </div>
      </div>
    </section>

    <!-- Listado -->
    <div class="max-w-7xl mx-auto px-4 py-8">
      <div v-if="loading && agencias.length === 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
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
        <div v-if="agencias.length" ref="cardsGrid" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <AgenciasAgenciaCard v-for="agencia in agencias" :key="agencia.id" :agencia="agencia" />
        </div>

        <div v-else class="text-center py-16">
          <i class="pi pi-building text-6xl text-gray-300 mb-3"></i>
          <p class="text-gray-800 font-semibold">No se encontraron agencias</p>
          <p class="text-sm text-gray-500">Prueba ajustando los filtros o revisa mas tarde.</p>
          <Button label="Limpiar filtros" icon="pi pi-filter-slash" class="mt-4" outlined @click="clearFilters" />
        </div>

        <div v-if="pagination.total > 0" class="mt-6 text-sm text-gray-500">
          Mostrando {{ agencias.length }} de {{ pagination.total }} agencias
        </div>

        <div class="flex justify-center mt-8">
          <Button
            v-if="canLoadMore"
            label="Cargar mas"
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
const { getAgencias, getDepartamentos, getCategorias } = useAgencias()

const agencias = ref<any[]>([])
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
  especialidad_id: null as number | null,
  licencia_turistica: null as string | null
})

const departamentos = ref<any[]>([])
const especialidades = ref<any[]>([])

const licenciaOptions = [
  { label: 'Con licencia', value: 'true' },
  { label: 'Sin licencia', value: 'false' }
]

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
  if (filters.value.departamento_id) params.departamento_id = String(filters.value.departamento_id)
  if (filters.value.especialidad_id) params.especialidad_id = String(filters.value.especialidad_id)
  if (filters.value.licencia_turistica) params.licencia_turistica = String(filters.value.licencia_turistica)
  return params
}

const reload = async () => {
  loading.value = true
  error.value = null
  pagination.value.page = 1
  try {
    const response: any = await getAgencias(buildParams() as any)
    if (response.success) {
      agencias.value = response.data?.agencias || []
      pagination.value = { ...pagination.value, ...(response.data?.pagination || {}) }
      await animateCards()
      return
    }
    error.value = response?.error?.message || 'No se pudieron cargar las agencias'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar las agencias'
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
    const response: any = await getAgencias(buildParams() as any)
    if (response.success) {
      const more = response.data?.agencias || []
      agencias.value = [...agencias.value, ...more]
      pagination.value = { ...pagination.value, ...(response.data?.pagination || {}) }
      await animateCards()
      return
    }
    error.value = response?.error?.message || 'No se pudieron cargar mas agencias'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar mas agencias'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loadingMore.value = false
  }
}

const clearFilters = () => {
  filters.value = {
    search: '',
    departamento_id: null,
    especialidad_id: null,
    licencia_turistica: null
  }
  reload()
}

const loadDepartamentos = async () => {
  try {
    const response: any = await getDepartamentos()
    if (response.success) {
      departamentos.value = response.data || []
    }
  } catch {
    // no-op
  }
}

const loadEspecialidades = async () => {
  try {
    const response: any = await getCategorias()
    if (response.success) {
      especialidades.value = response.data || []
    }
  } catch {
    // no-op
  }
}

onMounted(async () => {
  await Promise.all([loadDepartamentos(), loadEspecialidades()])
  await reload()
})
</script>
