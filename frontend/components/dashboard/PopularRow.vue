<template>
  <section class="max-w-7xl mx-auto px-4 py-12">
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-8">
      <div>
        <p class="text-xs uppercase tracking-[0.3em] text-white/50">Popular</p>
        <h2 class="text-2xl md:text-3xl font-semibold mt-2">Selecciones destacadas</h2>
      </div>
      <div class="inline-flex rounded-full border border-white/10 bg-white/5 p-1">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          type="button"
          class="px-4 py-2 rounded-full text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
          :class="activeTab === tab.value ? 'bg-white/15 text-white' : 'text-white/60 hover:text-white'"
          @click="activeTab = tab.value"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="n in 4"
        :key="n"
        class="h-56 rounded-[22px] border border-white/10 bg-white/5 animate-pulse"
      ></div>
    </div>

    <div v-else-if="showError" class="text-sm text-white/60">
      {{ errorMessage }}
    </div>

    <div v-else-if="showEmpty" class="text-sm text-white/60">
      {{ emptyMessage }}
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <button
        v-for="item in displayItems"
        :key="item.id"
        type="button"
        class="group relative overflow-hidden rounded-[22px] border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_20px_40px_rgba(0,0,0,0.45)] text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
        :aria-label="`Ver ${item.title}`"
        @click="handleItemClick(item)"
      >
        <div class="relative h-44">
          <img
            v-if="item.image"
            :src="resolveAssetUrl(item.image)"
            :alt="item.title"
            class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
          <div v-else class="w-full h-full bg-gradient-to-br from-white/10 to-white/5 flex items-center justify-center">
            <i class="pi pi-image text-3xl text-white/40"></i>
          </div>
          <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent"></div>
        </div>
        <div class="p-4 flex flex-col gap-3">
          <h3 class="text-base font-semibold text-white line-clamp-2">{{ item.title }}</h3>
          <div class="flex items-center justify-between text-xs text-white/60">
            <span>{{ item.meta }}</span>
            <span class="h-9 w-9 rounded-full border border-white/15 bg-white/10 flex items-center justify-center">
              <i class="pi pi-arrow-right text-xs"></i>
            </span>
          </div>
        </div>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useDashboardStore } from '~/stores/dashboard'

const authStore = useAuthStore()
const dashboardStore = useDashboardStore()
const { resolveAssetUrl } = useAssetUrl()

const tabs = [
  { label: 'Atracciones populares', value: 'atracciones' },
  { label: 'Paquetes recomendados', value: 'paquetes' }
]

const activeTab = ref<'atracciones' | 'paquetes'>('atracciones')

const displayItems = computed(() => {
  if (activeTab.value === 'atracciones') {
    return dashboardStore.popularAttractions.map((item) => ({
      id: item.id,
      title: item.name,
      meta: item.place,
      image: item.image,
      route: item.route,
      requiresLogin: false
    }))
  }
  return dashboardStore.popularPackages.map((item) => ({
    id: item.id,
    title: item.name,
    meta: `${item.duration} - Bs. ${item.priceFrom}`,
    image: item.image,
    route: item.route,
    requiresLogin: item.requiresLogin ?? true
  }))
})

const isLoading = computed(() =>
  activeTab.value === 'atracciones' ? dashboardStore.loadingAttractions : dashboardStore.loadingPackages
)
const showError = computed(() => {
  if (activeTab.value === 'atracciones') {
    return !dashboardStore.loadingAttractions && !!dashboardStore.attractionsError
  }
  return !dashboardStore.loadingPackages && !!dashboardStore.packagesError
})
const errorMessage = computed(() => {
  if (activeTab.value === 'atracciones') {
    return dashboardStore.attractionsError || 'No se pudieron cargar las atracciones'
  }
  return dashboardStore.packagesError || 'No se pudieron cargar los paquetes'
})
const showEmpty = computed(() => {
  if (activeTab.value === 'atracciones') {
    return !dashboardStore.loadingAttractions &&
      !dashboardStore.attractionsError &&
      displayItems.value.length === 0
  }
  return !dashboardStore.loadingPackages &&
    !dashboardStore.packagesError &&
    displayItems.value.length === 0
})
const emptyMessage = computed(() =>
  activeTab.value === 'atracciones'
    ? 'No hay atracciones disponibles por ahora.'
    : 'No hay paquetes disponibles por ahora.'
)

const handleItemClick = (item: {
  route: string
  requiresLogin?: boolean
}) => {
  if (item.requiresLogin && !authStore.isAuthenticated) {
    navigateTo('/login')
    return
  }
  navigateTo(item.route)
}
</script>
