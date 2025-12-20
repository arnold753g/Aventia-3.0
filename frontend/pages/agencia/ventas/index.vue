<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Ventas de paquetes</h1>
          <p class="muted mt-1">Salidas habilitadas, cupos y confirmación de pagos.</p>
        </div>
        <div class="flex gap-2">
          <Button
            label="Actualizar"
            icon="pi pi-refresh"
            severity="secondary"
            outlined
            :loading="refreshingSalidas"
            @click="refresh"
          />
          <Button
            label="Volver"
            icon="pi pi-arrow-left"
            severity="secondary"
            outlined
            @click="navigateTo('/agencia/dashboard')"
          />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Card class="surface-card">
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <IconField class="md:col-span-2">
              <InputIcon class="pi pi-search" />
              <InputText v-model="searchTerm" placeholder="Buscar por nombre de paquete" class="w-full" />
            </IconField>

            <div class="flex items-center gap-3">
              <InputSwitch v-model="soloPendientes" />
              <div class="min-w-0">
                <p class="text-sm font-medium text-gray-900">Solo con pendientes</p>
                <p class="text-xs text-gray-500">Salidas con cupos por confirmar</p>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <InputSwitch v-model="mostrarHistorico" />
              <div class="min-w-0">
                <p class="text-sm font-medium text-gray-900">Mostrar histórico</p>
                <p class="text-xs text-gray-500">Incluye completadas/canceladas</p>
              </div>
            </div>
          </div>

          <Divider class="my-4" />

          <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
            <div class="flex flex-wrap gap-3 text-sm text-gray-700">
              <span><span class="font-semibold">{{ groupedPaquetes.length }}</span> paquete(s)</span>
              <span><span class="font-semibold">{{ salidasFiltered.length }}</span> salida(s)</span>
              <span><span class="font-semibold text-amber-700">{{ totalReservados }}</span> por confirmar</span>
              <span><span class="font-semibold text-emerald-700">{{ totalConfirmados }}</span> confirmados</span>
            </div>
            <div class="text-xs text-gray-500">
              Se actualiza automáticamente al confirmar/rechazar pagos.
            </div>
          </div>
        </template>
      </Card>

      <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

      <div v-if="loadingSalidas" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card v-for="i in 6" :key="i" class="surface-card">
          <template #content>
            <Skeleton height="1.5rem" width="70%" class="mb-3" />
            <Skeleton height="5rem" class="mb-2" />
            <Skeleton height="5rem" class="mb-2" />
            <Skeleton height="5rem" />
          </template>
        </Card>
      </div>

      <Card v-else-if="groupedPaquetes.length === 0" class="surface-card">
        <template #content>
          <div class="text-center py-12">
            <i class="pi pi-calendar-times text-6xl muted mb-4 block"></i>
            <p class="text-xl font-semibold muted">Aún no hay salidas con ventas</p>
            <p class="text-sm muted mt-2">
              Cuando un turista compre un paquete, aparecerán aquí las salidas y sus cupos.
            </p>
          </div>
        </template>
      </Card>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card v-for="grupo in groupedPaquetes" :key="grupo.paquete_id" class="surface-card">
          <template #title>
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <p class="text-base font-bold text-gray-900 truncate">{{ grupo.paquete_nombre }}</p>
                <p class="text-xs text-gray-500 mt-1">
                  {{ frecuenciaLabel(grupo.paquete_frecuencia) }}
                  <span v-if="grupo.paquete_duracion_dias && grupo.paquete_duracion_dias > 1">
                    • {{ grupo.paquete_duracion_dias }} días
                  </span>
                </p>
                <div class="mt-2 flex flex-wrap gap-2">
                  <Tag v-if="grupo.total_reservados > 0" :value="`Pendientes: ${grupo.total_reservados}`" severity="warning" />
                  <Tag v-if="grupo.total_confirmados > 0" :value="`Confirmados: ${grupo.total_confirmados}`" severity="success" />
                  <Tag :value="`${grupo.salidas.length} salida(s)`" severity="secondary" />
                </div>
              </div>
              <div class="flex flex-col items-end gap-1">
                <Tag :value="frecuenciaLabel(grupo.paquete_frecuencia)" severity="info" />
                <Tag
                  v-if="grupo.paquete_duracion_dias && grupo.paquete_duracion_dias > 1"
                  value="Multi-día"
                  severity="secondary"
                />
              </div>
            </div>
          </template>

          <template #content>
            <div class="space-y-3">
              <div
                v-for="salida in visibleSalidas(grupo)"
                :key="salida.salida_id"
                class="rounded-xl border border-gray-200 bg-white p-3"
                :class="{ 'border-amber-300 bg-amber-50/40': salida.cupos_reservados > 0 }"
              >
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0 flex items-start gap-3">
                    <div class="w-12 shrink-0 rounded-lg border border-gray-200 bg-gray-50 text-center py-1">
                      <div class="text-lg font-bold leading-5 text-gray-900">{{ fechaBadge(salida.fecha_salida).day }}</div>
                      <div class="text-[11px] uppercase text-gray-600 mt-1">
                        {{ fechaBadge(salida.fecha_salida).month }}
                      </div>
                    </div>

                    <div class="min-w-0">
                      <p class="font-semibold text-gray-900 truncate">{{ formatFecha(salida.fecha_salida) }}</p>
                      <p class="text-xs text-gray-600 mt-1">
                        {{ tipoSalidaLabel(salida.tipo_salida) }} • Cupo máx: {{ salida.cupo_maximo }}
                      </p>

                      <div class="mt-2 flex flex-wrap gap-2">
                        <Tag
                          v-if="salida.cupos_reservados > 0"
                          :value="`${salida.cupos_reservados} por confirmar`"
                          severity="warning"
                        />
                        <Tag v-if="cuposDisponibles(salida) === 0" value="Sin cupos" severity="danger" />
                      </div>
                    </div>
                  </div>
                  <div class="flex flex-col items-end gap-2">
                    <Tag :value="estadoSalidaLabel(salida.estado)" :severity="estadoSalidaSeverity(salida.estado)" />
                    <Button
                      label="Detalle"
                      icon="pi pi-eye"
                      size="small"
                      severity="secondary"
                      outlined
                      @click="openDetalle(salida)"
                    />
                  </div>
                </div>

                <div class="mt-3 space-y-2">
                  <MeterGroup :value="meterValues(salida)" :max="safeCupoMaximo(salida.cupo_maximo)" />
                  <div class="flex flex-wrap gap-x-3 gap-y-1 text-xs text-gray-600">
                    <span class="text-emerald-700 font-semibold">Confirmados: {{ salida.cupos_confirmados }}</span>
                    <span class="text-amber-700 font-semibold">Por confirmar: {{ salida.cupos_reservados }}</span>
                    <span class="text-gray-700 font-semibold">Restantes: {{ cuposDisponibles(salida) }}</span>
                  </div>
                </div>
              </div>

              <div v-if="grupo.salidas.length > 3" class="flex justify-end">
                <Button
                  :label="expandedLabel(grupo)"
                  size="small"
                  severity="secondary"
                  text
                  @click="togglePaqueteExpanded(grupo.paquete_id)"
                />
              </div>
            </div>
          </template>
        </Card>
      </div>

      <Dialog
        v-model:visible="detalleDialog"
        modal
        :dismissableMask="true"
        :style="{ width: '95vw', maxWidth: '980px' }"
      >
        <template #header>
          <div class="flex flex-col">
            <span class="font-bold text-gray-900">
              {{ selectedSalida ? selectedSalida.paquete_nombre : 'Detalle de salida' }}
            </span>
            <span v-if="selectedSalida" class="text-sm text-gray-600">
              {{ formatFecha(selectedSalida.fecha_salida) }} • {{ tipoSalidaLabel(selectedSalida.tipo_salida) }}
            </span>
          </div>
        </template>

        <div class="space-y-4">
          <Message v-if="detalleError" severity="error" :closable="false">{{ detalleError }}</Message>

          <div v-if="detalleLoading" class="flex items-center justify-center py-10">
            <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="4" />
          </div>

          <div v-else-if="salidaDetalle" class="space-y-6">
            <div class="rounded-xl border border-gray-200 bg-gray-50 p-4 space-y-2">
              <div class="flex flex-wrap items-center justify-between gap-3">
                <div class="space-y-1">
                  <p class="text-xs uppercase tracking-wider text-gray-500 font-semibold">Cupos</p>
                  <p class="text-sm text-gray-700">
                    Máx: <span class="font-semibold">{{ salidaDetalle.cupo_maximo }}</span>
                    • Mín: <span class="font-semibold">{{ salidaDetalle.cupo_minimo }}</span>
                    • Estado: <span class="font-semibold">{{ estadoSalidaLabel(salidaDetalle.estado) }}</span>
                  </p>
                </div>
                <div class="text-sm text-gray-700">
                  <span class="font-semibold text-amber-700">{{ salidaDetalle.cupos_reservados }}</span> por confirmar •
                  <span class="font-semibold text-emerald-700">{{ salidaDetalle.cupos_confirmados }}</span> confirmados •
                  <span class="font-semibold">{{ cuposDisponibles(salidaDetalle) }}</span> restantes
                </div>
              </div>
              <MeterGroup :value="meterValues(salidaDetalle)" :max="safeCupoMaximo(salidaDetalle.cupo_maximo)" />
            </div>

            <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
              <SelectButton v-model="detalleTab" :options="detalleTabOptions" optionLabel="label" optionValue="value" />
              <Button
                label="Actualizar detalle"
                icon="pi pi-refresh"
                size="small"
                severity="secondary"
                outlined
                :loading="detalleLoading"
                @click="loadDetalleSalida({ keepTab: true })"
              />
            </div>

            <div class="space-y-6">
              <div v-show="detalleTab === 'pendientes'">
                <div class="flex items-center justify-between gap-3 mb-2">
                  <h3 class="text-base font-bold text-gray-900">
                    Pendientes de confirmación ({{ comprasPendientes.length }})
                  </h3>
                </div>

                <DataTable
                  v-if="comprasPendientes.length"
                  :value="comprasPendientes"
                  dataKey="compra_id"
                  class="p-datatable-sm"
                  :rows="10"
                  paginator
                  responsiveLayout="scroll"
                >
                  <Column header="Turista">
                    <template #body="{ data }">
                      <div class="min-w-0">
                        <p class="font-semibold text-gray-900">
                          {{ data.turista_nombre }} {{ data.turista_apellido_paterno }} {{ data.turista_apellido_materno }}
                        </p>
                        <p class="text-xs text-gray-600">
                          {{ data.turista_phone || 'Sin teléfono' }} • {{ data.turista_email || 'Sin email' }}
                        </p>
                      </div>
                    </template>
                  </Column>

                  <Column header="Personas" style="width: 90px">
                    <template #body="{ data }">
                      <span class="font-semibold">{{ data.total_participantes }}</span>
                    </template>
                  </Column>

                  <Column header="Monto" style="width: 130px">
                    <template #body="{ data }">
                      <span class="font-bold text-emerald-700">Bs. {{ formatMoney(data.precio_total) }}</span>
                    </template>
                  </Column>

                  <Column header="Pago" style="width: 180px">
                    <template #body="{ data }">
                      <div class="space-y-1">
                        <p class="text-sm font-semibold text-gray-900">{{ metodoPagoLabel(data.metodo_pago) }}</p>
                        <p class="text-xs text-gray-600">{{ formatFechaHora(data.fecha_pago) }}</p>
                      </div>
                    </template>
                  </Column>

                  <Column header="Acciones" style="width: 240px">
                    <template #body="{ data }">
                      <div class="flex flex-wrap gap-2 justify-end">
                        <Button
                          v-if="data.comprobante_foto"
                          icon="pi pi-image"
                          label="Comprobante"
                          size="small"
                          severity="secondary"
                          outlined
                          @click="openComprobante(data)"
                        />
                        <Button icon="pi pi-check" label="Confirmar" size="small" @click="openConfirmar(data)" />
                        <Button
                          icon="pi pi-times"
                          label="Rechazar"
                          size="small"
                          severity="danger"
                          outlined
                          @click="openRechazar(data)"
                        />
                      </div>
                    </template>
                  </Column>
                </DataTable>

                <p v-else class="text-sm text-gray-600">No hay pagos pendientes en esta salida.</p>
              </div>

              <div v-show="detalleTab === 'confirmados'">
                <h3 class="text-base font-bold text-gray-900 mb-2">Confirmados ({{ comprasConfirmadas.length }})</h3>

                <DataTable
                  v-if="comprasConfirmadas.length"
                  :value="comprasConfirmadas"
                  dataKey="compra_id"
                  class="p-datatable-sm"
                  :rows="10"
                  paginator
                  responsiveLayout="scroll"
                >
                  <Column header="Turista">
                    <template #body="{ data }">
                      <p class="font-semibold text-gray-900">
                        {{ data.turista_nombre }} {{ data.turista_apellido_paterno }} {{ data.turista_apellido_materno }}
                      </p>
                    </template>
                  </Column>

                  <Column header="Personas" style="width: 90px">
                    <template #body="{ data }">
                      <span class="font-semibold">{{ data.total_participantes }}</span>
                    </template>
                  </Column>

                  <Column header="Monto" style="width: 130px">
                    <template #body="{ data }">
                      <span class="font-bold text-emerald-700">Bs. {{ formatMoney(data.precio_total) }}</span>
                    </template>
                  </Column>

                  <Column header="Confirmación" style="width: 200px">
                    <template #body="{ data }">
                      <p class="text-sm text-gray-700">{{ formatFechaHora(data.fecha_confirmacion) }}</p>
                    </template>
                  </Column>
                </DataTable>

                <p v-else class="text-sm text-gray-600">Aún no hay pagos confirmados.</p>
              </div>

              <div v-show="detalleTab === 'sin_pago'">
                <h3 class="text-base font-bold text-gray-900 mb-2">Sin pago registrado ({{ comprasSinPago.length }})</h3>
                <div v-if="comprasSinPago.length" class="space-y-2">
                  <div
                    v-for="c in comprasSinPago"
                    :key="c.compra_id"
                    class="rounded-xl border border-gray-200 bg-white p-3 flex items-center justify-between gap-3"
                  >
                    <div>
                      <p class="font-semibold text-gray-900">
                        {{ c.turista_nombre }} {{ c.turista_apellido_paterno }} {{ c.turista_apellido_materno }}
                      </p>
                      <p class="text-xs text-gray-600">
                        {{ c.total_participantes }} persona(s) • Bs. {{ formatMoney(c.precio_total) }}
                      </p>
                    </div>
                    <Tag value="Sin pago" severity="warning" />
                  </div>
                </div>
                <p v-else class="text-sm text-gray-600">No hay compras sin pago en esta salida.</p>
              </div>

              <div v-show="detalleTab === 'rechazados'">
                <h3 class="text-base font-bold text-gray-900 mb-2">Rechazados ({{ comprasRechazadas.length }})</h3>
                <div v-if="comprasRechazadas.length" class="space-y-2">
                  <div v-for="c in comprasRechazadas" :key="c.compra_id" class="rounded-xl border border-gray-200 bg-white p-3">
                    <div class="flex items-start justify-between gap-3">
                      <div>
                        <p class="font-semibold text-gray-900">
                          {{ c.turista_nombre }} {{ c.turista_apellido_paterno }} {{ c.turista_apellido_materno }}
                        </p>
                        <p v-if="c.razon_rechazo" class="text-xs text-red-600 mt-1">{{ c.razon_rechazo }}</p>
                      </div>
                      <Tag value="Rechazado" severity="danger" />
                    </div>
                  </div>
                </div>
                <p v-else class="text-sm text-gray-600">No hay pagos rechazados en esta salida.</p>
              </div>
            </div>
          </div>
        </div>
      </Dialog>

      <Dialog v-model:visible="comprobanteDialog" modal header="Comprobante" :style="{ width: '95vw', maxWidth: '640px' }">
        <div v-if="selectedCompra?.comprobante_foto" class="space-y-3">
          <img
            :src="resolveAssetUrl(selectedCompra.comprobante_foto)"
            alt="Comprobante"
            class="w-full h-auto rounded-lg border border-gray-200"
          />
        </div>
        <p v-else class="text-sm text-gray-600">No hay comprobante disponible.</p>
      </Dialog>

      <Dialog v-model:visible="confirmDialog" modal header="Confirmar pago" :style="{ width: '95vw', maxWidth: '560px' }">
        <div class="space-y-4">
          <p class="text-sm text-gray-700">
            Confirma el pago seleccionado. Esto moverá los cupos de <span class="font-semibold">reservados</span> a
            <span class="font-semibold">confirmados</span>.
          </p>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">Notas (opcional)</label>
            <Textarea v-model="notasEncargado" rows="3" class="w-full" placeholder="Ej: recibido en oficina..." />
          </div>

          <div class="flex justify-end gap-2">
            <Button label="Cancelar" severity="secondary" outlined @click="confirmDialog = false" />
            <Button label="Confirmar" icon="pi pi-check" :loading="confirming" @click="confirmar" />
          </div>
        </div>
      </Dialog>

      <Dialog v-model:visible="rejectDialog" modal header="Rechazar pago" :style="{ width: '95vw', maxWidth: '560px' }">
        <div class="space-y-4">
          <p class="text-sm text-gray-700">Indica la razón del rechazo. Esto liberará los cupos reservados.</p>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">Razón</label>
            <Textarea v-model="razonRechazo" rows="3" class="w-full" placeholder="Ej: comprobante ilegible" />
            <Message v-if="rejectError" severity="error" :closable="false">{{ rejectError }}</Message>
          </div>

          <div class="flex justify-end gap-2">
            <Button label="Cancelar" severity="secondary" outlined @click="rejectDialog = false" />
            <Button label="Rechazar" icon="pi pi-times" severity="danger" :loading="rejecting" @click="rechazar" />
          </div>
        </div>
      </Dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
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
const { confirmarPago, rechazarPago } = usePago()

