<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Atracciones Turísticas</h1>
            <p class="text-gray-600 mt-1">Gestiona las atracciones de Bolivia</p>
          </div>
          <div class="flex gap-2">
            <Button
              label="Exportar"
              icon="pi pi-download"
              severity="secondary"
              @click="handleExport"
            />
            <Button
              label="Nueva Atracción"
              icon="pi pi-plus"
              @click="navigateTo('/admin/atracciones/crear')"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-4 py-8">
      <!-- Estadísticas -->
      <AtraccionesAtraccionStats :stats="stats" />

      <!-- Filtros -->
      <AtraccionesAtraccionFilters @filter-change="handleFilterChange" />

      <!-- Toggle vista -->
      <div class="mb-6 flex justify-center">
        <SelectButton
          v-model="viewMode"
          :options="viewOptions"
          optionLabel="label"
          optionValue="value"
        >
          <template #option="{ option }">
            <i :class="option.icon"></i>
          </template>
        </SelectButton>
      </div>

      <!-- Vista de Cards -->
      <div v-if="viewMode === 'cards'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <AtraccionesAtraccionCard
          v-for="atraccion in atracciones"
          :key="atraccion.id"
          :atraccion="atraccion"
        />
      </div>

      <!-- Vista de Tabla -->
      <Card v-else>
        <template #content>
          <DataTable
            :value="atracciones"
            :loading="loading"
            stripedRows
            paginator
            :rows="pagination.limit"
            :totalRecords="pagination.total"
            :lazy="true"
            @page="onPage"
            dataKey="id"
          >
            <template #empty>
              <div class="text-center py-8">
                <i class="pi pi-map-marker text-4xl text-gray-400 mb-4"></i>
                <p class="text-gray-600">No se encontraron atracciones</p>
              </div>
            </template>

            <Column field="id" header="ID" style="width: 80px" />

            <Column header="Atracción">
              <template #body="{ data }">
                <div class="flex items-center gap-3">
                  <img
                    v-if="data.fotos?.[0]"
                    :src="resolveFotoUrl(data.fotos[0].foto)"
                    class="w-12 h-12 rounded object-cover"
                  />
                  <div class="w-12 h-12 bg-gradient-to-br from-blue-100 to-green-100 rounded flex items-center justify-center" v-else>
                    <i class="pi pi-image text-gray-400"></i>
                  </div>
                  <div>
                    <p class="font-semibold text-gray-900">{{ data.nombre }}</p>
                    <p class="text-sm text-gray-500">
                      {{ data.provincia?.nombre }}
                    </p>
                  </div>
                </div>
              </template>
            </Column>

            <Column header="Categorías">
              <template #body="{ data }">
                <div class="flex flex-wrap gap-1">
                  <Chip
                    v-for="subcat in data.subcategorias?.slice(0, 2)"
                    :key="subcat.id"
                    :label="subcat.subcategoria?.nombre"
                    class="text-xs"
                  />
                  <Chip
                    v-if="data.subcategorias?.length > 2"
                    :label="`+${data.subcategorias.length - 2}`"
                    class="text-xs"
                  />
                </div>
              </template>
            </Column>

            <Column header="Dificultad">
              <template #body="{ data }">
                <span
                  v-if="data.nivel_dificultad"
                  :class="[
                    'px-3 py-1 rounded-full text-xs font-semibold',
                    `bg-${getNivelDificultadColor(data.nivel_dificultad)}-100`,
                    `text-${getNivelDificultadColor(data.nivel_dificultad)}-700`
                  ]"
                >
                  {{ getNivelDificultadLabel(data.nivel_dificultad) }}
                </span>
              </template>
            </Column>

            <Column header="Precio">
              <template #body="{ data }">
                {{ formatPrecioBoliviano(data.precio_entrada) }}
              </template>
            </Column>

            <Column header="Estado">
              <template #body="{ data }">
                <span
                  :class="[
                    'px-3 py-1 rounded-full text-xs font-semibold',
                    `bg-${getStatusAtraccionColor(data.status)}-100`,
                    `text-${getStatusAtraccionColor(data.status)}-700`
                  ]"
                >
                  {{ getStatusAtraccionLabel(data.status) }}
                </span>
              </template>
            </Column>

            <Column header="Acciones" style="width: 150px">
              <template #body="{ data }">
                <div class="flex gap-2">
                  <Button
                    icon="pi pi-eye"
                    severity="info"
                    text
                    rounded
                    @click="navigateTo(`/admin/atracciones/${data.id}`)"
                    v-tooltip.top="'Ver Detalle'"
                  />
                  <Button
                    icon="pi pi-pencil"
                    severity="warning"
                    text
                    rounded
                    @click="navigateTo(`/admin/atracciones/${data.id}/editar`)"
                    v-tooltip.top="'Editar'"
                  />
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    rounded
                    @click="confirmDelete(data)"
                    v-tooltip.top="'Desactivar'"
                  />
                </div>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>

    <!-- Dialog de confirmación -->
    <Dialog
      v-model:visible="showDeleteDialog"
      header="Confirmar Desactivación"
      :modal="true"
      :style="{ width: '450px' }"
    >
      <div class="flex items-start gap-4">
        <i class="pi pi-exclamation-triangle text-orange-500 text-3xl"></i>
        <div>
          <p class="mb-2">¿Está seguro de desactivar esta atracción?</p>
          <p class="text-sm font-semibold">{{ selectedAtraccion?.nombre }}</p>
          <p class="text-sm text-gray-500 mt-2">
            La atracción quedará inactiva y no será visible al público.
          </p>
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancelar"
          severity="secondary"
          @click="showDeleteDialog = false"
        />
        <Button
          label="Desactivar"
          severity="danger"
          @click="handleDelete"
          :loading="deleting"
        />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'
