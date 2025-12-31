<template>
  <Card class="surface-card">
    <template #title>
      <div class="flex items-center gap-2">
        <i class="pi pi-shopping-cart text-emerald-600"></i>
        <span>Comprar paquete</span>
      </div>
    </template>

    <template #content>
      <div class="space-y-5">
        <Message v-if="error" severity="error" :closable="false">
          {{ error }}
        </Message>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">Fecha</label>
            <DatePicker
              v-model="fecha"
              class="w-full"
              :minDate="minDate"
              :disabled="isSalidaUnica"
              dateFormat="dd/mm/yy"
              showIcon
            />
            <p v-if="isSalidaUnica" class="text-xs text-gray-500">
              Este paquete tiene una salida única: {{ fixedDateLabel }}
            </p>
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">Tipo de compra</label>
            <SelectButton
              v-model="tipoCompra"
              class="w-full"
              :options="tipoCompraOptions"
              optionLabel="label"
              optionValue="value"
              :disabled="!canChooseTipoCompra"
            />
            <p v-if="!canChooseTipoCompra" class="text-xs text-gray-500">
              Solo disponible: {{ tipoCompraLabel }}
            </p>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <Checkbox v-model="extranjero" :binary="true" inputId="extranjero" />
          <label for="extranjero" class="text-sm text-gray-700">Soy extranjero</label>
          <span v-if="precioExtranjero > 0" class="text-xs text-gray-500">
            (+ Bs. {{ formatMoney(precioExtranjero) }} por persona)
          </span>
        </div>

        <Divider />

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">Adultos</label>
            <InputNumber v-model="cantidadAdultos" class="w-full" :min="1" :useGrouping="false" />
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">Niños</label>
            <SelectButton
              v-model="tieneNinos"
              class="w-full"
              :options="tieneNinosOptions"
              optionLabel="label"
              optionValue="value"
            />
            <p class="text-xs text-gray-500">
              Los niños pagan desde {{ edadMinimaPago }} años.
            </p>
          </div>
        </div>

        <div v-if="tieneNinos === 'si'" class="space-y-3">
          <div class="rounded-xl border border-gray-200 bg-white p-4 space-y-3">
            <div class="flex items-center justify-between gap-2">
              <p class="text-sm font-semibold text-gray-900">Edades de los niños</p>
              <Button
                label="Agregar niño"
                icon="pi pi-plus"
                outlined
                :disabled="loading"
                @click="addNino"
              />
            </div>

            <div v-for="(_, index) in edadesNinos" :key="`nino-${index}`" class="flex items-end gap-2">
              <div class="flex-1 space-y-2">
                <label class="block text-sm font-medium text-gray-700">Edad niño {{ index + 1 }}</label>
                <InputNumber
                  v-model="edadesNinos[index]"
                  class="w-full"
                  :min="0"
                  :max="17"
                  :useGrouping="false"
                  :minFractionDigits="0"
                  :maxFractionDigits="0"
                />
              </div>
              <Button
                icon="pi pi-trash"
                severity="danger"
                text
                rounded
                :disabled="loading"
                @click="removeNino(index)"
              />
            </div>

            <Message v-if="ninosSinEdad > 0" severity="warn" :closable="false">
              Completa la edad de {{ ninosSinEdad }} niño(s) para calcular si pagan o no.
            </Message>

            <Message v-else-if="ninosEdadInvalida > 0" severity="warn" :closable="false">
              Las edades deben estar entre 0 y 17 años.
            </Message>

            <div v-else class="text-sm text-gray-700">
              <span class="font-semibold">{{ totalNinos }}</span> niño(s):
              <span class="font-semibold">{{ cantidadNinosPagan }}</span> pagan ·
              <span class="font-semibold">{{ cantidadNinosGratis }}</span> gratis
            </div>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <Checkbox v-model="tieneDiscapacidad" :binary="true" inputId="discapacidad" />
          <label for="discapacidad" class="text-sm text-gray-700">Algún participante tiene discapacidad</label>
        </div>

        <div v-if="tieneDiscapacidad" class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Descripción</label>
          <Textarea v-model="descripcionDiscapacidad" class="w-full" rows="3" autoResize />
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Notas</label>
          <Textarea v-model="notasTurista" class="w-full" rows="3" autoResize />
        </div>

        <Message v-if="tipoCompra === 'compartido' && checkingSalida" severity="info" :closable="false">
          Verificando salidas habilitadas para la fecha seleccionada...
        </Message>

        <Message v-else-if="tipoCompra === 'compartido' && salidaCheckError" severity="warn" :closable="false">
          {{ salidaCheckError }}
        </Message>

        <Message
          v-else-if="tipoCompra === 'compartido' && salidaExists === false && !cumpleCupoMinimo"
          severity="warn"
          :closable="false"
        >
          Para habilitar la primera salida en esta fecha debes registrar al menos {{ cupoMinimo }} participantes.
        </Message>

        <Message
          v-else-if="tipoCompra === 'compartido' && salidaExists === false && cumpleCupoMinimo"
          severity="info"
          :closable="false"
        >
          No existe una salida habilitada para esta fecha. Se creará una nueva salida al registrar la compra.
        </Message>

        <Divider />

        <div class="rounded-xl border border-gray-200 bg-gray-50 p-4 space-y-2">
          <div class="flex items-center justify-between text-sm">
            <span class="text-gray-600">Total personas</span>
            <span class="font-semibold">{{ totalPersonas }}</span>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-gray-600">Personas que pagan</span>
            <span class="font-semibold">{{ personasPagan }}</span>
          </div>
          <div v-if="totalNinos > 0" class="flex items-center justify-between text-sm">
            <span class="text-gray-600">Niños</span>
            <span class="font-semibold">{{ totalNinos }} ({{ cantidadNinosPagan }} pagan, {{ cantidadNinosGratis }} gratis)</span>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-gray-600">Precio unitario</span>
            <span class="font-semibold">Bs. {{ formatMoney(precioUnitario) }}</span>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-gray-600">Subtotal</span>
            <span class="font-semibold">Bs. {{ formatMoney(subtotal) }}</span>
          </div>
          <div v-if="recargoPrivado > 0" class="flex items-center justify-between text-sm">
            <span class="text-gray-600">Recargo privado ({{ recargoPrivadoPct }}%)</span>
            <span class="font-semibold">Bs. {{ formatMoney(recargoPrivado) }}</span>
          </div>
          <div class="flex items-center justify-between text-base pt-2 border-t border-gray-200">
            <span class="font-semibold text-gray-900">Total estimado</span>
            <span class="font-bold text-emerald-700">Bs. {{ formatMoney(totalEstimado) }}</span>
          </div>
          <p class="text-xs text-gray-500 mt-2">
            El total final se confirma al registrar la compra.
          </p>
        </div>

        <div class="flex justify-end gap-2">
          <Button
            label="Volver"
            icon="pi pi-arrow-left"
            severity="secondary"
            outlined
            :disabled="loading"
            @click="navigateBack"
          />
          <Button
            label="Continuar a pago"
            icon="pi pi-check"
            :loading="loading"
            :disabled="disableSubmit"
            @click="submit"
          />
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import type { CrearCompraRequest } from '~/types/compra'