const agenciaId = ref<number | null>(null)

const salidas = ref<AgenciaVentaSalida[]>([])
const loadingSalidas = ref(true)
const refreshingSalidas = ref(false)
const error = ref<string | null>(null)

const searchTerm = ref('')
const soloPendientes = ref(false)
const mostrarHistorico = ref(false)

const loadAgenciaId = async () => {
  if (agenciaId.value) return agenciaId.value
  const resp: any = await getMiAgencia()
  if (!resp?.success) throw new Error(resp?.error?.message || 'No se pudo cargar la agencia')
  agenciaId.value = Number(resp.data?.id)
  if (!agenciaId.value) throw new Error('No se pudo resolver el ID de tu agencia')
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

const formatMoney = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatFecha = (value?: any) => {
  if (!value) return ''
  const raw = String(value)
  const datePart = raw.split('T').shift() ?? raw
  const clean = datePart.split(' ').shift() ?? datePart
  const match = clean.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`
  return clean || raw
}

const formatFechaHora = (value?: any) => {
  if (!value) return ''
  const d = new Date(String(value))
  if (Number.isNaN(d.getTime())) return String(value)
  const date = d.toLocaleDateString('es-BO')
  const time = d.toLocaleTimeString('es-BO', { hour: '2-digit', minute: '2-digit' })
  return `${date} ${time}`
}

const frecuenciaLabel = (value?: string | null) => {
  const f = String(value || '').toLowerCase()
  if (f === 'salida_diaria') return 'Salida diaria'
  if (f === 'salida_unica') return 'Salida única'
  return value ? String(value) : 'N/D'
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

const metodoPagoLabel = (value?: string | null) => {
  const m = String(value || '').toLowerCase()
  const map: Record<string, string> = {
    efectivo: 'Efectivo',
    qr: 'QR',
    transferencia: 'Transferencia'
  }
  return map[m] || (value ? String(value) : 'Sin método')
}

const safeCupoMaximo = (value: any) => {
  const n = Number(value || 0)
  return n > 0 ? n : 1
}

const clamp = (value: number, min: number, max: number) => Math.min(max, Math.max(min, value))

const normalizeCupos = (salida: { cupo_maximo: number; cupos_reservados: number; cupos_confirmados: number }) => {
  const maxRaw = Number(salida.cupo_maximo || 0)
  const max = maxRaw > 0 ? maxRaw : 0
  const confirmadosRaw = Number(salida.cupos_confirmados || 0)
  const reservadosRaw = Number(salida.cupos_reservados || 0)

  const confirmados = clamp(confirmadosRaw, 0, max)
  const reservados = clamp(reservadosRaw, 0, Math.max(0, max - confirmados))
  const disponibles = Math.max(0, max - confirmados - reservados)

  return { max, confirmados, reservados, disponibles }
}

const cuposDisponibles = (salida: { cupo_maximo: number; cupos_reservados: number; cupos_confirmados: number }) => {
  return normalizeCupos(salida).disponibles
}

const meterValues = (salida: { cupo_maximo: number; cupos_reservados: number; cupos_confirmados: number }) => {
  const { confirmados, reservados, disponibles } = normalizeCupos(salida)
  return [
    { label: 'Confirmados', value: confirmados, color: '#22c55e' },
    { label: 'Por confirmar', value: reservados, color: '#f59e0b' },
    { label: 'Restantes', value: disponibles, color: '#9ca3af' }
  ]
}

const monthShortEs = ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic']

const fechaBadge = (value: string) => {
  const formatted = formatFecha(value)
  const [dayRaw, monthRaw] = formatted.split('/')
  const monthIndex = Number(monthRaw) - 1
  return {
    day: dayRaw || '',
    month: monthShortEs[monthIndex] || (monthRaw || '')
  }
}

type PaqueteGroup = {
  paquete_id: number
  paquete_nombre: string
  paquete_frecuencia: string
  paquete_duracion_dias: number | null
  paquete_horario: string | null
  salidas: AgenciaVentaSalida[]
  total_reservados: number
  total_confirmados: number
}

const salidasFiltered = computed(() => {
  const term = searchTerm.value.trim().toLowerCase()
  return salidas.value.filter((s) => {
    if (!mostrarHistorico.value) {
      const e = String(s.estado || '').toLowerCase()
      if (e === 'completada' || e === 'cancelada') return false
    }

    if (soloPendientes.value && Number(s.cupos_reservados || 0) <= 0) return false

    if (term) {
      const name = String(s.paquete_nombre || '').toLowerCase()
      if (!name.includes(term)) return false
    }

    return true
  })
})

const groupedPaquetes = computed<PaqueteGroup[]>(() => {
  const map = new Map<number, PaqueteGroup>()
  for (const s of salidasFiltered.value) {
    const reservados = Number(s.cupos_reservados || 0)
    const confirmados = Number(s.cupos_confirmados || 0)
    const existing = map.get(s.paquete_id)
    if (existing) {
      existing.salidas.push(s)
      existing.total_reservados += reservados
      existing.total_confirmados += confirmados
      continue
    }
    map.set(s.paquete_id, {
      paquete_id: s.paquete_id,
      paquete_nombre: s.paquete_nombre,
      paquete_frecuencia: s.paquete_frecuencia,
      paquete_duracion_dias: s.paquete_duracion_dias ?? null,
      paquete_horario: s.paquete_horario ?? null,
      salidas: [s],
      total_reservados: reservados,
      total_confirmados: confirmados
    })
  }

  const list = Array.from(map.values())
  for (const g of list) {
    g.salidas.sort((a, b) => {
      const da = String(a.fecha_salida)
      const db = String(b.fecha_salida)
      if (da === db) return String(a.tipo_salida).localeCompare(String(b.tipo_salida))
      return da.localeCompare(db)
    })
  }

  list.sort((a, b) => a.paquete_nombre.localeCompare(b.paquete_nombre))
  return list
})

const totalReservados = computed(() => salidasFiltered.value.reduce((acc, s) => acc + Number(s.cupos_reservados || 0), 0))
const totalConfirmados = computed(() => salidasFiltered.value.reduce((acc, s) => acc + Number(s.cupos_confirmados || 0), 0))

const expandedPaquetes = ref<Record<number, boolean>>({})

const isPaqueteExpanded = (paqueteId: number) => Boolean(expandedPaquetes.value[paqueteId])

const togglePaqueteExpanded = (paqueteId: number) => {
  expandedPaquetes.value = { ...expandedPaquetes.value, [paqueteId]: !expandedPaquetes.value[paqueteId] }
}

const visibleSalidas = (grupo: PaqueteGroup) => {
  if (grupo.salidas.length <= 3) return grupo.salidas
  if (isPaqueteExpanded(grupo.paquete_id)) return grupo.salidas
  return grupo.salidas.slice(0, 3)
}

const expandedLabel = (grupo: PaqueteGroup) => {
  if (grupo.salidas.length <= 3) return ''
  if (isPaqueteExpanded(grupo.paquete_id)) return 'Ocultar'
  return `Ver ${grupo.salidas.length - 3} más`
}

const loadSalidas = async (opts: { silent?: boolean } = {}) => {
  if (opts.silent) refreshingSalidas.value = true
  else loadingSalidas.value = true

  error.value = null
  try {
    const id = await loadAgenciaId()
    const response: any = await getVentasSalidas(id)

    if (!response?.success) {
      error.value = response?.error?.message || 'No se pudieron cargar las salidas'
      return
    }

    salidas.value = response.data?.salidas || []
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar las salidas'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loadingSalidas.value = false
    refreshingSalidas.value = false
  }
}

const refresh = async () => {
  await loadSalidas({ silent: true })
}

onMounted(() => loadSalidas())

const detalleDialog = ref(false)
const detalleLoading = ref(false)
const detalleError = ref<string | null>(null)

const selectedSalida = ref<AgenciaVentaSalida | null>(null)
const salidaDetalle = ref<VentaSalidaDetalle | null>(null)
const compras = ref<AgenciaVentaSalidaCompra[]>([])

const comprasPendientes = computed(() => compras.value.filter((c) => c.pago_id && c.estado === 'pendiente'))
const comprasConfirmadas = computed(() => compras.value.filter((c) => c.pago_id && c.estado === 'confirmado'))
const comprasRechazadas = computed(() => compras.value.filter((c) => c.pago_id && c.estado === 'rechazado'))
const comprasSinPago = computed(() => compras.value.filter((c) => !c.pago_id))

const detalleTab = ref<'pendientes' | 'confirmados' | 'sin_pago' | 'rechazados'>('pendientes')
const detalleTabOptions = computed(() => [
  { label: `Pendientes (${comprasPendientes.value.length})`, value: 'pendientes' },
  { label: `Confirmados (${comprasConfirmadas.value.length})`, value: 'confirmados' },
  { label: `Sin pago (${comprasSinPago.value.length})`, value: 'sin_pago' },
  { label: `Rechazados (${comprasRechazadas.value.length})`, value: 'rechazados' }
])

const pickBestTab = () => {
  if (comprasPendientes.value.length) return 'pendientes'
  if (comprasConfirmadas.value.length) return 'confirmados'
  if (comprasSinPago.value.length) return 'sin_pago'
  if (comprasRechazadas.value.length) return 'rechazados'
  return 'pendientes'
}

const loadDetalleSalida = async (opts: { keepTab?: boolean } = {}) => {
  if (!selectedSalida.value) return
  detalleLoading.value = true
  detalleError.value = null
  try {
    const id = await loadAgenciaId()
    const response: any = await getVentasSalidaCompras(id, selectedSalida.value.salida_id)
    if (!response?.success) {
      detalleError.value = response?.error?.message || 'No se pudo cargar el detalle de la salida'
      return
    }
    salidaDetalle.value = response.data?.salida || null
    compras.value = response.data?.compras || []
    if (!opts.keepTab) detalleTab.value = pickBestTab()
  } catch (err: any) {
    detalleError.value = err?.data?.error?.message || err?.message || 'No se pudo cargar el detalle de la salida'
    toast.add({ severity: 'error', summary: 'Error', detail: detalleError.value, life: 3000 })
  } finally {
    detalleLoading.value = false
  }
}

const openDetalle = async (salida: AgenciaVentaSalida) => {
  selectedSalida.value = salida
  detalleDialog.value = true
  await loadDetalleSalida()
}

const comprobanteDialog = ref(false)
const confirmDialog = ref(false)
const rejectDialog = ref(false)

const selectedCompra = ref<AgenciaVentaSalidaCompra | null>(null)
const notasEncargado = ref<string>('')
const razonRechazo = ref<string>('')
const rejectError = ref<string | null>(null)

const confirming = ref(false)
const rejecting = ref(false)

const openComprobante = (compra: AgenciaVentaSalidaCompra) => {
  selectedCompra.value = compra
  comprobanteDialog.value = true
}

const openConfirmar = (compra: AgenciaVentaSalidaCompra) => {
  selectedCompra.value = compra
  notasEncargado.value = ''
  confirmDialog.value = true
}

const openRechazar = (compra: AgenciaVentaSalidaCompra) => {
  selectedCompra.value = compra
  razonRechazo.value = ''
  rejectError.value = null
  rejectDialog.value = true
}

const confirmar = async () => {
  if (!selectedCompra.value?.pago_id) return
  confirming.value = true
  try {
    const notas = notasEncargado.value.trim()
    const payload = notas ? notas : undefined
    const resp: any = await confirmarPago(selectedCompra.value.pago_id, payload)
    if (!resp?.success) {
      toast.add({ severity: 'error', summary: 'Error', detail: resp?.error?.message || 'No se pudo confirmar el pago', life: 3500 })
      return
    }

    toast.add({ severity: 'success', summary: 'Confirmado', detail: resp?.message || 'Pago confirmado', life: 3000 })
    confirmDialog.value = false
    await Promise.all([loadDetalleSalida({ keepTab: true }), loadSalidas({ silent: true })])
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err?.data?.error?.message || err?.message || 'No se pudo confirmar el pago', life: 3500 })
  } finally {
    confirming.value = false
  }
}

const rechazar = async () => {
  if (!selectedCompra.value?.pago_id) return
  const reason = razonRechazo.value.trim()
  if (!reason) {
    rejectError.value = 'Debe ingresar una razón de rechazo.'
    return
  }

  rejecting.value = true
  rejectError.value = null
  try {
    const resp: any = await rechazarPago(selectedCompra.value.pago_id, reason)
    if (!resp?.success) {
      toast.add({ severity: 'error', summary: 'Error', detail: resp?.error?.message || 'No se pudo rechazar el pago', life: 3500 })
      return
    }

    toast.add({ severity: 'success', summary: 'Rechazado', detail: resp?.message || 'Pago rechazado', life: 3000 })
    rejectDialog.value = false
    await Promise.all([loadDetalleSalida({ keepTab: true }), loadSalidas({ silent: true })])
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err?.data?.error?.message || err?.message || 'No se pudo rechazar el pago', life: 3500 })
  } finally {
    rejecting.value = false
  }
}
</script>
