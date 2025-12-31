<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between flex-wrap gap-3">
          <div>
            <p class="text-sm text-gray-500">Panel de agencia</p>
            <h1 class="text-3xl font-bold text-gray-900">Paquetes turísticos</h1>
            <p class="muted mt-1">Crea y gestiona tus paquetes turísticos.</p>
          </div>
          <div class="flex gap-2">
            <Button label="Nuevo paquete" icon="pi pi-plus" @click="navigateTo('/agencia/paquetes/nuevo')" />
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Card class="surface-card">
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <IconField class="md:col-span-2">
              <InputIcon class="pi pi-search" />
              <InputText
                v-model="filters.search"
                placeholder="Buscar por nombre o descripción"
                class="w-full"
                @keyup.enter="applyFilters"
              />
            </IconField>

            <Select
              v-model="filters.status"
              :options="statusOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Estado"
              class="w-full"
              showClear
            />

            <Select
              v-model="filters.frecuencia"
              :options="frecuenciaOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Frecuencia"
              class="w-full"
              showClear
            />
          </div>

          <div class="flex justify-end gap-2 mt-4">
            <Button label="Limpiar" icon="pi pi-filter-slash" severity="secondary" outlined @click="resetFilters" />
            <Button label="Aplicar" icon="pi pi-filter" @click="applyFilters" />
          </div>
        </template>
      </Card>

      <!-- Toggle vista -->
      <div class="mb-6 flex justify-center">
        <SelectButton v-model="viewMode" :options="viewOptions" optionLabel="label" optionValue="value">
          <template #option="{ option }">
            <i :class="option.icon"></i>
          </template>
        </SelectButton>
      </div>

      <!-- Vista Cards -->
      <div v-if="viewMode === 'cards'">
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card v-for="n in 6" :key="n" class="surface-card">
          <template #content>
            <Skeleton height="180px" class="mb-4" />
            <Skeleton width="70%" height="1.4rem" class="mb-2" />
            <Skeleton width="45%" height="1rem" class="mb-2" />
            <Skeleton width="90%" height="1rem" />
          </template>
        </Card>
      </div>

      <div v-else-if="paquetes.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card
          v-for="p in paquetes"
          :key="p.id"
          class="surface-card overflow-hidden hover:shadow-lg transition-shadow"
        >
          <template #header>
            <div class="h-44 bg-gray-100 relative">
              <img :src="getPaqueteCover(p)" :alt="p.nombre" class="w-full h-full object-cover" />
              <div class="absolute top-3 left-3 flex gap-2">
                <Tag :value="getStatusLabel(p.status)" :severity="getStatusSeverity(p.status)" />
                <Tag :value="getFrecuenciaLabel(p.frecuencia)" severity="info" />
              </div>
              <div
                v-if="p.frecuencia === 'salida_unica' && p.fecha_salida_fija"
                class="absolute bottom-3 left-3 bg-white/90 backdrop-blur px-2 py-1 rounded text-xs border border-gray-200"
              >
                <i class="pi pi-calendar mr-1"></i>{{ formatFecha(p.fecha_salida_fija) }}
              </div>
            </div>
          </template>

          <template #content>
            <div class="space-y-2">
              <div>

                <h3 class="text-lg font-semibold text-gray-900 line-clamp-1">{{ p.nombre }}</h3>
              </div>

              <p class="text-sm text-gray-600 line-clamp-2">
                {{ p.descripcion || 'Sin descripción' }}
              </p>

              <div class="flex flex-wrap gap-2 text-xs text-gray-700">
                <span class="inline-flex items-center gap-1">
                  <i class="pi pi-clock"></i>{{ formatDuracion(p) }}
                </span>
                <span class="inline-flex items-center gap-1">
                  <i class="pi pi-users"></i>{{ p.cupo_minimo }}-{{ p.cupo_maximo }}
                </span>
                <span class="inline-flex items-center gap-1">
                  <i class="pi pi-tag"></i>Bs. {{ formatNumber(p.precio_base_nacionales) }}
                </span>
              </div>

              <div v-if="Number(p.precio_adicional_extranjeros || 0) > 0" class="text-xs text-gray-600">
                +Bs. {{ formatNumber(p.precio_adicional_extranjeros) }} extranjeros
              </div>
            </div>
          </template>

          <template #footer>
            <div class="flex gap-2">
              <Button
                label="Gestionar"
                icon="pi pi-pencil"
                outlined
                class="flex-1"
                @click="navigateTo(`/agencia/paquetes/${p.id}`)"
              />
              <Button icon="pi pi-trash" severity="danger" text @click="confirmDelete(p)" v-tooltip.top="'Eliminar'" />
            </div>
          </template>
        </Card>
      </div>

      <div v-else class="text-center py-16">
        <i class="pi pi-briefcase text-5xl text-gray-300 mb-3"></i>
        <p class="text-gray-700 font-semibold">No hay paquetes</p>
        <p class="text-sm text-gray-500">Crea tu primer paquete para empezar.</p>
        <Button label="Nuevo paquete" icon="pi pi-plus" class="mt-4" @click="navigateTo('/agencia/paquetes/nuevo')" />
      </div>

      <div v-if="pagination.total_pages > 1" class="flex items-center justify-between mt-6">
        <Button
          label="Anterior"
          icon="pi pi-angle-left"
          outlined
          :disabled="pagination.page <= 1 || loading"
          @click="goToPage(pagination.page - 1)"
        />
        <span class="text-sm text-gray-600">
          Página {{ pagination.page }} de {{ pagination.total_pages }} ({{ pagination.total }} total)
        </span>
        <Button
          label="Siguiente"
          icon="pi pi-angle-right"
          iconPos="right"
          outlined
          :disabled="pagination.page >= pagination.total_pages || loading"
          @click="goToPage(pagination.page + 1)"
        />
      </div>
      </div>

      <!-- Vista Tabla -->
      <Card v-else class="surface-card">
        <template #content>
          <DataTable
            :value="paquetes"
            :loading="loading"
            stripedRows
            paginator
            :rows="pagination.limit"
            :first="(pagination.page - 1) * pagination.limit"
            :totalRecords="pagination.total"
            :lazy="true"
            @page="onPage"
            dataKey="id"
          >
            <template #empty>
              <div class="text-center py-10 text-gray-500">No se encontraron paquetes</div>
            </template>

            <Column field="id" header="#" style="width: 80px" />

            <Column header="Paquete">
              <template #body="{ data }">
                <div class="flex items-center gap-3">
                  <img
                    v-if="data.fotos?.length"
                    :src="getPaqueteCover(data)"
                    class="w-12 h-12 rounded object-cover"
                    :alt="data.nombre"
                  />
                  <div
                    v-else
                    class="w-12 h-12 bg-gradient-to-br from-slate-100 to-emerald-100 rounded flex items-center justify-center"
                  >
                    <i class="pi pi-image text-gray-400"></i>
                  </div>
                  <div class="min-w-0">
                    <p class="font-semibold text-gray-900 truncate">{{ data.nombre }}</p>
                    <p class="text-sm text-gray-500 truncate">{{ data.descripcion || 'Sin descripción' }}</p>
                  </div>
                </div>
              </template>
            </Column>

            <Column header="Frecuencia" style="width: 170px">
              <template #body="{ data }">
                <div class="space-y-1">
                  <Tag :value="getFrecuenciaLabel(data.frecuencia)" severity="info" />
                  <div v-if="data.frecuencia === 'salida_unica' && data.fecha_salida_fija" class="text-xs text-gray-500">
                    <i class="pi pi-calendar mr-1"></i>{{ formatFecha(data.fecha_salida_fija) }}
                  </div>
                </div>
              </template>
            </Column>

            <Column header="Duración" style="width: 150px">
              <template #body="{ data }">{{ formatDuracion(data) }}</template>
            </Column>

            <Column header="Cupos" style="width: 120px">
              <template #body="{ data }">
                <span class="text-sm text-gray-700">{{ data.cupo_minimo }}-{{ data.cupo_maximo }}</span>
              </template>
            </Column>

            <Column header="Precio" style="width: 170px">
              <template #body="{ data }">
                <div>
                  <p class="font-semibold text-gray-900">Bs. {{ formatNumber(data.precio_base_nacionales) }}</p>
                  <p v-if="Number(data.precio_adicional_extranjeros || 0) > 0" class="text-xs text-gray-500">
                    +Bs. {{ formatNumber(data.precio_adicional_extranjeros) }} extranjeros
                  </p>
                </div>
              </template>
            </Column>

            <Column header="Estado" style="width: 140px">
              <template #body="{ data }">
                <Tag :value="getStatusLabel(data.status)" :severity="getStatusSeverity(data.status)" />
              </template>
            </Column>

            <Column header="Acciones" style="width: 170px">
              <template #body="{ data }">
                <div class="flex gap-2 justify-end">
                  <Button
                    icon="pi pi-pencil"
                    text
                    severity="warning"
                    @click="navigateTo(`/agencia/paquetes/${data.id}`)"
                    v-tooltip.top="'Gestionar'"
                  />
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    @click="confirmDelete(data)"
                    v-tooltip.top="'Eliminar'"
                  />
                </div>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>

    <Dialog v-model:visible="showDeleteDialog" header="Confirmar" :modal="true" :style="{ width: '420px' }">
      <div class="flex items-start gap-3">
        <i class="pi pi-exclamation-triangle text-orange-500 text-2xl"></i>
        <div>
          <p class="font-semibold">Eliminar paquete</p>
          <p class="text-sm text-gray-600">Esta seguro de eliminar este paquete ?? ya no sera visible para usted ni para el turista</p>
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" severity="secondary" @click="showDeleteDialog = false" />
        <Button label="Eliminar" severity="danger" :loading="deleting" @click="handleDelete" />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const { getMiAgencia } = useAgencias()
