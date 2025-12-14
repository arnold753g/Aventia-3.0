<template>
  <div class="page-shell">
    <!-- Header -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Gestión de Usuarios
            </h1>
            <p class="muted mt-1">Administra todos los usuarios del sistema</p>
          </div>
          <Button
            label="Nuevo Usuario"
            icon="pi pi-plus"
            @click="navigateTo('/admin/usuarios/crear')"
          />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card class="bg-gradient-to-br from-purple-50 to-purple-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-purple-600 font-medium">Total Usuarios</p>
                <p class="text-3xl font-bold text-purple-700">{{ stats.total || 0 }}</p>
              </div>
              <div class="w-12 h-12 rounded-full bg-purple-200 flex items-center justify-center">
                <i class="pi pi-users text-2xl text-purple-700"></i>
              </div>
            </div>
          </template>
        </Card>

        <Card class="bg-gradient-to-br from-green-50 to-green-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-green-600 font-medium">Activos</p>
                <p class="text-3xl font-bold text-green-700">{{ stats.active || 0 }}</p>
              </div>
              <div class="w-12 h-12 rounded-full bg-green-200 flex items-center justify-center">
                <i class="pi pi-check-circle text-2xl text-green-700"></i>
              </div>
            </div>
          </template>
        </Card>

        <Card class="bg-gradient-to-br from-blue-50 to-blue-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-blue-600 font-medium">Turistas</p>
                <p class="text-3xl font-bold text-blue-700">{{ stats.turistas || 0 }}</p>
              </div>
              <div class="w-12 h-12 rounded-full bg-blue-200 flex items-center justify-center">
                <i class="pi pi-user text-2xl text-blue-700"></i>
              </div>
            </div>
          </template>
        </Card>

        <Card class="bg-gradient-to-br from-orange-50 to-orange-100">
          <template #content>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-orange-600 font-medium">Administradores</p>
                <p class="text-3xl font-bold text-orange-700">{{ stats.admins || 0 }}</p>
              </div>
              <div class="w-12 h-12 rounded-full bg-orange-200 flex items-center justify-center">
                <i class="pi pi-shield text-2xl text-orange-700"></i>
              </div>
            </div>
          </template>
        </Card>
      </div>

      <!-- Filters -->
      <Card class="surface-card">
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <!-- Búsqueda con debounce -->
            <div class="md:col-span-2">
              <IconField>
                <InputIcon class="pi pi-search" />
                <InputText
                  v-model="searchQuery"
                  placeholder="Buscar por nombre, email o CI..."
                  class="w-full"
                />
              </IconField>
            </div>

            <!-- Filtro por rol -->
            <Select
              v-model="filters.rol"
              :options="roleOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Filtrar por rol"
              class="w-full"
              showClear
            />

            <!-- Filtro por estado -->
            <Select
              v-model="filters.status"
              :options="statusOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Filtrar por estado"
              class="w-full"
              showClear
            />
          </div>

          <!-- Chips de filtros activos -->
          <div v-if="hasActiveFilters" class="flex gap-2 mt-4">
            <Chip
              v-if="searchQuery"
              :label="`Búsqueda: ${searchQuery}`"
              removable
              @remove="searchQuery = ''"
            />
            <Chip
              v-if="filters.rol"
              :label="`Rol: ${getRolLabel(filters.rol)}`"
              removable
              @remove="filters.rol = ''"
            />
            <Chip
              v-if="filters.status"
              :label="`Estado: ${getStatusLabel(filters.status)}`"
              removable
              @remove="filters.status = ''"
            />
          </div>
        </template>
      </Card>

      <!-- Tabla de usuarios -->
      <Card class="surface-card">
        <template #content>
          <!-- Skeleton Loader -->
          <div v-if="loading" class="space-y-4">
            <Skeleton height="3rem" class="mb-2" v-for="i in 5" :key="i" />
          </div>

          <!-- Tabla de datos -->
          <DataTable
            v-else
            :value="usuarios"
            stripedRows
            paginator
            :rows="pagination.limit"
            :totalRecords="pagination.total"
            :lazy="true"
            dataKey="id"
            :first="(pagination.page - 1) * pagination.limit"
            @page="onPage"
            :rowHover="true"
            class="p-datatable-sm"
          >
            <template #empty>
              <div class="text-center py-12">
                <i class="pi pi-users text-6xl muted mb-4 block"></i>
                <p class="text-xl font-semibold muted">No se encontraron usuarios</p>
                <p class="text-sm muted mt-2">Intenta ajustar los filtros de búsqueda</p>
              </div>
            </template>

            <Column field="id" header="ID" :sortable="true" style="width: 80px">
              <template #body="{ data }">
                <span class="font-mono text-sm">#{{ data.id }}</span>
              </template>
            </Column>

            <Column header="Usuario" :sortable="true" sortField="nombre">
              <template #body="{ data }">
                <div class="flex items-center gap-3">
                  <UserAvatar
                    :src="resolvePhoto(data.profile_photo)"
                    :nombre="data.nombre"
                    :apellido="data.apellido_paterno"
                    :rol="data.rol"
                    :status="data.status"
                    size="sm"
                    :showStatus="true"
                  />
                  <div>
                    <p class="font-semibold" style="color: var(--color-text);">
                      {{ data.nombre }} {{ data.apellido_paterno }}
                    </p>
                    <p class="text-sm muted">{{ data.email }}</p>
                  </div>
                </div>
              </template>
            </Column>

            <Column field="ci" header="CI / Teléfono">
              <template #body="{ data }">
                <div>
                  <p class="font-mono text-sm">{{ data.ci }}-{{ data.expedido }}</p>
                  <p class="text-sm muted">{{ formatPhone(data.phone) }}</p>
                </div>
              </template>
            </Column>

            <Column header="Rol" field="rol" :sortable="true">
              <template #body="{ data }">
                <!-- Debug: {{ data.rol }} - {{ getRolLabel(data.rol) }} -->
                <UsuariosStatusBadge
                  type="rol"
                  :value="data.rol || ''"
                  :label="getRolLabel(data.rol || '')"
                  :icon="getRolIcon(data.rol || '')"
                />
              </template>
            </Column>

            <Column header="Estado" field="status" :sortable="true">
              <template #body="{ data }">
                <!-- Debug: {{ data.status }} - {{ getStatusLabel(data.status) }} -->
                <UsuariosStatusBadge
                  type="status"
                  :value="data.status || ''"
                  :label="getStatusLabel(data.status || '')"
                />
              </template>
            </Column>

            <Column header="Registro">
              <template #body="{ data }">
                <div class="text-sm">
                  <p>{{ formatDateShort(data.created_at) }}</p>
                  <p class="text-xs muted">{{ formatRelativeTime(data.created_at) }}</p>
                </div>
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
                    v-tooltip.top="'Ver detalles'"
                    @click="navigateTo(`/admin/usuarios/${data.id}`)"
                  />
                  <Button
                    icon="pi pi-pencil"
                    severity="warning"
                    text
                    rounded
                    v-tooltip.top="'Editar'"
                    @click="navigateTo(`/admin/usuarios/${data.id}/editar`)"
                  />
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    rounded
                    v-tooltip.top="'Desactivar'"
                    @click="confirmDeactivate(data)"
                    :disabled="data.id === authStore.user?.id"
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
      v-model:visible="showDeactivateDialog"
      header="Confirmar desactivación"
      :modal="true"
      :style="{ width: '450px' }"
    >
      <div class="flex items-start gap-4">
        <i class="pi pi-exclamation-triangle text-orange-500 text-3xl"></i>
        <div>
          <p class="mb-2">¿Está seguro de desactivar este usuario?</p>
          <p class="text-sm muted">
            <strong>{{ selectedUsuario?.nombre }} {{ selectedUsuario?.apellido_paterno }}</strong>
          </p>
          <p class="text-sm muted mt-2">
            El usuario no podrá iniciar sesión hasta que sea reactivado.
          </p>
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancelar"
          severity="secondary"
          @click="showDeactivateDialog = false"
        />
        <Button
          label="Desactivar"
          severity="danger"
          @click="handleDeactivate"
          :loading="deactivating"
        />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'
