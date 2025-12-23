<template>
  <div class="space-y-4">
    <HeroSlider :slides="heroSlides" />
    <QuickModules />
    <PopularRow />
    <RecentGrid />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import HeroSlider from '~/components/dashboard/HeroSlider.vue'
import QuickModules from '~/components/dashboard/QuickModules.vue'
import PopularRow from '~/components/dashboard/PopularRow.vue'
import RecentGrid from '~/components/dashboard/RecentGrid.vue'
import { useDashboardStore } from '~/stores/dashboard'

definePageMeta({
  layout: 'home'
})

const dashboardStore = useDashboardStore()
const { heroSlides } = storeToRefs(dashboardStore)

onMounted(() => {
  dashboardStore.loadAttractions()
  dashboardStore.loadPackages()
})
</script>
