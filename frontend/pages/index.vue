<template>
  <div class="page-shell">
    <div class="container mx-auto px-4 py-12">
      <div class="text-center mb-12">
        <h1 class="text-4xl md:text-5xl font-bold" style="color: var(--color-primary);">
          Sistema Andaria - Gestion Turistica
        </h1>
        <p class="text-lg mt-3 muted">
          Plataforma completa para agencias y turistas
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-6 mb-10">
        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-3">
              <i class="pi pi-database text-2xl" style="color: var(--color-primary)"></i>
              <span>Base de datos</span>
            </div>
          </template>
          <template #content>
            <div v-if="dbStatus.connected" style="color: var(--color-primary);">
              OK. Conectada correctamente
              <p class="text-sm muted mt-2">
                {{ dbStatus.tables }} tablas creadas
              </p>
            </div>
            <div v-else style="color: var(--color-contrast);">
              Error de conexion
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-3">
              <i class="pi pi-server text-2xl" style="color: var(--color-accent)"></i>
              <span>Backend API</span>
            </div>
          </template>
          <template #content>
            <div v-if="apiStatus.online" style="color: var(--color-primary);">
              OK. Servidor en linea
              <p class="text-sm muted mt-2">
                Puerto: {{ apiStatus.port }}
              </p>
            </div>
            <div v-else style="color: var(--color-contrast);">
              Servidor desconectado
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-3">
              <i class="pi pi-desktop text-2xl" style="color: var(--color-warm)"></i>
              <span>Frontend</span>
            </div>
          </template>
          <template #content>
            <div style="color: var(--color-primary);">
              OK. Aplicacion funcionando
              <p class="text-sm muted mt-2">
                Nuxt 3 + Vue 3 + PrimeVue
              </p>
            </div>
          </template>
        </Card>
      </div>

      <div class="text-center">
        <Button
          label="Ir al panel de administracion"
          icon="pi pi-sign-in"
          class="p-button-lg"
          @click="goToLogin"
        />
      </div>

      <div class="mt-14 surface-card p-8">
        <h2 class="text-2xl font-bold mb-6" style="color: var(--color-primary);">
          Estadisticas del sistema
        </h2>
        <div class="grid md:grid-cols-4 gap-6">
          <div class="text-center">
            <div class="text-3xl font-bold" style="color: var(--color-primary);">7</div>
            <div class="muted">Usuarios de prueba</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold" style="color: var(--color-accent);">3</div>
            <div class="muted">Agencias</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold" style="color: var(--color-contrast);">6</div>
            <div class="muted">Atracciones</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold" style="color: var(--color-warm);">4</div>
            <div class="muted">Paquetes</div>
          </div>
        </div>
      </div>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'

const toast = useToast()

const dbStatus = ref({
  connected: false,
  tables: 29
})

const apiStatus = ref({
  online: false,
  port: '5750'
})

const checkHealth = async () => {
  try {
    const response = await fetch('http://localhost:5750/health')
    if (response.ok) {
      apiStatus.value.online = true
      dbStatus.value.connected = true

      toast.add({
        severity: 'success',
        summary: 'Sistema operativo',
        detail: 'Todos los servicios estan funcionando correctamente',
        life: 3000
      })
    }
  } catch (error) {
    apiStatus.value.online = false

    toast.add({
      severity: 'error',
      summary: 'Error de conexion',
      detail: 'No se pudo conectar con el servidor backend',
      life: 5000
    })
  }
}

const goToLogin = () => {
  navigateTo('/login')
}

onMounted(() => {
  checkHealth()
})
</script>
