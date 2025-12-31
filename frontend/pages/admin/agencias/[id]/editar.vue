<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-6xl mx-auto px-4 py-6">
        <div class="flex items-center gap-4">
          <Button icon="pi pi-arrow-left" text rounded @click="navigateTo(`/admin/agencias/${id}`)" />
          <div>
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Editar Agencia
            </h1>
            <p class="muted mt-1">{{ form.nombre_comercial || 'Actualiza la información principal' }}</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="max-w-6xl mx-auto px-4 py-8">
      <Skeleton height="400px" />
    </div>

    <div v-else class="max-w-6xl mx-auto px-4 py-8">
      <Card class="surface-card">
        <template #content>
          <form class="space-y-8" @submit.prevent="handleSubmit">
            <!-- Informacion basica -->
            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-info-circle"></i>
                Informacion basica
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Nombre comercial *</label>
                  <InputText v-model="form.nombre_comercial" class="w-full" />
                  <Message v-if="errors.nombre_comercial" severity="error" size="small" variant="simple">
                    {{ errors.nombre_comercial }}
                  </Message>
                </div>

                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Descripcion</label>
                  <Textarea v-model="form.descripcion" class="w-full" rows="4" autoResize />
                  <Message v-if="errors.descripcion" severity="error" size="small" variant="simple">
                    {{ errors.descripcion }}
                  </Message>
                </div>
              </div>
            </div>

            <!-- Ubicacion -->
            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-map-marker"></i>
                Ubicacion
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium muted mb-2">Departamento</label>
                  <Select
                    v-model="form.departamento_id"
                    :options="departamentos"
                    optionLabel="nombre"
                    optionValue="id"
                    placeholder="Seleccione..."
                    class="w-full"
                    showClear
                  />
                  <Message v-if="errors.departamento_id" severity="error" size="small" variant="simple">
                    {{ errors.departamento_id }}
                  </Message>
                </div>

                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Direccion</label>
                  <InputText v-model="form.direccion" class="w-full" />
                  <Message v-if="errors.direccion" severity="error" size="small" variant="simple">
                    {{ errors.direccion }}
                  </Message>
                </div>

                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Ubicacion GPS (opcional)</label>
                  <ClientOnly>
                    <AgenciaMap
                      v-model:latitud="form.latitud"
                      v-model:longitud="form.longitud"
                      :editable="true"
                      :show-coordinate-inputs="false"
                      height="400px"
                    />
                    <template #fallback>
                      <div class="flex items-center justify-center h-96 bg-gray-100 rounded">
                        <i class="pi pi-spin pi-spinner text-4xl text-gray-400"></i>
                      </div>
                    </template>
                  </ClientOnly>
                  <Message v-if="errors.latitud" severity="error" size="small" variant="simple" class="mt-2">
                    {{ errors.latitud }}
                  </Message>
                  <Message v-if="errors.longitud" severity="error" size="small" variant="simple" class="mt-2">
                    {{ errors.longitud }}
                  </Message>
                </div>
              </div>
            </div>

            <!-- Contacto -->
            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-phone"></i>
                Contacto y redes
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium muted mb-2">Telefono</label>
                  <div class="grid grid-cols-1 sm:grid-cols-[140px,1fr] gap-3">
                    <Select
                      v-model="form.phone_prefix"
                      :options="phonePrefixes"
                      optionLabel="label"
                      optionValue="value"
                      class="w-full"
                      disabled
                    />
                    <InputText
                      v-model="form.phone_number"
                      class="w-full"
                      placeholder="Ej: 21234567"
                      maxlength="8"
                      inputmode="numeric"
                    />
                  </div>
                  <Message v-if="errors.phone_prefix" severity="error" size="small" variant="simple">
                    {{ errors.phone_prefix }}
                  </Message>
                  <Message v-if="errors.phone_number" severity="error" size="small" variant="simple">
                    {{ errors.phone_number }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Email</label>
                  <InputText v-model="form.email" class="w-full" />
                  <Message v-if="errors.email" severity="error" size="small" variant="simple">
                    {{ errors.email }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Sitio web</label>
                  <InputText v-model="form.sitio_web" class="w-full" />
                  <Message v-if="errors.sitio_web" severity="error" size="small" variant="simple">
                    {{ errors.sitio_web }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Facebook</label>
                  <InputText v-model="form.facebook" class="w-full" />
                  <Message v-if="errors.facebook" severity="error" size="small" variant="simple">
                    {{ errors.facebook }}
                  </Message>
                </div>

                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Instagram</label>
                  <InputText v-model="form.instagram" class="w-full" />
                  <Message v-if="errors.instagram" severity="error" size="small" variant="simple">
                    {{ errors.instagram }}
                  </Message>
                </div>
              </div>
            </div>

            <!-- Horarios y dias -->
            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-clock"></i>
                Horarios y dias
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium muted mb-2">Horario apertura</label>
                  <DatePicker
                    v-model="aperturaPickerValue"
                    timeOnly
                    hourFormat="24"
                    iconDisplay="input"
                    showIcon
                    icon="pi pi-clock"
                    fluid
                  />
                  <Message v-if="errors.horario_apertura" severity="error" size="small" variant="simple">
                    {{ errors.horario_apertura }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Horario cierre</label>
                  <DatePicker
                    v-model="cierrePickerValue"
                    timeOnly
                    hourFormat="24"
                    iconDisplay="input"
                    showIcon
                    icon="pi pi-clock"
                    fluid
                  />
                  <Message v-if="errors.horario_cierre" severity="error" size="small" variant="simple">
                    {{ errors.horario_cierre }}
                  </Message>
                </div>

                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Dias de atencion</label>
                  <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
                    <div v-for="dia in dias" :key="dia.id" class="flex items-center gap-2 border rounded px-3 py-2">
                      <Checkbox :inputId="`dia-${dia.id}`" :value="dia.id" v-model="form.dias_ids" />
                      <label :for="`dia-${dia.id}`" class="text-sm muted">{{ dia.nombre }}</label>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Operacion -->
            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-wallet"></i>
                Operacion
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium muted mb-2">Metodos de pago</label>
                  <div class="flex flex-col gap-2">
                    <label class="flex items-center gap-2">
                      <Checkbox v-model="form.acepta_qr" binary inputId="qr" />
                      <span class="text-sm muted">QR</span>
                    </label>
                    <label class="flex items-center gap-2">
                      <Checkbox v-model="form.acepta_transferencia" binary inputId="trans" />
                      <span class="text-sm muted">Transferencia</span>
                    </label>
                    <label class="flex items-center gap-2">
                      <Checkbox v-model="form.acepta_efectivo" binary inputId="cash" />
                      <span class="text-sm muted">Efectivo</span>
                    </label>
                  </div>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Licencia turistica</label>
                  <div class="flex items-center gap-3">
                    <InputSwitch v-model="form.licencia_turistica" />
                    <span class="text-sm muted">{{ form.licencia_turistica ? 'Si' : 'No' }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Gestion -->
            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-cog"></i>
                Gestion
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Encargado principal</label>
                  <Select
                    v-model="form.encargado_principal_id"
                    :options="encargados"
                    optionLabel="nombre_completo"
                    optionValue="id"
                    placeholder="Seleccione..."
                    class="w-full"
                    showClear
                    :disabled="!authStore.isAdmin"
                  >
                    <template #option="slotProps">
                      <div class="flex flex-col">
                        <span class="font-semibold">{{ slotProps.option.nombre_completo }}</span>
                        <span class="text-xs muted">{{ slotProps.option.email }}</span>
                      </div>
                    </template>
                  </Select>
                  <Message v-if="errors.encargado_principal_id" severity="error" size="small" variant="simple">
                    {{ errors.encargado_principal_id }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Estado</label>
                  <Select
                    v-model="form.status"
                    :options="statusOptions"
                    optionLabel="label"
                    optionValue="value"
                    class="w-full"
                    :disabled="!authStore.isAdmin"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Visibilidad</label>
                  <Select
                    v-model="form.visible_publico"
                    :options="visibleOptions"
                    optionLabel="label"
                    optionValue="value"
                    class="w-full"
                    :disabled="!authStore.isAdmin"
                  />
                </div>
              </div>
            </div>

            <div class="flex justify-end gap-4 pt-4 border-t">
              <Button
                label="Cancelar"
                severity="secondary"
                outlined
                type="button"
                @click="navigateTo(`/admin/agencias/${id}`)"
                :disabled="saving"
              />
              <Button
                label="Guardar cambios"
                icon="pi pi-save"
                type="submit"
                :loading="saving"
              />
            </div>
          </form>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { updateAgenciaSchema } from '~/utils/validations/agencia'
import { useAuthStore } from '~/stores/auth'
import { buildPhone, DEFAULT_PHONE_PREFIX, PHONE_PREFIXES, splitPhone } from '~/utils/phone'

const AgenciaMap = defineAsyncComponent(() =>
  import('~/components/atracciones/AtraccionMap.vue')
)

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
if (!authStore.isAdmin && !authStore.isEncargado) {
  navigateTo('/dashboard')
}

const route = useRoute()
const id = Number(route.params.id)
const toast = useToast()
const { getAgencia, updateAgencia, getDepartamentos, getEncargados, getDias } = useAgencias()

const DEFAULT_HORARIO_APERTURA = '08:30'
const DEFAULT_HORARIO_CIERRE = '18:30'

const form = ref<any>({
  nombre_comercial: '',
  descripcion: '',
  direccion: '',
  departamento_id: null,
  latitud: null,
  longitud: null,
  phone_prefix: DEFAULT_PHONE_PREFIX,
  phone_number: '',
  email: '',
  sitio_web: '',
  facebook: '',
  instagram: '',
  licencia_turistica: false,
  horario_apertura: DEFAULT_HORARIO_APERTURA,
  horario_cierre: DEFAULT_HORARIO_CIERRE,
  acepta_qr: true,
  acepta_transferencia: true,
  acepta_efectivo: true,
  encargado_principal_id: null,
  status: 'activa',
  visible_publico: true,
  dias_ids: [] as number[]
})

const errors = ref<Record<string, string>>({})
const departamentos = ref<any[]>([])
const encargados = ref<any[]>([])
const dias = ref<any[]>([])
const saving = ref(false)
const loading = ref(true)

const statusOptions = [
  { label: 'Activa', value: 'activa' },
  { label: 'Inactiva', value: 'inactiva' },
  { label: 'Suspendida', value: 'suspendida' },
  { label: 'En revision', value: 'en_revision' }
]

const visibleOptions = [
  { label: 'Publica', value: true },
  { label: 'Oculta', value: false }
]

const phonePrefixes = PHONE_PREFIXES.filter((prefix) => prefix.value === DEFAULT_PHONE_PREFIX)

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

const aperturaPickerValue = computed({
  get: () => timeStringToDate(form.value.horario_apertura),
  set: (val: Date | null) => {
    form.value.horario_apertura = dateToTimeString(val)
  }
})

const cierrePickerValue = computed({
  get: () => timeStringToDate(form.value.horario_cierre),
  set: (val: Date | null) => {
    form.value.horario_cierre = dateToTimeString(val)
  }
})

const handleSubmit = async () => {
  errors.value = {}
  const parsed = updateAgenciaSchema.safeParse({
    ...form.value,
    departamento_id: form.value.departamento_id ? Number(form.value.departamento_id) : undefined,
    encargado_principal_id: form.value.encargado_principal_id ? Number(form.value.encargado_principal_id) : undefined
  })

  if (!parsed.success) {
    parsed.error.issues.forEach((issue) => {
      const path = issue.path[0] as string
      errors.value[path] = issue.message
    })
    return
  }

  saving.value = true
  try {
    const payload = {
      ...parsed.data,
      telefono: buildPhone(parsed.data.phone_prefix, parsed.data.phone_number)
    }
    delete payload.phone_prefix
    delete payload.phone_number

    const response: any = await updateAgencia(id, payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Agencia actualizada', life: 3000 })
      navigateTo(`/admin/agencias/${id}`)
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo guardar', life: 3000 })
  } finally {
    saving.value = false
  }
}

const loadAgencia = async () => {
  try {
    const response: any = await getAgencia(id)
    if (response.success) {
      const data = response.data
      const { prefix, number } = splitPhone(data.telefono)
      form.value = {
        nombre_comercial: data.nombre_comercial,
        descripcion: data.descripcion || '',
        direccion: data.direccion || '',
        departamento_id: data.departamento_id || null,
        latitud: data.latitud ?? null,
        longitud: data.longitud ?? null,
        phone_prefix: prefix,
        phone_number: number,
        email: data.email || '',
        sitio_web: data.sitio_web || '',
        facebook: data.facebook || '',
        instagram: data.instagram || '',
        licencia_turistica: data.licencia_turistica || false,
        horario_apertura: normalizeTimeString(data.horario_apertura) || DEFAULT_HORARIO_APERTURA,
        horario_cierre: normalizeTimeString(data.horario_cierre) || DEFAULT_HORARIO_CIERRE,
        acepta_qr: data.acepta_qr,
        acepta_transferencia: data.acepta_transferencia,
        acepta_efectivo: data.acepta_efectivo,
        encargado_principal_id: data.encargado_principal_id || null,
        status: data.status,
        visible_publico: data.visible_publico,
        dias_ids: (data.dias || []).map((d: any) => d.id)
      }
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo cargar la agencia', life: 3000 })
  }
}

const loadDepartamentos = async () => {
  try {
    const response: any = await getDepartamentos()
    departamentos.value = response?.data || []
  } catch (error) {
    departamentos.value = []
  }
}

const loadEncargados = async () => {
  try {
    const response: any = await getEncargados({ only_unassigned: true, agencia_id: id })
    const list = response?.data || []
    encargados.value = list.map((u: any) => ({
      ...u,
      nombre_completo: [u.nombre, u.apellido_paterno, u.apellido_materno].filter(Boolean).join(' ')
    }))
  } catch (error) {
    encargados.value = []
  }
}

const loadDias = async () => {
  try {
    const response: any = await getDias()
    dias.value = response?.data || []
  } catch (error) {
    dias.value = []
  }
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([loadDepartamentos(), loadEncargados(), loadDias(), loadAgencia()])
  } finally {
    loading.value = false
  }
})
</script>
