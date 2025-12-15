<template>
  <div class="page-shell py-12 px-4">
    <div class="max-w-2xl mx-auto">
      <Card class="surface-card">
        <template #title>
          <div class="text-center">
            <h2 class="text-3xl font-bold" style="color: var(--color-primary);">
              Registro de usuario
            </h2>
            <p class="mt-2 muted">
              Crea tu cuenta en Sistema Andaria
            </p>
          </div>
        </template>
        <template #content>
          <Form v-slot="form" :initialValues="initialValues" :resolver="resolver" @submit="handleSubmit">
            <div class="space-y-6">
              <!-- Datos Personales -->
              <div class="grid md:grid-cols-3 gap-4">
                <FormField v-slot="$field" name="nombre">
                  <label class="block text-sm font-medium muted mb-2">Nombre *</label>
                  <InputText
                    v-model="formData.nombre"
                    placeholder="Ej: Carlos"
                    class="w-full"
                    :disabled="loading"
                  />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>

                <FormField v-slot="$field" name="apellido_paterno">
                  <label class="block text-sm font-medium muted mb-2">Apellido paterno *</label>
                  <InputText
                    v-model="formData.apellido_paterno"
                    placeholder="Ej: López"
                    class="w-full"
                    :disabled="loading"
                  />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>

                <FormField v-slot="$field" name="apellido_materno">
                  <label class="block text-sm font-medium muted mb-2">Apellido materno</label>
                  <InputText
                    v-model="formData.apellido_materno"
                    placeholder="Ej: Pérez"
                    class="w-full"
                    :disabled="loading"
                  />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
              </div>

              <!-- Identificación -->
              <div class="grid md:grid-cols-2 gap-4">
                <FormField v-slot="$field" name="ci">
                  <label class="block text-sm font-medium muted mb-2">
                    Carnet de identidad *
                    <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'5-15 caracteres alfanuméricos'"></i>
                  </label>
                  <InputText placeholder="Ej: 1234567" class="w-full" :disabled="loading" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                  <Message v-if="ciExists" severity="error" size="small" variant="simple">
                    Este CI ya está registrado
                  </Message>
                </FormField>

                <FormField v-slot="$field" name="expedido">
                  <label class="block text-sm font-medium muted mb-2">Expedido en *</label>
                  <Select
                    :options="departamentos"
                    optionLabel="label"
                    optionValue="value"
                    placeholder="Seleccione..."
                    class="w-full"
                    :disabled="loading"
                  />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
              </div>

              <!-- Fecha y Teléfono -->
              <div class="grid md:grid-cols-2 gap-4">
                <FormField v-slot="{ props, value, invalid, error }" name="fecha_nacimiento">
                  <label class="block text-sm font-medium muted mb-2">
                    Fecha de nacimiento *
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
                    :disabled="loading"
                  />
                  <Message v-if="invalid" severity="error" size="small" variant="simple">
                    {{ error }}
                  </Message>
                </FormField>

                <FormField v-slot="$field" name="phone">
                  <label class="block text-sm font-medium muted mb-2">
                    Teléfono *
                    <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'8 dígitos, inicia con 6 o 7'"></i>
                  </label>
                  <InputText placeholder="Ej: 71234567" class="w-full" maxlength="8" :disabled="loading" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
              </div>

              <!-- Email y Ciudad -->
              <div class="grid md:grid-cols-2 gap-4">
                <FormField v-slot="$field" name="email">
                  <label class="block text-sm font-medium muted mb-2">Email *</label>
                  <InputText type="email" placeholder="usuario@email.com" class="w-full" :disabled="loading" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                  <Message v-if="emailExists" severity="error" size="small" variant="simple">
                    Este email ya está registrado
                  </Message>
                </FormField>

                <FormField v-slot="$field" name="ciudad">
                  <label class="block text-sm font-medium muted mb-2">Ciudad</label>
                  <InputText placeholder="Ej: La Paz" class="w-full" :disabled="loading" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
              </div>

              <!-- Foto de Perfil -->
              <div>
                <label class="block text-sm font-medium muted mb-3">
                  <i class="pi pi-image mr-2"></i>
                  Foto de Perfil (Opcional)
                </label>
                <div class="flex flex-col md:flex-row items-start md:items-center gap-4">
                  <!-- Preview del avatar -->
                  <div class="flex-shrink-0">
                    <div v-if="preview" class="relative">
                      <img
                        :src="preview"
                        alt="Preview"
                        class="w-24 h-24 rounded-full object-cover border-2 border-gray-200"
                      />
                      <button
                        type="button"
                        @click="removePhoto"
                        class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center hover:bg-red-600 transition-colors"
                        :disabled="loading"
                      >
                        <i class="pi pi-times text-xs"></i>
                      </button>
                    </div>
                    <div v-else class="w-24 h-24 rounded-full bg-gray-100 flex items-center justify-center border-2 border-gray-200">
                      <i class="pi pi-user text-4xl text-gray-400"></i>
                    </div>
                  </div>

                  <!-- Controles de carga -->
                  <div class="flex-1">
                    <input
                      type="file"
                      accept="image/jpeg,image/jpg,image/png,image/webp"
                      @change="onFileChange"
                      class="hidden"
                      ref="fileInput"
                    />
                    <Button
                      type="button"
                      :label="preview ? 'Cambiar foto' : 'Seleccionar foto'"
                      icon="pi pi-upload"
                      severity="secondary"
                      outlined
                      @click="fileInput?.click()"
                      :disabled="loading"
                    />
                    <p class="text-xs muted mt-2">
                      Formatos: JPG, PNG, WebP. Tamaño máximo: 2MB
                    </p>
                  </div>
                </div>
              </div>

              <!-- Contraseñas -->
              <div class="grid md:grid-cols-2 gap-4">
                <FormField v-slot="$field" name="password">
                  <label class="block text-sm font-medium muted mb-2">
                    Contraseña *
                    <i class="pi pi-info-circle text-xs ml-1" v-tooltip.right="'Mín. 8 caracteres, incluye mayúscula, minúscula, número y carácter especial'"></i>
                  </label>
                  <Password placeholder="Contraseña segura" toggleMask :feedback="true" class="w-full" fluid :disabled="loading" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>

                <FormField v-slot="$field" name="confirmPassword">
                  <label class="block text-sm font-medium muted mb-2">Confirmar contraseña *</label>
                  <Password placeholder="Repite la contraseña" toggleMask :feedback="false" class="w-full" fluid :disabled="loading" />
                  <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                    {{ $field.error?.message }}
                  </Message>
                </FormField>
              </div>

              <!-- Botón de registro -->
              <Button
                type="submit"
                label="Registrarse"
                icon="pi pi-user-plus"
                class="w-full p-button-lg"
                :loading="loading"
                :disabled="!form.valid || emailExists || ciExists || checking"
              />
            </div>
          </Form>

          <div class="mt-6 text-center">
            <p class="text-sm muted">
              Ya tienes cuenta?
              <NuxtLink to="/login" class="font-semibold" style="color: var(--color-accent);">
                Inicia sesión aquí
              </NuxtLink>
            </p>
          </div>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'
