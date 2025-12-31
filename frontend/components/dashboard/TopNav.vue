<template>
  <header class="fixed top-0 inset-x-0 z-50">
    <div class="bg-black/40 backdrop-blur-xl border-b border-white/10">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center gap-3">
            <NuxtLink to="/" class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-2xl border border-white/15 bg-white/10 flex items-center justify-center">
                <i class="pi pi-compass text-lg text-white"></i>
              </div>
              <div class="leading-tight">
                <p class="text-sm uppercase tracking-[0.3em] text-white/50">Inicio</p>
                <p class="text-lg font-semibold text-white">ANDARIA</p>
              </div>
            </NuxtLink>
          </div>

          <nav class="hidden md:flex items-center gap-6 text-sm text-white/70">
            <NuxtLink
              v-for="link in navLinks"
              :key="link.to"
              :to="link.to"
              class="transition-colors hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
              :class="isActive(link.to) ? 'text-white' : ''"
            >
              {{ link.label }}
            </NuxtLink>
          </nav>

          <div class="flex items-center gap-3">

            <button
              type="button"
              class="hidden sm:flex items-center gap-3 px-3 py-2 rounded-full border border-white/10 bg-white/5 hover:border-white/25 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
              @click="handleAccountClick"
            >
              <div class="h-9 w-9 rounded-full bg-gradient-to-br from-white/30 to-white/10 flex items-center justify-center text-xs font-semibold uppercase">
                {{ userInitials }}
              </div>
              <div class="text-left leading-tight">
                <p class="text-sm text-white">{{ userLabel }}</p>
                <p class="text-xs text-white/50">{{ userRole }}</p>
              </div>
            </button>

            <button
              type="button"
              class="sm:hidden h-10 w-10 rounded-full border border-white/15 bg-white/5 flex items-center justify-center focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
              aria-label="Abrir menu"
              @click="mobileMenuOpen = true"
            >
              <i class="pi pi-bars text-white"></i>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 md:hidden transition-opacity"
      :class="mobileMenuOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'"
      @click="mobileMenuOpen = false"
    ></div>

    <div
      class="fixed top-16 right-0 w-72 max-w-[85vw] h-[calc(100vh-4rem)] z-50 md:hidden transform transition-transform"
      :class="mobileMenuOpen ? 'translate-x-0' : 'translate-x-full'"
    >
      <div class="h-full p-4 bg-[#0b0b0b] border-l border-white/10 shadow-2xl flex flex-col gap-5">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-full bg-gradient-to-br from-white/30 to-white/10 flex items-center justify-center text-xs font-semibold uppercase">
              {{ userInitials }}
            </div>
            <div>
              <p class="text-sm text-white">{{ userLabel }}</p>
              <p class="text-xs text-white/50">{{ userRole }}</p>
            </div>
          </div>
          <button
            type="button"
            class="h-9 w-9 rounded-full border border-white/10 bg-white/5 flex items-center justify-center focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
            aria-label="Cerrar menu"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-times text-white"></i>
          </button>
        </div>

        <nav class="flex-1 space-y-2">
          <NuxtLink
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="flex items-center justify-between px-3 py-2 rounded-xl border border-transparent text-white/70 hover:text-white hover:border-white/15 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
            @click="mobileMenuOpen = false"
          >
            <span>{{ link.label }}</span>
            <i class="pi pi-arrow-right text-xs"></i>
          </NuxtLink>
        </nav>

        <Button
          label="Ir a mi cuenta"
          icon="pi pi-user"
          class="w-full !bg-white !text-black hover:!bg-white/90 focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
          @click="handleAccountClick"
        />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const mobileMenuOpen = ref(false)
const searchQuery = ref('')
const isHydrated = ref(false)

const navLinks = [
  { label: 'Inicio', to: '/' },
  { label: 'Atracciones', to: '/atracciones' },
  { label: 'Paquetes', to: '/paquetes' },
  { label: 'Salidas', to: '/salidas' }
]

const userLabel = computed(() => {
  if (!isHydrated.value) return 'Invitado'
  if (authStore.user) return `${authStore.user.nombre} ${authStore.user.apellido_paterno}`
  return 'Invitado'
})

const userRole = computed(() => {
  if (!isHydrated.value) return 'Iniciar sesion'
  if (authStore.isAdmin) return 'Administrador'
  if (authStore.isEncargado) return 'Encargado'
  if (authStore.isTurista) return 'Turista'
  return 'Iniciar sesion'
})

const userInitials = computed(() => {
  if (!isHydrated.value) return 'A'
  if (!authStore.user) return 'A'
  const nombre = authStore.user.nombre?.[0] || ''
  const apellido = authStore.user.apellido_paterno?.[0] || ''
  return `${nombre}${apellido}`.toUpperCase() || 'A'
})

const isActive = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

const handleAccountClick = () => {
  mobileMenuOpen.value = false
  if (!authStore.isAuthenticated) {
    navigateTo('/login')
    return
  }
  if (authStore.isAdmin) {
    navigateTo('/dashboard')
    return
  }
  if (authStore.isEncargado) {
    navigateTo('/agencia/dashboard')
    return
  }
  navigateTo('/turista/dashboard')
}

onMounted(() => {
  isHydrated.value = true
})
</script>
