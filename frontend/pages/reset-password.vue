<template>
  <div class="page-shell py-12 px-4">
    <div class="max-w-md mx-auto">
      <Card class="surface-card">
        <template #title>
          <div class="text-center">
            <h2 class="text-2xl font-bold" style="color: var(--color-primary);">
              Restablecer contrasena
            </h2>
            <p class="mt-2 muted text-sm">
              Ingresa el codigo enviado a tu correo y tu nueva contrasena.
            </p>
          </div>
        </template>
        <template #content>
          <form class="space-y-6" @submit.prevent="handleSubmit">
            <div>
              <label for="code" class="block text-sm font-medium muted mb-2">
                Codigo de recuperacion
              </label>
              <InputText
                id="code"
                v-model="code"
                placeholder="000000"
                class="w-full text-center text-2xl tracking-widest"
                maxlength="6"
                :disabled="loading"
              />
              <Message v-if="errors.code" severity="error" size="small" variant="simple">
                {{ errors.code }}
              </Message>
            </div>

            <div>
              <label for="password" class="block text-sm font-medium muted mb-2">
                Nueva contrasena
              </label>
              <Password
                id="password"
                v-model="password"
                placeholder="Minimo 8 caracteres"
                toggleMask
                :feedback="false"
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
              label="Restablecer contrasena"
              :loading="loading"
              :disabled="!isFormValid"
              class="w-full"
            />

            <div class="text-center text-sm space-y-2">
              <button
                type="button"
                class="font-semibold disabled:opacity-60"
                style="color: var(--color-accent);"
                :disabled="cooldown > 0"
                @click="handleResend"
              >
                <span v-if="cooldown > 0">Reenviar codigo en {{ cooldown }}s</span>
                <span v-else>Reenviar codigo</span>
              </button>
              <div>
                <NuxtLink to="/login" class="font-semibold" style="color: var(--color-accent);">
                  Volver al inicio de sesion
                </NuxtLink>
              </div>
            </div>
          </form>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue'
import { useToast } from 'primevue/usetoast'
import { z } from 'zod'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const authAPI = useAuthAPI()

const email = ref((route.query.email as string) || '')
const code = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const cooldown = ref(0)
const errors = ref<{ code?: string; password?: string; confirmPassword?: string }>({})

let cooldownTimer: ReturnType<typeof setInterval> | null = null

const codeSchema = z.string().regex(/^\d{6}$/, 'Debe ser un codigo de 6 digitos')
const passwordSchema = z.string().min(8, 'La contrasena debe tener al menos 8 caracteres')

const schema = z.object({
  code: codeSchema,
  password: passwordSchema,
  confirmPassword: z.string()
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Las contrasenas no coinciden',
  path: ['confirmPassword']
})

const isFormValid = computed(() => {
  return code.value.length === 6 &&
    password.value.length >= 8 &&
    password.value === confirmPassword.value
})

const validateField = (field: 'code' | 'password' | 'confirmPassword') => {
  try {
    if (field === 'code') codeSchema.parse(code.value)
    if (field === 'password') passwordSchema.parse(password.value)
    if (field === 'confirmPassword') {
      if (password.value !== confirmPassword.value) {
        errors.value.confirmPassword = 'Las contrasenas no coinciden'
        return
      }
    }
    errors.value[field] = undefined
  } catch (err) {
    if (err instanceof z.ZodError) {
      errors.value[field] = err.errors[0]?.message
    }
  }
}

watch(code, () => validateField('code'))
watch(password, () => validateField('password'))
watch(confirmPassword, () => validateField('confirmPassword'))

const startCooldown = () => {
  if (cooldownTimer) clearInterval(cooldownTimer)
  cooldown.value = 60
  cooldownTimer = setInterval(() => {
    cooldown.value -= 1
    if (cooldown.value <= 0 && cooldownTimer) {
      clearInterval(cooldownTimer)
      cooldownTimer = null
    }
  }, 1000)
}

const handleSubmit = async () => {
  if (!email.value) return
  try {
    schema.parse({ code: code.value, password: password.value, confirmPassword: confirmPassword.value })
    errors.value = {}

    loading.value = true
    const response: any = await authAPI.resetPassword(email.value, code.value, password.value)

    toast.add({
      severity: 'success',
      summary: 'Contrasena restablecida',
      detail: response.message || 'Contrasena actualizada exitosamente',
      life: 5000
    })

    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.data?.error?.message || 'Error al restablecer la contrasena',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}

const handleResend = async () => {
  if (!email.value || cooldown.value > 0) return
  try {
    const response: any = await authAPI.forgotPassword(email.value)
    toast.add({
      severity: 'success',
      summary: 'Codigo reenviado',
      detail: response.message || 'Se envio un nuevo codigo a tu correo',
      life: 5000
    })
    startCooldown()
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.data?.error?.message || 'Error al reenviar el codigo',
      life: 5000
    })
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
    router.push('/forgot-password')
  }
})

onBeforeUnmount(() => {
  if (cooldownTimer) clearInterval(cooldownTimer)
})
</script>
