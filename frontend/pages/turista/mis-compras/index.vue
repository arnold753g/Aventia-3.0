<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Mis compras</h1>
          <p class="muted mt-1">Consulta el estado de tus compras y registra pagos pendientes.</p>
        </div>
        <div class="flex gap-2">
          <Button label="Explorar paquetes" icon="pi pi-briefcase" @click="navigateTo('/turista/paquetes')" />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
        <div>
          <p class="text-sm text-gray-600">{{ compras.length }} compra(s) cargadas · {{ pagination.total }} total</p>
          <p class="text-xs text-gray-500">Vista cronograma: agrupada por fecha y horario.</p>
        </div>
        <div class="flex gap-2">
          <Button label="Actualizar" icon="pi pi-refresh" severity="secondary" outlined :loading="loading" @click="refresh" />
        </div>
      </div>

      <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

      <div v-if="loading" class="space-y-4">
        <Card v-for="i in 3" :key="i" class="surface-card">
          <template #content>
            <div class="flex gap-4">
              <Skeleton width="96px" height="120px" />
              <div class="flex-1 space-y-3">
                <Skeleton height="1.2rem" width="35%" />
                <Skeleton height="6rem" />
              </div>
            </div>
          </template>
        </Card>
      </div>

      <div v-else class="space-y-4">
        <Card v-if="dias.length === 0" class="surface-card">
          <template #content>
            <div class="text-center py-12">
              <i class="pi pi-shopping-cart text-6xl muted mb-4 block"></i>
              <p class="text-xl font-semibold muted">Aún no tienes compras</p>
              <p class="text-sm muted mt-2">Explora paquetes turísticos y realiza tu primera compra.</p>
            </div>
          </template>
        </Card>

        <Card v-for="dia in dias" :key="dia.dateKey" class="surface-card overflow-hidden">
          <template #content>
            <div class="flex flex-col md:flex-row">
              <div
                class="md:w-28 w-full flex md:flex-col items-center justify-between md:justify-center gap-2 bg-emerald-50 border-b md:border-b-0 md:border-r border-emerald-100 px-4 py-3"
              >
                <div class="text-center leading-none">
                  <p class="text-3xl font-bold text-emerald-700">{{ dia.day }}</p>
                  <p class="text-xs uppercase tracking-widest text-emerald-700">{{ dia.month }}</p>
                </div>
                <p class="text-xs text-emerald-700/80">{{ dia.weekday }}</p>
              </div>

              <div class="flex-1 p-4">
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div class="rounded-xl border border-gray-200 bg-white p-3 space-y-3">
                    <div class="flex items-center justify-between">
                      <p class="text-xs uppercase tracking-wider text-gray-500 font-semibold">Mañana</p>
                      <Tag v-if="dia.slots.manana.length" :value="`${dia.slots.manana.length}`" severity="info" />
                    </div>

                    <p v-if="dia.slots.manana.length === 0" class="text-xs text-gray-500">Sin paquetes</p>

                    <div v-else class="space-y-2">
                      <div v-for="item in dia.slots.manana" :key="item.key" :class="entryCardClass(item)">
                        <div class="flex items-start justify-between gap-2">
                          <div class="min-w-0">
                            <p class="font-semibold text-gray-900 truncate">{{ item.compra.paquete?.nombre || `#${item.compra.paquete?.id}` }}</p>
                            <p class="text-xs text-gray-600 mt-1">Bs. {{ formatMoney(item.compra.precio_total) }} · {{ item.compra.total_participantes }} pax</p>
                          </div>
                          <Tag :value="statusLabel(item.compra.status)" :severity="statusSeverity(item.compra.status)" />
                        </div>

                        <div class="flex flex-wrap items-center justify-between gap-2 mt-3">
                          <div class="flex flex-wrap gap-2">
                            <Tag :value="item.compra.tipo_compra === 'privado' ? 'Privado' : 'Compartido'" :severity="item.compra.tipo_compra === 'privado' ? 'warning' : 'info'" />
                            <Tag v-if="item.totalDays > 1" :value="`Día ${item.dayIndex}/${item.totalDays}`" severity="secondary" />
                            <Tag v-if="item.compra.ultimo_pago?.estado === 'pendiente'" value="Pago pendiente" severity="warning" />
                          </div>

                          <Button
                            v-if="canPagar(item.compra)"
                            label="Pagar"
                            icon="pi pi-credit-card"
                            size="small"
                            @click="navigateTo(`/turista/compras/${item.compra.id}/pagar`)"
                          />
                          <Button
                            v-else
                            label="Ver"
                            icon="pi pi-eye"
                            severity="secondary"
                            outlined
                            size="small"
                            @click="navigateTo(`/turista/compras/${item.compra.id}/pagar`)"
                          />
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="rounded-xl border border-gray-200 bg-white p-3 space-y-3">
                    <div class="flex items-center justify-between">
                      <p class="text-xs uppercase tracking-wider text-gray-500 font-semibold">Tarde</p>
                      <Tag v-if="dia.slots.tarde.length" :value="`${dia.slots.tarde.length}`" severity="info" />
                    </div>

                    <p v-if="dia.slots.tarde.length === 0" class="text-xs text-gray-500">Sin paquetes</p>

                    <div v-else class="space-y-2">
                      <div v-for="item in dia.slots.tarde" :key="item.key" :class="entryCardClass(item)">
                        <div class="flex items-start justify-between gap-2">
                          <div class="min-w-0">
                            <p class="font-semibold text-gray-900 truncate">{{ item.compra.paquete?.nombre || `#${item.compra.paquete?.id}` }}</p>
                            <p class="text-xs text-gray-600 mt-1">Bs. {{ formatMoney(item.compra.precio_total) }} · {{ item.compra.total_participantes }} pax</p>
                          </div>
                          <Tag :value="statusLabel(item.compra.status)" :severity="statusSeverity(item.compra.status)" />
                        </div>

                        <div class="flex flex-wrap items-center justify-between gap-2 mt-3">
                          <div class="flex flex-wrap gap-2">
                            <Tag :value="item.compra.tipo_compra === 'privado' ? 'Privado' : 'Compartido'" :severity="item.compra.tipo_compra === 'privado' ? 'warning' : 'info'" />
                            <Tag v-if="item.totalDays > 1" :value="`Día ${item.dayIndex}/${item.totalDays}`" severity="secondary" />
                            <Tag v-if="item.compra.ultimo_pago?.estado === 'pendiente'" value="Pago pendiente" severity="warning" />
                          </div>

                          <Button
                            v-if="canPagar(item.compra)"
                            label="Pagar"
                            icon="pi pi-credit-card"
                            size="small"
                            @click="navigateTo(`/turista/compras/${item.compra.id}/pagar`)"
                          />
                          <Button
                            v-else
                            label="Ver"
                            icon="pi pi-eye"
                            severity="secondary"
                            outlined
                            size="small"
                            @click="navigateTo(`/turista/compras/${item.compra.id}/pagar`)"
                          />
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="rounded-xl border border-gray-200 bg-white p-3 space-y-3">
                    <div class="flex items-center justify-between">
                      <p class="text-xs uppercase tracking-wider text-gray-500 font-semibold">Todo el día</p>
                      <Tag v-if="dia.slots.todo_dia.length" :value="`${dia.slots.todo_dia.length}`" severity="info" />
                    </div>

                    <p v-if="dia.slots.todo_dia.length === 0" class="text-xs text-gray-500">Sin paquetes</p>

                    <div v-else class="space-y-2">
                      <div v-for="item in dia.slots.todo_dia" :key="item.key" :class="entryCardClass(item)">
                        <div class="flex items-start justify-between gap-2">
                          <div class="min-w-0">
                            <p class="font-semibold text-gray-900 truncate">{{ item.compra.paquete?.nombre || `#${item.compra.paquete?.id}` }}</p>
                            <p class="text-xs text-gray-600 mt-1">Bs. {{ formatMoney(item.compra.precio_total) }} · {{ item.compra.total_participantes }} pax</p>
                          </div>
                          <Tag :value="statusLabel(item.compra.status)" :severity="statusSeverity(item.compra.status)" />
                        </div>

                        <div class="flex flex-wrap items-center justify-between gap-2 mt-3">
                          <div class="flex flex-wrap gap-2">
                            <Tag :value="item.compra.tipo_compra === 'privado' ? 'Privado' : 'Compartido'" :severity="item.compra.tipo_compra === 'privado' ? 'warning' : 'info'" />
                            <Tag v-if="item.totalDays > 1" :value="`Día ${item.dayIndex}/${item.totalDays}`" severity="secondary" />
                            <Tag v-if="item.compra.ultimo_pago?.estado === 'pendiente'" value="Pago pendiente" severity="warning" />
                          </div>

                          <Button
                            v-if="canPagar(item.compra)"
                            label="Pagar"
                            icon="pi pi-credit-card"
                            size="small"
                            @click="navigateTo(`/turista/compras/${item.compra.id}/pagar`)"
                          />
                          <Button
                            v-else
                            label="Ver"
                            icon="pi pi-eye"
                            severity="secondary"
                            outlined
                            size="small"
                            @click="navigateTo(`/turista/compras/${item.compra.id}/pagar`)"
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </Card>

        <div v-if="hasMore" class="flex justify-center">
          <Button label="Cargar más" icon="pi pi-chevron-down" :loading="loadingMore" @click="loadMore" />
        </div>
      </div>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import type { CompraDetalle } from '~/types/compra'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