const { getPaquetes, deletePaquete } = usePaquetes()

const agencia = ref<any>(null)
const loading = ref(true)

const paquetes = ref<any[]>([])
const pagination = ref({ page: 1, limit: 12, total: 0, total_pages: 1 })

const filters = ref({
  search: '',
  status: null as string | null,
  frecuencia: null as string | null
})

const statusOptions = [
  { label: 'Activo', value: 'activo' },
  { label: 'Inactivo', value: 'inactivo' },
  { label: 'Borrador', value: 'borrador' }
]

const frecuenciaOptions = [
  { label: 'Salida diaria', value: 'salida_diaria' },
  { label: 'Salida única', value: 'salida_unica' }
]

const viewMode = ref<'cards' | 'table'>('cards')
const viewOptions = [
  { label: 'Cards', value: 'cards', icon: 'pi pi-th-large' },
  { label: 'Tabla', value: 'table', icon: 'pi pi-list' }
]

const showDeleteDialog = ref(false)
const deleting = ref(false)
const paqueteToDelete = ref<any>(null)

const agenciaId = computed(() => Number(agencia.value?.id || 0))

const loadAgencia = async () => {
  const response: any = await getMiAgencia()
  if (response.success) {
    agencia.value = response.data
  }
}

const loadPaquetes = async (page = 1) => {
  const id = agenciaId.value
  if (!id) return

  loading.value = true
  try {
    const params: any = {
      page,
      limit: pagination.value.limit
    }

    if (filters.value.search.trim()) params.search = filters.value.search.trim()
    if (filters.value.status) params.status = filters.value.status
    if (filters.value.frecuencia) params.frecuencia = filters.value.frecuencia

    const response: any = await getPaquetes(id, params)
    if (response.success) {
      paquetes.value = response.data?.paquetes || []
      pagination.value = response.data?.pagination || pagination.value
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudieron cargar los paquetes',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const applyFilters = async () => {
  await loadPaquetes(1)
}

const resetFilters = async () => {
  filters.value = { search: '', status: null, frecuencia: null }
  await loadPaquetes(1)
}

const goToPage = async (page: number) => {
  await loadPaquetes(page)
}

const onPage = async (event: any) => {
  await loadPaquetes(event.page + 1)
}

const confirmDelete = (paquete: any) => {
  paqueteToDelete.value = paquete
  showDeleteDialog.value = true
}

const handleDelete = async () => {
  const id = agenciaId.value
  const p = paqueteToDelete.value
  if (!id || !p?.id) return

  deleting.value = true
  try {
    const response: any = await deletePaquete(id, Number(p.id))
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Eliminado', detail: 'Paquete eliminado', life: 2500 })
      showDeleteDialog.value = false
      paqueteToDelete.value = null
      await loadPaquetes(pagination.value.page)
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo eliminar',
      life: 3000
    })
  } finally {
    deleting.value = false
  }
}

