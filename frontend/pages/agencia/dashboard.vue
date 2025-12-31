<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex flex-col gap-2 md:flex-row md:items-end md:justify-between">
          <div>
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Dashboard de agencia
            </h1>
            <p class="muted mt-1">
              {{ agencia?.nombre_comercial || 'Resumen ejecutivo del mes y estado de ventas.' }}
            </p>
          </div>
          <div class="text-sm text-gray-500">
            Mes actual: {{ monthLabel }}
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Message v-if="dashboardError" severity="error" :closable="false">{{ dashboardError }}</Message>

      <div v-if="loadingDashboard" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <Card v-for="i in 5" :key="`metric-skeleton-${i}`" class="surface-card">
          <template #content>
            <Skeleton height="1rem" width="60%" class="mb-3" />
            <Skeleton height="2rem" width="80%" />
          </template>
        </Card>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <Card class="surface-card">
          <template #content>
            <div class="flex items-start justify-between">
              <div>
                <p class="text-xs uppercase tracking-wide text-gray-500">Ingresos</p>
                <p class="text-2xl font-semibold text-gray-900">{{ formatCurrency(metrics.ingresos_mes) }}</p>
                <p class="text-xs" :class="ingresosTrendClass">
                  {{ formatPercent(comparatives.ingresos_percent) }} vs mes ant
                </p>
              </div>
              <div class="h-10 w-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
                <i class="pi pi-wallet"></i>
              </div>
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #content>
            <div class="flex items-start justify-between">
              <div>
                <p class="text-xs uppercase tracking-wide text-gray-500">Ventas confirmadas</p>
                <p class="text-2xl font-semibold text-gray-900">{{ metrics.ventas_confirmadas }}</p>
                <p class="text-xs" :class="ventasTrendClass">
                  {{ formatDiff(comparatives.ventas_diff) }} vs mes ant
                </p>
              </div>
              <div class="h-10 w-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
                <i class="pi pi-shopping-bag"></i>
              </div>
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #content>
            <div class="flex items-start justify-between">
              <div>
                <p class="text-xs uppercase tracking-wide text-gray-500">Turistas atendidos</p>
                <p class="text-2xl font-semibold text-gray-900">{{ metrics.turistas_atendidos }}</p>
                <p class="text-xs text-gray-500">Promedio: {{ metrics.turistas_promedio.toFixed(1) }}</p>
              </div>
              <div class="h-10 w-10 rounded-xl bg-purple-50 text-purple-600 flex items-center justify-center">
                <i class="pi pi-users"></i>
              </div>
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #content>
            <div class="flex items-start justify-between">
              <div>
                <p class="text-xs uppercase tracking-wide text-gray-500">Paquetes activos</p>
                <p class="text-2xl font-semibold text-gray-900">{{ metrics.paquetes_activos }}</p>
                <p class="text-xs text-gray-500">{{ metrics.paquetes_nuevos }} nuevos</p>
              </div>
              <div class="h-10 w-10 rounded-xl bg-orange-50 text-orange-600 flex items-center justify-center">
                <i class="pi pi-briefcase"></i>
              </div>
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #content>
            <div class="flex items-start justify-between">
              <div>
                <p class="text-xs uppercase tracking-wide text-gray-500">Pendientes de pago</p>
                <p class="text-2xl font-semibold text-gray-900">{{ metrics.pendientes_pago }}</p>
                <p class="text-xs text-gray-500">{{ formatCurrency(metrics.pendientes_pago_monto) }}</p>
              </div>
              <div class="h-10 w-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
                <i class="pi pi-exclamation-triangle"></i>
              </div>
            </div>
          </template>
        </Card>
      </div>

      <div v-if="loadingDashboard" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card v-for="i in 4" :key="`chart-skeleton-${i}`" class="surface-card">
          <template #content>
            <Skeleton height="280px" />
          </template>
        </Card>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card class="surface-card lg:col-span-2">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-chart-line" style="color: var(--color-primary);"></i>
              <span>Ventas por mes (6 meses)</span>
            </div>
          </template>
          <template #content>
            <div v-if="ventasMensuales.length === 0" class="text-sm text-gray-500">
              Sin ventas confirmadas en el periodo.
            </div>
            <div v-else class="h-72">
              <DashboardChart type="line" :data="ventasLineData" :options="ventasLineOptions" />
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-chart-pie" style="color: var(--color-primary);"></i>
              <span>Top 5 paquetes vendidos</span>
            </div>
          </template>
          <template #content>
            <div v-if="topPaquetes.length === 0" class="text-sm text-gray-500">
              Sin ventas confirmadas para mostrar.
            </div>
            <div v-else class="h-64">
              <DashboardChart type="pie" :data="topPaquetesData" :options="topPaquetesOptions" />
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-chart-bar" style="color: var(--color-primary);"></i>
              <span>Ingresos vs proyeccion</span>
            </div>
          </template>
          <template #content>
            <div class="h-64">
              <DashboardChart type="bar" :data="ingresosBarData" :options="ingresosBarOptions" />
            </div>
          </template>
        </Card>

        <Card class="surface-card lg:col-span-2">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-calendar" style="color: var(--color-primary);"></i>
              <span>Calendario de ocupacion</span>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-7 gap-2 text-xs text-gray-500 mb-2">
              <div v-for="day in weekDays" :key="day" class="text-center font-semibold">{{ day }}</div>
            </div>
            <div class="grid grid-cols-7 gap-2">
              <div v-for="(cell, index) in calendarCells" :key="`day-${index}`" class="min-h-[56px]">
                <div
                  v-if="cell"
                  class="h-full rounded-lg border p-2 flex flex-col gap-1"
                  :class="occupancyClass(cell.ocupacion)"
                >
                  <div class="flex items-center justify-between text-[11px]">
                    <span class="font-semibold">{{ cell.day }}</span>
                    <span v-if="cell.cupo_maximo > 0">{{ Math.round(cell.ocupacion * 100) }}%</span>
                  </div>
                  <div class="text-[10px] leading-tight">
                    <div>Conf: {{ cell.cupos_confirmados }}</div>
                    <div>Res: {{ cell.cupos_reservados }}</div>
                  </div>
                </div>
                <div v-else class="h-full rounded-lg border border-dashed border-gray-200 bg-gray-50"></div>
              </div>
            </div>
            <div class="mt-4 flex flex-wrap gap-3 text-xs text-gray-500">
              <div class="flex items-center gap-2"><span class="h-3 w-3 rounded bg-emerald-200"></span>0-25%</div>
              <div class="flex items-center gap-2"><span class="h-3 w-3 rounded bg-amber-200"></span>26-50%</div>
              <div class="flex items-center gap-2"><span class="h-3 w-3 rounded bg-orange-200"></span>51-75%</div>
              <div class="flex items-center gap-2"><span class="h-3 w-3 rounded bg-red-200"></span>76-100%</div>
            </div>
          </template>
        </Card>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-bell" style="color: var(--color-primary);"></i>
              <span>Alertas y acciones rapidas</span>
            </div>
          </template>
          <template #content>
            <div class="space-y-3 text-sm text-gray-700">
              <div class="flex items-start gap-3">
                <i class="pi pi-exclamation-triangle text-amber-500"></i>
                <p><span class="font-semibold">{{ alertas.pagos_pendientes }}</span> pagos pendientes de confirmacion.</p>
              </div>
              <div class="flex items-start gap-3">
                <i class="pi pi-calendar text-blue-500"></i>
                <p><span class="font-semibold">{{ alertas.salidas_proximas }}</span> salidas proximas (7 dias).</p>
              </div>
              <Button
                label="Ver calendario completo"
                icon="pi pi-calendar"
                severity="secondary"
                outlined
                @click="navigateTo('/agencia/calendario')"
              />
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-bolt" style="color: var(--color-primary);"></i>
              <span>Acciones rapidas</span>
            </div>
          </template>
          <template #content>
            <div class="space-y-4">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <p class="text-sm font-semibold text-gray-900">Mi agencia</p>
                  <p class="text-xs text-gray-500">Actualiza informacion principal.</p>
                </div>
                <div class="flex gap-2">
                  <Button label="Ver" icon="pi pi-eye" severity="secondary" outlined @click="navigateTo('/agencia/mi-agencia')" />
                  <Button label="Editar" icon="pi pi-pencil" severity="warning" @click="navigateTo('/agencia/mi-agencia/editar')" />
                </div>
              </div>

              <Divider />

              <div class="flex items-center justify-between gap-3">
                <div>
                  <p class="text-sm font-semibold text-gray-900">Paquetes turisticos</p>
                  <p class="text-xs text-gray-500">Gestiona tus paquetes y salidas.</p>
                </div>
                <div class="flex gap-2">
                  <Button label="Ver paquetes" icon="pi pi-eye" severity="secondary" outlined @click="navigateTo('/agencia/paquetes')" />
                  <Button label="Nuevo" icon="pi pi-plus" @click="navigateTo('/agencia/paquetes/nuevo')" />
                </div>
              </div>

              <Divider />

              <div class="flex items-center justify-between gap-3">
                <div>
                  <p class="text-sm font-semibold text-gray-900">Ventas de paquetes</p>
                  <p class="text-xs text-gray-500">Confirma pagos y revisa cupos.</p>
                </div>
                <div class="flex gap-2">
                  <Button label="Ver ventas" icon="pi pi-eye" severity="secondary" outlined @click="navigateTo('/agencia/ventas')" />
                </div>
              </div>
            </div>
          </template>
        </Card>
      </div>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import DashboardChart from '~/components/dashboard/DashboardChart.client.vue'

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const { getMiAgencia, getAgenciaDashboard } = useAgencias()