const props = defineProps<{
  paquete: any
}>()

const toast = useToast()
const { crearCompra } = useCompra()
const { getSalidas } = usePaquetesTuristicos()

const loading = ref(false)
const error = ref<string | null>(null)

const parseDateOnly = (value?: string | null) => {
  if (!value) return null
  const raw = String(value)
  const datePart = raw.split('T').shift() ?? raw
  const clean = datePart.split(' ').shift() ?? datePart
  const match = clean.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (!match) return null
  return new Date(Number(match[1]), Number(match[2]) - 1, Number(match[3]))
}

const isSalidaUnica = computed(() => props.paquete?.frecuencia === 'salida_unica' && !!props.paquete?.fecha_salida_fija)
const fixedDateLabel = computed(() => {
  const d = parseDateOnly(props.paquete?.fecha_salida_fija)
  return d ? formatDateBO(d) : ''
})

const diasPreviosCompra = computed(() => Number(props.paquete?.dias_previos_compra || 1))
const minDate = computed(() => {
  const d = new Date()
  d.setHours(0, 0, 0, 0)
  d.setDate(d.getDate() + Math.max(0, diasPreviosCompra.value))
  return d
})

const defaultFecha = computed(() => {
  if (isSalidaUnica.value) return parseDateOnly(props.paquete?.fecha_salida_fija)
  return minDate.value
})

const fecha = ref<Date | null>(defaultFecha.value)

const canChooseTipoCompra = computed(() => props.paquete?.frecuencia === 'salida_diaria' && !!props.paquete?.permite_privado)
const tipoCompraOptions = computed(() => {
  const base = [{ label: 'Compartido', value: 'compartido' }]
  if (canChooseTipoCompra.value) base.push({ label: 'Privado', value: 'privado' })
  return base
})