import DatePicker from 'primevue/datepicker'
import Form from '@primevue/forms/form'
import FormField from '@primevue/forms/formfield'
import { zodResolver } from '@primevue/forms/resolvers/zod'
import { registroPublicoSchema } from '~/utils/validations/usuario'

const toast = useToast()
const authStore = useAuthStore()
const runtimeConfig = useRuntimeConfig()

const resolver = zodResolver(registroPublicoSchema)

// Estado del formulario
const initialValues = {
  nombre: '',
  apellido_paterno: '',
  apellido_materno: '',
  ci: '',
  expedido: '',
  phone: '',
  fecha_nacimiento: '',
  email: '',
  ciudad: '',
  password: '',
  confirmPassword: ''
}

const loading = ref(false)
const emailExists = ref(false)
const ciExists = ref(false)
const checking = ref(false)

// Estado de la foto
const preview = ref<string | null>(null)
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

let emailTimer: ReturnType<typeof setTimeout> | null = null
let ciTimer: ReturnType<typeof setTimeout> | null = null

const apiBase = runtimeConfig.public.apiBase

// Estado del formulario para v-model
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

// Opciones de departamentos
const departamentos = [
  { label: 'La Paz', value: 'LP' },
  { label: 'Cochabamba', value: 'CB' },
  { label: 'Santa Cruz', value: 'SC' },
  { label: 'Potosí', value: 'PT' },
  { label: 'Oruro', value: 'OR' },
  { label: 'Tarija', value: 'TJ' },
  { label: 'Chuquisaca', value: 'CH' },
  { label: 'Beni', value: 'BE' },
  { label: 'Pando', value: 'PD' }
]

