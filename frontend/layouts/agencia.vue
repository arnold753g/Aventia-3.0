<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="fixed top-0 inset-x-0 bg-white shadow-lg border-b border-gray-200 z-50">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center gap-8">
            <NuxtLink to="/agencia/dashboard" class="flex items-center gap-2">
              <i class="pi pi-building text-3xl text-green-600"></i>
              <span class="text-2xl font-bold text-gray-900">ANDARIA</span>
            </NuxtLink>

            <div class="hidden md:flex items-center gap-1">
              <NuxtLink
                to="/agencia/dashboard"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-green-50 text-green-700 font-semibold"
              >
                <i class="pi pi-home mr-2"></i>
                Dashboard
              </NuxtLink>

              <NuxtLink
                to="/agencia/mi-agencia"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-green-50 text-green-700 font-semibold"
              >
                <i class="pi pi-building mr-2"></i>
                Mi agencia
              </NuxtLink>

              <NuxtLink
                to="/agencia/paquetes"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-green-50 text-green-700 font-semibold"
              >
                <i class="pi pi-briefcase mr-2"></i>
                Paquetes
              </NuxtLink>

              <NuxtLink
                to="/agencia/ventas"
                class="px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
                active-class="bg-green-50 text-green-700 font-semibold"
              >
                <i class="pi pi-shopping-bag mr-2"></i>
                Ventas
              </NuxtLink>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <div class="hidden md:flex items-center gap-2">
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

            <div class="hidden md:block">
              <Button
                label="Cerrar sesion"
                icon="pi pi-sign-out"
                severity="danger"
                outlined
                @click="handleLogout"
              />
            </div>
            <Button
              icon="pi pi-bars"
              text
              rounded
              class="md:hidden"
              @click="mobileMenuOpen = true"
            />
          </div>
        </div>
      </div>
    </nav>
    <div
      class="fixed inset-0 bg-black/40 z-40 md:hidden transition-opacity"
      :class="mobileMenuOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'"
      @click="mobileMenuOpen = false"
    ></div>
    <div
      class="fixed top-16 right-0 w-72 max-w-[85vw] h-[calc(100vh-4rem)] bg-white shadow-xl z-50 md:hidden transform transition-transform"
      :class="mobileMenuOpen ? 'translate-x-0' : 'translate-x-full'"
    >
      <div class="p-4 h-full flex flex-col gap-4">
        <div class="flex items-start justify-between border-b border-gray-200 pb-3">
          <div class="flex items-center gap-3">
            <UserAvatar
              :nombre="authStore.user?.nombre || ''"
              :apellido="authStore.user?.apellido_paterno || ''"
              :rol="authStore.user?.rol"
              size="sm"
              showStatus
              :status="authStore.user?.status"
            />
            <div>
              <p class="text-sm font-semibold text-gray-900">
                {{ authStore.user?.nombre }} {{ authStore.user?.apellido_paterno }}
              </p>
              <p class="text-xs text-gray-500">
                {{ getRolLabel(authStore.user?.rol || '') }}
              </p>
            </div>
          </div>
          <Button icon="pi pi-times" text rounded @click="mobileMenuOpen = false" />
        </div>
        <div class="flex-1 space-y-1">
          <NuxtLink
            to="/agencia/dashboard"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors text-gray-700"
            active-class="bg-green-50 text-green-700 font-semibold"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-home"></i>
            Dashboard
          </NuxtLink>
          <NuxtLink
            to="/agencia/mi-agencia"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors text-gray-700"
            active-class="bg-green-50 text-green-700 font-semibold"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-building"></i>
            Mi agencia
          </NuxtLink>
          <NuxtLink
            to="/agencia/paquetes"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors text-gray-700"
            active-class="bg-green-50 text-green-700 font-semibold"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-briefcase"></i>
            Paquetes
          </NuxtLink>
          <NuxtLink
            to="/agencia/ventas"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors text-gray-700"
            active-class="bg-green-50 text-green-700 font-semibold"
            @click="mobileMenuOpen = false"
          >
            <i class="pi pi-shopping-bag"></i>
            Ventas
          </NuxtLink>
        </div>
        <Button
          label="Cerrar sesion"
          icon="pi pi-sign-out"
          severity="danger"
          outlined
          class="w-full"
          @click="handleLogout"
        />
      </div>
    </div>

    <main class="pt-16">
      <slot />
    </main>

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
import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { getRolLabel } from '~/utils/formatters'
import UserAvatar from '~/components/usuarios/UserAvatar.vue'

const authStore = useAuthStore()
const mobileMenuOpen = ref(false)

const handleLogout = () => {
  authStore.logout()
  navigateTo('/login')
}
</script>