type SlotKey = 'manana' | 'tarde' | 'todo_dia'

type ScheduleEntry = {
  key: string
  dateKey: string
  slot: SlotKey
  compra: CompraDetalle
  dayIndex: number
  totalDays: number
}

type DayGroup = {
  dateKey: string
  day: number
  month: string
  weekday: string
  slots: Record<SlotKey, ScheduleEntry[]>
}

const toast = useToast()
const { listarMisCompras } = useCompra()

const compras = ref<CompraDetalle[]>([])
const loading = ref(true)
const loadingMore = ref(false)
const error = ref<string | null>(null)

const pagination = ref({
  page: 1,
  page_size: 20,
  total: 0,
  total_pages: 0
})

const formatMoney = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const statusLabel = (status: string) => {
  const map: Record<string, string> = {
    pendiente_confirmacion: 'Pendiente de confirmación',
    confirmada: 'Confirmada',
    rechazada: 'Rechazada',
    cancelada: 'Cancelada',
    completada: 'Completada'
  }
  return map[status] || status
}

const statusSeverity = (status: string) => {
  const map: Record<string, any> = {
    pendiente_confirmacion: 'warning',
    confirmada: 'success',
    rechazada: 'danger',
    cancelada: 'secondary',
    completada: 'info'
  }
  return map[status] || 'secondary'
}

