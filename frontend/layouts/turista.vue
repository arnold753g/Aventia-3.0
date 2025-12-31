<template>
  <div class="panel-shell panel-accent-blue">
    <!-- Listener de notificaciones WebSocket -->
    <NotificationListener />

    <nav class="fixed top-0 inset-x-0 z-50 panel-nav">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center gap-8">
            <NuxtLink to="/turista/dashboard" class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-2xl border border-white/15 bg-white/10 flex items-center justify-center">
                <i class="pi pi-compass text-lg text-white"></i>
              </div>
              <div class="leading-tight">
                <p class="text-sm uppercase tracking-[0.3em] text-white/50">Turista</p>
                <p class="text-lg font-semibold text-white">ANDARIA</p>
              </div>
            </NuxtLink>

            <div class="hidden md:flex items-center gap-1 text-sm">
              <NuxtLink
                to="/turista/dashboard"
                class="panel-nav-link px-4 py-2 rounded-lg transition-colors"
                active-class="panel-nav-link-active"
              >

                Dashboard
              </NuxtLink>

              <NuxtLink
                to="/turista/atracciones"
                class="panel-nav-link px-4 py-2 rounded-lg transition-colors"
                active-class="panel-nav-link-active"
              >

                Atracciones
              </NuxtLink>

              <NuxtLink
                to="/turista/paquetes"
                class="panel-nav-link px-4 py-2 rounded-lg transition-colors"
                active-class="panel-nav-link-active"
              >

                Paquetes
              </NuxtLink>

              <NuxtLink
                to="/turista/agencias"
                class="panel-nav-link px-4 py-2 rounded-lg transition-colors"
                active-class="panel-nav-link-active"
              >

                Agencias
              </NuxtLink>

              <NuxtLink
                to="/turista/salidas"
                class="panel-nav-link px-4 py-2 rounded-lg transition-colors"
                active-class="panel-nav-link-active"
              >

                Salidas
              </NuxtLink>

              <NuxtLink
                to="/turista/mis-compras"
                class="panel-nav-link px-4 py-2 rounded-lg transition-colors"
                active-class="panel-nav-link-active"
              >
                
                Mis compras
              </NuxtLink>

            </div>
          </div>

          <div class="flex items-center gap-4">
            <!-- Notificaciones -->
            <div class="relative dark">
              <NotificationBadge @toggle="notificationPanelOpen = !notificationPanelOpen" />
              <NotificationPanel
                :isOpen="notificationPanelOpen"
                @close="notificationPanelOpen = false"
              />
            </div>

            <NuxtLink
              to="/turista/perfil"
              class="hidden md:flex items-center gap-3 px-3 py-2 rounded-full border border-white/10 bg-white/5 hover:border-white/25 transition-colors"
              title="Mi perfil"
            >
              <UserAvatar
                :nombre="authStore.user?.nombre || ''"
                :apellido="authStore.user?.apellido_paterno || ''"
                :rol="authStore.user?.rol"
                size="sm"
                showStatus
                :status="authStore.user?.status"
              />
              <div class="hidden md:block leading-tight">
                <p class="text-sm font-semibold text-white">
                  {{ authStore.user?.nombre }} {{ authStore.user?.apellido_paterno }}
                </p>
                <p class="text-xs text-white/50">
                  {{ getRolLabel(authStore.user?.rol || '') }}
                </p>
              </div>
            </NuxtLink>

            <div class="hidden md:flex">
              <button
                type="button"
                class="flex items-center gap-2 px-3 py-2 rounded-full border border-white/10 bg-white/5 text-white/70 hover:text-white hover:border-white/25 transition-colors"
                @click="handleLogout"
              >
                <i class="pi pi-sign-out"></i>
                <span class="text-sm">Cerrar sesion</span>
              </button>
            </div>
            <button
              type="button"
              class="md:hidden h-10 w-10 rounded-full border border-white/15 bg-white/5 flex items-center justify-center text-white/80 hover:text-white hover:border-white/25 transition-colors"
              @click="mobileMenuOpen = true"
            >
              <i class="pi pi-bars"></i>
            </button>
          </div>
        </div>
      </div>
    </nav>
    <div
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 md:hidden transition-opacity"
      :class="mobileMenuOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'"
      @click="mobileMenuOpen = false"
    ></div>
    <div
      class="fixed top-16 right-0 w-72 max-w-[85vw] h-[calc(100vh-4rem)] panel-menu shadow-2xl z-50 md:hidden transform transition-transform"
      :class="mobileMenuOpen ? 'translate-x-0' : 'translate-x-full'"
    >
      <div class="p-4 h-full flex flex-col gap-4">
        <div class="flex items-start justify-between border-b border-white/10 pb-3">
          <NuxtLink
            to="/turista/perfil"
            class="flex items-center gap-3 rounded-xl border border-transparent px-2 py-1 hover:border-white/20 transition-colors"
            title="Mi perfil"
            @click="mobileMenuOpen = false"
          >
            <UserAvatar
              :nombre="authStore.user?.nombre || ''"
              :apellido="authStore.user?.apellido_paterno || ''"
              :rol="authStore.user?.rol"
              size="sm"
              showStatus
              :status="authStore.user?.status"
            />
            <div>
              <p class="text-sm font-semibold text-white">
                {{ authStore.user?.nombre }} {{ authStore.user?.apellido_paterno }}
              </p>
              <p class="text-xs text-white/50">
                {{ getRolLabel(authStore.user?.rol || '') }}
              </p>
            </div>
          </NuxtLink>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="h-9 w-9 rounded-full border border-white/10 bg-white/5 flex items-center justify-center text-white/70 hover:text-white hover:border-white/25 transition-colors"
              @click="mobileMenuOpen = false"
            >
              <i class="pi pi-times"></i>
            </button>
          </div>
        </div>
        <div class="flex-1 space-y-1">
          <NuxtLink
            to="/turista/dashboard"
            class="panel-nav-link flex items-center gap-3 px-3 py-2 rounded-xl border border-transparent transition-colors hover:border-white/15"
            active-class="panel-nav-link-active"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-home"></i>
            Dashboard
          </NuxtLink>
          <NuxtLink
            to="/turista/atracciones"
            class="panel-nav-link flex items-center gap-3 px-3 py-2 rounded-xl border border-transparent transition-colors hover:border-white/15"
            active-class="panel-nav-link-active"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-map-marker"></i>
            Atracciones
          </NuxtLink>
          <NuxtLink
            to="/turista/paquetes"
            class="panel-nav-link flex items-center gap-3 px-3 py-2 rounded-xl border border-transparent transition-colors hover:border-white/15"
            active-class="panel-nav-link-active"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-briefcase"></i>
            Paquetes
          </NuxtLink>
          <NuxtLink
            to="/turista/agencias"
            class="panel-nav-link flex items-center gap-3 px-3 py-2 rounded-xl border border-transparent transition-colors hover:border-white/15"
            active-class="panel-nav-link-active"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-building"></i>
            Agencias
          </NuxtLink>
          <NuxtLink
            to="/turista/salidas"
            class="panel-nav-link flex items-center gap-3 px-3 py-2 rounded-xl border border-transparent transition-colors hover:border-white/15"
            active-class="panel-nav-link-active"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-calendar"></i>
            Salidas
          </NuxtLink>
          <NuxtLink
            to="/turista/mis-compras"
            class="panel-nav-link flex items-center gap-3 px-3 py-2 rounded-xl border border-transparent transition-colors hover:border-white/15"
            active-class="panel-nav-link-active"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-shopping-cart"></i>
            Mis compras
          </NuxtLink>
        </div>
        <button
          type="button"
          class="w-full px-4 py-3 rounded-full bg-white text-black font-semibold hover:bg-white/90 transition-colors"
          @click="handleLogout"
        >
          Cerrar sesion
        </button>
      </div>
    </div>

    <main class="pt-16">
      <slot />
    </main>

    <footer class="panel-footer relative z-10 mt-16 text-white">
      <div class="max-w-7xl mx-auto px-4 py-12 grid grid-cols-1 md:grid-cols-4 gap-8 text-sm text-white/70">
        <div>
          <p class="text-white font-semibold text-lg">ANDARIA</p>
          <p class="mt-3 text-white/60">
            Plataforma de experiencias turisticas en Bolivia.
          </p>
        </div>
        <div>
          <p class="text-white font-semibold mb-3">Explora</p>
          <div class="space-y-2">
            <NuxtLink to="/atracciones" class="block hover:text-white">Atracciones</NuxtLink>
            <NuxtLink to="/paquetes" class="block hover:text-white">Paquetes</NuxtLink>
            <NuxtLink to="/salidas" class="block hover:text-white">Salidas</NuxtLink>
          </div>
        </div>
        <div>
          <p class="text-white font-semibold mb-3">Recursos</p>
          <div class="space-y-2">
            <a href="#" class="block hover:text-white">Soporte</a>
            <a href="#" class="block hover:text-white">Preguntas frecuentes</a>
            <a href="#" class="block hover:text-white">Politicas</a>
          </div>
        </div>
        <div>
          <p class="text-white font-semibold mb-3">Follow us</p>
          <div class="flex items-center gap-3 text-white/70">
            <button
              type="button"
              class="h-10 w-10 rounded-full border border-white/15 bg-white/5 hover:border-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
              aria-label="Facebook"
            >
              <i class="pi pi-facebook"></i>
            </button>
            <button
              type="button"
              class="h-10 w-10 rounded-full border border-white/15 bg-white/5 hover:border-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
              aria-label="Instagram"
            >
              <i class="pi pi-instagram"></i>
            </button>
            <button
              type="button"
              class="h-10 w-10 rounded-full border border-white/15 bg-white/5 hover:border-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
              aria-label="Twitter"
            >
              <i class="pi pi-twitter"></i>
            </button>
          </div>
        </div>
      </div>
      <div class="border-t border-white/10 py-4 text-center text-xs text-white/50">
        (c) 2025 Andaria. Todos los derechos reservados.
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { getRolLabel } from '~/utils/formatters'
import UserAvatar from '~/components/usuarios/UserAvatar.vue'

const authStore = useAuthStore()
const mobileMenuOpen = ref(false)
const notificationPanelOpen = ref(false)

const handleLogout = () => {
  authStore.logout()
  navigateTo('/login')
}
</script>
