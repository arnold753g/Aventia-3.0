<template>
  <div class="page-shell">
    <nav class="bg-white border-b border-gray-200">
      <div class="container mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold" style="color: var(--color-primary);">
            Sistema Andaria
          </h1>
          <div class="flex items-center gap-4">
            <span class="muted">
              {{ authStore.fullName }}
              <span class="text-xs" style="color: var(--color-muted);">({{ authStore.user?.rol }})</span>
            </span>
            <Button
              label="Cerrar sesion"
              icon="pi pi-sign-out"
              @click="handleLogout"
            />
          </div>
        </div>
      </div>
    </nav>

    <div class="container mx-auto px-4 py-8">
      <h2 class="text-3xl font-bold mb-8" style="color: var(--color-primary);">
        Dashboard
      </h2>

      <Card class="mb-8 surface-card">
        <template #content>
          <div class="text-center py-8">
            <h3 class="text-2xl font-bold mb-3" style="color: var(--color-primary);">
              Bienvenido {{ authStore.user?.nombre }}!
            </h3>
            <p class="muted">
              Has iniciado sesion exitosamente en el sistema
            </p>
            <div class="mt-4 text-sm muted">
              Email: {{ authStore.user?.email }}
            </div>
          </div>
        </template>
      </Card>

      <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-users" style="color: var(--color-primary);"></i>
              <span>Usuarios</span>
            </div>
          </template>
          <template #content>
            <div class="text-3xl font-bold" style="color: var(--color-primary);">7</div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-building" style="color: var(--color-accent);"></i>
              <span>Agencias</span>
            </div>
          </template>
          <template #content>
            <div class="text-3xl font-bold" style="color: var(--color-accent);">3</div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-box" style="color: var(--color-contrast);"></i>
              <span>Paquetes</span>
            </div>
          </template>
          <template #content>
            <div class="text-3xl font-bold" style="color: var(--color-contrast);">4</div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-shopping-cart" style="color: var(--color-warm);"></i>
              <span>Compras</span>
            </div>
          </template>
          <template #content>
            <div class="text-3xl font-bold" style="color: var(--color-warm);">1</div>
          </template>
        </Card>
      </div>

      <div class="text-center space-x-4">
        <Button
          label="Volver al inicio"
          icon="pi pi-home"
          @click="navigateTo('/')"
        />
        <Button
          label="Ver perfil"
          icon="pi pi-user"
          text
          @click="getProfile"
        />
      </div>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'

const toast = useToast()
const authStore = useAuthStore()

if (!authStore.isAuthenticated) {
  navigateTo('/login')
}

const handleLogout = () => {
  toast.add({
    severity: 'info',
    summary: 'Sesion cerrada',
    detail: 'Has cerrado sesion correctamente',
    life: 3000
  })

  authStore.logout()
}

const getProfile = async () => {
  const success = await authStore.getProfile()
  if (success) {
    toast.add({
      severity: 'success',
      summary: 'Perfil actualizado',
      detail: 'Informacion del perfil cargada correctamente',
      life: 3000
    })
  } else {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'No se pudo cargar el perfil',
      life: 3000
    })
  }
}
</script>