const normalizeDateKey = (value?: string) => {
  if (!value) return ''
  const raw = String(value)
  const datePart = raw.split('T').shift() ?? raw
  return datePart.split(' ').shift() ?? datePart
}

const dateFromKey = (key: string) => {
  const [year, month, day] = key.split('-').map((n) => Number(n))
  return new Date(year, (month || 1) - 1, day || 1)
}

const toDateKey = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const monthLabel = (date: Date) =>
  date
    .toLocaleString('es-BO', { month: 'short' })
    .replace(/\./g, '')
    .toUpperCase()

const weekdayLabel = (date: Date) =>
  date
    .toLocaleString('es-BO', { weekday: 'short' })
    .replace(/\./g, '')
    .toUpperCase()

const entryCardClass = (item: ScheduleEntry) => {
  const base = 'rounded-lg border p-3'
  const needsPayment =
    item.compra.status === 'pendiente_confirmacion' && item.compra.ultimo_pago?.estado !== 'pendiente'
  if (needsPayment) return [base, 'border-red-200 bg-red-50']
  if (item.totalDays > 1) return [base, 'border-indigo-200 bg-indigo-50']
  return [base, 'border-gray-100 bg-gray-50']
}

const normalizeSlot = (horario?: string | null): SlotKey => {
  const h = String(horario || '').toLowerCase()
  if (h === 'mañana' || h === 'manana') return 'manana'
  if (h === 'tarde') return 'tarde'
  return 'todo_dia'
}