import { getRolLabel, getStatusLabel, getRolIcon, formatPhone, formatDateShort } from '~/utils/formatters'
import UserAvatar from '~/components/usuarios/UserAvatar.vue'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const toast = useToast()
const authStore = useAuthStore()
const { getUsuarios, deactivateUsuario, getStats } = useUsuarios()
const apiOrigin = new URL(useRuntimeConfig().public.apiBase).origin

// Estado
const usuarios = ref<any[]>([])
const stats = ref<Record<string, any>>({})
const loading = ref(false)
const deactivating = ref(false)
const showDeactivateDialog = ref(false)
const selectedUsuario = ref<any>(null)
const searchQuery = ref('')

const pagination = ref({
  page: 1,
  limit: 10,
  total: 0,
  total_pages: 0
})

const filters = ref({
  rol: '',
  status: ''
})

// Opciones de filtros
const roleOptions = [
  { label: 'Administrador', value: 'admin' },
  { label: 'Turista', value: 'turista' },
  { label: 'Encargado de Agencia', value: 'encargado_agencia' }
]

const statusOptions = [
  { label: 'Activo', value: 'active' },
  { label: 'Inactivo', value: 'inactive' },
  { label: 'Suspendido', value: 'suspended' }
]

// Computed
const hasActiveFilters = computed(() => {
  return searchQuery.value || filters.value.rol || filters.value.status
})

