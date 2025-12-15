<template>
  <div class="page-shell flex items-center justify-center px-4 py-10">
    <Card class="w-full max-w-md surface-card">
      <template #title>
        <div class="text-center">
          <h2 class="text-3xl font-bold" style="color: var(--color-primary);">
            Iniciar sesion
          </h2>
          <p class="mt-2 muted">
            Sistema Andaria
          </p>
        </div>
      </template>
      <template #content>
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label for="email" class="block text-sm font-medium muted mb-2">
              Email
            </label>
            <InputText
              id="email"
              v-model="email"
              type="email"
              placeholder="tu@email.com"
              class="w-full"
              required
              :disabled="loading"
            />
          </div>

          <div>
            <label for="password" class="block text-sm font-medium muted mb-2">
              Contrasena
            </label>
            <InputText
              id="password"
              v-model="password"
              type="password"
              placeholder="********"
              class="w-full"
              required
              :disabled="loading"
            />
          </div>

          <Button
            type="submit"
            label="Ingresar"
            icon="pi pi-sign-in"
            class="w-full p-button-lg"
            :loading="loading"
          />
        </form>

        <div class="mt-6 text-center">
          <p class="text-sm muted">
            No tienes cuenta?
            <NuxtLink to="/registro" class="font-semibold" style="color: var(--color-accent);">
              Registrate aqui
            </NuxtLink>
          </p>
        </div>

        <div class="mt-8 p-4 rounded-lg" style="background: var(--color-neutral);">
          <p class="text-sm font-semibold" style="color: var(--color-primary); margin-bottom: 0.5rem;">
            Usuarios de prueba:
          </p>
          <div class="space-y-2">
            <div class="flex items-center justify-between text-xs">
              <span class="muted">Admin: admin@andaria.bo</span>
              <Button
                label="Usar"
                size="small"
                text
                @click="fillCredentials('admin@andaria.bo', 'admin123')"
                :disabled="loading"
              />
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="muted">Turista: juan.perez@email.com</span>
              <Button
                label="Usar"
                size="small"
                text
                @click="fillCredentials('juan.perez@email.com', 'turista123')"
                :disabled="loading"
              />
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="muted">Agencia: maria.lopez@agencia.com</span>
              <Button
                label="Usar"
                size="small"
                text
                @click="fillCredentials('maria.lopez@agencia.com', 'agencia123')"
                :disabled="loading"
              />
            </div>
          </div>
        </div>
      </template>
    </Card>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'

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
