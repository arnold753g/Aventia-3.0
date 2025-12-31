<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Calendario de salidas</h1>
          <p class="muted mt-1">Visualiza cupos, estados y gestiona salidas por fecha.</p>
        </div>
        <div class="text-sm text-gray-500">
          {{ rangeLabel }}
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

      <Card class="surface-card">
        <template #content>
          <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div class="flex flex-wrap items-center gap-2">
              <Button icon="pi pi-angle-left" severity="secondary" outlined @click="prevRange" />
              <div class="min-w-[220px]">
                <p class="text-xs uppercase tracking-wide text-gray-500">{{ viewLabel }}</p>
                <p class="text-lg font-semibold text-gray-900">{{ rangeLabel }}</p>
              </div>
              <Button icon="pi pi-angle-right" severity="secondary" outlined @click="nextRange" />
              <Button label="Hoy" severity="secondary" text @click="goToToday" />
            </div>

            <div class="flex flex-wrap items-center gap-3">
              <SelectButton
                v-model="viewMode"
                :options="viewOptions"
                optionLabel="label"
                optionValue="value"
              />
              <Dropdown
                v-model="paqueteFilter"
                :options="paqueteOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Todos los paquetes"
                class="min-w-[220px]"
                showClear
              />
              <Button
                label="Exportar iCal"
                icon="pi pi-calendar-plus"
                severity="secondary"
                outlined
                @click="exportICal"
              />
              <Button
                label="Habilitar salida"
                icon="pi pi-plus"
                @click="openCrearSalida"
              />
            </div>
          </div>
        </template>
      </Card>

      <Card v-if="crearDialog" class="surface-card">
        <template #content>
          <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
            <div>
              <p class="text-sm text-gray-500">Gestion de salidas</p>
              <h2 class="text-xl font-bold text-gray-900">Habilitar salida</h2>
              <p class="text-sm text-gray-600">Define el paquete, fecha y estado inicial.</p>
            </div>
            <Button label="Cerrar" severity="secondary" outlined @click="crearDialog = false" />
          </div>

          <Divider class="my-4" />

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="space-y-2 md:col-span-2">
              <label class="block text-sm font-medium text-gray-700">Paquete</label>
              <Dropdown
                v-model="crearForm.paquete_id"
                :options="paqueteSelectOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Selecciona un paquete"
                class="w-full"
              />
            </div>
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700">Estado inicial</label>
              <Dropdown
                v-model="crearForm.estado"
                :options="estadoSalidaOptions"
                optionLabel="label"
                optionValue="value"
                class="w-full"
              />
            </div>
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700">Fecha de salida</label>
              <Calendar v-model="crearForm.fecha" :showIcon="true" dateFormat="yy-mm-dd" class="w-full" />
            </div>
          </div>

          <div class="flex justify-end gap-2 mt-4">
            <Button label="Cancelar" severity="secondary" outlined @click="crearDialog = false" />
            <Button label="Habilitar" icon="pi pi-check" :loading="crearLoading" @click="crearSalida" />
          </div>
        </template>
      </Card>

      <Card v-if="dayDialog" class="surface-card">
        <template #content>
          <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-3">
            <div>
              <p class="text-sm text-gray-500">Salidas del dia</p>
              <h2 class="text-xl font-bold text-gray-900">Salidas del {{ selectedDateLabel }}</h2>
              <p class="text-xs text-gray-600">{{ selectedDaySalidas.length }} salida(s)</p>
            </div>
            <Button label="Cerrar" severity="secondary" outlined @click="closeDayPanel" />
          </div>

          <Divider class="my-4" />

          <div class="space-y-4">
            <div v-if="selectedDaySalidas.length === 0" class="text-sm text-gray-500">No hay salidas para este dia.</div>

            <div v-for="salida in selectedDaySalidas" :key="`salida-${salida.salida_id}`" class="rounded-xl border border-gray-200 bg-white p-4">
              <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-3">
                <div class="min-w-0">
                  <p class="text-base font-semibold text-gray-900">{{ salida.paquete_nombre }}</p>
                  <p class="text-xs text-gray-600 mt-1">
                    {{ tipoSalidaLabel(salida.tipo_salida) }} ¶ú {{ salida.paquete_horario || 'Horario por definir' }}
                  </p>
                  <p class="text-xs text-gray-500 mt-1">Estado: {{ estadoSalidaLabel(salida.estado) }}</p>
                </div>
                <div class="flex flex-wrap gap-2">
                  <Tag :value="cupoEstadoLabel(salida)" :severity="cupoEstadoSeverity(salida)" />
                  <Tag :value="estadoSalidaLabel(salida.estado)" :severity="estadoSalidaSeverity(salida.estado)" />
                </div>
              </div>

              <div class="mt-3 text-xs text-gray-600">
                Cupos: {{ salida.cupos_confirmados }} confirmados ¶ú {{ salida.cupos_reservados }} reservados ¶ú
                {{ cuposDisponibles(salida) }} disponibles
              </div>

              <div class="mt-4 flex flex-wrap gap-2 justify-end">
                <Button label="Ver detalle" icon="pi pi-eye" severity="secondary" outlined @click="openDetalle(salida)" />
                <Button label="Editar logistica" icon="pi pi-pencil" severity="warning" @click="openLogistica(salida)" />
                <Button
                  v-if="salida.estado !== 'cancelada' && salida.estado !== 'completada'"
                  label="Cancelar"
                  icon="pi pi-times"
                  severity="danger"
                  outlined
                  @click="openCancelar(salida)"
                />
              </div>
            </div>
          </div>
        </template>
      </Card>

      <Card v-if="detalleDialog" class="surface-card">
        <template #content>
          <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-3">
            <div>
              <p class="text-sm text-gray-500">Detalle de salida</p>
              <h2 class="text-xl font-bold text-gray-900">{{ detalleSalida?.paquete_nombre || 'Detalle de salida' }}</h2>
              <p class="text-xs text-gray-600">
                {{ detalleSalida?.fecha_salida ? formatFecha(detalleSalida?.fecha_salida) : '' }}
              </p>
            </div>
            <Button label="Cerrar detalle" severity="secondary" outlined @click="closeDetalle" />
          </div>

          <Divider class="my-4" />

          <div class="space-y-4">
            <Message v-if="detalleError" severity="error" :closable="false">{{ detalleError }}</Message>
            <div v-if="detalleLoading" class="flex items-center justify-center py-10">
              <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="4" />
            </div>
            <div v-else-if="detalleSalida" class="space-y-4">
              <div class="rounded-xl border border-gray-200 bg-gray-50 p-4 text-sm text-gray-700">
                <p><span class="font-semibold">Tipo:</span> {{ tipoSalidaLabel(detalleSalida.tipo_salida) }}</p>
                <p><span class="font-semibold">Estado:</span> {{ estadoSalidaLabel(detalleSalida.estado) }}</p>
                <p>
                  <span class="font-semibold">Cupos:</span>
                  {{ detalleSalida.cupos_confirmados }} confirmados ¶ú {{ detalleSalida.cupos_reservados }} reservados ¶ú
                  {{ cuposDisponibles(detalleSalida) }} disponibles
                </p>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-3 gap-3 text-sm text-gray-700">
                <div class="rounded-lg border border-gray-200 bg-white p-3">
                  <p class="text-xs text-gray-500 uppercase">Compras</p>
                  <p class="text-lg font-semibold text-gray-900">{{ detalleTotales.compras }}</p>
                </div>
                <div class="rounded-lg border border-gray-200 bg-white p-3">
                  <p class="text-xs text-gray-500 uppercase">Participantes</p>
                  <p class="text-lg font-semibold text-gray-900">{{ detalleTotales.participantes }}</p>
                </div>
                <div class="rounded-lg border border-gray-200 bg-white p-3">
                  <p class="text-xs text-gray-500 uppercase">Ingresos</p>
                  <p class="text-lg font-semibold text-emerald-700">Bs {{ formatMoney(detalleTotales.ingresos) }}</p>
                </div>
              </div>

              <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                <div class="flex flex-wrap gap-2">
                  <Button
                    label="Editar logistica"
                    icon="pi pi-pencil"
                    severity="warning"
                    size="small"
                    :disabled="!detalleBaseSalida"
                    @click="openLogisticaDetalle"
                  />
                  <Button
                    v-if="detalleSalida.estado !== 'cancelada' && detalleSalida.estado !== 'completada'"
                    label="Cancelar salida"
                    icon="pi pi-times"
                    severity="danger"
                    outlined
                    size="small"
                    :disabled="!detalleBaseSalida"
                    @click="openCancelarDetalle"
                  />
                </div>
                <Button
                  label="Actualizar detalle"
                  icon="pi pi-refresh"
                  size="small"
                  severity="secondary"
                  outlined
                  :loading="detalleLoading"
                  @click="refreshDetalle"
                />
              </div>

              <div v-if="logisticaDialog" class="rounded-xl border border-emerald-200 bg-emerald-50/40 p-4 space-y-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <label class="block text-sm font-medium text-gray-700">Punto de encuentro</label>
                    <InputText v-model="logisticaForm.punto_encuentro" class="w-full" />
                  </div>
                  <div class="space-y-2">
                    <label class="block text-sm font-medium text-gray-700">Hora de encuentro</label>
                    <InputText v-model="logisticaForm.hora_encuentro" class="w-full" placeholder="08:30" />
                  </div>
                  <div class="space-y-2">
                    <label class="block text-sm font-medium text-gray-700">Guia</label>
                    <InputText v-model="logisticaForm.guia_nombre" class="w-full" />
                  </div>
                  <div class="space-y-2">
                    <label class="block text-sm font-medium text-gray-700">Telefono guia</label>
                    <InputText v-model="logisticaForm.guia_telefono" class="w-full" />
                  </div>
                  <div class="space-y-2 md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700">Notas logistica</label>
                    <Textarea v-model="logisticaForm.notas_logistica" rows="3" class="w-full" />
                  </div>
                  <div class="space-y-2 md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700">Instrucciones para turistas</label>
                    <Textarea v-model="logisticaForm.instrucciones_turistas" rows="3" class="w-full" />
                  </div>
                </div>
                <div class="flex justify-end gap-2">
                  <Button label="Cancelar" severity="secondary" outlined @click="logisticaDialog = false" />
                  <Button label="Guardar" icon="pi pi-save" :loading="logisticaLoading" @click="guardarLogistica" />
                </div>
              </div>

              <div v-if="cancelDialog" class="rounded-xl border border-red-200 bg-red-50/40 p-4 space-y-4">
                <div class="flex items-center justify-between gap-2">
                  <p class="text-sm font-semibold text-gray-900">Cancelar salida</p>
                  <Button label="Cerrar" severity="secondary" outlined size="small" @click="cancelDialog = false" />
                </div>
                <p class="text-sm text-gray-700">Indica la razon para cancelar esta salida.</p>
                <Textarea v-model="cancelReason" rows="3" class="w-full" placeholder="Ej: no se alcanzo el cupo minimo" />
                <div class="flex justify-end gap-2">
                  <Button label="Volver" severity="secondary" outlined @click="cancelDialog = false" />
                  <Button label="Cancelar salida" icon="pi pi-times" severity="danger" :loading="cancelLoading" @click="confirmarCancelacion" />
                </div>
              </div>

              <DataTable
                :value="detalleCompras"
                dataKey="compra_id"
                class="p-datatable-sm"
                :rows="10"
                paginator
                responsiveLayout="scroll"
              >
                <Column header="Turista">
                  <template #body="{ data }">
                    <div>
                      <p class="font-semibold text-gray-900">
                        {{ data.turista_nombre }} {{ data.turista_apellido_paterno }} {{ data.turista_apellido_materno }}
                      </p>
                      <p class="text-xs text-gray-600">{{ data.turista_phone || 'Sin telefono' }}</p>
                    </div>
                  </template>
                </Column>
                <Column header="Participantes" style="width: 120px">
                  <template #body="{ data }">
                    <span class="font-semibold">{{ data.total_participantes }}</span>
                  </template>
                </Column>
                <Column header="Monto" style="width: 140px">
                  <template #body="{ data }">
                    <span class="font-semibold text-emerald-700">Bs {{ formatMoney(data.precio_total) }}</span>
                  </template>
                </Column>
                <Column header="Confirmacion" style="width: 170px">
                  <template #body="{ data }">
                    <span class="text-xs text-gray-600">{{ formatFechaHora(data.fecha_confirmacion) || 'ƒ?"' }}</span>
                  </template>
                </Column>
                <Column header="Notas">
                  <template #body="{ data }">
                    <span class="text-xs text-gray-600" :title="data.notas_turista || ''">
                      {{ data.notas_turista ? truncateText(data.notas_turista, 40) : 'ƒ?"' }}
                    </span>
                  </template>
                </Column>
                <Column header="Estado pago" style="width: 140px">
                  <template #body="{ data }">
                    <Tag :value="estadoPagoLabel(data.estado)" :severity="estadoPagoSeverity(data.estado)" />
                  </template>
                </Column>
                <Column header="Comprobante" style="width: 160px">
                  <template #body="{ data }">
                    <Button
                      v-if="data.comprobante_foto"
                      icon="pi pi-image"
                      label="Ver"
                      size="small"
                      severity="secondary"
                      outlined
                      @click="openComprobante(data)"
                    />
                    <span v-else class="text-xs text-gray-400">Sin comprobante</span>
                  </template>
                </Column>
              </DataTable>

              <div v-if="comprobanteDialog" class="rounded-xl border border-gray-200 bg-white p-4 space-y-3">
                <div class="flex items-center justify-between gap-2">
                  <p class="text-sm font-semibold text-gray-900">Comprobante</p>
                  <Button label="Cerrar" severity="secondary" outlined size="small" @click="comprobanteDialog = false" />
                </div>
                <div v-if="selectedCompra?.comprobante_foto" class="space-y-3">
                  <img
                    :src="resolveAssetUrl(selectedCompra.comprobante_foto)"
                    alt="Comprobante"
                    class="w-full h-auto rounded-lg border border-gray-200"
                  />
                </div>
                <p v-else class="text-sm text-gray-600">No hay comprobante disponible.</p>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card class="surface-card lg:col-span-2">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-calendar text-green-600"></i>
              <span>{{ viewLabel }}</span>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-7 gap-2 text-xs text-gray-500 mb-2">
              <div v-for="day in weekDays" :key="day" class="text-center font-semibold">{{ day }}</div>
            </div>

            <div v-if="loading" class="grid grid-cols-7 gap-2">
              <div v-for="i in 35" :key="`skeleton-${i}`" class="h-20 rounded-lg border border-gray-200 bg-gray-50"></div>
            </div>

            <div v-else-if="viewMode === 'mes'" class="grid grid-cols-7 gap-2">
              <div v-for="(cell, index) in calendarCells" :key="`day-${index}`" class="min-h-[80px]">
                <div
                  v-if="cell"
                  class="h-full rounded-lg border p-2 flex flex-col gap-1 cursor-pointer transition-colors"
                  :class="dayCellClass(cell)"
                  @click="openDay(cell.date)"
                >
                  <div class="flex items-center justify-between text-[11px]">
                    <span class="font-semibold">{{ cell.day }}</span>
                    <span v-if="cell.stats.total > 0">{{ cell.stats.total }}</span>
                  </div>
                  <div class="text-[10px] leading-tight">
                    <div v-if="cell.stats.confirmadas">Conf: {{ cell.stats.confirmadas }}</div>
                    <div v-if="cell.stats.pendientes">Pend: {{ cell.stats.pendientes }}</div>
                    <div v-if="cell.stats.llenas">Llenas: {{ cell.stats.llenas }}</div>
                  </div>
                </div>
                <div v-else class="h-full rounded-lg border border-dashed border-gray-200 bg-gray-50"></div>
              </div>
            </div>

            <div v-else class="grid grid-cols-1 gap-4">
              <div class="grid grid-cols-7 gap-2">
                <div
                  v-for="day in weekDates"
                  :key="day.date"
                  class="rounded-lg border border-gray-200 bg-white p-3 min-h-[120px] cursor-pointer"
                  @click="openDay(day.date)"
                >
                  <div class="flex items-center justify-between text-xs text-gray-500">
                    <span class="font-semibold text-gray-900">{{ day.label }}</span>
                    <span>{{ day.count }}</span>
                  </div>
                  <div v-if="day.items.length" class="mt-2 space-y-1 text-[11px] text-gray-600">
                    <div v-for="item in day.items.slice(0, 3)" :key="item.salida_id" class="truncate">
                      {{ item.paquete_nombre }}
                    </div>
                    <div v-if="day.items.length > 3" class="text-[10px] text-gray-400">
                      +{{ day.items.length - 3 }} mas
                    </div>
                  </div>
                  <div v-else class="mt-6 text-[11px] text-gray-400">Sin salidas</div>
                </div>
              </div>
            </div>

            <div class="mt-4 flex flex-wrap gap-3 text-xs text-gray-500">
              <div class="flex items-center gap-2"><span class="h-3 w-3 rounded bg-emerald-200"></span>Cupo minimo alcanzado</div>
              <div class="flex items-center gap-2"><span class="h-3 w-3 rounded bg-amber-200"></span>Esperando cupo minimo</div>
              <div class="flex items-center gap-2"><span class="h-3 w-3 rounded bg-red-200"></span>Capacidad maxima</div>
            </div>
          </template>
        </Card>

        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-bell text-amber-500"></i>
              <span>Alertas de bajo cupo</span>
            </div>
          </template>
          <template #content>
            <div v-if="alertasBajoCupo.length === 0" class="text-sm text-gray-500">
              No hay salidas con bajo cupo en este rango.
            </div>
            <div v-else class="space-y-3 text-sm">
              <div v-for="salida in alertasBajoCupo" :key="`alert-${salida.salida_id}`" class="rounded-lg border border-amber-200 bg-amber-50 p-3">
                <p class="font-semibold text-gray-900">{{ salida.paquete_nombre }}</p>
                <p class="text-xs text-gray-600">{{ formatFecha(salida.fecha_salida) }} · {{ tipoSalidaLabel(salida.tipo_salida) }}</p>
                <p class="text-xs text-amber-700 mt-1">
                  {{ cuposActuales(salida) }} / {{ salida.cupo_minimo }} cupo minimo
                </p>
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
import { computed, onMounted, ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import type { AgenciaVentaSalida, AgenciaVentaSalidaCompra, VentaSalidaDetalle } from '~/types/ventas'

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const config = useRuntimeConfig()
const assetsBase = String(config.public.apiBase || '').replace(/\/api\/v1\/?$/, '')
const { getMiAgencia } = useAgencias()
const { getVentasSalidas, getVentasSalidaCompras } = useVentasAgencia()
const { getPaquetes, createPaqueteSalida, updatePaqueteSalida } = usePaquetes()

const agenciaId = ref<number | null>(null)
const salidas = ref<AgenciaVentaSalida[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const viewOptions = [
  { label: 'Mes', value: 'mes' },
  { label: 'Semana', value: 'semana' }
]
const viewMode = ref<'mes' | 'semana'>('mes')
const currentDate = ref(new Date())
const paqueteFilter = ref<number | null>(null)
const paquetesDisponibles = ref<any[]>([])
const crearDialog = ref(false)
const crearLoading = ref(false)
const crearForm = ref({
  paquete_id: null as number | null,
  fecha: null as Date | null,
  estado: 'pendiente'
})

const estadoSalidaOptions = [
  { label: 'Pendiente', value: 'pendiente' },
  { label: 'Activa', value: 'activa' }
]

const weekDays = ['L', 'M', 'M', 'J', 'V', 'S', 'D']
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

const resolveAgenciaId = async () => {
  if (agenciaId.value) return agenciaId.value
  const resp: any = await getMiAgencia()
  if (!resp?.success) throw new Error(resp?.error?.message || 'No se pudo cargar la agencia')
  agenciaId.value = Number(resp.data?.id)
  if (!agenciaId.value) throw new Error('No se pudo resolver el ID de la agencia')
  return agenciaId.value
}

const resolveAssetUrl = (path?: string | null) => {
  if (!path) return ''
  let normalized = String(path).replace(/\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const formatDateKey = (date: Date) => {
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

const normalizeDateKey = (value: string) => {
  if (!value) return ''
  const raw = value.split('T')[0]
  return raw.split(' ')[0]
}

const monthLabel = computed(() => {
  const date = currentDate.value
  return `${monthLabels[date.getMonth()]} ${date.getFullYear()}`
})

const rangeLabel = computed(() => {
  if (viewMode.value === 'mes') return monthLabel.value
  const start = weekStart.value
  const end = weekEnd.value
  return `Semana del ${formatFecha(formatDateKey(start))} al ${formatFecha(formatDateKey(end))}`
})

const viewLabel = computed(() => (viewMode.value === 'mes' ? 'Vista mensual' : 'Vista semanal'))

const weekStart = computed(() => {
  const date = currentDate.value
  const day = date.getDay()
  const diff = (day + 6) % 7
  return new Date(date.getFullYear(), date.getMonth(), date.getDate() - diff)
})

const weekEnd = computed(() => {
  const start = weekStart.value
  return new Date(start.getFullYear(), start.getMonth(), start.getDate() + 6)
})

const rangeStart = computed(() => {
  if (viewMode.value === 'semana') return weekStart.value
  return new Date(currentDate.value.getFullYear(), currentDate.value.getMonth(), 1)
})

const rangeEnd = computed(() => {
  if (viewMode.value === 'semana') return weekEnd.value
  return new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 0)
})

const paqueteOptions = computed(() => {
  const map = new Map<number, string>()
  salidas.value.forEach((salida) => {
    map.set(salida.paquete_id, salida.paquete_nombre)
  })
  const options = Array.from(map.entries())
    .map(([value, label]) => ({ label, value }))
    .sort((a, b) => a.label.localeCompare(b.label))
  return options
})

const paqueteSelectOptions = computed(() => {
  return (paquetesDisponibles.value || [])
    .map((paquete: any) => ({
      label: paquete.nombre,
      value: paquete.id
    }))
    .sort((a, b) => a.label.localeCompare(b.label))
})

const salidasByDate = computed(() => {
  const map = new Map<string, AgenciaVentaSalida[]>()
  salidas.value.forEach((salida) => {
    const key = normalizeDateKey(salida.fecha_salida)
    if (!key) return
    if (!map.has(key)) map.set(key, [])
    map.get(key)?.push(salida)
  })
  return map
})

const calendarCells = computed(() => {
  if (viewMode.value !== 'mes') return []
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth() + 1
  const daysInMonth = new Date(year, month, 0).getDate()
  const firstDay = new Date(year, month - 1, 1)
  const offset = (firstDay.getDay() + 6) % 7

  const cells: any[] = []
  for (let i = 0; i < offset; i += 1) {
    cells.push(null)
  }
  for (let day = 1; day <= daysInMonth; day += 1) {
    const dateKey = formatDateKey(new Date(year, month - 1, day))
    const items = salidasByDate.value.get(dateKey) || []
    cells.push({
      date: dateKey,
      day,
      items,
      stats: buildDayStats(items)
    })
  }
  return cells
})

const weekDates = computed(() => {
  if (viewMode.value !== 'semana') return []
  const start = weekStart.value
  return Array.from({ length: 7 }, (_, index) => {
    const date = new Date(start.getFullYear(), start.getMonth(), start.getDate() + index)
    const key = formatDateKey(date)
    const items = salidasByDate.value.get(key) || []
    return {
      date: key,
      label: `${weekDays[index]} ${date.getDate()}`,
      items,
      count: items.length
    }
  })
})

const buildDayStats = (items: AgenciaVentaSalida[]) => {
  const stats = { total: items.length, confirmadas: 0, pendientes: 0, llenas: 0 }
  items.forEach((salida) => {
    const total = cuposActuales(salida)
    if (total >= salida.cupo_maximo) stats.llenas += 1
    else if (total >= salida.cupo_minimo) stats.confirmadas += 1
    else stats.pendientes += 1
  })
  return stats
}

const dayCellClass = (cell: any) => {
  if (!cell || cell.stats.total === 0) return 'bg-gray-50 text-gray-400 border-gray-100'
  if (cell.stats.llenas > 0) return 'bg-red-50 text-red-700 border-red-100'
  if (cell.stats.pendientes > 0) return 'bg-amber-50 text-amber-700 border-amber-100'
  return 'bg-emerald-50 text-emerald-700 border-emerald-100'
}

const cuposActuales = (salida: { cupos_confirmados: number; cupos_reservados: number }) => {
  return Number(salida.cupos_confirmados || 0) + Number(salida.cupos_reservados || 0)
}

const cuposDisponibles = (salida: { cupo_maximo: number; cupos_confirmados: number; cupos_reservados: number }) => {
  return Math.max(0, Number(salida.cupo_maximo || 0) - cuposActuales(salida))
}

const tipoSalidaLabel = (value?: string | null) => {
  const t = String(value || '').toLowerCase()
  if (t === 'privado') return 'Privado'
  if (t === 'compartido') return 'Compartido'
  return value ? String(value) : 'N/D'
}

const estadoSalidaLabel = (value?: string | null) => {
  const e = String(value || '').toLowerCase()
  const map: Record<string, string> = {
    pendiente: 'Pendiente',
    activa: 'Activa',
    completada: 'Completada',
    cancelada: 'Cancelada'
  }
  return map[e] || (value ? String(value) : 'N/D')
}

const estadoSalidaSeverity = (value?: string | null) => {
  const e = String(value || '').toLowerCase()
  const map: Record<string, any> = {
    pendiente: 'warning',
    activa: 'success',
    completada: 'secondary',
    cancelada: 'danger'
  }
  return map[e] || 'secondary'
}

const cupoEstadoLabel = (salida: AgenciaVentaSalida) => {
  const total = cuposActuales(salida)
  if (total >= salida.cupo_maximo) return 'Capacidad maxima'
  if (total >= salida.cupo_minimo) return 'Cupo minimo OK'
  return 'Bajo cupo'
}

const cupoEstadoSeverity = (salida: AgenciaVentaSalida) => {
  const total = cuposActuales(salida)
  if (total >= salida.cupo_maximo) return 'danger'
  if (total >= salida.cupo_minimo) return 'success'
  return 'warning'
}

const formatFecha = (value?: string) => {
  if (!value) return ''
  const clean = normalizeDateKey(value)
  const match = clean.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`
  return clean
}

const formatFechaHora = (value?: string | null) => {
  if (!value) return ''
  const date = new Date(String(value))
  if (Number.isNaN(date.getTime())) return String(value)
  const fecha = date.toLocaleDateString('es-BO')
  const hora = date.toLocaleTimeString('es-BO', { hour: '2-digit', minute: '2-digit' })
  return `${fecha} ${hora}`
}

const truncateText = (value: string, max: number) => {
  if (!value) return ''
  if (value.length <= max) return value
  return `${value.slice(0, max - 1)}…`
}

const formatMoney = (value: number) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const estadoPagoLabel = (value?: string | null) => {
  const e = String(value || '').toLowerCase()
  if (!e) return 'Sin pago'
  if (e === 'pendiente') return 'Pendiente'
  if (e === 'confirmado') return 'Confirmado'
  if (e === 'rechazado') return 'Rechazado'
  return value ? String(value) : 'Sin pago'
}

const estadoPagoSeverity = (value?: string | null) => {
  const e = String(value || '').toLowerCase()
  if (!e) return 'secondary'
  if (e === 'pendiente') return 'warning'
  if (e === 'confirmado') return 'success'
  if (e === 'rechazado') return 'danger'
  return 'secondary'
}

const alertasBajoCupo = computed(() => {
  const todayKey = formatDateKey(new Date())
  return salidas.value
    .filter((salida) => {
      if (salida.estado === 'cancelada' || salida.estado === 'completada') return false
      return cuposActuales(salida) < salida.cupo_minimo && normalizeDateKey(salida.fecha_salida) >= todayKey
    })
    .sort((a, b) => normalizeDateKey(a.fecha_salida).localeCompare(normalizeDateKey(b.fecha_salida)))
    .slice(0, 6)
})

const loadSalidas = async () => {
  loading.value = true
  error.value = null
  try {
    const id = await resolveAgenciaId()
    const params: any = {
      desde: formatDateKey(rangeStart.value),
      hasta: formatDateKey(rangeEnd.value)
    }
    if (paqueteFilter.value) params.paquete_id = paqueteFilter.value
    const resp: any = await getVentasSalidas(id, params)
    if (!resp?.success) {
      error.value = resp?.error?.message || 'No se pudieron cargar las salidas'
      return
    }
    salidas.value = resp.data?.salidas || []
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar las salidas'
  } finally {
    loading.value = false
  }
}

const loadPaquetes = async () => {
  try {
    const id = await resolveAgenciaId()
    const resp: any = await getPaquetes(id, { page: 1, limit: 200, status: 'activo' })
    if (resp?.success) {
      paquetesDisponibles.value = resp.data?.paquetes || []
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err?.data?.error?.message || err?.message || 'No se pudieron cargar los paquetes', life: 3000 })
  }
}

const openCrearSalida = async () => {
  dayDialog.value = false
  closeDetalle()
  await loadPaquetes()
  if (!crearForm.value.fecha) {
    crearForm.value.fecha = new Date()
  }
  crearDialog.value = true
}

const crearSalida = async () => {
  if (!crearForm.value.paquete_id || !crearForm.value.fecha) {
    toast.add({ severity: 'warn', summary: 'Datos incompletos', detail: 'Selecciona paquete y fecha', life: 2500 })
    return
  }
  crearLoading.value = true
  try {
    const id = await resolveAgenciaId()
    const payload = {
      fecha_salida: formatDateKey(crearForm.value.fecha),
      estado: crearForm.value.estado
    }
    const resp: any = await createPaqueteSalida(id, crearForm.value.paquete_id, payload)
    if (!resp?.success) {
      toast.add({ severity: 'error', summary: 'Error', detail: resp?.error?.message || 'No se pudo habilitar la salida', life: 3000 })
      return
    }
    toast.add({ severity: 'success', summary: 'Salida habilitada', detail: resp?.message || 'Salida creada', life: 2500 })
    crearDialog.value = false
    crearForm.value = {
      paquete_id: null,
      fecha: null,
      estado: 'pendiente'
    }
    await loadSalidas()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err?.data?.error?.message || err?.message || 'No se pudo habilitar la salida', life: 3000 })
  } finally {
    crearLoading.value = false
  }
}

const prevRange = () => {
  const date = currentDate.value
  if (viewMode.value === 'mes') {
    currentDate.value = new Date(date.getFullYear(), date.getMonth() - 1, 1)
  } else {
    currentDate.value = new Date(date.getFullYear(), date.getMonth(), date.getDate() - 7)
  }
}

const nextRange = () => {
  const date = currentDate.value
  if (viewMode.value === 'mes') {
    currentDate.value = new Date(date.getFullYear(), date.getMonth() + 1, 1)
  } else {
    currentDate.value = new Date(date.getFullYear(), date.getMonth(), date.getDate() + 7)
  }
}

const goToToday = () => {
  currentDate.value = new Date()
}

const dayDialog = ref(false)
const selectedDate = ref<string | null>(null)

const selectedDaySalidas = computed(() => {
  if (!selectedDate.value) return []
  return salidasByDate.value.get(selectedDate.value) || []
})

const selectedDateLabel = computed(() => (selectedDate.value ? formatFecha(selectedDate.value) : ''))

const closeDayPanel = () => {
  dayDialog.value = false
  selectedDate.value = null
}

const resetDetallePanels = () => {
  logisticaDialog.value = false
  cancelDialog.value = false
  comprobanteDialog.value = false
  selectedCompra.value = null
  cancelReason.value = ''
  logisticaSalida.value = null
  cancelSalida.value = null
}

const closeDetalle = () => {
  detalleDialog.value = false
  detalleBaseSalida.value = null
  detalleSalida.value = null
  detalleCompras.value = []
  detalleError.value = null
  detalleTotalesServer.value = null
  resetDetallePanels()
}

const openDay = (dateKey: string) => {
  closeDetalle()
  crearDialog.value = false
  selectedDate.value = dateKey
  dayDialog.value = true
}

const detalleDialog = ref(false)
const detalleLoading = ref(false)
const detalleError = ref<string | null>(null)
const detalleBaseSalida = ref<AgenciaVentaSalida | null>(null)
const detalleSalida = ref<(VentaSalidaDetalle & { paquete_nombre?: string }) | null>(null)
const detalleCompras = ref<AgenciaVentaSalidaCompra[]>([])
const detalleTotalesServer = ref<{ ingresos: number; participantes: number; compras: number } | null>(null)

const detalleTotales = computed(() => {
  if (detalleTotalesServer.value) return detalleTotalesServer.value
  const ingresos = detalleCompras.value.reduce((acc, compra) => acc + Number(compra.precio_total || 0), 0)
  const participantes = detalleCompras.value.reduce((acc, compra) => acc + Number(compra.total_participantes || 0), 0)
  return { ingresos, participantes, compras: detalleCompras.value.length }
})

const loadDetalle = async (salida: AgenciaVentaSalida) => {
  detalleLoading.value = true
  detalleError.value = null
  try {
    const id = await resolveAgenciaId()
    const resp: any = await getVentasSalidaCompras(id, salida.salida_id)
    if (!resp?.success) {
      detalleError.value = resp?.error?.message || 'No se pudo cargar el detalle'
      return
    }
    detalleSalida.value = {
      ...(resp.data?.salida || {}),
      paquete_nombre: salida.paquete_nombre
    }
    detalleCompras.value = resp.data?.compras || []
    detalleTotalesServer.value = resp.data?.totales || null
  } catch (err: any) {
    detalleError.value = err?.data?.error?.message || err?.message || 'No se pudo cargar el detalle'
  } finally {
    detalleLoading.value = false
  }
}

const openDetalle = async (salida: AgenciaVentaSalida) => {
  dayDialog.value = false
  crearDialog.value = false
  resetDetallePanels()
  detalleBaseSalida.value = salida
  detalleDialog.value = true
  await loadDetalle(salida)
}

const refreshDetalle = async () => {
  if (!detalleBaseSalida.value) return
  await loadDetalle(detalleBaseSalida.value)
}

const openLogisticaDetalle = async () => {
  if (!detalleBaseSalida.value) return
  await openLogistica(detalleBaseSalida.value)
}

const openCancelarDetalle = async () => {
  if (!detalleBaseSalida.value) return
  await openCancelar(detalleBaseSalida.value)
}

const comprobanteDialog = ref(false)
const selectedCompra = ref<AgenciaVentaSalidaCompra | null>(null)

const openComprobante = (compra: AgenciaVentaSalidaCompra) => {
  selectedCompra.value = compra
  comprobanteDialog.value = true
  logisticaDialog.value = false
  cancelDialog.value = false
}

const logisticaDialog = ref(false)
const logisticaLoading = ref(false)
const logisticaSalida = ref<AgenciaVentaSalida | null>(null)

const logisticaForm = ref({
  punto_encuentro: '',
  hora_encuentro: '',
  notas_logistica: '',
  instrucciones_turistas: '',
  guia_nombre: '',
  guia_telefono: ''
})

const fillLogisticaForm = (salida: VentaSalidaDetalle | null) => {
  logisticaForm.value = {
    punto_encuentro: salida?.punto_encuentro || '',
    hora_encuentro: salida?.hora_encuentro || '',
    notas_logistica: salida?.notas_logistica || '',
    instrucciones_turistas: salida?.instrucciones_turistas || '',
    guia_nombre: salida?.guia_nombre || '',
    guia_telefono: salida?.guia_telefono || ''
  }
}

const openLogistica = async (salida: AgenciaVentaSalida) => {
  dayDialog.value = false
  crearDialog.value = false
  resetDetallePanels()
  detalleBaseSalida.value = salida
  detalleDialog.value = true
  logisticaSalida.value = salida
  await loadDetalle(salida)
  fillLogisticaForm(detalleSalida.value as VentaSalidaDetalle)
  logisticaDialog.value = true
}

const guardarLogistica = async () => {
  if (!logisticaSalida.value) return
  logisticaLoading.value = true
  try {
    const id = await resolveAgenciaId()
    const payload = {
      punto_encuentro: logisticaForm.value.punto_encuentro,
      hora_encuentro: logisticaForm.value.hora_encuentro,
      notas_logistica: logisticaForm.value.notas_logistica,
      instrucciones_turistas: logisticaForm.value.instrucciones_turistas,
      guia_nombre: logisticaForm.value.guia_nombre,
      guia_telefono: logisticaForm.value.guia_telefono
    }
    const resp: any = await updatePaqueteSalida(id, logisticaSalida.value.paquete_id, logisticaSalida.value.salida_id, payload)
    if (!resp?.success) {
      toast.add({ severity: 'error', summary: 'Error', detail: resp?.error?.message || 'No se pudo actualizar', life: 3000 })
      return
    }
    toast.add({ severity: 'success', summary: 'Actualizado', detail: 'Logistica guardada', life: 2500 })
    logisticaDialog.value = false
    await Promise.all([loadSalidas(), refreshDetalle()])
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err?.data?.error?.message || err?.message || 'No se pudo actualizar', life: 3000 })
  } finally {
    logisticaLoading.value = false
  }
}

const cancelDialog = ref(false)
const cancelLoading = ref(false)
const cancelReason = ref('')
const cancelSalida = ref<AgenciaVentaSalida | null>(null)

const openCancelar = async (salida: AgenciaVentaSalida) => {
  dayDialog.value = false
  crearDialog.value = false
  resetDetallePanels()
  detalleBaseSalida.value = salida
  detalleDialog.value = true
  cancelSalida.value = salida
  cancelReason.value = ''
  cancelDialog.value = true
  await loadDetalle(salida)
}

const confirmarCancelacion = async () => {
  if (!cancelSalida.value) return
  const reason = cancelReason.value.trim()
  if (!reason) {
    toast.add({ severity: 'warn', summary: 'Razon requerida', detail: 'Ingresa la razon de cancelacion', life: 2500 })
    return
  }
  cancelLoading.value = true
  try {
    const id = await resolveAgenciaId()
    const payload = { estado: 'cancelada', razon_cancelacion: reason }
    const resp: any = await updatePaqueteSalida(id, cancelSalida.value.paquete_id, cancelSalida.value.salida_id, payload)
    if (!resp?.success) {
      toast.add({ severity: 'error', summary: 'Error', detail: resp?.error?.message || 'No se pudo cancelar', life: 3000 })
      return
    }
    toast.add({ severity: 'success', summary: 'Salida cancelada', detail: 'Se actualizo el estado', life: 2500 })
    cancelDialog.value = false
    await Promise.all([loadSalidas(), refreshDetalle()])
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err?.data?.error?.message || err?.message || 'No se pudo cancelar', life: 3000 })
  } finally {
    cancelLoading.value = false
  }
}

const exportICal = () => {
  const items = salidas.value.filter((salida) => salida.estado !== 'cancelada')
  if (items.length === 0) {
    toast.add({ severity: 'warn', summary: 'Sin salidas', detail: 'No hay salidas para exportar', life: 2500 })
    return
  }
  const now = new Date()
  const dtstamp = now.toISOString().replace(/[-:]/g, '').split('.')[0] + 'Z'
  const lines = [
    'BEGIN:VCALENDAR',
    'VERSION:2.0',
    'PRODID:-//ANDARIA//Calendario de salidas//ES',
    'CALSCALE:GREGORIAN'
  ]
  items.forEach((salida) => {
    const dateKey = normalizeDateKey(salida.fecha_salida)
    const dtstart = dateKey.replace(/-/g, '')
    const summary = `Salida: ${salida.paquete_nombre} (${tipoSalidaLabel(salida.tipo_salida)})`
    lines.push('BEGIN:VEVENT')
    lines.push(`UID:salida-${salida.salida_id}@andaria`)
    lines.push(`DTSTAMP:${dtstamp}`)
    lines.push(`DTSTART;VALUE=DATE:${dtstart}`)
    lines.push(`SUMMARY:${summary.replace(/,/g, '\\,')}`)
    lines.push('END:VEVENT')
  })
  lines.push('END:VCALENDAR')
  const content = lines.join('\r\n')
  if (!process.client) return
  const blob = new Blob([content], { type: 'text/calendar;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `salidas_${formatDateKey(rangeStart.value)}_${formatDateKey(rangeEnd.value)}.ics`
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

watch([viewMode, currentDate, paqueteFilter], () => {
  loadSalidas()
})

onMounted(() => {
  loadSalidas()
})
</script>
