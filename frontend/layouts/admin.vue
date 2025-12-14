<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navbar -->
    <nav class="bg-white shadow-lg border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <!-- Logo -->
          <div class="flex items-center gap-8">
            <NuxtLink to="/dashboard" class="flex items-center gap-2">
              <i class="pi pi-compass text-3xl text-blue-600"></i>
              <span class="text-2xl font-bold text-gray-900">ANDARIA</span>
            </NuxtLink>

            <!-- Navigation Links -->
            <div class="hidden md:flex items-center gap-1">
              <NuxtLink
                to="/dashboard"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-blue-50 text-blue-600 font-semibold"
              >
                <i class="pi pi-home mr-2"></i>
                Dashboard
              </NuxtLink>

              <NuxtLink
                to="/admin/usuarios"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-blue-50 text-blue-600 font-semibold"
              >
                <i class="pi pi-users mr-2"></i>
                Usuarios
              </NuxtLink>

              <NuxtLink
                to="/admin/atracciones"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-blue-50 text-blue-600 font-semibold"
              >
                <i class="pi pi-map-marker mr-2"></i>
                Atracciones
              </NuxtLink>

              <NuxtLink
                to="/admin/agencias"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-blue-50 text-blue-600 font-semibold"
              >
                <i class="pi pi-building mr-2"></i>
                Agencias
              </NuxtLink>
            </div>
          </div>

          <!-- User Menu -->
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <UserAvatar
                :nombre="authStore.user?.nombre || ''"
                :apellido="authStore.user?.apellido_paterno || ''"
                :rol="authStore.user?.rol"
                size="sm"
                showStatus
                :status="authStore.user?.status"
              />
              <div class="hidden md:block">
                <p class="text-sm font-semibold text-gray-900">
                  {{ authStore.user?.nombre }} {{ authStore.user?.apellido_paterno }}
                </p>
                <p class="text-xs text-gray-500">
                  {{ getRolLabel(authStore.user?.rol || '') }}
                </p>
              </div>
            </div>

            <Button
              icon="pi pi-sign-out"
              severity="danger"
              text
              rounded
              @click="handleLogout"
              v-tooltip.bottom="'Cerrar Sesión'"
            />
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main>
      <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 mt-12">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="text-center text-sm text-gray-600">
          <p>&copy; 2024 ANDARIA - Sistema de Gestión Turística de Bolivia</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { getRolLabel } from '~/utils/formatters'
import UserAvatar from '~/components/usuarios/UserAvatar.vue'

const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  navigateTo('/login')
}
</script>
