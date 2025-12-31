<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-2 md:flex-row md:items-end md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Reportes</h1>
          <p class="muted mt-1">Genera reportes de ventas, ocupacion, finanzas y turistas.</p>
        </div>
        <div class="text-sm text-gray-500">
          Descargas en PDF, CSV o JSON.
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

      <Card class="surface-card">
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div class="space-y-2">
              <label class="text-sm font-semibold text-gray-700">Tipo de reporte</label>
              <SelectButton v-model="reportType" :options="reportOptions" optionLabel="label" optionValue="value" />
            </div>

            <div class="space-y-2">
              <label class="text-sm font-semibold text-gray-700">Formato</label>
              <Dropdown v-model="formato" :options="formatOptions" optionLabel="label" optionValue="value" class="w-full" />
            </div>

            <div class="space-y-2">
              <label class="text-sm font-semibold text-gray-700">Fecha inicio</label>
              <Calendar v-model="fechaInicio" :showIcon="true" dateFormat="yy-mm-dd" class="w-full" />
            </div>

            <div class="space-y-2">
              <label class="text-sm font-semibold text-gray-700">Fecha fin</label>
              <Calendar v-model="fechaFin" :showIcon="true" dateFormat="yy-mm-dd" class="w-full" />
            </div>

            <div v-if="showTipoCompra" class="space-y-2">
              <label class="text-sm font-semibold text-gray-700">Tipo de compra</label>
              <Dropdown v-model="tipoCompra" :options="tipoCompraOptions" optionLabel="label" optionValue="value" class="w-full" />
            </div>
          </div>

          <Divider class="my-4" />

          <div class="flex flex-wrap items-center gap-3 justify-end">
            <Button label="Limpiar" severity="secondary" outlined @click="resetForm" />
            <Button label="Generar reporte" icon="pi pi-file" :loading="loading" @click="generarReporte" />
          </div>
        </template>
      </Card>

      <Card v-if="showJson" class="surface-card">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-database text-green-600"></i>
            <span>Resultado JSON</span>
          </div>
        </template>
        <template #content>
          <pre class="text-xs text-gray-700 whitespace-pre-wrap">{{ prettyJson }}</pre>
        </template>
      </Card>
    </div>
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
const { generarReporteVentas, generarReporteOcupacion, generarReporteFinanciero, generarReporteTuristas } = useReportes()

const reportOptions = [
  { label: 'Ventas', value: 'ventas' },
  { label: 'Ocupacion', value: 'ocupacion' },
  { label: 'Financiero', value: 'financiero' },
  { label: 'Turistas', value: 'turistas' }
]

const formatOptions = [
  { label: 'PDF', value: 'pdf' },
  { label: 'CSV', value: 'csv' },
  { label: 'JSON', value: 'json' }
]

const tipoCompraOptions = [
  { label: 'Todos', value: '' },
  { label: 'Compartido', value: 'compartido' },
  { label: 'Privado', value: 'privado' }
]

const agenciaId = ref<number | null>(null)
const reportType = ref<'ventas' | 'ocupacion' | 'financiero' | 'turistas'>('ventas')
const formato = ref<'pdf' | 'csv' | 'json'>('pdf')
const paqueteId = ref('')
const tipoCompra = ref('')
const error = ref<string | null>(null)
const loading = ref(false)
const reportData = ref<any>(null)

const now = new Date()
const fechaInicio = ref<Date | null>(new Date(now.getFullYear(), now.getMonth(), 1))
const fechaFin = ref<Date | null>(new Date(now.getFullYear(), now.getMonth() + 1, 0))

const showPaqueteFilter = computed(() => reportType.value !== 'turistas')
const showTipoCompra = computed(() => reportType.value === 'ventas')
const showJson = computed(() => formato.value === 'json' && reportData.value)
const prettyJson = computed(() => JSON.stringify(reportData.value, null, 2))

const formatDate = (value: Date | null) => {
  if (!value) return ''
  const yyyy = value.getFullYear()
  const mm = String(value.getMonth() + 1).padStart(2, '0')
  const dd = String(value.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

const downloadBlob = (blob: Blob, filename: string) => {
  if (!process.client) return
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename || `reporte_${reportType.value}.${formato.value}`
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

const resolveAgenciaId = async () => {
  if (agenciaId.value) return agenciaId.value
  const resp: any = await getMiAgencia()
  if (!resp?.success) throw new Error(resp?.error?.message || 'No se pudo cargar la agencia')
  agenciaId.value = Number(resp.data?.id)
  if (!agenciaId.value) throw new Error('No se pudo resolver el ID de la agencia')
  return agenciaId.value
}

const resetForm = () => {
  reportType.value = 'ventas'
  formato.value = 'pdf'
  paqueteId.value = ''
  tipoCompra.value = ''
  reportData.value = null
  error.value = null
  const current = new Date()
  fechaInicio.value = new Date(current.getFullYear(), current.getMonth(), 1)
  fechaFin.value = new Date(current.getFullYear(), current.getMonth() + 1, 0)
}

const generarReporte = async () => {
  loading.value = true
  error.value = null
  reportData.value = null
  try {
    const id = await resolveAgenciaId()
    const params: any = {
      fechaInicio: formatDate(fechaInicio.value),
      fechaFin: formatDate(fechaFin.value),
      formato: formato.value
    }

    if (showPaqueteFilter.value && paqueteId.value) {
      const parsed = Number(paqueteId.value)
      if (!Number.isFinite(parsed)) {
        throw new Error('Paquete ID invalido')
      }
      params.paqueteId = parsed
    }
    if (showTipoCompra.value && tipoCompra.value) {
      params.tipoCompra = tipoCompra.value as 'compartido' | 'privado'
    }

    let response: any
    if (reportType.value === 'ventas') {
      response = await generarReporteVentas(id, params)
    } else if (reportType.value === 'ocupacion') {
      response = await generarReporteOcupacion(id, params)
    } else if (reportType.value === 'financiero') {
      response = await generarReporteFinanciero(id, params)
    } else {
      response = await generarReporteTuristas(id, params)
    }

    if (formato.value === 'json') {
      reportData.value = response?.data ?? response
      return
    }

    const filename = response?.filename || ''
    if (response?.blob) {
      downloadBlob(response.blob, filename)
      toast.add({ severity: 'success', summary: 'Reporte generado', detail: 'Archivo descargado.', life: 2500 })
      return
    }

    toast.add({ severity: 'warn', summary: 'Sin archivo', detail: 'No se recibio archivo para descargar.', life: 2500 })
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudo generar el reporte'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  resolveAgenciaId().catch(() => null)
})
</script>
