<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-6xl mx-auto px-4 py-6 flex items-center gap-4">
        <Button icon="pi pi-arrow-left" text rounded @click="navigateTo(`/admin/atracciones/${route.params.id}`)" />
        <div>
          <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
            Editar Atracción
          </h1>
          <p class="muted mt-1">Actualiza la información principal</p>
        </div>
      </div>
    </div>

    <div v-if="loading" class="max-w-6xl mx-auto px-4 py-8">
      <Skeleton height="400px" />
    </div>

    <div v-else-if="initialValues" class="max-w-6xl mx-auto px-4 py-8">
      <Card class="surface-card">
        <template #content>
          <Form v-slot="form" :initialValues="initialValues" :resolver="resolver" @submit="handleSubmit">
            <div class="space-y-8">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField v-slot="$field" name="nombre">
                  <label class="block text-sm font-medium muted mb-2">Nombre</label>
                  <InputText v-model="formData.nombre" class="w-full" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
                <FormField v-slot="$field" name="precio_entrada">
                  <label class="block text-sm font-medium muted mb-2">Precio (Bs)</label>
                  <InputText type="number" step="0.5" class="w-full" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
                <FormField v-slot="$field" name="status">
                  <label class="block text-sm font-medium muted mb-2">Estado</label>
                  <Dropdown
                    v-model="status"
                    :options="statusOptions"
                    optionLabel="label"
                    optionValue="value"
                    class="w-full"
                    @change="form.setFieldValue('status', status || '')"
                  />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
                <FormField v-slot="$field" name="nivel_dificultad">
                  <label class="block text-sm font-medium muted mb-2">Dificultad</label>
                  <Dropdown
                    v-model="nivelDificultad"
                    :options="nivelOptions"
                    optionLabel="label"
                    optionValue="value"
                    class="w-full"
                    showClear
                    @change="form.setFieldValue('nivel_dificultad', nivelDificultad || '')"
                  />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
              </div>

              <FormField v-slot="$field" name="descripcion">
                <label class="block text-sm font-medium muted mb-2">Descripción</label>
                <Textarea rows="4" class="w-full" autoResize />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-4">
                  <FormField v-slot="$field" name="provincia_id">
                    <label class="block text-sm font-medium muted mb-2">Provincia</label>
                    <Dropdown
                      v-model="provinciaId"
                      :options="provincias"
                      optionLabel="nombre"
                      optionValue="id"
                      class="w-full"
                      placeholder="Seleccione provincia"
                      @change="form.setFieldValue('provincia_id', provinciaId || '')"
                    />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                  <label class="block text-sm font-medium muted mb-2">Filtrar por departamento</label>
                  <Dropdown
                    v-model="departamentoId"
                    :options="departamentos"
                    optionLabel="nombre"
                    optionValue="id"
                    class="w-full"
                    placeholder="Departamento"
                    showClear
                    @change="loadProvincias"
                  />

                  <FormField v-slot="$field" name="direccion">
                    <label class="block text-sm font-medium muted mb-2">Dirección</label>
                    <InputText class="w-full" />
                  </FormField>
                </div>

                <div>
                  <AtraccionMap
                    :latitud="latitud"
                    :longitud="longitud"
                    editable
                    height="300px"
                    @update:latitud="val => { latitud.value = val; form.setFieldValue('latitud', val) }"
                    @update:longitud="val => { longitud.value = val; form.setFieldValue('longitud', val) }"
                  />
                  <div class="grid grid-cols-2 gap-4 mt-4">
                    <FormField v-slot="$field" name="horario_apertura">
                      <label class="block text-sm font-medium muted mb-2">Apertura</label>
                      <InputText type="time" class="w-full" />
                    </FormField>
                    <FormField v-slot="$field" name="horario_cierre">
                      <label class="block text-sm font-medium muted mb-2">Cierre</label>
                      <InputText type="time" class="w-full" />
                    </FormField>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField v-slot="$field" name="telefono">
                  <label class="block text-sm font-medium muted mb-2">Teléfono</label>
                  <InputText class="w-full" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
                <FormField v-slot="$field" name="email">
                  <label class="block text-sm font-medium muted mb-2">Email</label>
                  <InputText type="email" class="w-full" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
                <FormField v-slot="$field" name="sitio_web">
                  <label class="block text-sm font-medium muted mb-2">Sitio web</label>
                  <InputText class="w-full" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
                <FormField v-slot="$field" name="facebook">
                  <label class="block text-sm font-medium muted mb-2">Facebook</label>
                  <InputText class="w-full" />
                </FormField>
                <FormField v-slot="$field" name="instagram">
                  <label class="block text-sm font-medium muted mb-2">Instagram</label>
                  <InputText class="w-full" />
                </FormField>
              </div>

              <div class="flex items-center gap-4">
                <Checkbox v-model="requiereAgencia" inputId="requiere_agencia" binary @change="form.setFieldValue('requiere_agencia', !!requiereAgencia)" />
                <label for="requiere_agencia" class="text-sm">Requiere agencia</label>

                <Checkbox v-model="accesoParticular" inputId="acceso_particular" binary @change="form.setFieldValue('acceso_particular', !!accesoParticular)" />
                <label for="acceso_particular" class="text-sm">Acceso particular</label>

                <Checkbox v-model="visiblePublico" inputId="visible_publico" binary @change="form.setFieldValue('visible_publico', !!visiblePublico)" />
                <label for="visible_publico" class="text-sm">Visible al público</label>
              </div>

              <div class="flex justify-end gap-4 pt-4 border-t">
                <Button
                  label="Cancelar"
                  severity="secondary"
                  outlined
                  @click="navigateTo(`/admin/atracciones/${route.params.id}`)"
                  :disabled="saving"
                />
                <Button
                  label="Guardar cambios"
                  icon="pi pi-save"
                  type="submit"
                  :loading="saving"
                  :disabled="!form.valid"
                />
              </div>
            </div>
          </Form>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { zodResolver } from '@primevue/forms/resolvers/zod'
