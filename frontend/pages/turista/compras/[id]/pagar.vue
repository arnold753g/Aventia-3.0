<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Pagar compra</h1>
          <p class="muted mt-1">Registra tu pago para que el encargado pueda confirmarlo.</p>
        </div>
        <div class="flex gap-2">
          <Button
            label="Mis compras"
            icon="pi pi-list"
            severity="secondary"
            outlined
            :disabled="loading"
            @click="navigateTo('/turista/mis-compras')"
          />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <div v-if="loading">
        <Skeleton height="220px" class="mb-4" />
        <Skeleton height="320px" />
      </div>

      <div v-else-if="error" class="text-center space-y-4 py-12">
        <i class="pi pi-exclamation-triangle text-5xl text-orange-500"></i>
        <h2 class="text-2xl font-bold text-gray-900">No se pudo cargar la compra</h2>
        <p class="muted">{{ error }}</p>
        <div class="flex justify-center gap-2">
          <Button label="Volver" icon="pi pi-arrow-left" severity="secondary" outlined @click="navigateTo('/turista/mis-compras')" />
          <Button label="Reintentar" icon="pi pi-refresh" @click="loadAll" />
        </div>
      </div>

      <div v-else-if="compra" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
          <Card class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-receipt text-emerald-600"></i>
                <span>Resumen</span>
              </div>
            </template>
            <template #content>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p class="text-xs text-gray-500">Paquete</p>
                  <p class="font-semibold text-gray-900">{{ compra.paquete?.nombre || `#${compra.paquete?.id}` }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Estado</p>
                  <Tag :value="statusLabel" :severity="statusSeverity" />
                </div>
                <div>
                  <p class="text-xs text-gray-500">Tipo de compra</p>
                  <p class="font-semibold text-gray-900">{{ tipoCompraLabel }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha de compra</p>
                  <p class="font-semibold text-gray-900">{{ formatFechaHora(compra.fecha_compra) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha seleccionada</p>
                  <p class="font-semibold text-gray-900">{{ formatFecha(compra.fecha_seleccionada) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Participantes</p>
                  <p class="font-semibold text-gray-900">{{ compra.total_participantes }}</p>
                </div>
                <div class="md:col-span-2">
                  <p class="text-xs text-gray-500">Total</p>
                  <p class="text-2xl font-bold text-emerald-700">Bs. {{ formatMoney(compra.precio_total) }}</p>
                </div>
              </div>
            </template>
          </Card>

          <Card v-if="paquete" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-map-marker text-sky-600"></i>
                <span>Detalle del paquete</span>
              </div>
            </template>
            <template #content>
              <div class="flex flex-col md:flex-row gap-4">
                <img
                  v-if="paqueteFotoUrl"
                  :src="paqueteFotoUrl"
                  alt="Foto del paquete"
                  class="w-full md:w-56 h-40 object-cover rounded-lg border border-gray-200"
                  loading="lazy"
                />
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 flex-1">
                  <div>
                    <p class="text-xs text-gray-500">Agencia</p>
                    <p class="font-semibold text-gray-900">{{ agenciaNombre }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Ubicacion</p>
                    <p class="font-semibold text-gray-900">{{ agenciaUbicacion }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Duracion</p>
                    <p class="font-semibold text-gray-900">{{ duracionLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Horario</p>
                    <p class="font-semibold text-gray-900">{{ horarioLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Hora salida</p>
                    <p class="font-semibold text-gray-900">{{ horaSalidaLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Dificultad</p>
                    <p class="font-semibold text-gray-900">{{ dificultadLabel }}</p>
                  </div>
                </div>
              </div>
            </template>
          </Card>

          <Message v-if="compra.status !== 'pendiente_confirmacion'" severity="info" :closable="false">
            Esta compra ya no admite registrar pagos. Estado: {{ statusLabel }}.
          </Message>

          <Message v-else-if="compra.ultimo_pago?.estado === 'pendiente'" severity="success" :closable="false">
            Ya registraste un pago pendiente para esta compra. Espera la confirmación del encargado.
          </Message>

          <SubirComprobante v-else :compra-id="compraId" :monto="Number(compra.precio_total || 0)" @pago-registrado="loadAll" />
        </div>

        <div class="space-y-6">
          <Card v-if="paquete" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-building text-blue-600"></i>
                <span>Datos de pago</span>
              </div>
            </template>
            <template #content>
              <div v-if="paquete.agencia_datos_pago" class="space-y-3">
                <div v-if="paquete.agencia_datos_pago.nombre_banco" class="text-sm">
                  <span class="text-gray-500">Banco:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ paquete.agencia_datos_pago.nombre_banco }}</span>
                </div>
                <div v-if="paquete.agencia_datos_pago.numero_cuenta" class="text-sm">
                  <span class="text-gray-500">Cuenta:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ paquete.agencia_datos_pago.numero_cuenta }}</span>
                </div>
                <div v-if="paquete.agencia_datos_pago.nombre_titular" class="text-sm">
                  <span class="text-gray-500">Titular:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ paquete.agencia_datos_pago.nombre_titular }}</span>
                </div>

                <div v-if="paquete.agencia_datos_pago.qr_pago_foto" class="pt-3">
                  <p class="text-xs text-gray-500 mb-2">QR</p>
                  <img
                    :src="resolveAssetUrl(paquete.agencia_datos_pago.qr_pago_foto)"
                    alt="QR de pago"
                    class="w-full rounded-lg border border-gray-200"
                    loading="lazy"
                  />
                </div>
              </div>

              <div v-else class="text-sm text-gray-600">
                La agencia no tiene datos de pago configurados.
              </div>
            </template>
          </Card>

          <Card v-if="compra?.ultimo_pago" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-receipt text-emerald-600"></i>
                <span>Pago registrado</span>
              </div>
            </template>
            <template #content>
              <div class="space-y-3">
                <div class="text-sm">
                  <span class="text-gray-500">Metodo:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ metodoPagoLabel }}</span>
                </div>
                <div class="text-sm">
                  <span class="text-gray-500">Monto:</span>
                  <span class="font-semibold text-gray-900 ml-2">Bs. {{ formatMoney(compra.ultimo_pago?.monto) }}</span>
                </div>
                <div class="text-sm flex items-center gap-2">
                  <span class="text-gray-500">Estado:</span>
                  <Tag :value="pagoStatusLabel" :severity="pagoStatusSeverity" />
                </div>
                <div v-if="compra.ultimo_pago?.fecha_confirmacion" class="text-sm">
                  <span class="text-gray-500">Confirmado:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ formatFechaHora(compra.ultimo_pago?.fecha_confirmacion) }}</span>
                </div>

                <div v-if="compra.ultimo_pago?.comprobante_foto" class="pt-3">
                  <p class="text-xs text-gray-500 mb-2">Comprobante</p>
                  <img
                    :src="resolveAssetUrl(compra.ultimo_pago?.comprobante_foto)"
                    alt="Comprobante"
                    class="w-full rounded-lg border border-gray-200"
                    loading="lazy"
                  />
                </div>
                <p v-else class="text-sm text-gray-600">No hay comprobante adjunto.</p>
              </div>
            </template>
          </Card>

          <Card class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-info-circle text-gray-600"></i>
                <span>Recordatorio</span>
              </div>
            </template>
            <template #content>
              <ul class="text-sm text-gray-700 space-y-2">
                <li>El monto debe coincidir exactamente con el total.</li>
                <li>El encargado confirmará o rechazará tu pago.</li>
                <li>Si se rechaza, se liberan cupos reservados.</li>
              </ul>
            </template>
          </Card>
        </div>
      </div>
    </div>
    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import SubirComprobante from '~/components/paquetes/SubirComprobante.vue'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const route = useRoute()
const config = useRuntimeConfig()
const assetsBase = String(config.public.apiBase || '').replace(/\/api\/v1\/?$/, '')

const { obtenerDetalleCompra } = useCompra()
const { getPaquete } = usePaquetesTuristicos()

const compraId = Number(route.params.id)

const compra = ref<any>(null)
const paquete = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const statusLabel = computed(() => {
  const map: Record<string, string> = {
    pendiente_confirmacion: 'Pendiente de confirmación',
    confirmada: 'Confirmada',
    rechazada: 'Rechazada',
    cancelada: 'Cancelada',
    completada: 'Completada'
  }
  return map[compra.value?.status] || compra.value?.status || 'N/D'
})

const statusSeverity = computed(() => {
  const map: Record<string, any> = {
    pendiente_confirmacion: 'warning',
    confirmada: 'success',
    rechazada: 'danger',
    cancelada: 'secondary',
    completada: 'info'
  }
  return map[compra.value?.status] || 'secondary'
})

const tipoCompraLabel = computed(() => {
  const map: Record<string, string> = {
    compartido: 'Compartido',
    privado: 'Privado'
  }
  return map[compra.value?.tipo_compra] || compra.value?.tipo_compra || 'N/D'
})

const metodoPagoLabel = computed(() => {
  const map: Record<string, string> = {
    efectivo: 'Efectivo',
    qr: 'QR',
    transferencia: 'Transferencia'
  }
  return map[compra.value?.ultimo_pago?.metodo_pago] || compra.value?.ultimo_pago?.metodo_pago || 'N/D'
})

const pagoStatusLabel = computed(() => {
  const map: Record<string, string> = {
    pendiente: 'Pendiente',
    confirmado: 'Confirmado',
    rechazado: 'Rechazado'
  }
  return map[compra.value?.ultimo_pago?.estado] || compra.value?.ultimo_pago?.estado || 'N/D'
})

const pagoStatusSeverity = computed(() => {
  const map: Record<string, any> = {
    pendiente: 'warning',
    confirmado: 'success',
    rechazado: 'danger'
  }
  return map[compra.value?.ultimo_pago?.estado] || 'secondary'
})

const frecuenciaLabel = computed(() => {
  const map: Record<string, string> = {
    salida_diaria: 'Salida diaria',
    salida_unica: 'Salida unica'
  }
  return map[paquete.value?.frecuencia] || paquete.value?.frecuencia || 'N/D'
})

const duracionLabel = computed(() => {
  const value = Number(paquete.value?.duracion_dias || 0)
  if (!value) return 'N/D'
  return value === 1 ? '1 dia' : `${value} dias`
})

const horarioLabel = computed(() => {
  const raw = paquete.value?.horario
  if (!raw) return 'N/D'
  const map: Record<string, string> = {
    manana: 'Manana',
    mañana: 'Manana',
    tarde: 'Tarde',
    todo_dia: 'Todo el dia'
  }
  return map[String(raw)] || String(raw)
})

const dificultadLabel = computed(() => {
  const map: Record<string, string> = {
    facil: 'Facil',
    medio: 'Medio',
    dificil: 'Dificil',
    extremo: 'Extremo'
  }
  return map[paquete.value?.nivel_dificultad] || paquete.value?.nivel_dificultad || 'N/D'
})

const horaSalidaLabel = computed(() => {
  const raw = paquete.value?.hora_salida
  if (!raw) return 'N/D'
  const value = String(raw).trim()
  const match = value.match(/(\d{1,2}):(\d{2})/)
  if (match) return `${match[1].padStart(2, '0')}:${match[2]}`
  const parsed = new Date(value)
  if (!Number.isNaN(parsed.getTime())) {
    const hh = String(parsed.getHours()).padStart(2, '0')
    const mm = String(parsed.getMinutes()).padStart(2, '0')
    return `${hh}:${mm}`
  }
  return value
})

const agenciaNombre = computed(() => paquete.value?.agencia?.nombre_comercial || 'N/D')
const agenciaUbicacion = computed(() => paquete.value?.agencia?.departamento?.nombre || paquete.value?.agencia?.direccion || 'N/D')

const paqueteFotoUrl = computed(() => {
  const fotos = (paquete.value?.fotos || []).slice()
  fotos.sort((a: any, b: any) => {
    if (!!a.es_principal !== !!b.es_principal) return a.es_principal ? -1 : 1
    return (a.orden || 0) - (b.orden || 0)
  })
  const principal = fotos[0]
  return principal?.foto ? resolveAssetUrl(principal.foto) : ''
})

const formatMoney = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatFechaHora = (value?: any) => {
  if (!value) return ''
  const raw = String(value)
  const parsed = new Date(raw)
  if (!Number.isNaN(parsed.getTime())) {
    const d = String(parsed.getDate()).padStart(2, '0')
    const m = String(parsed.getMonth() + 1).padStart(2, '0')
    const y = parsed.getFullYear()
    const hh = String(parsed.getHours()).padStart(2, '0')
    const mm = String(parsed.getMinutes()).padStart(2, '0')
    return `${d}/${m}/${y} ${hh}:${mm}`
  }
  return formatFecha(raw)
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

const resolveAssetUrl = (path?: string) => {
  if (!path) return ''
  let normalized = String(path).replace(/\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const loadAll = async () => {
  loading.value = true
  error.value = null
  try {
    const response: any = await obtenerDetalleCompra(compraId)
    if (!response.success) {
      error.value = response?.error?.message || 'No se pudo cargar la compra'
      return
    }
    compra.value = response.data

    const paqueteId = Number(compra.value?.paquete?.id)
    if (paqueteId) {
      const paqueteResp: any = await getPaquete(paqueteId)
      if (paqueteResp.success) {
        paquete.value = paqueteResp.data
      }
    }
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudo cargar la compra'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>
