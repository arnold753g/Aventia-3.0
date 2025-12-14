<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-5xl mx-auto px-4 py-6">
        <div class="flex items-center gap-4">
          <Button icon="pi pi-arrow-left" text rounded @click="navigateTo('/admin/agencias')" />
          <div>
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Crear Agencia (Completa)
            </h1>
            <p class="muted mt-1">Complete el formulario para registrar una nueva agencia</p>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-5xl mx-auto px-4 py-8">
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
                  <InputText v-model="form.nombre_comercial" placeholder="Ej: Andaria Tours" class="w-full" />
                  <Message v-if="errors.nombre_comercial" severity="error" size="small" variant="simple">
                    {{ errors.nombre_comercial }}
                  </Message>
                </div>

                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Descripcion</label>
                  <Textarea v-model="form.descripcion" rows="4" autoResize class="w-full" placeholder="Descripcion de la agencia (opcional)" />
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
                  <label class="block text-sm font-medium muted mb-2">Departamento *</label>
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
                  <label class="block text-sm font-medium muted mb-2">Direccion *</label>
                  <InputText v-model="form.direccion" placeholder="Ej: Av. Principal #123" class="w-full" />
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
                  <label class="block text-sm font-medium muted mb-2">Telefono *</label>
                  <InputMask v-model="form.telefono" mask="99999999" class="w-full" placeholder="8 digitos (2-7)" />
                  <Message v-if="errors.telefono" severity="error" size="small" variant="simple">
                    {{ errors.telefono }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Email *</label>
                  <InputText v-model="form.email" class="w-full" placeholder="contacto@agencia.com" />
                  <Message v-if="errors.email" severity="error" size="small" variant="simple">
                    {{ errors.email }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Sitio web</label>
                  <InputText v-model="form.sitio_web" class="w-full" placeholder="https://..." />
                  <Message v-if="errors.sitio_web" severity="error" size="small" variant="simple">
                    {{ errors.sitio_web }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Facebook</label>
                  <InputText v-model="form.facebook" class="w-full" placeholder="https://facebook.com/..." />
                  <Message v-if="errors.facebook" severity="error" size="small" variant="simple">
                    {{ errors.facebook }}
                  </Message>
                </div>

                <div class="md:col-span-2">
                  <label class="block text-sm font-medium muted mb-2">Instagram</label>
                  <InputText v-model="form.instagram" class="w-full" placeholder="https://instagram.com/..." />
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
                  <InputMask v-model="form.horario_apertura" mask="99:99" slotChar="HH:MM" class="w-full" placeholder="08:00" />
                  <Message v-if="errors.horario_apertura" severity="error" size="small" variant="simple">
                    {{ errors.horario_apertura }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Horario cierre</label>
                  <InputMask v-model="form.horario_cierre" mask="99:99" slotChar="HH:MM" class="w-full" placeholder="18:00" />
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
                  <label class="block text-sm font-medium muted mb-2">Encargado principal *</label>
                  <Select
                    v-model="form.encargado_principal_id"
                    :options="encargados"
                    optionLabel="nombre_completo"
                    optionValue="id"
                    placeholder="Seleccione..."
                    class="w-full"
                    showClear
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
                  <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Visibilidad</label>
                  <Select v-model="form.visible_publico" :options="visibleOptions" optionLabel="label" optionValue="value" class="w-full" />
                </div>
              </div>
            </div>

            <!-- Botones -->
            <div class="flex justify-end gap-4 pt-4 border-t">
              <Button
                label="Cancelar"
                severity="secondary"
                outlined
                type="button"
                @click="navigateTo('/admin/agencias')"
                :disabled="submitting"
              />
              <Button
                label="Crear Agencia"
                icon="pi pi-save"
                type="submit"
                :loading="submitting"
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
import { ref, onMounted, defineAsyncComponent } from 'vue'
import { useToast } from 'primevue/usetoast'
import { createAgenciaCompletaSchema } from '~/utils/validations/agencia'
import { useAuthStore } from '~/stores/auth'

const AgenciaMap = defineAsyncComponent(() =>
  import('~/components/atracciones/AtraccionMap.vue')
)

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
if (!authStore.isAdmin) {
  navigateTo('/dashboard')
}

const toast = useToast()
const { createAgenciaCompleta, getDepartamentos, getEncargados, getDias } = useAgencias()

const form = ref({
  nombre_comercial: '',
  descripcion: '',
  direccion: '',
  departamento_id: null as number | null,
  latitud: null as number | null,
  longitud: null as number | null,
  telefono: '',
  email: '',
  sitio_web: '',
  facebook: '',
  instagram: '',
  licencia_turistica: false,
  horario_apertura: '08:00',
  horario_cierre: '18:00',
  acepta_qr: true,
  acepta_transferencia: true,
  acepta_efectivo: true,
  encargado_principal_id: null as number | null,
  status: 'activa',
  visible_publico: true,
  dias_ids: [] as number[]
})

const errors = ref<Record<string, string>>({})
const departamentos = ref<any[]>([])
const encargados = ref<any[]>([])
const dias = ref<any[]>([])
const submitting = ref(false)

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

const handleSubmit = async () => {
  errors.value = {}
  const parsed = createAgenciaCompletaSchema.safeParse({
    ...form.value,
    departamento_id: Number(form.value.departamento_id),
    encargado_principal_id: Number(form.value.encargado_principal_id)
  })

  if (!parsed.success) {
    parsed.error.issues.forEach((issue) => {
      const path = issue.path[0] as string
      errors.value[path] = issue.message
    })
    return
  }

  submitting.value = true
  try {
    const response: any = await createAgenciaCompleta(parsed.data)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Creada', detail: 'Agencia creada correctamente', life: 3000 })
      navigateTo(`/admin/agencias/${response.data.id}`)
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo crear', life: 3000 })
  } finally {
    submitting.value = false
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
    const response: any = await getEncargados({ only_unassigned: true })
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

onMounted(() => {
  loadDepartamentos()
  loadEncargados()
  loadDias()
})
</script>
