<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-4xl mx-auto px-4 py-6">
        <div class="flex items-center gap-4">
          <Button
            icon="pi pi-arrow-left"
            text
            rounded
            @click="navigateTo('/admin/usuarios')"
          />
          <div>
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Crear Nuevo Usuario
            </h1>
            <p class="muted mt-1">Complete todos los campos requeridos</p>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-4xl mx-auto px-4 py-8">
      <Card class="surface-card">
        <template #content>
          <Form v-slot="form" :initialValues="initialValues" :resolver="resolver" @submit="handleSubmit">
            <div class="space-y-8">
              <!-- Datos Personales -->
              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-user"></i>
                  Datos Personales
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <FormField v-slot="$field" name="nombre">
                    <label class="block text-sm font-medium muted mb-2">
                      Nombre *
                      <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'Ingrese el nombre completo'"></i>
                    </label>
                    <InputText v-model="formData.nombre" placeholder="Ej: Juan" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                  <FormField v-slot="$field" name="apellido_paterno">
                    <label class="block text-sm font-medium muted mb-2">Apellido Paterno *</label>
                    <InputText v-model="formData.apellido_paterno" placeholder="Ej: Pérez" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                  <FormField v-slot="$field" name="apellido_materno">
                    <label class="block text-sm font-medium muted mb-2">Apellido Materno *</label>
                    <InputText v-model="formData.apellido_materno" placeholder="Ej: López" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>
                </div>
              </div>

              <!-- Identificación y Contacto -->
              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-id-card"></i>
                  Identificación y Contacto
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField v-slot="$field" name="ci">
                    <label class="block text-sm font-medium muted mb-2">
                      CI *
                      <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'Carnet de Identidad boliviano (5-15 caracteres)'"></i>
                    </label>
                    <InputText placeholder="Ej: 1234567" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                  <FormField v-slot="$field" name="expedido">
                    <label class="block text-sm font-medium muted mb-2">Expedido en *</label>
                    <Select :options="departamentos" optionLabel="label" optionValue="value" placeholder="Seleccione..." class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                                    <div>
                    <label class="block text-sm font-medium muted mb-2">
                      Telefono *
                      <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'Seleccione prefijo y numero'"></i>
                    </label>
                    <div class="grid grid-cols-1 sm:grid-cols-[140px,1fr] gap-3">
                      <FormField v-slot="$prefixField" name="phone_prefix">
                        <Select :options="phonePrefixes" optionLabel="label" optionValue="value" placeholder="Prefijo" class="w-full" />
                        <Message v-if="$prefixField?.invalid" severity="error" size="small" variant="simple">
                          {{ $prefixField.error?.message }}
                        </Message>
                      </FormField>

                      <FormField v-slot="$numberField" name="phone_number">
                        <InputText placeholder="Ej: 71234567" class="w-full" maxlength="14" inputmode="numeric" />
                        <Message v-if="$numberField?.invalid" severity="error" size="small" variant="simple">
                          {{ $numberField.error?.message }}
                        </Message>
                      </FormField>
                    </div>
                  </div>

                  <FormField v-slot="$field" name="email">
                    <label class="block text-sm font-medium muted mb-2">Email *</label>
                    <InputText type="email" placeholder="usuario@ejemplo.com" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>
                </div>
              </div>

              <!-- Información Adicional -->
              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-calendar"></i>
                  Información Adicional
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField v-slot="{ props, value, invalid, error }" name="fecha_nacimiento">
                    <label class="block text-sm font-medium muted mb-2">
                      Fecha de Nacimiento *
                      <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'Debe ser mayor de 18 años'"></i>
                    </label>
                    <DatePicker
                      v-bind="{ ...props, onChange: undefined }"
                      :modelValue="toDate(value as string)"
                      @update:modelValue="(val) => props?.onChange?.({ value: formatDateYMD(val as Date | null) })"
                      showIcon
                      dateFormat="yy-mm-dd"
                      inputClass="w-full"
                      class="w-full"
                    />
                    <Message v-if="invalid" severity="error" size="small" variant="simple">
                      {{ error }}
                    </Message>
                  </FormField>

                  <FormField v-slot="$field" name="ciudad">
                    <label class="block text-sm font-medium muted mb-2">Ciudad</label>
                    <InputText placeholder="Ej: La Paz" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>
                </div>
              </div>

              <!-- Foto de Perfil -->
              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-image"></i>
                  Foto de Perfil
                </h3>
                <div class="flex items-center gap-4">
                  <!-- Preview con imagen o placeholder -->
                  <div class="flex-shrink-0">
                    <div v-if="preview" class="relative">
                      <img
                        :src="preview"
                        alt="Preview"
                        class="w-24 h-24 rounded-full object-cover border-2 border-gray-200 shadow-lg"
                      />
                      <button
                        type="button"
                        @click="removePhoto"
                        class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center hover:bg-red-600 transition-colors shadow-md"
                      >
                        <i class="pi pi-times text-xs"></i>
                      </button>
                    </div>
                    <div v-else class="w-24 h-24 rounded-full bg-gray-100 flex items-center justify-center border-2 border-gray-200 shadow-lg">
                      <i class="pi pi-user text-4xl text-gray-400"></i>
                    </div>
                  </div>

                  <!-- Controles de carga -->
                  <div class="flex-1">
                    <input
                      type="file"
                      accept="image/*"
                      @change="onFileChange"
                      class="hidden"
                      ref="fileInput"
                    />
                    <Button
                      type="button"
                      :label="preview ? 'Cambiar foto' : 'Seleccionar Foto'"
                      icon="pi pi-upload"
                      severity="secondary"
                      outlined
                      @click="fileInput?.click()"
                    />
                    <p class="text-xs muted mt-2">Formatos aceptados: JPG, PNG. Tamaño máximo: 2MB</p>
                  </div>
                </div>
              </div>

              <!-- Credenciales de Acceso -->
              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-lock"></i>
                  Credenciales de Acceso
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField v-slot="$field" name="password">
                    <label class="block text-sm font-medium muted mb-2">
                      Contraseña *
                      <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'Mín. 8 caracteres, incluye mayúscula, minúscula, número y carácter especial'"></i>
                    </label>
                    <Password placeholder="Contraseña segura" toggleMask :feedback="true" class="w-full" fluid />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                  <FormField v-slot="$field" name="rol">
                    <label class="block text-sm font-medium muted mb-2">Rol *</label>
                    <Select :options="roles" optionLabel="label" optionValue="value" placeholder="Seleccione un rol" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>
                </div>
              </div>

              <!-- Botones de Acción -->
              <div class="flex justify-end gap-4 pt-4 border-t">
                <Button
                  label="Cancelar"
                  severity="secondary"
                  outlined
                  @click="navigateTo('/admin/usuarios')"
                  :disabled="loading"
                />
                <Button
                  label="Crear Usuario"
                  icon="pi pi-save"
                  type="submit"
                  :loading="loading"
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
import { ref, computed, reactive, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import DatePicker from 'primevue/datepicker'
import Form from '@primevue/forms/form'
import FormField from '@primevue/forms/formfield'
import { zodResolver } from '@primevue/forms/resolvers/zod'
import { createUsuarioSchema } from '~/utils/validations/usuario'
import { buildPhone, DEFAULT_PHONE_PREFIX, PHONE_PREFIXES } from '~/utils/phone'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const toast = useToast()
const { createUsuario } = useUsuarios()

const loading = ref(false)
const preview = ref<string | null>(null)
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

// Valores iniciales
const initialValues = {
  nombre: '',
  apellido_paterno: '',
  apellido_materno: '',
  ci: '',
  expedido: '',
  phone_prefix: DEFAULT_PHONE_PREFIX,
  phone_number: '',
  email: '',
  fecha_nacimiento: '',
  ciudad: '',
  password: '',
  rol: ''
}

// Resolver de Zod
const resolver = zodResolver(createUsuarioSchema)

// Estado del formulario para v-model con capitalización
const formData = reactive({
  nombre: '',
  apellido_paterno: '',
  apellido_materno: ''
})

// Función para capitalizar primera letra
const capitalize = (t: string) => t ? t.charAt(0).toUpperCase() + t.slice(1).toLowerCase() : ''

// Watchers para capitalizar automáticamente
watch(() => formData.nombre, (newVal) => {
  const capitalized = capitalize(newVal)
  if (capitalized !== newVal) {
    formData.nombre = capitalized
    initialValues.nombre = capitalized
  }
})

watch(() => formData.apellido_paterno, (newVal) => {
  const capitalized = capitalize(newVal)
  if (capitalized !== newVal) {
    formData.apellido_paterno = capitalized
    initialValues.apellido_paterno = capitalized
  }
})

watch(() => formData.apellido_materno, (newVal) => {
  const capitalized = capitalize(newVal)
  if (capitalized !== newVal) {
    formData.apellido_materno = capitalized
    initialValues.apellido_materno = capitalized
  }
})

// Opciones
const departamentos = [
  { label: 'La Paz', value: 'LP' },
  { label: 'Cochabamba', value: 'CB' },
  { label: 'Santa Cruz', value: 'SC' },
  { label: 'Potosí', value: 'PT' },
  { label: 'Oruro', value: 'OR' },
  { label: 'Tarija', value: 'TJ' },
  { label: 'Chuquisaca', value: 'CH' },
  { label: 'Beni', value: 'BN' },
  { label: 'Pando', value: 'PD' }
]

const roles = [
  { label: 'Turista', value: 'turista' },
  { label: 'Encargado de Agencia', value: 'encargado_agencia' },
  { label: 'Administrador', value: 'admin' }
]

const phonePrefixes = PHONE_PREFIXES

// Fecha máxima (18 años atrás)
const maxBirthDate = computed(() => {
  const date = new Date()
  date.setFullYear(date.getFullYear() - 18)
  date.setHours(0, 0, 0, 0)
  return date
})

const formatDateYMD = (date: Date | null) => {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const toDate = (value?: string) => {
  if (!value) return null
  const date = new Date(`${value}T00:00:00`)
  return Number.isNaN(date.getTime()) ? null : date
}

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const file = target.files[0]

  // Validar tamaño (máximo 2MB)
  if (file.size > 2 * 1024 * 1024) {
    toast.add({
      severity: 'warn',
      summary: 'Archivo muy grande',
      detail: 'El tamaño máximo permitido es 2MB',
      life: 3000
    })
    return
  }

  // Validar tipo
  if (!file.type.startsWith('image/')) {
    toast.add({
      severity: 'warn',
      summary: 'Tipo de archivo inválido',
      detail: 'Solo se permiten imágenes (JPG, PNG)',
      life: 3000
    })
    return
  }

  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = () => {
    preview.value = reader.result as string
    console.log('Preview creado exitosamente')
  }
  reader.readAsDataURL(file)
}

const removePhoto = () => {
  selectedFile.value = null
  preview.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
  console.log('Foto removida')
}

const handleSubmit = async ({ valid, values }: any) => {
  if (!valid) {
    toast.add({
      severity: 'warn',
      summary: 'Formulario incompleto',
      detail: 'Por favor complete todos los campos requeridos correctamente',
      life: 3000
    })
    return
  }

  loading.value = true
  try {
    const payload: any = {
      ...values,
      phone: buildPhone(values.phone_prefix, values.phone_number)
    }
    delete payload.phone_prefix
    delete payload.phone_number

    // Agregar foto si existe
    if (selectedFile.value) {
      console.log('Archivo seleccionado:', {
        name: selectedFile.value.name,
        size: selectedFile.value.size,
        type: selectedFile.value.type
      })
      payload.profile_photo = selectedFile.value
    } else {
      console.log('No se seleccionó ninguna foto')
    }

    console.log('Enviando payload para crear usuario...')
    const response: any = await createUsuario(payload)
    console.log('Respuesta del servidor:', response)

    if (response.success) {
      toast.add({
        severity: 'success',
        summary: '¡Usuario Creado!',
        detail: 'El usuario ha sido creado exitosamente',
        life: 3000
      })

      setTimeout(() => {
        navigateTo('/admin/usuarios')
      }, 1200)
    } else {
      console.error('Error en respuesta:', response)
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: response.error?.message || 'Error al crear usuario',
        life: 5000
      })
    }
  } catch (error: any) {
    console.error('Error creating user:', error)
    console.error('Error details:', {
      message: error.message,
      data: error.data,
      statusCode: error.statusCode,
      response: error.response
    })

    const errorMessage =
      error.data?.error?.message ||
      error.data?.message ||
      error.message ||
      'Error al crear usuario'

    toast.add({
      severity: 'error',
      summary: 'Error al crear usuario',
      detail: errorMessage,
      life: 5000
    })
  } finally {
    loading.value = false
  }
}
</script>