const agencia = ref<any>(null)
const dashboard = ref<any>(null)
const loadingDashboard = ref(true)
const dashboardError = ref<string | null>(null)

const monthLabels = [
  'Enero',
  'Febrero',
  'Marzo',
  'Abril',
  'Mayo',
  'Junio',
  'Julio',
  'Agosto',
  'Septiembre',
  'Octubre',
  'Noviembre',
  'Diciembre'
]

const weekDays = ['L', 'M', 'M', 'J', 'V', 'S', 'D']

const now = new Date()
const selectedMonth = ref(now.getMonth() + 1)
const selectedYear = ref(now.getFullYear())

const monthLabel = computed(() => `${monthLabels[selectedMonth.value - 1]} ${selectedYear.value}`)

const metrics = computed(() => dashboard.value?.metrics ?? {
  ingresos_mes: 0,
  ventas_confirmadas: 0,
  turistas_atendidos: 0,
  turistas_promedio: 0,
  paquetes_activos: 0,
  paquetes_nuevos: 0,
  pendientes_pago: 0,
  pendientes_pago_monto: 0
})

const comparatives = computed(() => dashboard.value?.comparatives ?? {
  ingresos_percent: 0,
  ventas_diff: 0
})

const ventasMensuales = computed(() => dashboard.value?.series?.ventas_mensuales ?? [])
const topPaquetes = computed(() => dashboard.value?.series?.top_paquetes ?? [])
const ingresosProyeccion = computed(() => dashboard.value?.series?.ingresos_vs_proyeccion ?? {
  confirmados: 0,
  pendientes: 0
})
const calendarioOcupacion = computed(() => dashboard.value?.calendario_ocupacion ?? [])
const alertas = computed(() => dashboard.value?.alertas ?? { pagos_pendientes: 0, salidas_proximas: 0 })

