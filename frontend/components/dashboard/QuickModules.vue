<template>
  <section class="max-w-7xl mx-auto px-4 py-14">
    <div class="flex items-center justify-between mb-8">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-white/50">Accesos rapidos</p>
        <h2 class="text-2xl md:text-3xl font-semibold mt-2">Explora lo mejor de Andaria</h2>
      </div>
      <span class="hidden md:inline text-sm text-white/50">Curado para viajeros</span>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <button
        v-for="module in modules"
        :key="module.id"
        type="button"
        class="group relative overflow-hidden rounded-[26px] border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_24px_60px_rgba(0,0,0,0.45)] text-left transition-transform hover:-translate-y-1.5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
        :aria-label="`Ir a ${module.title}`"
        @click="handleModuleClick(module)"
      >
        <div class="absolute inset-0">
          <img
            :src="resolveAssetUrl(module.image)"
            :alt="module.title"
            class="w-full h-full object-cover"
            loading="lazy"
          />
          <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/40 to-transparent"></div>
        </div>
        <div class="relative z-10 p-6 flex flex-col h-full justify-between">
          <div class="h-12 w-12 rounded-2xl border border-white/20 bg-white/10 flex items-center justify-center">
            <i :class="module.icon" class="text-lg text-white"></i>
          </div>
          <div>
            <h3 class="text-xl font-semibold text-white mt-6">{{ module.title }}</h3>
            <p class="text-sm text-white/70 mt-2">{{ module.description }}</p>
          </div>
          <div class="mt-6 inline-flex items-center gap-2 text-sm text-white/80">
            <span>Explorar</span>
            <i class="pi pi-arrow-right text-xs"></i>
          </div>
        </div>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useDashboardStore } from '~/stores/dashboard'
import type { DashboardModule } from '~/types/dashboard'

const authStore = useAuthStore()
const dashboardStore = useDashboardStore()
const { resolveAssetUrl } = useAssetUrl()

const modules = computed(() => dashboardStore.modules)

const handleModuleClick = (module: DashboardModule) => {
  if (module.requiresLogin && !authStore.isAuthenticated) {
    navigateTo('/login')
    return
  }
  navigateTo(module.route)
}
</script>
