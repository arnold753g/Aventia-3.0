<template>
  <div class="page-shell py-12 px-4">
    <div class="max-w-md mx-auto">
      <Card class="surface-card">
        <template #title>
          <div class="text-center">
            <h2 class="text-2xl font-bold" style="color: var(--color-primary);">
              Establece tu contrasena
            </h2>
            <p class="mt-2 muted text-sm">
              Tu correo ya fue verificado. Ahora crea tu contrasena para acceder.
            </p>
          </div>
        </template>
        <template #content>
          <form class="space-y-6" @submit.prevent="handleSubmit">
            <div>
              <label for="password" class="block text-sm font-medium muted mb-2">
                Nueva contrasena
              </label>
              <Password
                id="password"
                v-model="password"
                placeholder="Minimo 8 caracteres"
                toggleMask
                :feedback="true"
                class="w-full"
                :disabled="loading"
              />
              <Message v-if="errors.password" severity="error" size="small" variant="simple">
                {{ errors.password }}
              </Message>
            </div>

            <div>
              <label for="confirmPassword" class="block text-sm font-medium muted mb-2">
                Confirmar contrasena
              </label>
              <Password
                id="confirmPassword"
                v-model="confirmPassword"
                placeholder="Repite tu contrasena"
                toggleMask
                :feedback="false"
                class="w-full"
                :disabled="loading"
              />
              <Message v-if="errors.confirmPassword" severity="error" size="small" variant="simple">
                {{ errors.confirmPassword }}
              </Message>
            </div>

            <Button
              type="submit"
              label="Establecer contrasena"
              :loading="loading"
              :disabled="!isFormValid"
              class="w-full"
            />
          </form>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { z } from 'zod'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const authAPI = useAuthAPI()

const email = ref((route.query.email as string) || '')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errors = ref<{ password?: string; confirmPassword?: string }>({})

const schema = z.object({
  password: z.string().min(8, 'La contrasena debe tener al menos 8 caracteres'),
  confirmPassword: z.string()
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Las contrasenas no coinciden',
  path: ['confirmPassword']
})

const isFormValid = computed(() => {
  return password.value.length >= 8 && password.value === confirmPassword.value
})

watch([password, confirmPassword], () => {
  try {
    schema.parse({ password: password.value, confirmPassword: confirmPassword.value })
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
    schema.parse({ password: password.value, confirmPassword: confirmPassword.value })

    loading.value = true
    const response: any = await authAPI.setInitialPassword(email.value, password.value)

    toast.add({
      severity: 'success',
      summary: 'Contrasena establecida',
      detail: response.message || 'Contrasena creada exitosamente',
      life: 5000
    })

    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.data?.error?.message || 'Error al establecer la contrasena',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (!email.value) {
    toast.add({
      severity: 'warn',
      summary: 'Email requerido',
      detail: 'Debes proporcionar un email',
      life: 3000
    })
    router.push('/login')
  }
})
</script>
