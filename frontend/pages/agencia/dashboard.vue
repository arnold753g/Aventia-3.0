<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
          Panel de agencia
        </h1>
        <p class="muted mt-1">Gestiona tu agencia y su informacion.</p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 grid gap-6 md:grid-cols-2">
      <Card class="surface-card md:col-span-2">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-building" style="color: var(--color-primary);"></i>
            <span>Mi agencia</span>
          </div>
        </template>
        <template #content>
          <div v-if="loading">
            <Skeleton height="140px" />
          </div>

          <div v-else-if="error" class="space-y-3">
            <Message severity="warn" :closable="false">{{ error }}</Message>
            <Button label="Reintentar" icon="pi pi-refresh" outlined @click="loadAgencia" />
          </div>

          <div v-else class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <div class="text-xl font-semibold text-gray-900">{{ agencia?.nombre_comercial }}</div>
              <div class="text-sm text-gray-600">{{ agencia?.departamento?.nombre || 'Departamento N/D' }}</div>
              <div class="text-sm text-gray-600">
                Fotos: {{ agencia?.fotos?.length || 0 }} | Especialidades: {{ agencia?.especialidades?.length || 0 }}
              </div>
            </div>
            <div class="flex gap-2">
              <Button label="Ver" icon="pi pi-eye" outlined @click="navigateTo('/agencia/mi-agencia')" />
              <Button label="Editar" icon="pi pi-pencil" severity="warning" @click="navigateTo('/agencia/mi-agencia/editar')" />
            </div>
          </div>
        </template>
      </Card>

      <Card class="surface-card md:col-span-2">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-briefcase" style="color: var(--color-primary);"></i>
            <span>Paquetes turísticos</span>
          </div>
        </template>
        <template #content>
          <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <div class="text-sm text-gray-600">
                Crea y gestiona tus paquetes (salida diaria o salida única) y revisa sus salidas habilitadas.
              </div>
            </div>
            <div class="flex gap-2">
              <Button label="Ver paquetes" icon="pi pi-eye" outlined @click="navigateTo('/agencia/paquetes')" />
              <Button label="Nuevo paquete" icon="pi pi-plus" @click="navigateTo('/agencia/paquetes/nuevo')" />
            </div>
          </div>
        </template>
      </Card>

      <Card class="surface-card md:col-span-2">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-shopping-bag" style="color: var(--color-primary);"></i>
            <span>Ventas de paquetes</span>
          </div>
        </template>
        <template #content>
          <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <div class="text-sm text-gray-600">
                Revisa los pagos registrados por turistas y confirma o rechaza cada pago.
              </div>
            </div>
            <div class="flex gap-2">
              <Button label="Ver ventas" icon="pi pi-eye" outlined @click="navigateTo('/agencia/ventas')" />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const { getMiAgencia } = useAgencias()

const agencia = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const loadAgencia = async () => {
  loading.value = true
  error.value = null
  try {
    const response: any = await getMiAgencia()
    if (response.success) {
      agencia.value = response.data
      return
    }
    error.value = response?.error?.message || 'No se pudo cargar la agencia'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudo cargar la agencia'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAgencia()
})
</script>
