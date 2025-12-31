<template>
  <section class="max-w-7xl mx-auto px-4 py-12">
    <div class="flex items-center justify-between mb-8">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-white/50">Recientes</p>
        <h2 class="text-2xl md:text-3xl font-semibold mt-2">Novedades para ti</h2>
      </div>
      <span class="hidden md:inline text-sm text-white/50">Actualizado hoy</span>
    </div>

    <div v-if="isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="n in 3"
        :key="n"
        class="h-64 rounded-[22px] border border-white/10 bg-white/5 animate-pulse"
      ></div>
    </div>

    <div v-else-if="showError" class="text-sm text-white/60">
      {{ errorMessage }}
    </div>

    <div v-else-if="showEmpty" class="text-sm text-white/60">
      No hay novedades disponibles por ahora.
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <button
          v-for="item in items"
          :key="item.id"
          type="button"
          class="group overflow-hidden rounded-[22px] border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_20px_40px_rgba(0,0,0,0.45)] text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
          :aria-label="`Ver ${item.title}`"
          @click="handleItemClick(item)"
      >
        <div class="relative h-48 overflow-hidden">
          <img
              v-if="item.image"
              :src="resolveAssetUrl(item.image)"
              :alt="item.title"
              class="absolute inset-0 w-full h-full object-cover transition-transform duration-300 ease-out transform-gpu will-change-transform group-hover:scale-105"
              loading="lazy"
          />

          <div
              class="absolute -inset-px bg-gradient-to-t from-black/70 via-black/30 to-transparent pointer-events-none"
          ></div>

          <span class="absolute top-3 left-3 text-xs px-2 py-1 rounded-full border border-white/20 bg-white/10 text-white/80">
        {{ badgeLabel(item.type) }}
      </span>
        </div>

        <div class="p-4 space-y-3">
          <p class="text-xs text-white/60">{{ item.date }}</p>
          <h3 class="text-base font-semibold text-white line-clamp-2">{{ item.title }}</h3>
        </div>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useDashboardStore } from '~/stores/dashboard'
import type { RecentItem } from '~/types/dashboard'

const authStore = useAuthStore()
const dashboardStore = useDashboardStore()
const { resolveAssetUrl } = useAssetUrl()

const items = computed(() => dashboardStore.recentItems)
const isLoading = computed(() => dashboardStore.loadingAttractions && items.value.length === 0)
const showError = computed(() => !dashboardStore.loadingAttractions && !!dashboardStore.attractionsError)
const errorMessage = computed(() => dashboardStore.attractionsError || 'No se pudieron cargar las novedades')
const showEmpty = computed(() => !dashboardStore.loadingAttractions && !dashboardStore.attractionsError && items.value.length === 0)

const badgeLabel = (type: RecentItem['type']) => {
  if (type === 'paquete') return 'Paquete'
  if (type === 'salida') return 'Salida'
  return 'Atraccion'
}

const handleItemClick = (item: RecentItem) => {
  if (item.requiresLogin && !authStore.isAuthenticated) {
    navigateTo('/login')
    return
  }
  navigateTo(item.route)
}
</script>
