<template>
  <div class="page-shell py-12 px-4">
    <div class="max-w-md mx-auto">
      <Card class="surface-card">
        <template #title>
          <div class="text-center">
            <h2 class="text-2xl font-bold" style="color: var(--color-primary);">
              Verifica tu correo
            </h2>
            <p class="mt-2 muted text-sm">
              Ingresa el codigo de 6 digitos enviado a<br>
              <span class="font-semibold">{{ email }}</span>
            </p>
          </div>
        </template>
        <template #content>
          <form class="space-y-6" @submit.prevent="handleVerify">
            <div>
              <label for="code" class="block text-sm font-medium muted mb-2">
                Codigo de verificacion
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

            <div class="space-y-3">
              <Button
                type="submit"
                label="Verificar"
                :loading="loading"
                :disabled="!isCodeValid"
                class="w-full"
              />

              <Button
                type="button"
                :label="cooldown > 0 ? `Reenviar en ${cooldown}s` : 'Reenviar codigo'"
                severity="secondary"
                outlined
                :loading="resending"
                :disabled="cooldown > 0"
                class="w-full"
                @click="handleResend"
              />
            </div>

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
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useToast } from 'primevue/usetoast'
import { z } from 'zod'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const authAPI = useAuthAPI()

const email = ref((route.query.email as string) || '')
const code = ref('')
const loading = ref(false)
const resending = ref(false)
const cooldown = ref(0)
const errors = ref<{ code?: string }>({})

let cooldownTimer: ReturnType<typeof setInterval> | null = null

const codeSchema = z.string().regex(/^\d{6}$/, 'Debe ser un codigo de 6 digitos')

const isCodeValid = computed(() => code.value.length === 6 && !errors.value.code)

watch(code, (value) => {
  if (!value) {
    errors.value.code = undefined
    return
  }
  try {
    codeSchema.parse(value)
    errors.value.code = undefined
  } catch (err) {
    if (err instanceof z.ZodError) {
      errors.value.code = err.errors[0]?.message
    }
  }
})

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

const handleVerify = async () => {
  if (!email.value) {
    toast.add({
      severity: 'warn',
      summary: 'Email requerido',
      detail: 'Debes proporcionar un email para verificar',
      life: 3000
    })
    router.push('/login')
    return
  }

  try {
    codeSchema.parse(code.value)
    errors.value = {}

    loading.value = true
    const response: any = await authAPI.verifyEmail(email.value, code.value)

    toast.add({
      severity: 'success',
      summary: 'Correo verificado',
      detail: response.message || 'Correo verificado exitosamente',
      life: 5000
    })

    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.data?.error?.message || 'Error al verificar el codigo',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}

const handleResend = async () => {
  if (!email.value || cooldown.value > 0) return

  try {
    resending.value = true
    const response: any = await authAPI.resendCode(email.value)

    toast.add({
      severity: 'success',
      summary: 'Codigo enviado',
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
  } finally {
    resending.value = false
  }
}

onMounted(() => {
  if (!email.value) {
    toast.add({
      severity: 'warn',
      summary: 'Email requerido',
      detail: 'Debes proporcionar un email para verificar',
      life: 3000
    })
    router.push('/login')
  }
})

onBeforeUnmount(() => {
  if (cooldownTimer) clearInterval(cooldownTimer)
})
</script>
