<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between flex-wrap gap-3">
          <div>
            <p class="text-sm text-gray-500">Panel de agencia</p>
            <h1 class="text-3xl font-bold text-gray-900">Nuevo paquete turístico</h1>
            <p class="muted mt-1">Crea una nueva plantilla de paquete para tu agencia.</p>
          </div>
          <div class="flex gap-2">
            <Button label="Volver" icon="pi pi-arrow-left" outlined @click="navigateTo('/agencia/paquetes')" />
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Card class="surface-card">
        <template #content>
          <div v-if="loadingAgencia" class="space-y-3">
            <Skeleton height="2rem" width="50%" />
            <Skeleton height="10rem" />
          </div>

          <div v-else-if="!agencia" class="space-y-3">
            <Message severity="warn" :closable="false">No se pudo cargar tu agencia.</Message>
            <Button label="Reintentar" icon="pi pi-refresh" outlined :loading="loadingAgencia" @click="loadAgencia" />
          </div>

          <form v-else class="space-y-6" @submit.prevent="handleSubmit">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Nombre</label>
                <InputText v-model="form.nombre" class="w-full" placeholder="Ej: City Tour Tarija" />
              </div>

              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                <Textarea
                  v-model="form.descripcion"
                  class="w-full"
                  rows="4"
                  autoResize
                  placeholder="Describe el paquete, experiencia, recomendaciones..."
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Estado</label>
                <Select
                  v-model="form.status"
                  :options="statusOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full"
                />
              </div>

              <div class="flex items-center gap-3 pt-6">
                <InputSwitch v-model="form.visible_publico" />
                <div>
                  <p class="text-sm font-medium text-gray-700">Visible al público</p>
                  <p class="text-xs text-gray-500">Si está desactivado, solo tú lo verás.</p>
                </div>
              </div>
            </div>

            <Divider />

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Frecuencia</label>
                <SelectButton
                  v-model="form.frecuencia"
                  :options="frecuenciaOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full"
                />
                <small class="text-gray-500">
                  <span v-if="form.frecuencia === 'salida_diaria'">El turista elige la fecha al comprar.</span>
                  <span v-else>Existe una única fecha (fija) para la salida.</span>
                </small>
              </div>

              <div v-if="form.frecuencia === 'salida_unica'" class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Fecha de salida fija</label>
                <DatePicker v-model="form.fecha_salida_fija" class="w-full" dateFormat="yy-mm-dd" />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Duración (días)</label>
                <InputNumber v-model="form.duracion_dias" class="w-full" :min="1" :max="30" />
                <small class="text-gray-500" v-if="isMultiDay">Noches: {{ duracionNoches }}</small>
                <small class="text-gray-500" v-else>Paquete de un día.</small>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Días previos de compra</label>
                <InputNumber v-model="form.dias_previos_compra" class="w-full" :min="1" :max="60" />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nivel de dificultad</label>
                <Select
                  v-model="form.nivel_dificultad"
                  :options="dificultadOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full"
                  showClear
                  placeholder="(Opcional)"
                />
              </div>
            </div>

            <div v-if="!isMultiDay" class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Horario</label>
                <Select
                  v-model="form.horario"
                  :options="horarioOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full"
                  showClear
                  placeholder="(Opcional)"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Hora de salida</label>
                <DatePicker
                  v-model="horaSalidaPickerValue"
                  timeOnly
                  hourFormat="24"
                  iconDisplay="input"
                  showIcon
                  icon="pi pi-clock"
                  fluid
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Duración (horas)</label>
                <InputNumber v-model="form.duracion_horas_num" class="w-full" :min="1" :max="24" />
              </div>
            </div>

            <Divider />

            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Cupo mínimo</label>
                <InputNumber v-model="form.cupo_minimo" class="w-full" :min="1" :max="500" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Cupo máximo</label>
                <InputNumber v-model="form.cupo_maximo" class="w-full" :min="1" :max="500" />
              </div>
              <div class="md:col-span-2 flex items-center gap-3 pt-6">
                <Checkbox v-model="form.permite_privado" binary inputId="permitePrivado" />
                <label for="permitePrivado" class="text-sm text-gray-700">Permite contratar como privado</label>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Precio base (nacionales)</label>
                <InputNumber v-model="form.precio_base_nacionales" class="w-full" :min="0" :max="999999" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Adicional (extranjeros)</label>
                <InputNumber v-model="form.precio_adicional_extranjeros" class="w-full" :min="0" :max="999999" />
              </div>
            </div>

            <Divider />

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <StringListEditor
                v-model="form.incluye"
                label="Incluye"
                placeholder="Ej: Transporte"
                emptyText="Sin elementos"
              />
              <StringListEditor
                v-model="form.no_incluye"
                label="No incluye"
                placeholder="Ej: Almuerzo"
                emptyText="Sin elementos"
              />
              <StringListEditor
                v-model="form.que_llevar"
                label="Qué llevar"
                placeholder="Ej: Bloqueador solar"
                emptyText="Sin elementos"
              />
            </div>

            <div class="flex justify-end gap-2 pt-2">
              <Button label="Cancelar" severity="secondary" outlined @click="navigateTo('/agencia/paquetes')" />
              <Button label="Crear paquete" icon="pi pi-check" type="submit" :loading="saving" />
            </div>
          </form>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import StringListEditor from '~/components/paquetes/StringListEditor.vue'

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const { getMiAgencia } = useAgencias()
const { createPaquete } = usePaquetes()