const resolveFotoUrl = (path?: string) => {
  if (!path) return '/images/placeholder.svg'
  if (path.startsWith('http')) return path
  const clean = path.replace(/^\.?\//, '')

  const apiBase = useRuntimeConfig().public.apiBase as unknown as string
  let origin = ''
  if (typeof apiBase === 'string' && apiBase.startsWith('http')) {
    origin = new URL(apiBase).origin
  } else if (typeof window !== 'undefined') {
    origin = window.location.origin
  }

  return origin ? `${origin}/${clean}` : `/${clean}`
}

const getPaqueteCover = (p: any) => {
  const fotos = p?.fotos || []
  const principal = fotos.find((f: any) => f.es_principal) || fotos[0]
  return principal?.foto ? resolveFotoUrl(principal.foto) : '/images/placeholder.svg'
}

const getStatusLabel = (status?: string) => {
  const map: Record<string, string> = {
    activo: 'Activo',
    inactivo: 'Inactivo',
    borrador: 'Borrador',
    eliminado: 'Eliminado'
  }
  return map[status || ''] || (status || 'N/D')
}

const getStatusSeverity = (status?: string) => {
  const map: Record<string, string> = {
    activo: 'success',
    inactivo: 'warning',
    borrador: 'info',
    eliminado: 'danger'
  }
  return map[status || ''] || 'secondary'
}

const getFrecuenciaLabel = (frecuencia?: string) => {
  const map: Record<string, string> = {
    salida_diaria: 'Salida diaria',
    salida_unica: 'Salida única'
  }
  return map[frecuencia || ''] || (frecuencia || 'N/D')
}

const formatFecha = (value?: any) => {
  if (!value) return ''

  if (value instanceof Date) {
    const d = String(value.getUTCDate()).padStart(2, '0')
    const m = String(value.getUTCMonth() + 1).padStart(2, '0')
    const y = value.getUTCFullYear()
    return `${d}/${m}/${y}`
  }

  const raw = String(value)
  const datePart = raw.split('T')[0].split(' ')[0]
  const match = datePart.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`

  const parsed = new Date(raw)
  if (!Number.isNaN(parsed.getTime())) {
    const d = String(parsed.getUTCDate()).padStart(2, '0')
    const m = String(parsed.getUTCMonth() + 1).padStart(2, '0')
    const y = parsed.getUTCFullYear()
    return `${d}/${m}/${y}`
  }

  return datePart || raw
}

const formatDuracion = (p: any) => {
  const dias = Number(p?.duracion_dias || 1)
  if (dias > 1) {
    const noches = Number(p?.duracion_noches || (dias - 1))
    return `${dias} días / ${noches} noches`
  }
  return '1 día'
}

const formatNumber = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

onMounted(async () => {
  try {
    await loadAgencia()
    await loadPaquetes(1)
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo cargar el módulo',
      life: 3000
    })
    loading.value = false
  }
})
</script>