const tipoCompra = ref<'compartido' | 'privado'>(canChooseTipoCompra.value ? 'compartido' : 'compartido')

const tipoCompraLabel = computed(() => (canChooseTipoCompra.value ? 'Compartido/Privado' : 'Compartido'))

const extranjero = ref(false)
const cantidadAdultos = ref(1)
const tieneNinos = ref<'no' | 'si'>('no')
const tieneNinosOptions = [
  { label: 'No', value: 'no' },
  { label: 'Sí', value: 'si' }
]
const edadesNinos = ref<Array<number | null>>([])
const tieneDiscapacidad = ref(false)
const descripcionDiscapacidad = ref<string>('')
const notasTurista = ref<string>('')

const precioBase = computed(() => Number(props.paquete?.precio_base_nacionales || 0))
const precioExtranjero = computed(() => Number(props.paquete?.precio_adicional_extranjeros || 0))
const recargoPrivadoPct = computed(() => Number(props.paquete?.politicas?.recargo_privado_porcentaje || 0))
const edadMinimaPago = computed(() => {
  const value = props.paquete?.politicas?.edad_minima_pago
  if (value === null || value === undefined) return 6
  const parsed = Number(value)
  if (!Number.isFinite(parsed)) return 6
  return Math.max(0, parsed)
})

watch(tieneNinos, (value) => {
  if (value === 'si') {
    if (edadesNinos.value.length === 0) edadesNinos.value = [null]
    return
  }
  edadesNinos.value = []
})

const totalNinos = computed(() => (tieneNinos.value === 'si' ? edadesNinos.value.length : 0))
const ninosSinEdad = computed(() => (tieneNinos.value === 'si' ? edadesNinos.value.filter((edad) => edad === null).length : 0))
const ninosEdadInvalida = computed(() => {
  if (tieneNinos.value !== 'si') return 0
  return edadesNinos.value.filter((edad) => typeof edad === 'number' && (!Number.isFinite(edad) || edad < 0 || edad > 17)).length
})

const cantidadNinosPagan = computed(() => {
  if (tieneNinos.value !== 'si') return 0
  return edadesNinos.value.filter((edad) => typeof edad === 'number' && Number.isFinite(edad) && edad >= 0 && edad <= 17 && edad >= edadMinimaPago.value).length
})

const cantidadNinosGratis = computed(() => {
  if (tieneNinos.value !== 'si') return 0
  return edadesNinos.value.filter((edad) => typeof edad === 'number' && Number.isFinite(edad) && edad >= 0 && edad <= 17 && edad < edadMinimaPago.value).length
})

const addNino = () => {
  edadesNinos.value.push(null)
}

const removeNino = (index: number) => {
  edadesNinos.value.splice(index, 1)
  if (edadesNinos.value.length === 0) {
    tieneNinos.value = 'no'
  }
}

const precioUnitario = computed(() => precioBase.value + (extranjero.value ? precioExtranjero.value : 0))
const totalPersonas = computed(() => (cantidadAdultos.value || 0) + totalNinos.value)
const personasPagan = computed(() => (cantidadAdultos.value || 0) + (cantidadNinosPagan.value || 0))
const subtotal = computed(() => precioUnitario.value * personasPagan.value)
const recargoPrivado = computed(() => (tipoCompra.value === 'privado' ? subtotal.value * (recargoPrivadoPct.value / 100) : 0))
const totalEstimado = computed(() => subtotal.value + recargoPrivado.value)

const cupoMinimo = computed(() => {
  const n = Number(props.paquete?.cupo_minimo || 1)
  if (!Number.isFinite(n) || n < 1) return 1
  return Math.floor(n)
})

const salidaExists = ref<boolean | null>(null)
const checkingSalida = ref(false)
const salidaCheckError = ref<string | null>(null)

const requiereCupoMinimo = computed(() => tipoCompra.value === 'compartido' && salidaExists.value === false)
const cumpleCupoMinimo = computed(() => totalPersonas.value >= cupoMinimo.value)

const disableSubmit = computed(() => {
  if (loading.value) return true
  if (!fecha.value) return true
  if ((cantidadAdultos.value || 0) < 1) return true
  if (tieneNinos.value === 'si' && (ninosSinEdad.value > 0 || ninosEdadInvalida.value > 0)) return true
  if (tipoCompra.value === 'compartido' && checkingSalida.value) return true
  if (requiereCupoMinimo.value && !cumpleCupoMinimo.value) return true
  return false
})