const agencia = ref<any>(null)
const loadingAgencia = ref(true)

const DEFAULT_HORA_SALIDA = '08:30'

const form = ref({
  nombre: '',
  descripcion: '',
  frecuencia: 'salida_diaria',
  fecha_salida_fija: null as Date | null,
  duracion_dias: 1,
  dias_previos_compra: 1,
  nivel_dificultad: null as string | null,
  horario: null as string | null,
  hora_salida: DEFAULT_HORA_SALIDA,
  duracion_horas_num: 4,
  cupo_minimo: 1,
  cupo_maximo: 10,
  permite_privado: true,
  precio_base_nacionales: 0,
  precio_adicional_extranjeros: 0,
  incluye: [] as string[],
  no_incluye: [] as string[],
  que_llevar: [] as string[],
  status: 'borrador',
  visible_publico: false
})

const statusOptions = [
  { label: 'Borrador', value: 'borrador' },
  { label: 'Activo', value: 'activo' },
  { label: 'Inactivo', value: 'inactivo' }
]

const frecuenciaOptions = [
  { label: 'Salida diaria', value: 'salida_diaria' },
  { label: 'Salida única', value: 'salida_unica' }
]

const dificultadOptions = [
  { label: 'Fácil', value: 'facil' },
  { label: 'Medio', value: 'medio' },
  { label: 'Difícil', value: 'dificil' },
  { label: 'Extremo', value: 'extremo' }
]

const horarioOptions = [
  { label: 'Mañana', value: 'mañana' },
  { label: 'Tarde', value: 'tarde' },
  { label: 'Todo el día', value: 'todo_dia' }
]

const normalizeTimeString = (value?: string | null) => {
  if (!value) return ''
  const match = String(value).match(/(\d{1,2}):(\d{2})/)
  if (!match) return ''
  return `${match[1].padStart(2, '0')}:${match[2]}`
}

const timeStringToDate = (timeStr?: string | null) => {
  const normalized = normalizeTimeString(timeStr)
  if (!normalized) return null
  const [hours, minutes] = normalized.split(':').map(Number)
  if (Number.isNaN(hours) || Number.isNaN(minutes)) return null
  const date = new Date()
  date.setHours(hours, minutes, 0, 0)
  return date
}

const dateToTimeString = (date: Date | null) => {
  if (!date) return ''
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  return `${hh}:${mm}`
}

