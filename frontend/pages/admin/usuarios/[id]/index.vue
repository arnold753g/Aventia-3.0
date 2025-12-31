<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-5xl mx-auto px-4 py-6">
        <div class="flex items-center gap-4">
          <Button
            icon="pi pi-arrow-left"
            text
            rounded
            @click="navigateTo('/admin/usuarios')"
          />
          <div class="flex-1">
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Detalles del Usuario
            </h1>
            <p class="muted mt-1">Información completa del usuario</p>
          </div>
          <Button
            label="Editar"
            icon="pi pi-pencil"
            severity="warning"
            @click="navigateTo(`/admin/usuarios/${route.params.id}/editar`)"
          />
        </div>
      </div>
    </div>

    <div v-if="loading" class="max-w-5xl mx-auto px-4 py-8">
      <Skeleton height="200px" class="mb-4" />
      <Skeleton height="400px" />
    </div>

    <div v-else-if="usuario" class="max-w-5xl mx-auto px-4 py-8 space-y-6">
      <!-- Perfil -->
      <Card>
        <template #content>
          <div class="flex items-start gap-6">
            <UserAvatar
              :src="resolvePhoto(usuario.profile_photo)"
              :nombre="usuario.nombre"
              :apellido="usuario.apellido_paterno"
              :rol="usuario.rol"
              :status="usuario.status"
              size="xl"
              :showStatus="true"
              :showRolIcon="true"
            />
            <div class="flex-1">
              <h2 class="text-2xl font-bold mb-2">
                {{ getFullName(usuario.nombre, usuario.apellido_paterno, usuario.apellido_materno) }}
              </h2>
              <div class="flex gap-2 mb-4">
                <StatusBadge type="rol" :value="resolveRol(usuario)" :label="getRolLabel(resolveRol(usuario))" :icon="getRolIcon(resolveRol(usuario))" />
                <StatusBadge type="status" :value="resolveStatus(usuario)" :label="getStatusLabel(resolveStatus(usuario))" />
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <p class="text-sm muted">Email</p>
                  <p class="font-medium">{{ usuario.email }}</p>
                </div>
                <div>
                  <p class="text-sm muted">Teléfono</p>
                  <p class="font-medium">{{ formatPhone(usuario.phone) }}</p>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Información Personal -->
      <Card>
        <template #header>
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold flex items-center gap-2">
              <i class="pi pi-user"></i>
              Información Personal
            </h3>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <p class="text-sm muted mb-1">CI</p>
              <p class="font-medium text-lg">{{ usuario.ci }}-{{ usuario.expedido }}</p>
            </div>
            <div>
              <p class="text-sm muted mb-1">Fecha de Nacimiento</p>
              <p class="font-medium text-lg">{{ formatDate(usuario.fecha_nacimiento) }}</p>
            </div>
            <div>
              <p class="text-sm muted mb-1">Ciudad</p>
              <p class="font-medium text-lg">{{ usuario.ciudad || 'No especificado' }}</p>
            </div>
            <div>
              <p class="text-sm muted mb-1">Nacionalidad</p>
              <p class="font-medium text-lg">{{ usuario.nationality || 'Bolivia' }}</p>
            </div>
          </div>
        </template>
      </Card>

      <!-- Información del Sistema -->
      <Card>
        <template #header>
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold flex items-center gap-2">
              <i class="pi pi-cog"></i>
              Información del Sistema
            </h3>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <p class="text-sm muted mb-1">Fecha de Registro</p>
              <p class="font-medium text-lg">{{ formatDateTime(usuario.created_at) }}</p>
            </div>
            <div>
              <p class="text-sm muted mb-1">Última Actualización</p>
              <p class="font-medium text-lg">{{ formatDateTime(usuario.updated_at) }}</p>
            </div>
            <div>
              <p class="text-sm muted mb-1">ID del Usuario</p>
              <p class="font-mono font-medium text-lg">#{{ usuario.id }}</p>
            </div>
            <div>
              <p class="text-sm muted mb-1">Estado de Verificación</p>
              <Tag :value="usuario.email_verified ? 'Verificado' : 'Pendiente'" :severity="usuario.email_verified ? 'success' : 'warning'" />
            </div>
          </div>
        </template>
      </Card>

      <!-- Acciones Administrativas -->
      <Card>
        <template #header>
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold flex items-center gap-2">
              <i class="pi pi-shield"></i>
              Acciones Administrativas
            </h3>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h4 class="font-semibold mb-2">Cambiar Rol</h4>
              <div class="flex gap-2">
                <Select v-model="nuevoRol" :options="roleOptions" optionLabel="label" optionValue="value" placeholder="Seleccionar rol" class="flex-1" />
                <Button label="Actualizar" icon="pi pi-refresh" :loading="updatingRol" @click="handleUpdateRol" />
              </div>
            </div>
            <div>
              <h4 class="font-semibold mb-2">Cambiar Estado</h4>
              <div class="flex gap-2">
                <Select v-model="nuevoStatus" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="Seleccionar estado" class="flex-1" />
                <Button label="Actualizar" icon="pi pi-check" severity="success" :loading="updatingStatus" @click="handleUpdateStatus" />
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { getFullName, getRolLabel, getStatusLabel, getRolIcon, formatPhone, formatDate, formatDateTime } from '~/utils/formatters'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const route = useRoute()
const toast = useToast()
const { getUsuario, updateRol, updateStatus } = useUsuarios()
const apiOrigin = new URL(useRuntimeConfig().public.apiBase).origin

