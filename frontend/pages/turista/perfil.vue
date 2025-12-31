<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-4xl mx-auto px-4 py-6">
        <div class="flex items-center gap-4">
          <Button icon="pi pi-arrow-left" text rounded @click="navigateTo('/turista/dashboard')" />
          <div>
            <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
              Mi perfil
            </h1>
            <p class="muted mt-1">Actualiza tu información personal.</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="max-w-4xl mx-auto px-4 py-8">
      <Skeleton height="420px" />
    </div>

    <div v-else-if="initialValues" class="max-w-4xl mx-auto px-4 py-8">
      <Card class="surface-card">
        <template #content>
          <Form v-slot="form" :initialValues="initialValues" :resolver="resolver" @submit="handleSubmit">
            <div class="space-y-8">
              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-user"></i>
                  Datos personales
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <FormField v-slot="$field" name="nombre">
                    <label class="block text-sm font-medium muted mb-2">Nombre</label>
                    <InputText placeholder="Ej: Juan" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                  <FormField v-slot="$field" name="apellido_paterno">
                    <label class="block text-sm font-medium muted mb-2">Apellido paterno</label>
                    <InputText placeholder="Ej: Pérez" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>

                  <FormField v-slot="$field" name="apellido_materno">
                    <label class="block text-sm font-medium muted mb-2">Apellido materno</label>
                    <InputText placeholder="Ej: López" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>
                </div>
              </div>

              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-phone"></i>
                  Contacto
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div>
                    <label class="block text-sm font-medium muted mb-2">Telefono</label>
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

                  <FormField v-slot="$field" name="ciudad">
                    <label class="block text-sm font-medium muted mb-2">Ciudad</label>
                    <InputText placeholder="Ej: Tarija" class="w-full" />
                    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                      {{ $field.error?.message }}
                    </Message>
                  </FormField>
                </div>

                <div class="mt-4">
                  <label class="block text-sm font-medium muted mb-2">Email</label>
                  <InputText :modelValue="usuario?.email" class="w-full" disabled />
                </div>
              </div>

              <div>
                <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" style="color: var(--color-primary);">
                  <i class="pi pi-image"></i>
                  Foto de perfil
                </h3>
                <div class="flex items-center gap-4">
                  <UserAvatar
                    :src="preview || resolvePhoto(usuario?.profile_photo)"
                    :nombre="initialValues.nombre"
                    :apellido="initialValues.apellido_paterno"
                    size="xl"
                  />
                  <div class="flex-1">
                    <input
                      ref="fileInput"
                      type="file"
                      accept="image/*"
                      class="hidden"
                      @change="onFileChange"
                    />
                    <div class="flex flex-wrap gap-2">
                      <Button
                        label="Elegir foto"
                        icon="pi pi-upload"
                        severity="secondary"
                        outlined
                        @click="fileInput?.click()"
                      />
                      <Button
                        v-if="selectedFile"
                        label="Quitar"
                        icon="pi pi-times"
                        severity="danger"
                        text
                        @click="clearSelectedFile"
                      />
                    </div>
                    <p class="text-xs muted mt-2">Formatos aceptados: JPG, PNG, WEBP. Tamaño máximo: 2MB</p>
                  </div>
                </div>
              </div>

              <div class="flex justify-end gap-4 pt-4 border-t">
                <Button
                  label="Cancelar"
                  severity="secondary"
                  outlined
                  @click="navigateTo('/turista/dashboard')"
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
import { onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import Form from '@primevue/forms/form'
import FormField from '@primevue/forms/formfield'
import { zodResolver } from '@primevue/forms/resolvers/zod'
import { updateUsuarioSchema } from '~/utils/validations/usuario'
import { useAuthStore } from '~/stores/auth'
import UserAvatar from '~/components/usuarios/UserAvatar.vue'
import { buildPhone, PHONE_PREFIXES, splitPhone } from '~/utils/phone'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const authStore = useAuthStore()
const { getUsuario, updateUsuario } = useUsuarios()
const apiOrigin = new URL(useRuntimeConfig().public.apiBase).origin

const loading = ref(false)
const saving = ref(false)
const usuario = ref<any>(null)
const initialValues = ref<any>(null)

const preview = ref<string | null>(null)
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const resolver = zodResolver(updateUsuarioSchema)

const phonePrefixes = PHONE_PREFIXES

const resolvePhoto = (path?: string) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${apiOrigin}/${path.replace(/^\//, '')}`
}

const loadUsuario = async () => {
  const userId = authStore.user?.id
  if (!userId) {
    await authStore.getProfile()
  }

  const id = authStore.user?.id
  if (!id) {
    navigateTo('/login')
    return
  }

  loading.value = true
  try {
    const response: any = await getUsuario(id)
    if (response.success) {
      usuario.value = response.data
      const { prefix, number } = splitPhone(response.data.phone)
      initialValues.value = {
        nombre: response.data.nombre,
        apellido_paterno: response.data.apellido_paterno,
        apellido_materno: response.data.apellido_materno,
        phone_prefix: prefix,
        phone_number: number,
        ciudad: response.data.ciudad || ''
      }
      return
    }
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: response?.error?.message || 'No se pudo cargar tu perfil',
      life: 3000
    })
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || error.message || 'No se pudo cargar tu perfil',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const file = target.files[0]

  if (file.size > 2 * 1024 * 1024) {
    toast.add({
      severity: 'warn',
      summary: 'Archivo muy grande',
      detail: 'El tamaño máximo permitido es 2MB',
      life: 3000
    })
    return
  }

  if (!file.type.startsWith('image/')) {
    toast.add({
      severity: 'warn',
      summary: 'Tipo de archivo inválido',
      detail: 'Solo se permiten imágenes',
      life: 3000
    })
    return
  }

  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = () => {
    preview.value = reader.result as string
  }
  reader.readAsDataURL(file)
}

const clearSelectedFile = () => {
  selectedFile.value = null
  preview.value = null
  if (fileInput.value) fileInput.value.value = ''
}

const handleSubmit = async ({ valid, values }: any) => {
  if (!valid) {
    toast.add({
      severity: 'warn',
      summary: 'Formulario incompleto',
      detail: 'Por favor complete los campos correctamente',
      life: 3000
    })
    return
  }

  const id = authStore.user?.id
  if (!id) return

  saving.value = true
  try {
    const payload: any = {
      ...values,
      phone: buildPhone(values.phone_prefix, values.phone_number)
    }
    delete payload.phone_prefix
    delete payload.phone_number
    if (selectedFile.value) payload.profile_photo = selectedFile.value

    const response: any = await updateUsuario(id, payload)
    if (response.success) {
      toast.add({
        severity: 'success',
        summary: 'Guardado',
        detail: 'Tu perfil fue actualizado exitosamente',
        life: 3000
      })

      await authStore.getProfile()
      usuario.value = { ...usuario.value, ...response.data }
      clearSelectedFile()
      return
    }

    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: response?.error?.message || 'No se pudo actualizar tu perfil',
      life: 5000
    })
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || error.message || 'No se pudo actualizar tu perfil',
      life: 5000
    })
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadUsuario()
})
</script>
