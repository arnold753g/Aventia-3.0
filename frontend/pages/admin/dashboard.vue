<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <h1 class="text-3xl font-bold" style="color: var(--color-primary);">Panel de administracion</h1>
        <p class="muted mt-1">Resumen rapido y accesos a modulos</p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 grid gap-6 md:grid-cols-2">
      <Card class="surface-card">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-users" style="color: var(--color-primary);"></i>
            <span>Usuarios</span>
          </div>
        </template>
        <template #content>
          <p class="muted mb-4">Gestiona altas, ediciones, roles y estados de usuario.</p>
          <Button
            label="Ir a usuarios"
            icon="pi pi-arrow-right"
            @click="navigateTo('/admin/usuarios')"
          />
        </template>
      </Card>

      <Card class="surface-card">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-users" style="color: var(--color-primary);"></i>
            <span>Usuarios</span>
          </div>
        </template>
        <template #content>
          <p class="muted mb-4">Gestiona altas, ediciones, roles y estados de usuario.</p>
          <Button
              label="Ir a usuarios"
              icon="pi pi-arrow-right"
              @click="navigateTo('/admin/usuarios')"
          />
        </template>
      </Card>


      <Card class="surface-card">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-home" style="color: var(--color-accent);"></i>
            <span>Dashboard general</span>
          </div>
        </template>
        <template #content>
          <p class="muted mb-4">Volver al dashboard general de la aplicacion.</p>
          <Button
            label="Ir al dashboard"
            icon="pi pi-home"
            severity="secondary"
            @click="navigateTo('/dashboard')"
          />
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})
</script>