const pad2 = (n: number) => String(n).padStart(2, '0')
const toDateOnly = (d: Date) => `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())}`

let salidaCheckSeq = 0
const checkSalidaHabilitada = async () => {
  salidaCheckError.value = null

  if (tipoCompra.value !== 'compartido') {
    salidaExists.value = null
    checkingSalida.value = false
    return
  }

  if (!fecha.value) {
    salidaExists.value = null
    checkingSalida.value = false
    return
  }

  const paqueteId = Number(props.paquete?.id)
  if (!Number.isFinite(paqueteId) || paqueteId <= 0) {
    salidaExists.value = null
    checkingSalida.value = false
    return
  }

  const seq = ++salidaCheckSeq
  checkingSalida.value = true
  salidaExists.value = null

  try {
    const response: any = await getSalidas(paqueteId, { fecha: toDateOnly(fecha.value), tipo: 'compartido' })
    if (seq !== salidaCheckSeq) return

    if (response?.success) {
      const salidas = response.data?.salidas
      salidaExists.value = Array.isArray(salidas) && salidas.length > 0
      return
    }

    salidaCheckError.value = response?.error?.message || 'No se pudo verificar salidas habilitadas'
  } catch (err: any) {
    if (seq !== salidaCheckSeq) return
    salidaCheckError.value = err?.data?.error?.message || err?.message || 'No se pudo verificar salidas habilitadas'
  } finally {
    if (seq === salidaCheckSeq) checkingSalida.value = false
  }
}

watch(
  () => [fecha.value ? toDateOnly(fecha.value) : null, tipoCompra.value],
  () => {
    void checkSalidaHabilitada()
  },
  { immediate: true }
)

const formatDateBO = (d: Date) => `${pad2(d.getDate())}/${pad2(d.getMonth() + 1)}/${d.getFullYear()}`

const formatMoney = (value: number) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const navigateBack = () => {
  navigateTo(`/turista/paquetes/${props.paquete?.id}`)
}

const submit = async () => {
  error.value = null

  if (!fecha.value) {
    error.value = 'Debe seleccionar una fecha'
    return
  }

  if (tieneNinos.value === 'si') {
    if (ninosSinEdad.value > 0) {
      error.value = 'Complete la edad de todos los niños'
      return
    }

    if (ninosEdadInvalida.value > 0) {
      error.value = 'Las edades de los niños deben estar entre 0 y 17 años'
      return
    }
  }

  if (totalPersonas.value < 1) {
    error.value = 'Debe registrar al menos 1 participante'
    return
  }

  if (tipoCompra.value === 'compartido') {
    if (checkingSalida.value) {
      error.value = 'Espere un momento mientras verificamos las salidas disponibles'
      return
    }

    if (salidaExists.value === false && !cumpleCupoMinimo.value) {
      error.value = `Para habilitar la primera salida en esta fecha debe registrar al menos ${cupoMinimo.value} participantes`
      return
    }
  }

  loading.value = true
  try {
    const payload: CrearCompraRequest = {
      paquete_id: Number(props.paquete?.id),
      fecha_seleccionada: toDateOnly(fecha.value),
      tipo_compra: tipoCompra.value,
      extranjero: extranjero.value,
      cantidad_adultos: Number(cantidadAdultos.value || 0),
      cantidad_ninos_pagan: Number(cantidadNinosPagan.value || 0),
      cantidad_ninos_gratis: Number(cantidadNinosGratis.value || 0),
      tiene_discapacidad: tieneDiscapacidad.value,
      descripcion_discapacidad: tieneDiscapacidad.value ? (descripcionDiscapacidad.value || null) : null,
      notas_turista: notasTurista.value || null
    }

    console.log('📦 Payload a enviar:', JSON.stringify(payload, null, 2))
    const response: any = await crearCompra(payload)
    if (response?.success) {
      toast.add({ severity: 'success', summary: 'Compra creada', detail: response.message || 'Compra registrada', life: 3000 })
      const compraId = response.data?.compra_id
      if (compraId) {
        navigateTo(`/turista/compras/${compraId}/pagar`)
        return
      }
      error.value = 'Compra creada, pero no se pudo obtener el ID'
      return
    }

    error.value = response?.error?.message || 'No se pudo registrar la compra'
  } catch (err: any) {
    console.error('❌ Error al crear compra:', err)
    console.error('❌ Detalles del error:', JSON.stringify(err?.data || err, null, 2))
    error.value = err?.data?.error?.message || err?.message || 'No se pudo registrar la compra'
  } finally {
    loading.value = false
  }
}
</script>