import { extractHora } from '~/utils/formatters-atraccion'
import { updateAtraccionSchema } from '~/utils/validations/atraccion'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const route = useRoute()
const toast = useToast()
const authStore = useAuthStore()
const { getAtraccion, updateAtraccion, getDepartamentos, getProvincias } = useAtracciones()

const loading = ref(false)
const saving = ref(false)
const initialValues = ref<any>(null)

const departamentos = ref<any[]>([])
const provincias = ref<any[]>([])
const provinciaId = ref<number | null>(null)
const departamentoId = ref<number | null>(null)
const latitud = ref<number | null>(null)
const longitud = ref<number | null>(null)
const nivelDificultad = ref<string | null>(null)
const status = ref<string | null>(null)
const requiereAgencia = ref(false)
const accesoParticular = ref(true)
const visiblePublico = ref(true)

const resolver = zodResolver(updateAtraccionSchema)

const nivelOptions = [
  { label: 'Fácil', value: 'facil' },
  { label: 'Medio', value: 'medio' },
  { label: 'Difícil', value: 'dificil' },
  { label: 'Extremo', value: 'extremo' }
]

const statusOptions = [
  { label: 'Activa', value: 'activa' },
  { label: 'Inactiva', value: 'inactiva' },
  { label: 'Mantenimiento', value: 'mantenimiento' },
  { label: 'Fuera de temporada', value: 'fuera_temporada' }
]

const formData = reactive({
  nombre: ''
})

// Función para capitalizar palabras
const capitalize = (text: string) => {
  if (!text) return text
  return text
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

// Watcher para capitalizar nombre en tiempo real
watch(() => formData.nombre, (newVal) => {
  const capitalized = capitalize(newVal)
  if (capitalized !== newVal) {
    formData.nombre = capitalized
    if (initialValues.value) {
      initialValues.value.nombre = capitalized
    }
  }
})

const loadDepartamentos = async () => {
  const response: any = await getDepartamentos()
  if (response?.success) {
    departamentos.value = response.data
  }
}

const loadProvincias = async () => {
  const response: any = await getProvincias(departamentoId.value || undefined)
  if (response?.success) {
    provincias.value = response.data
  }
}

const loadAtraccion = async () => {
  loading.value = true
  try {
    const response: any = await getAtraccion(Number(route.params.id))
    if (response.success) {
      const a = response.data
      provinciaId.value = a.provincia_id
      departamentoId.value = a.provincia?.departamento?.id || null
      latitud.value = a.latitud
      longitud.value = a.longitud
      nivelDificultad.value = a.nivel_dificultad || null
      status.value = a.status || null
      requiereAgencia.value = !!a.requiere_agencia
      accesoParticular.value = !!a.acceso_particular
      visiblePublico.value = !!a.visible_publico
      formData.nombre = a.nombre

      initialValues.value = {
        nombre: a.nombre,
        descripcion: a.descripcion,
        provincia_id: a.provincia_id,
        direccion: a.direccion,
        latitud: a.latitud,
        longitud: a.longitud,
        horario_apertura: extractHora(a.horario_apertura),
        horario_cierre: extractHora(a.horario_cierre),
        precio_entrada: a.precio_entrada,
        nivel_dificultad: a.nivel_dificultad || '',
        requiere_agencia: !!a.requiere_agencia,
        acceso_particular: !!a.acceso_particular,
        status: a.status || '',
        visible_publico: !!a.visible_publico,
        telefono: a.telefono || '',
        email: a.email || '',
        sitio_web: a.sitio_web || '',
        facebook: a.facebook || '',
        instagram: a.instagram || ''
      }
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo cargar la atracción',
      life: 3000
    })
    navigateTo('/admin/atracciones')
  } finally {
    loading.value = false
  }
}

const handleSubmit = async ({ valid, values }: any) => {
  if (!valid) {
    toast.add({
      severity: 'warn',
      summary: 'Formulario incompleto',
      detail: 'Revisa los campos requeridos',
      life: 3000
    })
    return
  }

  saving.value = true
  try {
    const payload: any = {
      ...values,
      nombre: formData.nombre,
      provincia_id: values.provincia_id ? Number(values.provincia_id) : undefined,
      latitud: values.latitud ? Number(values.latitud) : undefined,
      longitud: values.longitud ? Number(values.longitud) : undefined,
      requiere_agencia: !!requiereAgencia.value,
      acceso_particular: !!accesoParticular.value,
      visible_publico: !!visiblePublico.value
    }

    const response: any = await updateAtraccion(Number(route.params.id), payload)

    if (response.success) {
      toast.add({
        severity: 'success',
        summary: 'Cambios guardados',
        detail: 'Atracción actualizada correctamente',
        life: 3000
      })
      navigateTo(`/admin/atracciones/${route.params.id}`)
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: response.error?.message || 'No se pudo actualizar',
        life: 4000
      })
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || error.message || 'No se pudo actualizar',
      life: 4000
    })
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (!authStore.isAdmin) {
    navigateTo('/dashboard')
    return
  }
  await loadDepartamentos()
  await loadProvincias()
  await loadAtraccion()
})
</script>
