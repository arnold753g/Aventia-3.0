<template>
  <section class="relative min-h-[68vh] lg:min-h-[72vh] overflow-hidden">
    <div class="absolute inset-0">
      <img
        v-if="activeSlide"
        :src="resolveAssetUrl(activeSlide.image)"
        :alt="activeSlide.title"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <div class="absolute inset-0 bg-gradient-to-r from-black/80 via-black/55 to-black/30"></div>
      <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/25 to-transparent"></div>
    </div>

    <div class="relative z-10 max-w-7xl mx-auto px-4 py-16 lg:py-20">
      <div class="flex flex-col lg:flex-row gap-10 lg:gap-16 items-start lg:items-stretch lg:min-h-[520px]">
      <div class="flex-1 max-w-2xl">
          <p class="text-xs uppercase tracking-[0.4em] text-white/60">Atracciones turisticas</p>

          <Transition name="fade-slide" mode="out-in">
            <div v-if="activeSlide" :key="activeSlide.id" class="mt-4 space-y-5">
              <h1 class="text-4xl md:text-6xl lg:text-7xl font-semibold leading-tight">
                {{ activeSlide.title }}
              </h1>
              <p class="text-white/70 text-lg md:text-xl max-w-xl">
                {{ activeSlide.description }}
              </p>
              <div class="flex flex-wrap items-center gap-4 text-sm text-white/70">
                <template v-if="hasMeta">
                  <div v-for="item in metaItems" :key="item.text" class="flex items-center gap-2">
                    <i :class="item.icon"></i>
                    <span>{{ item.text }}</span>
                  </div>
                </template>
                <template v-else>
                  <div v-if="activeSlide.location" class="flex items-center gap-2">
                    <i class="pi pi-map-marker"></i>
                    <span>{{ activeSlide.location }}</span>
                  </div>
                  <div v-if="activeSlide.departureDate" class="flex items-center gap-2">
                    <i class="pi pi-calendar"></i>
                    <span>{{ activeSlide.departureDate }}</span>
                  </div>
                  <div v-if="activeSlide.seats && activeSlide.seats > 0" class="flex items-center gap-2">
                    <i class="pi pi-users"></i>
                    <span>{{ activeSlide.seats }} cupos</span>
                  </div>
                </template>
              </div>
              <div class="flex flex-wrap items-center gap-4 pt-2">
                <Button
                  :label="ctaLabel"
                  icon="pi pi-arrow-right"
                  class="rounded-full px-6 focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
                  @click="handlePrimaryAction"
                />
                <div v-if="showPrice" class="text-white/70 text-sm">
                  <span v-if="priceLabel">{{ priceLabel }} </span>
                  <span class="text-white font-semibold">{{ priceValue }}</span>
                </div>
              </div>
            </div>
            <div v-else class="mt-4 space-y-4">
              <h1 class="text-3xl md:text-5xl font-semibold leading-tight text-white">
                Sin atracciones disponibles
              </h1>
              <p class="text-white/70 text-lg max-w-xl">
                Pronto publicaremos nuevas atracciones para tu viaje.
              </p>
            </div>
          </Transition>
        </div>

        <div v-if="slides.length" class="w-full lg:w-[420px] lg:self-end">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm text-white/60">Explora mas atracciones</p>
            <span class="text-sm text-white/70">{{ slidePosition }}</span>
          </div>

          <div class="flex gap-4 overflow-x-auto pb-2 pr-2 snap-x snap-mandatory">
            <button
              v-for="(slide, index) in slides"
              :key="slide.id"
              type="button"
              class="group min-w-[220px] snap-start rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_20px_40px_rgba(0,0,0,0.45)] overflow-hidden transition-transform focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
              :class="index === activeIndex ? 'scale-[1.03] border-white/30' : 'hover:scale-[1.02]'"
              :aria-label="`Ir a ${slide.title}`"
              @click="goToSlide(index)"
            >
              <div class="relative h-28">
                <img
                  :src="resolveAssetUrl(slide.image)"
                  :alt="slide.title"
                  class="w-full h-full object-cover"
                  loading="lazy"
                />
                <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/30 to-transparent"></div>
              </div>
              <div class="p-4 text-left space-y-2">
                <p class="text-sm font-semibold text-white">{{ slide.title }}</p>
                <p class="text-xs text-white/60">{{ slide.location }}</p>
              </div>
            </button>
          </div>

          <div class="mt-6 flex items-center justify-between">
            <div class="flex items-center gap-2">
              <button
                type="button"
                class="h-10 w-10 rounded-full border border-white/15 bg-white/5 flex items-center justify-center hover:border-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
                aria-label="Slide anterior"
                @click="prevSlide"
              >
                <i class="pi pi-chevron-left text-white"></i>
              </button>
              <button
                type="button"
                class="h-10 w-10 rounded-full border border-white/15 bg-white/5 flex items-center justify-center hover:border-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 focus-visible:ring-offset-2 focus-visible:ring-offset-black/80"
                aria-label="Slide siguiente"
                @click="nextSlide"
              >
                <i class="pi pi-chevron-right text-white"></i>
              </button>
            </div>
            <div class="flex items-center gap-3">
              <div class="h-1.5 w-32 rounded-full bg-white/20 overflow-hidden">
                <div class="h-full bg-white/70 transition-all" :style="{ width: progressWidth }"></div>
              </div>
              <span class="text-xs text-white/60">{{ slidePosition }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { HeroSlide } from '~/types/dashboard'
import { useAuthStore } from '~/stores/auth'

const props = defineProps<{
  slides: HeroSlide[]
}>()

const authStore = useAuthStore()
const { resolveAssetUrl } = useAssetUrl()

const activeIndex = ref(0)

const slides = computed(() => props.slides || [])
const activeSlide = computed(() => slides.value[activeIndex.value])
const metaItems = computed(() => activeSlide.value?.meta || [])
const hasMeta = computed(() => metaItems.value.length > 0)
const ctaLabel = computed(() => activeSlide.value?.ctaLabel || 'Ver atraccion')
const priceLabel = computed(() => {
  if (activeSlide.value?.priceLabel) return activeSlide.value.priceLabel
  return activeSlide.value?.priceFrom ? 'Desde' : ''
})
const priceValue = computed(() => {
  if (activeSlide.value?.priceValue) return activeSlide.value.priceValue
  if (activeSlide.value?.priceFrom) return `Bs. ${activeSlide.value.priceFrom}`
  return ''
})
const showPrice = computed(() => Boolean(priceValue.value))

const slidePosition = computed(() => {
  if (!slides.value.length) return ''
  const current = String(activeIndex.value + 1).padStart(2, '0')
  const total = String(slides.value.length).padStart(2, '0')
  return `${current} / ${total}`
})

const progressWidth = computed(() => {
  if (!slides.value.length) return '0%'
  return `${((activeIndex.value + 1) / slides.value.length) * 100}%`
})

const goToSlide = (index: number) => {
  if (!slides.value.length) return
  const total = slides.value.length
  activeIndex.value = (index + total) % total
}

const nextSlide = () => goToSlide(activeIndex.value + 1)
const prevSlide = () => goToSlide(activeIndex.value - 1)

const handlePrimaryAction = () => {
  if (!activeSlide.value) return
  if (activeSlide.value.requiresLogin && !authStore.isAuthenticated) {
    navigateTo('/login')
    return
  }
  navigateTo(activeSlide.value.route)
}

const handleKeyDown = (event: KeyboardEvent) => {
  const target = event.target as HTMLElement | null
  if (target?.tagName === 'INPUT' || target?.tagName === 'TEXTAREA' || target?.isContentEditable) {
    return
  }
  if (event.key === 'ArrowRight') {
    nextSlide()
  }
  if (event.key === 'ArrowLeft') {
    prevSlide()
  }
}

watch(
  () => slides.value.length,
  (len) => {
    if (len && activeIndex.value >= len) activeIndex.value = 0
  }
)

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
:global(.fade-slide-enter-active),
:global(.fade-slide-leave-active) {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

:global(.fade-slide-enter-from),
:global(.fade-slide-leave-to) {
  opacity: 0;
  transform: translateY(12px);
}
</style>