const scheduleEntries = computed<ScheduleEntry[]>(() => {
  const out: ScheduleEntry[] = []

  for (const compra of compras.value) {
    const startKey = normalizeDateKey(compra.fecha_seleccionada)
    if (!startKey) continue

    const totalDays = Math.max(1, Number(compra.paquete?.duracion_dias || 1))
    const startDate = dateFromKey(startKey)

    if (totalDays > 1) {
      for (let i = 0; i < totalDays; i++) {
        const d = new Date(startDate)
        d.setDate(startDate.getDate() + i)
        const key = toDateKey(d)
        out.push({
          key: `${compra.id}-${key}-${i + 1}`,
          dateKey: key,
          slot: 'todo_dia',
          compra,
          dayIndex: i + 1,
          totalDays
        })
      }
      continue
    }

    const slot = normalizeSlot(compra.paquete?.horario)
    out.push({
      key: `${compra.id}-${startKey}-1`,
      dateKey: startKey,
      slot,
      compra,
      dayIndex: 1,
      totalDays: 1
    })
  }

  return out
})

const dias = computed<DayGroup[]>(() => {
  const map = new Map<string, DayGroup>()

  for (const entry of scheduleEntries.value) {
    let group = map.get(entry.dateKey)
    if (!group) {
      const date = dateFromKey(entry.dateKey)
      group = {
        dateKey: entry.dateKey,
        day: date.getDate(),
        month: monthLabel(date),
        weekday: weekdayLabel(date),
        slots: { manana: [], tarde: [], todo_dia: [] }
      }
      map.set(entry.dateKey, group)
    }

    group.slots[entry.slot].push(entry)
  }

  const groups = Array.from(map.values())
  groups.sort((a, b) => a.dateKey.localeCompare(b.dateKey))

  for (const group of groups) {
    for (const slot of Object.keys(group.slots) as SlotKey[]) {
      group.slots[slot].sort((a, b) => String(b.compra.fecha_compra || '').localeCompare(String(a.compra.fecha_compra || '')))
    }
  }

  return groups
})

const canPagar = (compra: CompraDetalle) => compra.status === 'pendiente_confirmacion' && compra.ultimo_pago?.estado !== 'pendiente'

const hasMore = computed(() => pagination.value.page < pagination.value.total_pages)

const mergeUnique = (incoming: CompraDetalle[]) => {
  const ids = new Set(compras.value.map((c) => c.id))
  for (const item of incoming) {
    if (ids.has(item.id)) continue
    compras.value.push(item)
    ids.add(item.id)
  }
}

const loadPage = async (page: number, opts: { append: boolean }) => {
  if (opts.append) {
    loadingMore.value = true
  } else {
    loading.value = true
  }

  error.value = null

  try {
    const response: any = await listarMisCompras({ page, page_size: pagination.value.page_size })
    if (!response?.success) {
      error.value = response?.error?.message || 'No se pudieron cargar las compras'
      return
    }

    const list: CompraDetalle[] = response.data?.compras || []
    const pag = response.data?.pagination || {}
    pagination.value = { ...pagination.value, ...pag }

    if (opts.append) {
      mergeUnique(list)
      return
    }

    compras.value = list
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar las compras'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const refresh = async () => {
  await loadPage(1, { append: false })
}

const loadMore = async () => {
  if (loadingMore.value) return
  if (!hasMore.value) return
  await loadPage(pagination.value.page + 1, { append: true })
}

onMounted(refresh)
</script>
