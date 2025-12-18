<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between flex-wrap gap-3">
          <div>
            <p class="text-sm text-gray-500">Modulo de agencias turisticas</p>
            <h1 class="text-3xl font-bold text-gray-900">Agencias</h1>
          </div>
          <div class="flex gap-2">
            <Button
              label="Crear completa"
              icon="pi pi-list"
              severity="secondary"
              @click="navigateTo('/admin/agencias/crear-completa')"
            />
            <Button
              label="Crear rapida"
              icon="pi pi-plus"
              @click="navigateTo('/admin/agencias/crear')"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <!-- Stats -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card class="bg-gradient-to-br from-blue-50 to-blue-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-blue-600 font-medium">Total</p>
                <p class="text-3xl font-bold text-blue-700">{{ stats.total || 0 }}</p>
              </div>
              <span class="w-12 h-12 rounded-full bg-blue-200 flex items-center justify-center">
                <i class="pi pi-building text-2xl text-blue-700"></i>
              </span>
            </div>
          </template>
        </Card>

        <Card class="bg-gradient-to-br from-green-50 to-green-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-green-600 font-medium">Activas</p>
                <p class="text-3xl font-bold text-green-700">{{ stats.activas || 0 }}</p>
              </div>
              <span class="w-12 h-12 rounded-full bg-green-200 flex items-center justify-center">
                <i class="pi pi-check-circle text-2xl text-green-700"></i>
              </span>
            </div>
          </template>
        </Card>

        <Card class="bg-gradient-to-br from-orange-50 to-orange-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-orange-600 font-medium">En revision</p>
                <p class="text-3xl font-bold text-orange-700">{{ stats.en_revision || 0 }}</p>
              </div>
              <span class="w-12 h-12 rounded-full bg-orange-200 flex items-center justify-center">
                <i class="pi pi-clock text-2xl text-orange-700"></i>
              </span>
            </div>
          </template>
        </Card>

        <Card class="bg-gradient-to-br from-indigo-50 to-indigo-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-indigo-600 font-medium">Con licencia</p>
                <p class="text-3xl font-bold text-indigo-700">{{ stats.con_licencia || 0 }}</p>
              </div>
              <span class="w-12 h-12 rounded-full bg-indigo-200 flex items-center justify-center">
                <i class="pi pi-id-card text-2xl text-indigo-700"></i>
              </span>
            </div>
          </template>
        </Card>
      </div>

      <!-- Filters -->
      <Card class="surface-card">
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <IconField class="md:col-span-2">
              <InputIcon class="pi pi-search" />
              <InputText v-model="filters.search" placeholder="Buscar por nombre o direccion" class="w-full" />
            </IconField>

            <Select
              v-model="filters.departamento_id"
              :options="departamentos"
              optionLabel="nombre"
              optionValue="id"
              placeholder="Departamento"
              class="w-full"
              showClear
            />

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
              v-model="filters.licencia_turistica"
              :options="licenciaOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Licencia turistica"
              class="w-full"
              showClear
            />
          </div>
        </template>
      </Card>

      <!-- Toggle view -->
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

      <!-- Cards view -->
      <div v-if="viewMode === 'cards'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <Card v-for="agencia in agencias" :key="agencia.id" class="surface-card h-full">
          <template #content>
            <div class="flex justify-between items-start">
              <div>
                
                <h3 class="text-xl font-semibold text-gray-900">{{ agencia.nombre_comercial }}</h3>
                <p class="text-sm text-gray-600">{{ agencia.direccion }}</p>
                <p class="text-sm text-gray-600 mt-1">
                  {{ agencia.departamento?.nombre || 'Sin departamento' }}
                </p>
              </div>
              <Tag
                :value="getStatusLabel(agencia.status)"
                :severity="getStatusSeverity(agencia.status)"
              />
            </div>

            <div class="mt-3 flex flex-wrap gap-2 text-sm text-gray-700">
              <span class="inline-flex items-center gap-1">
                <i class="pi pi-phone"></i>{{ agencia.telefono || 'N/D' }}
              </span>
              <span class="inline-flex items-center gap-1">
                <i class="pi pi-envelope"></i>{{ agencia.email || 'N/D' }}
              </span>
            </div>

            <div class="mt-4 flex gap-2">
              <Button label="Ver" size="small" icon="pi pi-eye" outlined @click="navigateTo(`/admin/agencias/${agencia.id}`)" />
              <Button label="Editar" size="small" icon="pi pi-pencil" outlined severity="warning" @click="navigateTo(`/admin/agencias/${agencia.id}/editar`)" />
              <Button label="Eliminar" size="small" icon="pi pi-trash" outlined severity="danger" @click="confirmDelete(agencia)" />
            </div>
          </template>
        </Card>
      </div>

      <!-- Table view -->
      <Card v-else>
        <template #content>
          <DataTable
            :value="agencias"
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
              <div class="text-center py-10 text-gray-500">
                No se encontraron agencias
              </div>
            </template>

            <Column field="id" header="ID" style="width: 70px" />
            <Column header="Agencia">
              <template #body="{ data }">
                <div>
                  <p class="font-semibold text-gray-900">{{ data.nombre_comercial }}</p>
                  <p class="text-sm text-gray-600">{{ data.direccion }}</p>
                  <p class="text-xs text-gray-500">{{ data.departamento?.nombre }}</p>
                </div>
              </template>
            </Column>
            <Column header="Contacto">
              <template #body="{ data }">
                <p class="text-sm">{{ data.telefono }}</p>
                <p class="text-xs text-gray-500">{{ data.email }}</p>
              </template>
            </Column>
            <Column header="Licencia">
              <template #body="{ data }">
                <Tag :value="data.licencia_turistica ? 'Si' : 'No'" :severity="data.licencia_turistica ? 'success' : 'warning'" />
              </template>
            </Column>
            <Column header="Estado">
              <template #body="{ data }">
                <Tag :value="getStatusLabel(data.status)" :severity="getStatusSeverity(data.status)" />
              </template>
            </Column>
            <Column header="Acciones" style="width: 180px">
              <template #body="{ data }">
                <div class="flex gap-2">
                  <Button icon="pi pi-eye" text @click="navigateTo(`/admin/agencias/${data.id}`)" v-tooltip.top="'Ver detalle'" />
                  <Button icon="pi pi-pencil" text severity="warning" @click="navigateTo(`/admin/agencias/${data.id}/editar`)" v-tooltip.top="'Editar'" />
                  <Button icon="pi pi-trash" text severity="danger" @click="confirmDelete(data)" v-tooltip.top="'Eliminar'" />
                </div>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>

    <!-- Delete dialog -->
    <Dialog v-model:visible="showDeleteDialog" header="Confirmar" :modal="true" :style="{ width: '420px' }">
      <div class="flex items-start gap-3">
        <i class="pi pi-exclamation-triangle text-orange-500 text-2xl"></i>
        <div>
          <p class="font-semibold">Eliminar agencia</p>
          <p class="text-sm text-gray-600">Esta accion desactivara la agencia y la ocultara del publico.</p>
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
import { ref, watch, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
if (!authStore.isAdmin) {
  navigateTo('/dashboard')
}

const toast = useToast()
const { getAgencias, deleteAgencia, getStats, getDepartamentos } = useAgencias()

const agencias = ref<any[]>([])
const departamentos = ref<any[]>([])
const stats = ref<Record<string, any>>({})
const loading = ref(false)
const deleting = ref(false)
const showDeleteDialog = ref(false)
const selectedAgencia = ref<any>(null)

const pagination = ref({
  page: 1,
  limit: 12,
  total: 0,
  total_pages: 0
})

const filters = ref({
  search: '',
  departamento_id: '',
  status: '',
  licencia_turistica: ''
})

const viewMode = ref<'cards' | 'table'>('cards')
const viewOptions = [
  { label: 'Cards', value: 'cards', icon: 'pi pi-th-large' },
  { label: 'Tabla', value: 'table', icon: 'pi pi-list' }
]

const statusOptions = [
  { label: 'Activa', value: 'activa' },
  { label: 'En revision', value: 'en_revision' },
  { label: 'Inactiva', value: 'inactiva' },
  { label: 'Suspendida', value: 'suspendida' }
]

const licenciaOptions = [
  { label: 'Con licencia', value: 'true' },
  { label: 'Sin licencia', value: 'false' }
]

const loadAgencias = async () => {
  loading.value = true
  try {
    const response: any = await getAgencias({
      page: pagination.value.page,
      limit: pagination.value.limit,
      ...filters.value
    })
    if (response.success) {
      agencias.value = response.data.agencias
      pagination.value = {
        ...pagination.value,
        ...response.data.pagination
      }
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo cargar agencias',
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
    // silent
  }
}

const loadDepartamentos = async () => {
  try {
    const response: any = await getDepartamentos()
    departamentos.value = response?.data || []
  } catch (error) {
    departamentos.value = []
  }
}

watch(filters, () => {
  pagination.value.page = 1
  loadAgencias()
}, { deep: true })

const onPage = (event: any) => {
  pagination.value.page = event.page + 1
  loadAgencias()
}

const confirmDelete = (agencia: any) => {
  selectedAgencia.value = agencia
  showDeleteDialog.value = true
}

const handleDelete = async () => {
  if (!selectedAgencia.value) return
  deleting.value = true
  try {
    const response: any = await deleteAgencia(selectedAgencia.value.id)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Eliminada', detail: 'Agencia desactivada', life: 2500 })
      showDeleteDialog.value = false
      loadAgencias()
      loadStats()
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo eliminar', life: 3000 })
  } finally {
    deleting.value = false
  }
}

const getStatusLabel = (status?: string) => {
  switch (status) {
    case 'activa': return 'Activa'
    case 'en_revision': return 'En revision'
    case 'inactiva': return 'Inactiva'
    case 'suspendida': return 'Suspendida'
    default: return 'N/D'
  }
}

const getStatusSeverity = (status?: string) => {
  switch (status) {
    case 'activa': return 'success'
    case 'en_revision': return 'warning'
    case 'inactiva': return 'secondary'
    case 'suspendida': return 'danger'
    default: return 'info'
  }
}

onMounted(() => {
  loadDepartamentos()
  loadStats()
  loadAgencias()
})
</script>
