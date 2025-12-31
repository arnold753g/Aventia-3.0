<template>
  <div
    class="home-shell min-h-screen relative overflow-x-hidden"
    :class="isDark ? 'dark bg-[#050505] text-white' : 'bg-white text-gray-900'"
  >
    <div v-if="isDark" class="pointer-events-none absolute inset-0">
      <div class="absolute -top-32 -right-32 h-96 w-96 bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.12),transparent_60%)]"></div>
      <div class="absolute top-1/3 -left-32 h-[28rem] w-[28rem] bg-[radial-gradient(circle_at_center,rgba(59,130,246,0.12),transparent_60%)]"></div>
      <div class="absolute bottom-0 right-0 h-80 w-80 bg-[radial-gradient(circle_at_center,rgba(16,185,129,0.08),transparent_60%)]"></div>
    </div>

    <TopNav />

    <main class="relative z-10 pt-16">
      <slot />
    </main>

    <footer class="relative z-10 border-t border-white/10 mt-16 bg-[#050505] text-white">
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
import { computed } from 'vue'
import TopNav from '~/components/dashboard/TopNav.vue'

const route = useRoute()
const isDark = computed(() => !route.path.startsWith('/atracciones'))
</script>