const loading = ref(false)
const usuario = ref<any>(null)
const nuevoRol = ref('')
const nuevoStatus = ref('')
const updatingRol = ref(false)
const updatingStatus = ref(false)

const roleOptions = [
  { label: 'Turista', value: 'turista' },
  { label: 'Encargado de Agencia', value: 'encargado_agencia' },
  { label: 'Administrador', value: 'admin' }
]

const statusOptions = [
  { label: 'Activo', value: 'active' },
  { label: 'Inactivo', value: 'inactive' },
  { label: 'Suspendido', value: 'suspended' }
]

const resolveRol = (usuario: any) => {
  const raw =
    usuario?.rol ??
    usuario?.role ??
    usuario?.roles ??
    usuario?.user_role ??
    usuario?.rol_name ??
    usuario?.role_name ??
    usuario?.roleSlug ??
    usuario?.role_slug ??
    usuario?.roleName ??
    usuario?.type ??
    usuario?.user_type

  if (typeof raw === 'string') return raw

  if (Array.isArray(raw) && raw.length) {
    const first = raw[0]
    if (typeof first === 'string') return first
    if (first && typeof first === 'object') return first.slug ?? first.name ?? first.value ?? first.tipo ?? ''
  }

  if (raw && typeof raw === 'object') {
    return raw.slug ?? raw.name ?? raw.value ?? raw.tipo ?? raw.label ?? ''
  }

  return ''
}

const resolveStatus = (usuario: any) => {
  const raw =
    usuario?.status ??
    usuario?.estado ??
    usuario?.state ??
    usuario?.estatus ??
    usuario?.active ??
    usuario?.is_active ??
    usuario?.isActive ??
    usuario?.enabled

  if (typeof raw === 'string') return raw
  if (typeof raw === 'boolean') return raw ? 'active' : 'inactive'

  if (raw && typeof raw === 'object') {
    return raw.slug ?? raw.name ?? raw.value ?? raw.estado ?? raw.label ?? ''
  }

  return ''
}

const loadUsuario = async () => {
  loading.value = true
  try {
    const response: any = await getUsuario(Number(route.params.id))
    if (response.success) {
      usuario.value = response.data
      nuevoRol.value = response.data.rol
      nuevoStatus.value = response.data.status
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'Error al cargar usuario',
      life: 3000
    })
    navigateTo('/admin/usuarios')
  } finally {
    loading.value = false
  }
}

const handleUpdateRol = async () => {
  if (nuevoRol.value === usuario.value.rol) {
    toast.add({
      severity: 'warn',
      summary: 'Sin cambios',
      detail: 'El rol seleccionado es el mismo que el actual',
      life: 3000
    })
    return
  }

  updatingRol.value = true
  try {
    const response: any = await updateRol(Number(route.params.id), nuevoRol.value)
    if (response.success) {
      toast.add({
        severity: 'success',
        summary: 'Rol actualizado',
        detail: 'El rol del usuario ha sido actualizado correctamente',
        life: 3000
      })
      loadUsuario()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'Error al actualizar rol',
      life: 3000
    })
  } finally {
    updatingRol.value = false
  }
}

const handleUpdateStatus = async () => {
  if (nuevoStatus.value === usuario.value.status) {
    toast.add({
      severity: 'warn',
      summary: 'Sin cambios',
      detail: 'El estado seleccionado es el mismo que el actual',
      life: 3000
    })
    return
  }

  updatingStatus.value = true
  try {
    const response: any = await updateStatus(Number(route.params.id), nuevoStatus.value)
    if (response.success) {
      toast.add({
        severity: 'success',
        summary: 'Estado actualizado',
        detail: 'El estado del usuario ha sido actualizado correctamente',
        life: 3000
      })
      loadUsuario()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'Error al actualizar estado',
      life: 3000
    })
  } finally {
    updatingStatus.value = false
  }
}

const resolvePhoto = (path?: string) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${apiOrigin}/${path.replace(/^\//, '')}`
}

onMounted(() => {
  loadUsuario()
})
</script>
