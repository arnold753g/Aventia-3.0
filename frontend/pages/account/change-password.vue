<template>
  <div class="page-shell py-12 px-4">
    <div class="max-w-2xl mx-auto">
      <Card class="surface-card">
        <template #title>
          <div class="text-center">
            <h2 class="text-2xl font-bold" style="color: var(--color-primary);">
              Cambiar contrasena
            </h2>
            <p class="mt-2 muted text-sm">
              Actualiza tu contrasena para mantener tu cuenta segura.
            </p>
          </div>
        </template>
        <template #content>
          <form class="space-y-6" @submit.prevent="handleSubmit">
            <div>
              <label for="currentPassword" class="block text-sm font-medium muted mb-2">
                Contrasena actual
              </label>
              <Password
                id="currentPassword"
                v-model="currentPassword"
                placeholder="Tu contrasena actual"
                toggleMask
                :feedback="false"
                class="w-full"
                :disabled="loading"
              />
              <Message v-if="errors.currentPassword" severity="error" size="small" variant="simple">
                {{ errors.currentPassword }}
              </Message>
            </div>

            <div>
              <label for="newPassword" class="block text-sm font-medium muted mb-2">
                Nueva contrasena
              </label>
              <Password
                id="newPassword"
                v-model="newPassword"
                placeholder="Minimo 8 caracteres"
                toggleMask
                :feedback="true"
                class="w-full"
                :disabled="loading"
              />
              <Message v-if="errors.newPassword" severity="error" size="small" variant="simple">
                {{ errors.newPassword }}
              </Message>
            </div>

            <div>
              <label for="confirmPassword" class="block text-sm font-medium muted mb-2">
                Confirmar nueva contrasena
              </label>
              <Password
                id="confirmPassword"
                v-model="confirmPassword"
                placeholder="Repite tu nueva contrasena"
                toggleMask
                :feedback="false"
                class="w-full"
                :disabled="loading"
              />
              <Message v-if="errors.confirmPassword" severity="error" size="small" variant="simple">
                {{ errors.confirmPassword }}
              </Message>
            </div>

            <div class="flex flex-col sm:flex-row gap-3">
              <Button
                type="submit"
                label="Cambiar contrasena"
                :loading="loading"
                :disabled="!isFormValid"
              />
              <Button
                type="button"
                label="Cancelar"
                severity="secondary"
                outlined
                @click="router.back()"
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
import { ref, computed, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { z } from 'zod'

definePageMeta({
  middleware: 'auth',
  layout: 'home'
})

const router = useRouter()
const toast = useToast()
const authAPI = useAuthAPI()

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errors = ref<{ currentPassword?: string; newPassword?: string; confirmPassword?: string }>({})

const schema = z.object({
  currentPassword: z.string().min(1, 'Contrasena actual requerida'),
  newPassword: z.string().min(8, 'La nueva contrasena debe tener al menos 8 caracteres'),
  confirmPassword: z.string()
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: 'Las contrasenas no coinciden',
  path: ['confirmPassword']
}).refine((data) => data.currentPassword !== data.newPassword, {
  message: 'La nueva contrasena debe ser diferente a la actual',
  path: ['newPassword']
})

const isFormValid = computed(() => {
  return currentPassword.value &&
    newPassword.value.length >= 8 &&
    newPassword.value === confirmPassword.value &&
    currentPassword.value !== newPassword.value
})

watch([currentPassword, newPassword, confirmPassword], () => {
  try {
    schema.parse({
      currentPassword: currentPassword.value,
      newPassword: newPassword.value,
      confirmPassword: confirmPassword.value
    })
    errors.value = {}
  } catch (err) {
    if (err instanceof z.ZodError) {
      errors.value = {}
      err.errors.forEach((error) => {
        const field = error.path[0] as keyof typeof errors.value
        errors.value[field] = error.message
      })
    }
  }
})

const handleSubmit = async () => {
  try {
    schema.parse({
      currentPassword: currentPassword.value,
      newPassword: newPassword.value,
      confirmPassword: confirmPassword.value
    })

    loading.value = true
    const response: any = await authAPI.changePassword(currentPassword.value, newPassword.value)

    toast.add({
      severity: 'success',
      summary: 'Contrasena actualizada',
      detail: response.message || 'Contrasena actualizada exitosamente',
      life: 8000
    })

    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''

    setTimeout(() => {
      router.back()
    }, 2000)
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.data?.error?.message || 'Error al cambiar la contrasena',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}
</script>
