<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Comprar paquete</h1>
          <p class="muted mt-1">Completa los datos para registrar tu compra.</p>
        </div>
        <div class="flex gap-2">
          <Button
            label="Volver al detalle"
            icon="pi pi-arrow-left"
            severity="secondary"
            outlined
            :disabled="loading"
            @click="navigateTo(`/turista/paquetes/${paqueteId}`)"
          />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8">
      <div v-if="loading">
        <Skeleton height="280px" />
      </div>

      <div v-else-if="error" class="text-center space-y-4 py-12">
        <i class="pi pi-exclamation-triangle text-5xl text-orange-500"></i>
        <h2 class="text-2xl font-bold text-gray-900">No se pudo cargar el paquete</h2>
        <p class="muted">{{ error }}</p>
        <div class="flex justify-center gap-2">
          <Button label="Volver" icon="pi pi-arrow-left" severity="secondary" outlined @click="navigateTo('/turista/paquetes')" />
          <Button label="Reintentar" icon="pi pi-refresh" @click="loadPaquete" />
        </div>
      </div>

      <div v-else-if="paquete">
        <div class="mb-6">
          <Card class="surface-card">
            <template #content>
              <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                <div>
                  <p class="text-xs uppercase tracking-widest text-gray-500">Paquete</p>
                  <h2 class="text-xl font-bold text-gray-900">{{ paquete.nombre }}</h2>
                  <p v-if="paquete.agencia?.nombre_comercial" class="text-sm text-gray-600 mt-1">
                    <i class="pi pi-building mr-2"></i>{{ paquete.agencia.nombre_comercial }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-xs text-gray-500">Desde</p>
                  <p class="text-2xl font-bold text-emerald-700">Bs. {{ formatMoney(paquete.precio_base_nacionales) }}</p>
                  <p v-if="Number(paquete.precio_adicional_extranjeros || 0) > 0" class="text-xs text-gray-500">
                    + Bs. {{ formatMoney(paquete.precio_adicional_extranjeros) }} extranjeros
                  </p>
                </div>
              </div>
            </template>
          </Card>
        </div>

        <CompraForm :paquete="paquete" />
      </div>
    </div>
    <Toast />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import CompraForm from '~/components/paquetes/CompraForm.vue'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const route = useRoute()
const { getPaquete } = usePaquetesTuristicos()

const paqueteId = Number(route.params.id)

const paquete = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const formatMoney = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const loadPaquete = async () => {
  loading.value = true
  error.value = null
  try {
    const response: any = await getPaquete(paqueteId)
    if (response.success) {
      paquete.value = response.data
      return
    }
    error.value = response?.error?.message || 'No se pudo cargar el paquete'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudo cargar el paquete'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadPaquete)
</script>
