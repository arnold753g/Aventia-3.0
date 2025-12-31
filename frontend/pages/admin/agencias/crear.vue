<template>
  <div class="max-w-4xl mx-auto px-4 py-10">
    <Card>
      <template #title>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">Creacion rapida</p>
            <h1 class="text-2xl font-bold text-gray-900">Nueva agencia</h1>
          </div>
          <Button label="Volver" icon="pi pi-arrow-left" outlined @click="navigateTo('/admin/agencias')" />
        </div>
      </template>
      <template #content>
        <div class="p-4 rounded-lg bg-orange-50 border border-orange-100 text-orange-800 mb-6">
          Las agencias creadas por esta via quedan en estado <strong>en revision</strong> y con visibilidad deshabilitada hasta completar sus datos.
        </div>

        <form class="space-y-5" @submit.prevent="handleSubmit">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Nombre comercial</label>
              <InputText v-model="form.nombre_comercial" class="w-full" placeholder="Ej: Andaria Tours" />
              <small v-if="errors.nombre_comercial" class="text-red-500">{{ errors.nombre_comercial }}</small>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Telefono</label>
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
              <small v-if="errors.phone_prefix" class="text-red-500">{{ errors.phone_prefix }}</small>
              <small v-if="errors.phone_number" class="text-red-500">{{ errors.phone_number }}</small>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Departamento</label>
              <Select
                v-model="form.departamento_id"
                :options="departamentos"
                optionLabel="nombre"
                optionValue="id"
                placeholder="Seleccione"
                class="w-full"
                showClear
              />
              <small v-if="errors.departamento_id" class="text-red-500">{{ errors.departamento_id }}</small>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Encargado principal</label>
              <Select
                v-model="form.encargado_principal_id"
                :options="encargados"
                optionLabel="nombre_completo"
                optionValue="id"
                placeholder="Seleccione"
                class="w-full"
                showClear
              >
                <template #option="slotProps">
                  <div class="flex flex-col">
                    <span class="font-semibold">{{ slotProps.option.nombre_completo }}</span>
                    <span class="text-xs text-gray-500">{{ slotProps.option.email }}</span>
                  </div>
                </template>
              </Select>
              <small v-if="errors.encargado_principal_id" class="text-red-500">{{ errors.encargado_principal_id }}</small>
            </div>
          </div>

          <div class="flex justify-end gap-2">
            <Button label="Cancelar" severity="secondary" type="button" @click="navigateTo('/admin/agencias')" />
            <Button label="Crear" icon="pi pi-check" type="submit" :loading="submitting" />
          </div>
        </form>
      </template>
    </Card>
    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { createAgenciaRapidaSchema } from '~/utils/validations/agencia'
import { useAuthStore } from '~/stores/auth'
import { buildPhone, DEFAULT_PHONE_PREFIX, PHONE_PREFIXES } from '~/utils/phone'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
if (!authStore.isAdmin) {
  navigateTo('/dashboard')
}

const toast = useToast()
const { createAgenciaRapida, getDepartamentos, getEncargados } = useAgencias()

const form = ref({
  nombre_comercial: '',
  departamento_id: null as number | null,
  phone_prefix: DEFAULT_PHONE_PREFIX,
  phone_number: '',
  encargado_principal_id: null as number | null
})

const errors = ref<Record<string, string>>({})
const departamentos = ref<any[]>([])
const encargados = ref<any[]>([])
const submitting = ref(false)

const phonePrefixes = PHONE_PREFIXES.filter((prefix) => prefix.value === DEFAULT_PHONE_PREFIX)

const handleSubmit = async () => {
  errors.value = {}
  const parsed = createAgenciaRapidaSchema.safeParse(form.value)
  if (!parsed.success) {
    parsed.error.issues.forEach((issue) => {
      const path = issue.path[0] as string
      errors.value[path] = issue.message
    })
    return
  }

  submitting.value = true
  try {
    const payload = {
      ...parsed.data,
      telefono: buildPhone(parsed.data.phone_prefix, parsed.data.phone_number)
    }
    delete payload.phone_prefix
    delete payload.phone_number

    const response: any = await createAgenciaRapida(payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Creada', detail: 'Agencia creada en revision', life: 3000 })
      navigateTo(`/admin/agencias/${response.data.id}/editar`)
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

onMounted(() => {
  loadDepartamentos()
  loadEncargados()
})
</script>