const formatCurrency = (value: number) => {
  const safe = Number(value || 0)
  return `Bs ${safe.toLocaleString('es-BO', { minimumFractionDigits: 0, maximumFractionDigits: 0 })}`
}

const formatDiff = (value: number) => `${value >= 0 ? '+' : ''}${value}`

const formatPercent = (value: number) => {
  const safe = Number.isFinite(value) ? value : 0
  return `${safe >= 0 ? '+' : ''}${safe.toFixed(1)}%`
}

const ingresosTrendClass = computed(() =>
  comparatives.value.ingresos_percent >= 0 ? 'text-emerald-600' : 'text-red-600'
)

const ventasTrendClass = computed(() =>
  comparatives.value.ventas_diff >= 0 ? 'text-emerald-600' : 'text-red-600'
)

const ventasLineData = computed(() => {
  const labels = ventasMensuales.value.map((row: any) => {
    const label = monthLabels[row.month - 1] || ''
    return `${label.slice(0, 3)} ${row.year}`
  })

  return {
    labels,
    datasets: [
      {
        label: 'Ventas',
        data: ventasMensuales.value.map((row: any) => row.ventas),
        borderColor: '#2563eb',
        backgroundColor: 'rgba(37, 99, 235, 0.15)',
        tension: 0.35,
        yAxisID: 'y'
      },
      {
        label: 'Ingresos',
        data: ventasMensuales.value.map((row: any) => row.ingresos),
        borderColor: '#16a34a',
        backgroundColor: 'rgba(22, 163, 74, 0.15)',
        tension: 0.35,
        yAxisID: 'y1'
      }
    ]
  }
})

const ventasLineOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom',
      labels: {
        color: '#6b7280'
      }
    },
    tooltip: {
      callbacks: {
        label(context: any) {
          if (context.dataset.label === 'Ingresos') {
            return `Ingresos: ${formatCurrency(context.parsed.y)}`
          }
          return `Ventas: ${context.parsed.y}`
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: { color: '#6b7280' },
      grid: { color: '#e5e7eb' }
    },
    y1: {
      beginAtZero: true,
      position: 'right',
      ticks: { color: '#6b7280' },
      grid: { drawOnChartArea: false }
    },
    x: {
      ticks: { color: '#6b7280' },
      grid: { display: false }
    }
  }
}

const topPaquetesData = computed(() => {
  const colors = ['#2563eb', '#16a34a', '#f59e0b', '#ef4444', '#8b5cf6']
  return {
    labels: topPaquetes.value.map((row: any) => row.nombre),
    datasets: [
      {
        data: topPaquetes.value.map((row: any) => row.ventas),
        backgroundColor: colors
      }
    ]
  }
})

const topPaquetesOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom',
      labels: { color: '#6b7280' }
    },
    tooltip: {
      callbacks: {
        label(context: any) {
          const index = context.dataIndex
          const item = topPaquetes.value[index]
          if (!item) return ''
          return `${item.nombre}: ${item.ventas} ventas (Bs ${Number(item.ingresos || 0).toLocaleString('es-BO')})`
        }
      }
    }
  }
}

const ingresosBarData = computed(() => ({
  labels: ['Confirmados', 'Pendientes'],
  datasets: [
    {
      data: [ingresosProyeccion.value.confirmados, ingresosProyeccion.value.pendientes],
      backgroundColor: ['#16a34a', '#f59e0b'],
      borderRadius: 8
    }
  ]
}))

const ingresosBarOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label(context: any) {
          return formatCurrency(context.parsed.y)
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: { color: '#6b7280' },
      grid: { color: '#e5e7eb' }
    },
    x: {
      ticks: { color: '#6b7280' },
      grid: { display: false }
    }
  }
}

const occupancyMap = computed(() => {
  const map = new Map()
  calendarioOcupacion.value.forEach((row: any) => {
    map.set(row.fecha, row)
  })
  return map
})

const formatDateKey = (year: number, month: number, day: number) => {
  const mm = String(month).padStart(2, '0')
  const dd = String(day).padStart(2, '0')
  return `${year}-${mm}-${dd}`
}

const calendarCells = computed(() => {
  if (!dashboard.value) return []

  const year = dashboard.value.anio || selectedYear.value
  const month = dashboard.value.mes || selectedMonth.value
  const daysInMonth = new Date(year, month, 0).getDate()
  const firstDay = new Date(year, month - 1, 1)
  const offset = (firstDay.getDay() + 6) % 7

  const cells: any[] = []
  for (let i = 0; i < offset; i += 1) {
    cells.push(null)
  }

  for (let day = 1; day <= daysInMonth; day += 1) {
    const key = formatDateKey(year, month, day)
    const info = occupancyMap.value.get(key) || {}
    cells.push({
      day,
      date: key,
      cupo_maximo: info.cupo_maximo || 0,
      cupos_confirmados: info.cupos_confirmados || 0,
      cupos_reservados: info.cupos_reservados || 0,
      ocupacion: info.ocupacion || 0
    })
  }

  return cells
})

const occupancyClass = (value: number) => {
  if (!value) return 'bg-gray-50 text-gray-400 border-gray-100'
  if (value < 0.26) return 'bg-emerald-50 text-emerald-700 border-emerald-100'
  if (value < 0.51) return 'bg-amber-50 text-amber-700 border-amber-100'
  if (value < 0.76) return 'bg-orange-50 text-orange-700 border-orange-100'
  return 'bg-red-50 text-red-700 border-red-100'
}

const loadDashboard = async () => {
  loadingDashboard.value = true
  dashboardError.value = null

  try {
    const agenciaResponse: any = await getMiAgencia()
    if (!agenciaResponse?.success) {
      dashboardError.value = agenciaResponse?.error?.message || 'No se pudo cargar la agencia'
      return
    }

    agencia.value = agenciaResponse.data

    const response: any = await getAgenciaDashboard(agencia.value.id, {
      mes: selectedMonth.value,
      anio: selectedYear.value
    })

    if (response.success) {
      dashboard.value = response.data
      return
    }

    dashboardError.value = response?.error?.message || 'No se pudo cargar el dashboard'
  } catch (err: any) {
    dashboardError.value = err?.data?.error?.message || err?.message || 'No se pudo cargar el dashboard'
    toast.add({ severity: 'error', summary: 'Error', detail: dashboardError.value, life: 3000 })
  } finally {
    loadingDashboard.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>