const resolveRol = (usuario: any) => {
  // Debug: console.log para verificar qué está recibiendo
  console.log('Usuario completo:', usuario)
  console.log('Rol directo:', usuario?.rol)

  // Primero intentar acceder directamente al campo 'rol'
  if (usuario?.rol && typeof usuario.rol === 'string') {
    return usuario.rol
  }

  // Intentar alternativas comunes
  const raw =
    usuario?.role ??
    usuario?.roles ??
    usuario?.user_role ??
    usuario?.type

  if (typeof raw === 'string') return raw

  if (Array.isArray(raw) && raw.length) {
    const first = raw[0]
    if (typeof first === 'string') return first
    if (first && typeof first === 'object') return first.slug ?? first.name ?? first.value ?? ''
  }

  if (raw && typeof raw === 'object') {
    return raw.slug ?? raw.name ?? raw.value ?? raw.label ?? ''
  }

  console.warn('No se pudo resolver el rol para:', usuario)
  return ''
}

const resolveStatus = (usuario: any) => {
  // Debug: console.log para verificar qué está recibiendo
  console.log('Status directo:', usuario?.status)

  // Primero intentar acceder directamente al campo 'status'
  if (usuario?.status && typeof usuario.status === 'string') {
    return usuario.status
  }

  // Intentar alternativas comunes
  const raw =
    usuario?.estado ??
    usuario?.state ??
    usuario?.active ??
    usuario?.is_active

  if (typeof raw === 'string') return raw
  if (typeof raw === 'boolean') return raw ? 'active' : 'inactive'

  if (raw && typeof raw === 'object') {
    return raw.slug ?? raw.name ?? raw.value ?? raw.label ?? ''
  }

  console.warn('No se pudo resolver el status para:', usuario)
  return ''
}

// Búsqueda con debounce
let searchTimeout: NodeJS.Timeout
watch(searchQuery, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    loadUsuarios()
  }, 500)
})

// Filtros reactivos
watch(filters, () => {
  pagination.value.page = 1
  loadUsuarios()
}, { deep: true })

// Verificar permisos
if (!authStore.isAdmin) {
  navigateTo('/dashboard')
}

const loadUsuarios = async () => {
  loading.value = true
  try {
    const response: any = await getUsuarios({
      page: pagination.value.page,
      limit: pagination.value.limit,
      search: searchQuery.value,
      ...filters.value
    })

    if (response.success) {
      usuarios.value = response.data.usuarios
      pagination.value = {
        ...pagination.value,
        ...response.data.pagination
      }
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'Error al cargar usuarios',
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

const onPage = (event: any) => {
  pagination.value.page = event.page + 1
  loadUsuarios()
}

const confirmDeactivate = (usuario: any) => {
  selectedUsuario.value = usuario
  showDeactivateDialog.value = true
}

const handleDeactivate = async () => {
  if (!selectedUsuario.value) return

  deactivating.value = true
  try {
    const response: any = await deactivateUsuario(selectedUsuario.value.id)

    if (response.success) {
      toast.add({
        severity: 'success',
        summary: 'Usuario desactivado',
        detail: 'El usuario ha sido desactivado exitosamente',
        life: 3000
      })

      showDeactivateDialog.value = false
      loadUsuarios()
      loadStats()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'Error al desactivar usuario',
      life: 3000
    })
  } finally {
    deactivating.value = false
  }
}

const resolvePhoto = (path?: string) => {
  if (!path) {
    //console.log('resolvePhoto: No path provided')
    return ''
  }
  if (path.startsWith('http')) {
    //console.log('resolvePhoto: Full URL:', path)
    return path
  }

  const resolvedUrl = `${apiOrigin}/${path.replace(/^\//, '')}`
  //console.log('resolvePhoto: Constructed URL:', resolvedUrl, 'from path:', path)
  return resolvedUrl
}

const formatRelativeTime = (date: string) => {
  const now = new Date()
  const then = new Date(date)
  const diff = now.getTime() - then.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return 'Hoy'
  if (days === 1) return 'Ayer'
  if (days < 7) return `Hace ${days} días`
  if (days < 30) return `Hace ${Math.floor(days / 7)} semanas`
  if (days < 365) return `Hace ${Math.floor(days / 30)} meses`
  return `Hace ${Math.floor(days / 365)} años`
}

onMounted(() => {
  loadUsuarios()
  loadStats()
})
</script>
