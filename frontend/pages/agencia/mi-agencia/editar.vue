<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-6xl mx-auto px-4 py-6">
        <div class="flex items-center gap-4">
          <Button icon="pi pi-arrow-left" text rounded @click="navigateTo('/agencia/mi-agencia')" />
          <div>
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Editar mi agencia
            </h1>
            <p class="muted mt-1">{{ form.nombre_comercial || 'Actualiza la informacion principal' }}</p>
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

            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-phone"></i>
                Contacto y redes
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium muted mb-2">Telefono *</label>
                  <InputText v-model="form.telefono" class="w-full" />
                  <Message v-if="errors.telefono" severity="error" size="small" variant="simple">
                    {{ errors.telefono }}
                  </Message>
                </div>

                <div>
                  <label class="block text-sm font-medium muted mb-2">Email *</label>
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

                <div>
                  <label class="block text-sm font-medium muted mb-2">Instagram</label>
                  <InputText v-model="form.instagram" class="w-full" />
                  <Message v-if="errors.instagram" severity="error" size="small" variant="simple">
                    {{ errors.instagram }}
                  </Message>
                </div>
              </div>
            </div>

            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-clock"></i>
                Horarios
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium muted mb-2">Horario apertura (HH:MM)</label>
                  <InputText v-model="form.horario_apertura" class="w-full" placeholder="08:00" />
                  <Message v-if="errors.horario_apertura" severity="error" size="small" variant="simple">
                    {{ errors.horario_apertura }}
                  </Message>
                </div>
                <div>
                  <label class="block text-sm font-medium muted mb-2">Horario cierre (HH:MM)</label>
                  <InputText v-model="form.horario_cierre" class="w-full" placeholder="18:00" />
                  <Message v-if="errors.horario_cierre" severity="error" size="small" variant="simple">
                    {{ errors.horario_cierre }}
                  </Message>
                </div>
              </div>
            </div>

            <div>
              <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                <i class="pi pi-wallet"></i>
                Metodos de pago
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div class="flex items-center gap-2 border rounded px-3 py-2">
                  <Checkbox v-model="form.acepta_qr" binary inputId="pago_qr" />
                  <label for="pago_qr" class="text-sm muted">QR</label>
                </div>
                <div class="flex items-center gap-2 border rounded px-3 py-2">
                  <Checkbox v-model="form.acepta_transferencia" binary inputId="pago_trans" />
                  <label for="pago_trans" class="text-sm muted">Transferencia</label>
                </div>
                <div class="flex items-center gap-2 border rounded px-3 py-2">
                  <Checkbox v-model="form.acepta_efectivo" binary inputId="pago_cash" />
                  <label for="pago_cash" class="text-sm muted">Efectivo</label>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-end gap-3">
              <Button label="Cancelar" severity="secondary" type="button" outlined @click="navigateTo('/agencia/mi-agencia')" />
              <Button label="Guardar" icon="pi pi-save" type="submit" :loading="saving" />
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
import { updateAgenciaSchema } from '~/utils/validations/agencia'

const AgenciaMap = defineAsyncComponent(() =>
  import('~/components/atracciones/AtraccionMap.vue')
)

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const { getMiAgencia, updateAgencia, getDepartamentos, getDias } = useAgencias()

const agenciaId = ref<number | null>(null)
const form = ref<any>({
  nombre_comercial: '',
  descripcion: '',
  direccion: '',
  departamento_id: null,
  latitud: null,
  longitud: null,
  telefono: '',
  email: '',
  sitio_web: '',
  facebook: '',
  instagram: '',
  horario_apertura: '',
  horario_cierre: '',
  acepta_qr: true,
  acepta_transferencia: true,
  acepta_efectivo: true,
  dias_ids: [] as number[]
})

const errors = ref<Record<string, string>>({})
const departamentos = ref<any[]>([])
const dias = ref<any[]>([])
const saving = ref(false)
const loading = ref(true)

const handleSubmit = async () => {
  errors.value = {}
  const parsed = updateAgenciaSchema.safeParse({
    ...form.value,
    departamento_id: form.value.departamento_id ? Number(form.value.departamento_id) : undefined
  })

  if (!parsed.success) {
    parsed.error.issues.forEach((issue) => {
      const path = issue.path[0] as string
      errors.value[path] = issue.message
    })
    return
  }

  if (!agenciaId.value) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Agencia no cargada', life: 3000 })
    return
  }

  saving.value = true
  try {
    const payload = { ...parsed.data }
    const response: any = await updateAgencia(agenciaId.value, payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Agencia actualizada', life: 3000 })
      navigateTo('/agencia/mi-agencia')
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo guardar', life: 3000 })
  } finally {
    saving.value = false
  }
}

const loadAgencia = async () => {
  try {
    const response: any = await getMiAgencia()
    if (response.success) {
      const data = response.data
      agenciaId.value = Number(data.id)
      form.value = {
        nombre_comercial: data.nombre_comercial,
        descripcion: data.descripcion || '',
        direccion: data.direccion || '',
        departamento_id: data.departamento_id || null,
        latitud: data.latitud ?? null,
        longitud: data.longitud ?? null,
        telefono: data.telefono || '',
        email: data.email || '',
        sitio_web: data.sitio_web || '',
        facebook: data.facebook || '',
        instagram: data.instagram || '',
        horario_apertura: data.horario_apertura || '',
        horario_cierre: data.horario_cierre || '',
        acepta_qr: data.acepta_qr,
        acepta_transferencia: data.acepta_transferencia,
        acepta_efectivo: data.acepta_efectivo,
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
  } catch {
    departamentos.value = []
  }
}

const loadDias = async () => {
  try {
    const response: any = await getDias()
    dias.value = response?.data || []
  } catch {
    dias.value = []
  }
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([loadDepartamentos(), loadDias(), loadAgencia()])
  } finally {
    loading.value = false
  }
})
</script>

