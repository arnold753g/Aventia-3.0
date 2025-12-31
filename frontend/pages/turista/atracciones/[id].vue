<template>
  <div class="page-shell">
    <div v-if="loading" class="max-w-7xl mx-auto px-4 py-10">
      <Skeleton height="360px" class="mb-6" />
      <Skeleton height="2rem" width="55%" class="mb-3" />
      <Skeleton height="1rem" width="100%" class="mb-2" />
      <Skeleton height="1rem" width="80%" />
    </div>

    <div v-else-if="error" class="max-w-7xl mx-auto px-4 py-16 text-center space-y-4">
      <i class="pi pi-exclamation-triangle text-5xl text-orange-500"></i>
      <h2 class="text-2xl font-bold text-gray-900">No se pudo cargar la atracción</h2>
      <p class="muted">{{ error }}</p>
      <div class="flex justify-center gap-2">
        <Button label="Volver" icon="pi pi-arrow-left" severity="secondary" outlined @click="navigateTo('/turista/atracciones')" />
        <Button label="Reintentar" icon="pi pi-refresh" @click="loadAtraccion" />
      </div>
    </div>

    <div v-else-if="atraccion">
      <!-- Hero -->
      <section class="relative bg-white border-b border-gray-200">
        <div class="relative h-[58vh] min-h-[360px] overflow-hidden">
          <ClientOnly>
            <Carousel
              v-if="fotosOrdenadas.length"
              :items-to-show="1"
              :wrap-around="true"
              :transition="600"
              class="h-full"
            >
              <Slide v-for="(foto, idx) in fotosOrdenadas" :key="foto.id || idx">
                <div class="relative h-[58vh] min-h-[360px]">
                  <img
                    :src="resolveFotoUrl(foto.foto)"
                    :alt="atraccion.nombre"
                    class="w-full h-full object-cover"
                    loading="lazy"
                  />
                </div>
              </Slide>
              <template #addons>
                <Navigation />
                <Pagination />
              </template>
            </Carousel>
          </ClientOnly>

          <div v-if="!fotosOrdenadas.length" class="w-full h-full bg-gradient-to-br from-blue-100 to-green-100 flex items-center justify-center">
            <i class="pi pi-image text-6xl text-gray-400"></i>
          </div>

          <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/15 to-transparent pointer-events-none" />

          <Button
            icon="pi pi-arrow-left"
            class="!absolute top-4 left-4 z-20 !bg-white/80 backdrop-blur !border-gray-200"
            rounded
            @click="navigateTo('/turista/atracciones')"
          />

          <div class="absolute bottom-0 left-0 right-0 z-10">
            <div class="max-w-7xl mx-auto px-4 pb-8">
              <div class="flex flex-wrap gap-2 mb-3">
                <Tag
                  v-if="atraccion.nivel_dificultad"
                  :value="getNivelDificultadLabel(atraccion.nivel_dificultad)"
                  :severity="nivelSeverity"
                />
                <Tag v-if="atraccion.precio_entrada === 0" value="Gratis" severity="success" icon="pi pi-check" />
                <Tag v-if="atraccion.requiere_agencia" value="Requiere agencia" severity="warning" icon="pi pi-users" />
              </div>

              <h1 ref="heroTitle" class="text-4xl md:text-5xl font-bold text-white leading-tight">
                {{ atraccion.nombre }}
              </h1>

              <div class="flex flex-wrap items-center gap-4 text-white/90 mt-3">
                <div class="flex items-center gap-2">
                  <i class="pi pi-map-marker"></i>
                  <span>
                    {{ atraccion.provincia?.nombre || 'Provincia N/D' }}
                    <span v-if="atraccion.provincia?.departamento?.nombre">
                      , {{ atraccion.provincia?.departamento?.nombre }}
                    </span>
                  </span>
                </div>
                <div v-if="atraccion.precio_entrada > 0" class="flex items-center gap-2">
                  <i class="pi pi-tag"></i>
                  <span class="font-semibold">Bs. {{ atraccion.precio_entrada.toFixed(2) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Contenido -->
      <div class="max-w-7xl mx-auto px-4 py-10">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Principal -->
          <div class="lg:col-span-2 space-y-6">
            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-info-circle" style="color: var(--color-primary);"></i>
                  <span>Descripción</span>
                </div>
              </template>
              <template #content>
                <p class="text-gray-700 whitespace-pre-line">
                  {{ atraccion.descripcion || 'No hay descripción disponible.' }}
                </p>
              </template>
            </Card>

            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-clock" style="color: var(--color-accent);"></i>
                  <span>Horarios y temporada</span>
                </div>
              </template>
              <template #content>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <p class="text-sm text-gray-500 mb-1">Horario</p>
                    <p class="text-gray-900 font-semibold">
                      {{ formatHorario(atraccion.horario_apertura, atraccion.horario_cierre) }}
                    </p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500 mb-1">Mejor época</p>
                    <p class="text-gray-900 font-semibold">
                      {{ formatMejorEpoca(atraccion.mes_inicio, atraccion.mes_fin) }}
                    </p>
                  </div>
                </div>
              </template>
            </Card>

            <Card v-if="atraccion.latitud && atraccion.longitud" class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-map" style="color: var(--color-contrast);"></i>
                  <span>Ubicación</span>
                </div>
              </template>
              <template #content>
                <div class="space-y-4">
                  <div v-if="atraccion.direccion">
                    <p class="text-sm text-gray-500 mb-1">Dirección</p>
                    <p class="text-gray-800">{{ atraccion.direccion }}</p>
                  </div>
                  <AtraccionesAtraccionMap
                    :latitud="atraccion.latitud"
                    :longitud="atraccion.longitud"
                    :editable="false"
                    :showCoordinateInputs="false"
                    height="420px"
                  />
                </div>
              </template>
            </Card>

            <Card v-if="fotosOrdenadas.length" class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-images" style="color: var(--color-primary);"></i>
                  <span>Galería</span>
                </div>
              </template>
              <template #content>
                <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
                  <button
                    v-for="(foto, idx) in fotosOrdenadas"
                    :key="foto.id || idx"
                    type="button"
                    class="relative group rounded-xl overflow-hidden border border-gray-200"
                    @click="openFoto(resolveFotoUrl(foto.foto))"
                  >
                    <img
                      :src="resolveFotoUrl(foto.foto)"
                      :alt="atraccion.nombre"
                      class="w-full h-40 object-cover group-hover:scale-105 transition-transform duration-300"
                      loading="lazy"
                    />
                    <div class="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors" />
                  </button>
                </div>
              </template>
            </Card>
          </div>

          <!-- Sidebar -->
          <div class="space-y-6">
            <Card class="surface-card sticky top-24">
              <template #title>Información rápida</template>
              <template #content>
                <div class="space-y-4">
                  <div>
                    <p class="text-xs text-gray-500 mb-1">Precio</p>
                    <p v-if="atraccion.precio_entrada > 0" class="text-2xl font-bold text-blue-600">
                      Bs. {{ atraccion.precio_entrada.toFixed(2) }}
                    </p>
                    <Tag v-else value="Entrada gratuita" severity="success" icon="pi pi-check" />
                  </div>

                  <Divider />

                  <div v-if="atraccion.nivel_dificultad">
                    <p class="text-xs text-gray-500 mb-2">Dificultad</p>
                    <Tag
                      :value="getNivelDificultadLabel(atraccion.nivel_dificultad)"
                      :severity="nivelSeverity"
                      class="w-full justify-center"
                    />
                  </div>

                  <Divider />

                  <div>
                    <p class="text-xs text-gray-500 mb-2">Requisitos</p>
                    <div class="space-y-2 text-sm text-gray-700">
                      <div class="flex items-start gap-2">
                        <i :class="['pi', atraccion.requiere_agencia ? 'pi-check-circle text-green-600' : 'pi-times-circle text-gray-400']"></i>
                        <span>Requiere agencia de turismo</span>
                      </div>
                      <div class="flex items-start gap-2">
                        <i :class="['pi', atraccion.acceso_particular ? 'pi-check-circle text-green-600' : 'pi-times-circle text-gray-400']"></i>
                        <span>Acceso particular permitido</span>
                      </div>
                    </div>
                  </div>

                  <Divider v-if="hasContacto" />

                  <div v-if="hasContacto" class="space-y-3">
                    <p class="text-xs text-gray-500">Contacto</p>
                    <Button
                      v-if="atraccion.telefono"
                      :label="atraccion.telefono"
                      icon="pi pi-phone"
                      class="w-full"
                      @click="llamar"
                    />
                    <Button
                      v-if="atraccion.email"
                      label="Enviar correo"
                      icon="pi pi-envelope"
                      severity="secondary"
                      outlined
                      class="w-full"
                      @click="enviarCorreo"
                    />
                  </div>
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>
    </div>

    <Dialog v-model:visible="fotoDialog" header="Foto" :modal="true" :style="{ width: 'min(920px, 95vw)' }">
      <img v-if="fotoSeleccionada" :src="fotoSeleccionada" alt="Foto" class="w-full max-h-[80vh] object-contain rounded-lg" />
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { Carousel, Slide, Pagination, Navigation } from 'vue3-carousel'
import 'vue3-carousel/dist/carousel.css'
import { getNivelDificultadLabel, formatHorario, formatMejorEpoca } from '~/utils/formatters-atraccion'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const route = useRoute()
const toast = useToast()
const { getAtraccion } = useAtracciones()
const config = useRuntimeConfig()
const assetsBase = config.public.apiBase.replace(/\/api\/v1\/?$/, '')

const loading = ref(true)
const error = ref<string | null>(null)
const atraccion = ref<any>(null)

const fotoDialog = ref(false)
const fotoSeleccionada = ref<string | null>(null)
const heroTitle = ref<HTMLElement | null>(null)

const resolveFotoUrl = (path?: string) => {
  if (!path) return ''
  let normalized = path.replace(/\\\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const fotosOrdenadas = computed(() => {
  const fotos = (atraccion.value?.fotos || []).slice()
  fotos.sort((a: any, b: any) => {
    if (!!a.es_principal !== !!b.es_principal) return a.es_principal ? -1 : 1
    return (a.orden || 0) - (b.orden || 0)
  })
  return fotos
})

const hasContacto = computed(() => !!(atraccion.value?.telefono || atraccion.value?.email))

const nivelSeverity = computed(() => {
  const nivel = atraccion.value?.nivel_dificultad
  if (!nivel) return 'secondary'
  const map: Record<string, string> = {
    facil: 'success',
    medio: 'info',
    dificil: 'warning',
    extremo: 'danger'
  }
  return map[nivel] || 'secondary'
})

const openFoto = (url: string) => {
  fotoSeleccionada.value = url
  fotoDialog.value = true
}

const llamar = () => {
  if (!atraccion.value?.telefono) return
  window.location.href = `tel:${atraccion.value.telefono}`
}

const enviarCorreo = () => {
  if (!atraccion.value?.email) return
  window.location.href = `mailto:${atraccion.value.email}`
}

const loadAtraccion = async () => {
  loading.value = true
  error.value = null
  atraccion.value = null

  try {
    const id = Number(route.params.id)
    const response: any = await getAtraccion(id)
    if (response.success) {
      atraccion.value = response.data

      try {
        const { gsap } = await import('gsap')
        if (heroTitle.value) {
          gsap.from(heroTitle.value, { opacity: 0, y: 10, duration: 0.5, ease: 'power2.out' })
        }
      } catch {
        // no-op
      }
      return
    }

    error.value = response?.error?.message || 'No se pudo cargar la atracción'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudo cargar la atracción'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAtraccion()
})
</script>