const timeStringToTimestamp = (timeStr?: string | null) => {
  const normalized = normalizeTimeString(timeStr)
  if (!normalized) return ''
  const [hours, minutes] = normalized.split(':').map(Number)
  const date = new Date()
  date.setHours(hours, minutes, 0, 0)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d} ${normalized}:00`
}

const isMultiDay = computed(() => Number(form.value.duracion_dias || 1) > 1)
const duracionNoches = computed(() => Math.max(0, Number(form.value.duracion_dias || 1) - 1))
const horaSalidaPickerValue = computed({
  get: () => timeStringToDate(form.value.hora_salida),
  set: (val: Date | null) => {
    form.value.hora_salida = dateToTimeString(val)
  }
})

watch(
  () => form.value.frecuencia,
  (next) => {
    if (next === 'salida_diaria') form.value.fecha_salida_fija = null
  }
)

watch(
  () => form.value.duracion_dias,
  (next) => {
    const safe = Number(next || 1)
    if (!Number.isFinite(safe) || safe < 1) form.value.duracion_dias = 1
    if (safe > 1) {
      form.value.horario = null
      form.value.hora_salida = ''
      return
    }
    if (!form.value.hora_salida) {
      form.value.hora_salida = DEFAULT_HORA_SALIDA
    }
  }
)

const agenciaId = computed(() => Number(agencia.value?.id || 0))

const loadAgencia = async () => {
  loadingAgencia.value = true
  try {
    const response: any = await getMiAgencia()
    if (response.success) {
      agencia.value = response.data
      return
    }
    agencia.value = null
  } catch {
    agencia.value = null
  } finally {
    loadingAgencia.value = false
  }
}

const saving = ref(false)

const formatDate = (date: Date) => {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const handleSubmit = async () => {
  const id = agenciaId.value
  if (!id) return

  const nombre = form.value.nombre.trim()
  if (!nombre) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'El nombre es obligatorio', life: 3000 })
    return
  }

  if (form.value.frecuencia === 'salida_unica' && !form.value.fecha_salida_fija) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'La fecha fija es obligatoria', life: 3000 })
    return
  }

  const cupoMin = Number(form.value.cupo_minimo || 0)
  const cupoMax = Number(form.value.cupo_maximo || 0)
  if (cupoMin <= 0 || cupoMax <= 0 || cupoMax < cupoMin) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'Revisa cupo mínimo y máximo', life: 3000 })
    return
  }

  const duracionDias = Math.max(1, Number(form.value.duracion_dias || 1))
  if (duracionDias > 1 && duracionNoches.value !== duracionDias - 1) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'Duración inválida', life: 3000 })
    return
  }

  const payload: any = {
    nombre,
    descripcion: form.value.descripcion?.trim() || '',
    frecuencia: form.value.frecuencia,
    duracion_dias: duracionDias,
    duracion_noches: duracionDias > 1 ? duracionDias - 1 : undefined,
    fecha_salida_fija: form.value.frecuencia === 'salida_unica' && form.value.fecha_salida_fija
      ? formatDate(form.value.fecha_salida_fija)
      : undefined,
    dias_previos_compra: Math.max(1, Number(form.value.dias_previos_compra || 1)),
    nivel_dificultad: form.value.nivel_dificultad || undefined,
    cupo_minimo: cupoMin,
    cupo_maximo: cupoMax,
    permite_privado: !!form.value.permite_privado,
    precio_base_nacionales: Number(form.value.precio_base_nacionales || 0),
    precio_adicional_extranjeros: Number(form.value.precio_adicional_extranjeros || 0),
    incluye: form.value.incluye || [],
    no_incluye: form.value.no_incluye || [],
    que_llevar: form.value.que_llevar || [],
    status: form.value.status,
    visible_publico: !!form.value.visible_publico
  }

  if (duracionDias <= 1) {
    payload.horario = form.value.horario || undefined
    const horaSalidaPayload = timeStringToTimestamp(form.value.hora_salida)
    payload.hora_salida = horaSalidaPayload || undefined
    payload.duracion_horas = form.value.duracion_horas_num ? `${Number(form.value.duracion_horas_num)} hours` : undefined
  }

  saving.value = true
  try {
    const response: any = await createPaquete(id, payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Creado', detail: 'Paquete creado', life: 2500 })
      await navigateTo(`/agencia/paquetes/${response.data.id}`)
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo crear el paquete',
      life: 3500
    })
  } finally {
    saving.value = false
  }
}

onMounted(loadAgencia)
</script>