import {
  getNivelDificultadLabel,
  getNivelDificultadColor,
  getStatusAtraccionLabel,
  getStatusAtraccionColor,
  formatPrecioBoliviano
} from '~/utils/formatters-atraccion'
import { exportToExcel } from '~/utils/export'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const toast = useToast()
const authStore = useAuthStore()
const { getAtracciones, deleteAtraccion, getStats } = useAtracciones()
const config = useRuntimeConfig()
const assetsBase = config.public.apiBase.replace(/\/api\/v1\/?$/, '')

const atracciones = ref<any[]>([])
const stats = ref<Record<string, any>>({})
const loading = ref(false)
const deleting = ref(false)
const showDeleteDialog = ref(false)
const selectedAtraccion = ref<any>(null)
const viewMode = ref('cards')

const viewOptions = [
  { label: 'Cards', value: 'cards', icon: 'pi pi-th-large' },
  { label: 'Tabla', value: 'table', icon: 'pi pi-list' }
]

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

const pagination = ref({
  page: 1,
  limit: 12,
  total: 0,
  total_pages: 0
})

const filters = ref({
  search: '',
  departamento_id: '',
  provincia_id: '',
  categoria_id: '',
  subcategoria_id: '',
  nivel_dificultad: '',
  status: ''
})

// Verificar que sea admin
if (!authStore.isAdmin) {
  navigateTo('/dashboard')
}

const loadAtracciones = async () => {
  loading.value = true
  try {
    const response: any = await getAtracciones({
      page: pagination.value.page,
      limit: pagination.value.limit,
      ...filters.value
    })

    if (response.success) {
      atracciones.value = response.data.atracciones
      pagination.value = {
        ...pagination.value,
        ...response.data.pagination
      }
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Error al cargar atracciones',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response: any = await getStats()
    if (response.success) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('Error al cargar estadísticas')
  }
}

const handleFilterChange = (newFilters: any) => {
  filters.value = newFilters
  pagination.value.page = 1
  loadAtracciones()
}

const onPage = (event: any) => {
  pagination.value.page = event.page + 1
  loadAtracciones()
}

const confirmDelete = (atraccion: any) => {
  selectedAtraccion.value = atraccion
  showDeleteDialog.value = true
}

const handleDelete = async () => {
  if (!selectedAtraccion.value) return

  deleting.value = true
  try {
    const response: any = await deleteAtraccion(selectedAtraccion.value.id)

    if (response.success) {
      toast.add({
        severity: 'success',
        summary: 'Atracción Desactivada',
        detail: 'La atracción ha sido desactivada exitosamente',
        life: 3000
      })

      showDeleteDialog.value = false
      loadAtracciones()
      loadStats()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'Error al desactivar atracción',
      life: 3000
    })
  } finally {
    deleting.value = false
  }
}

const handleExport = () => {
  const data = atracciones.value.map((a: any) => ({
    ID: a.id,
    Nombre: a.nombre,
    Provincia: a.provincia?.nombre,
    Departamento: a.provincia?.departamento?.nombre,
    Precio: a.precio_entrada,
    Dificultad: a.nivel_dificultad,
    Estado: a.status
  }))

  exportToExcel(data, 'atracciones')

  toast.add({
    severity: 'success',
    summary: 'Exportado',
    detail: 'Datos exportados a Excel',
    life: 3000
  })
}

onMounted(() => {
  loadAtracciones()
  loadStats()
})
</script>