// Funciones para manejo de fechas
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

// Manejo de archivo de foto
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
      detail: 'Solo se permiten imágenes (JPG, PNG, WebP)',
      life: 3000
    })
    return
  }

  selectedFile.value = file

  // Crear preview
  const reader = new FileReader()
  reader.onload = (e) => {
    preview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)

  console.log('Foto seleccionada:', {
    name: file.name,
    size: file.size,
    type: file.type
  })
}

const removePhoto = () => {
  selectedFile.value = null
  preview.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
  console.log('Foto removida')
}

// Verificar duplicados en el backend
const checkDuplicados = async (email?: string, ci?: string) => {
  if (!email && !ci) return

  try {
    checking.value = true
    const query = new URLSearchParams()
    if (email) query.append('email', email)
    if (ci) query.append('ci', ci)

    const res: any = await $fetch(`${apiBase}/usuarios/check?${query.toString()}`)
    emailExists.value = !!res.data?.emailExists
    ciExists.value = !!res.data?.ciExists
  } catch (error) {
    // Silencioso para no bloquear la UI
  } finally {
    checking.value = false
  }
}

// Watchers para verificación en tiempo real
watch(
  () => initialValues.email,
  (val) => {
    if (emailTimer) clearTimeout(emailTimer)
    emailTimer = setTimeout(() => {
      if (val) {
        checkDuplicados(val, undefined)
      } else {
        emailExists.value = false
      }
    }, 400)
  }
)

watch(
  () => initialValues.ci,
  (val) => {
    if (ciTimer) clearTimeout(ciTimer)
    ciTimer = setTimeout(() => {
      if (val) {
        checkDuplicados(undefined, val)
      } else {
        ciExists.value = false
      }
    }, 400)
  }
)

// Manejo del submit
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

  if (emailExists.value || ciExists.value) {
    toast.add({
      severity: 'error',
      summary: 'Datos duplicados',
      detail: 'El email o CI ya están registrados',
      life: 3000
    })
    return
  }

  loading.value = true

  try {
    // SEGURIDAD: El rol siempre es 'turista' para registros públicos
    // No se permite especificar el rol desde el frontend
    const payload: any = {
      nombre: values.nombre,
      apellido_paterno: values.apellido_paterno,
      apellido_materno: values.apellido_materno || '',
      ci: values.ci,
      expedido: values.expedido,
      fecha_nacimiento: values.fecha_nacimiento,
      phone: values.phone,
      email: values.email,
      ciudad: values.ciudad || '',
      password: values.password,
      rol: 'turista' // HARDCODED por seguridad
    }

    // Agregar foto si fue seleccionada
    if (selectedFile.value) {
      console.log('Agregando foto al payload:', {
        name: selectedFile.value.name,
        size: selectedFile.value.size,
        type: selectedFile.value.type
      })
      payload.profile_photo = selectedFile.value
    }

    console.log('Enviando registro con datos:', Object.keys(payload))
    const result = await authStore.register(payload)

    if (result.success) {
      toast.add({
        severity: 'success',
        summary: 'Registro exitoso',
        detail: 'Tu cuenta ha sido creada correctamente',
        life: 3000
      })

      setTimeout(() => {
        navigateTo('/turista/dashboard')
      }, 1500)
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error de registro',
        detail: result.error || 'Error al crear la cuenta',
        life: 5000
      })
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.message || 'Error al registrar usuario',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}
</script>
