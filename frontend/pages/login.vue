<template>
  <section class="relative min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 py-12">
    <div class="pointer-events-none absolute inset-0">
      <div class="absolute -top-24 left-1/4 h-72 w-72 bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.18),transparent_60%)]"></div>
      <div class="absolute bottom-16 right-10 h-80 w-80 bg-[radial-gradient(circle_at_center,rgba(59,130,246,0.18),transparent_60%)]"></div>
    </div>

    <div class="relative w-full max-w-6xl grid gap-8 lg:grid-cols-[1.1fr_0.9fr]">
      <div class="relative overflow-hidden rounded-[28px] border border-white/10 bg-white/5 backdrop-blur-xl p-8 shadow-[0_30px_60px_rgba(0,0,0,0.45)]">
        <div class="absolute inset-0 bg-gradient-to-br from-white/10 via-transparent to-transparent"></div>
        <div class="relative space-y-6">
          <div class="flex items-center gap-4">
            <div class="h-12 w-12 rounded-2xl border border-white/20 bg-white/10 flex items-center justify-center">
              <i class="pi pi-compass text-xl text-white"></i>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.3em] text-white/60">Andaria</p>
              <p class="text-2xl font-semibold text-white">Bienvenido</p>
            </div>
          </div>

          <p class="text-white/70 text-lg">
            Accede para gestionar reservas, compras de paquetes y salidas confirmadas.
          </p>

          <div class="space-y-4 text-sm text-white/60">
            <div class="flex items-start gap-3">
              <span class="mt-1 h-2 w-2 rounded-full bg-emerald-400"></span>
              <p>Explora atracciones turisticas de Tarija.</p>
            </div>
            <div class="flex items-start gap-3">
              <span class="mt-1 h-2 w-2 rounded-full bg-sky-400"></span>
              <p>Paquetes turisticos de agencias verificadas.</p>
            </div>
            <div class="flex items-start gap-3">
              <span class="mt-1 h-2 w-2 rounded-full bg-amber-400"></span>
              <p>Salidas habilitadas y promociones.</p>
            </div>
          </div>

          <div class="flex flex-wrap gap-3 pt-2">
            <NuxtLink
              to="/atracciones"
              class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-4 py-2 text-xs text-white/80 hover:border-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
            >
              Ver atracciones
              <i class="pi pi-arrow-right text-xs"></i>
            </NuxtLink>
            <NuxtLink
              to="/paquetes"
              class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-4 py-2 text-xs text-white/80 hover:border-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
            >
              Ver paquetes
              <i class="pi pi-arrow-right text-xs"></i>
            </NuxtLink>
          </div>
        </div>
      </div>

      <div class="rounded-[28px] border border-white/10 bg-white/5 backdrop-blur-xl p-8 shadow-[0_30px_60px_rgba(0,0,0,0.45)]">
        <div class="text-center space-y-2">
          <p class="text-xs uppercase tracking-[0.3em] text-white/60">Cuenta</p>
          <h2 class="text-3xl md:text-4xl font-semibold text-white">Iniciar sesion</h2>
          <p class="text-sm text-white/60">Sistema Andaria</p>
        </div>

        <form @submit.prevent="handleLogin" class="mt-8 space-y-6">
          <div>
            <label for="email" class="block text-sm text-white/70 mb-2">
              Email
            </label>
            <InputText
              id="email"
              v-model="email"
              type="email"
              placeholder="tu@email.com"
              autocomplete="email"
              class="w-full bg-white/5 border border-white/10 text-white placeholder-white/40 focus:border-white/30 focus:ring-2 focus:ring-white/20"
              required
              :disabled="loading"
            />
          </div>

          <div>
            <label for="password" class="block text-sm text-white/70 mb-2">
              Contrasena
            </label>
            <InputText
              id="password"
              v-model="password"
              type="password"
              placeholder="********"
              autocomplete="current-password"
              class="w-full bg-white/5 border border-white/10 text-white placeholder-white/40 focus:border-white/30 focus:ring-2 focus:ring-white/20"
              required
              :disabled="loading"
            />
          </div>

          <div class="text-right text-sm">
            <NuxtLink to="/forgot-password" class="text-white/70 hover:text-white">
              ¿Olvidaste tu contrasena?
            </NuxtLink>
          </div>

          <Button
            type="submit"
            label="Ingresar"
            icon="pi pi-sign-in"
            class="w-full !bg-white !text-black hover:!bg-white/90 focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
            :loading="loading"
          />
        </form>

        <div class="mt-6 text-center text-sm text-white/60">
          No tienes cuenta?
          <NuxtLink to="/registro" class="font-semibold text-white hover:text-white/90">
            Registrate aqui
          </NuxtLink>
        </div>
        
      </div>
    </div>

    <Toast />
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'home'
})

const toast = useToast()
const authStore = useAuthStore()
const route = useRoute()

const email = ref('')
const password = ref('')
const loading = ref(false)

const fillCredentials = (emailVal: string, passVal: string) => {
  email.value = emailVal
  password.value = passVal
}

const getDefaultHome = () => {
  if (authStore.isAdmin) return '/admin/dashboard'
  if (authStore.isEncargado) return '/agencia/dashboard'
  if (authStore.isTurista) return '/turista/dashboard'
  return '/dashboard'
}

const getSafeRedirect = () => {
  const redirect = route.query.redirect
  if (typeof redirect !== 'string') return null
  if (!redirect.startsWith('/') || redirect.startsWith('//')) return null
  if (redirect.startsWith('/login') || redirect.startsWith('/registro')) return null
  if (redirect.startsWith('/admin') && !authStore.isAdmin) return null
  if (redirect.startsWith('/agencia') && !authStore.isEncargado) return null
  if (redirect.startsWith('/turista') && !authStore.isTurista) return null
  return redirect
}

const handleLogin = async () => {
  loading.value = true

  try {
    const result = await authStore.login(email.value, password.value)

    if (result.success) {
      toast.add({
        severity: 'success',
        summary: 'Bienvenido',
        detail: 'Has iniciado sesion exitosamente',
        life: 3000
      })

      const target = getSafeRedirect() || getDefaultHome()
      setTimeout(() => {
        navigateTo(target)
      }, 500)
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error de autenticacion',
        detail: result.error || 'Credenciales invalidas',
        life: 5000
      })
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.message || 'Error al iniciar sesion',
      life: 5000
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (authStore.isAuthenticated) {
    const target = getSafeRedirect() || getDefaultHome()
    navigateTo(target)
  }
})
</script>
