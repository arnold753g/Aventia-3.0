<template>
  <div class="page-shell py-12 px-4">
    <div class="max-w-md mx-auto">
      <Card class="surface-card">
        <template #title>
          <div class="text-center">
            <h2 class="text-2xl font-bold" style="color: var(--color-primary);">
              ¿Olvidaste tu contrasena?
            </h2>
            <p class="mt-2 muted text-sm">
              Ingresa tu correo y te enviaremos un codigo para restablecerla.
            </p>
          </div>
        </template>
        <template #content>
          <form class="space-y-6" @submit.prevent="handleSubmit">
            <div>
              <label for="email" class="block text-sm font-medium muted mb-2">
                Correo electronico
              </label>
              <InputText
                id="email"
                v-model="email"
                type="email"
                placeholder="tu@correo.com"
                class="w-full"
                :disabled="loading"
              />
              <Message v-if="errors.email" severity="error" size="small" variant="simple">
                {{ errors.email }}
              </Message>
            </div>

            <Button
              type="submit"
              label="Enviar codigo"
              :loading="loading"
              :disabled="!email"
              class="w-full"
            />

            <div class="text-center text-sm">
              <NuxtLink to="/login" class="font-semibold" style="color: var(--color-accent);">
                Volver al inicio de sesion
              </NuxtLink>
            </div>
          </form>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { z } from 'zod'

const router = useRouter()
const toast = useToast()
const authAPI = useAuthAPI()

const email = ref('')
const loading = ref(false)
const errors = ref<{ email?: string }>({})

const emailSchema = z.string().email('Correo electronico invalido')

watch(email, (value) => {
  if (!value) {
    errors.value.email = undefined
    return
  }
  try {
    emailSchema.parse(value)
    errors.value.email = undefined
  } catch (err) {
    if (err instanceof z.ZodError) {
      errors.value.email = err.errors[0]?.message
    }
  }
})

const handleSubmit = async () => {
  try {
    emailSchema.parse(email.value)
    errors.value = {}

    loading.value = true
    const response: any = await authAPI.forgotPassword(email.value)

    toast.add({
      severity: 'success',
      summary: 'Codigo enviado',
      detail: response.message || 'Revisa tu correo para el codigo de recuperacion',
      life: 5000
    })

    setTimeout(() => {
      router.push({
        path: '/reset-password',
        query: { email: email.value }
      })
    }, 1500)
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.data?.error?.message || 'Error al enviar el codigo',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}
</script>
