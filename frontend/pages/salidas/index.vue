<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4 py-20">
    <div class="max-w-xl w-full text-center rounded-[26px] border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_24px_60px_rgba(0,0,0,0.45)] p-8">
      <p class="text-xs uppercase tracking-[0.3em] text-white/50">Salidas habilitadas</p>
      <h1 class="text-2xl md:text-3xl font-semibold mt-3">Acceso solo con sesion</h1>
      <p class="text-white/70 mt-3">
        Para reservar una salida confirmada necesitas iniciar sesion como turista.
      </p>
      <div class="mt-8 flex flex-col sm:flex-row gap-3 justify-center">
        <Button
          label="Iniciar sesion"
          icon="pi pi-sign-in"
          class="focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
          @click="goToLogin"
        />
        <Button
          v-if="authStore.isAuthenticated"
          label="Ver salidas"
          icon="pi pi-arrow-right"
          severity="secondary"
          outlined
          class="focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
          @click="goToSalidas"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'home'
})

const authStore = useAuthStore()

const goToLogin = () => {
  navigateTo('/login')
}

const goToSalidas = () => {
  if (authStore.isAdmin) {
    navigateTo('/dashboard')
    return
  }
  if (authStore.isEncargado) {
    navigateTo('/agencia/ventas')
    return
  }
  navigateTo('/turista/salidas')
}

const redirectIfAuthenticated = () => {
  if (authStore.isAuthenticated) {
    goToSalidas()
  }
}

onMounted(() => {
  redirectIfAuthenticated()
})

watch(
  () => authStore.isAuthenticated,
  (isAuth) => {
    if (isAuth) redirectIfAuthenticated()
  }
)
</script>
