<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold" style="color: var(--color-primary);">
            Panel del turista
          </h1>
          <p class="muted mt-1">Explora atracciones turísticas y gestiona tu perfil.</p>
        </div>
        <div class="flex gap-2">
          <Button
            label="Explorar atracciones"
            icon="pi pi-map-marker"
            @click="navigateTo('/turista/atracciones')"
          />
          <Button
            label="Mi perfil"
            icon="pi pi-user"
            severity="secondary"
            outlined
            @click="navigateTo('/turista/perfil')"
          />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <Card class="surface-card">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-star" style="color: var(--color-accent);"></i>
            <span>Atracciones Turisticas Recomendadas</span>
          </div>
        </template>
        <template #content>
          <div v-if="loading">
            <Skeleton height="240px" />
          </div>

          <div v-else-if="error" class="space-y-3">
            <Message severity="warn" :closable="false">{{ error }}</Message>
            <Button label="Reintentar" icon="pi pi-refresh" outlined @click="loadRecomendadas" />
          </div>

          <div v-else-if="recomendadas.length === 0" class="text-center py-10 muted">
            No hay atracciones disponibles en este momento.
          </div>

          <ClientOnly>
            <Carousel
              v-if="recomendadas.length"
              :items-to-show="1.15"
              :wrap-around="true"
              :transition="550"
              class="py-2"
            >
              <Slide v-for="atraccion in recomendadas" :key="atraccion.id">
                <div class="px-2 w-full">
                  <div
                    class="relative h-60 rounded-2xl overflow-hidden border border-gray-200 cursor-pointer"
                    @click="navigateTo(`/turista/atracciones/${atraccion.id}`)"
                  >
                    <img
                      v-if="resolvePrincipalFoto(atraccion)"
                      :src="resolvePrincipalFoto(atraccion)"
                      :alt="atraccion.nombre"
                      class="w-full h-full object-cover"
                      loading="lazy"
                    />
                    <div v-else class="w-full h-full bg-gradient-to-br from-blue-100 to-green-100 flex items-center justify-center">
                      <i class="pi pi-image text-5xl text-gray-400"></i>
                    </div>

                    <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/15 to-transparent" />

                    <div class="absolute bottom-0 left-0 right-0 p-4">
                      <div class="flex items-center gap-2 text-xs text-white/90">
                        <i class="pi pi-map-marker" />
                        <span>
                          {{ atraccion.provincia?.nombre || 'Provincia N/D' }}
                          <span v-if="atraccion.provincia?.departamento?.nombre">
                            , {{ atraccion.provincia?.departamento?.nombre }}
                          </span>
                        </span>
                      </div>
                      <h3 class="text-lg font-semibold text-white leading-tight mt-1 line-clamp-1">
                        {{ atraccion.nombre }}
                      </h3>
                      <p class="text-xs text-white/80 mt-1 line-clamp-2">
                        {{ atraccion.descripcion || 'Descubre esta atracción turística.' }}
                      </p>
                    </div>
                  </div>
                </div>
              </Slide>

              <template #addons>
                <Navigation />
                <Pagination />
              </template>
            </Carousel>
          </ClientOnly>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { Carousel, Slide, Pagination, Navigation } from 'vue3-carousel'
import 'vue3-carousel/dist/carousel.css'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const { getAtracciones } = useAtracciones()
const config = useRuntimeConfig()
const assetsBase = config.public.apiBase.replace(/\/api\/v1\/?$/, '')

const recomendadas = ref<any[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

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

const resolvePrincipalFoto = (atraccion: any) => {
  const principal = atraccion?.fotos?.find((f: any) => f.es_principal)
  const path = principal?.foto || atraccion?.fotos?.[0]?.foto
  return resolveFotoUrl(path)
}

const loadRecomendadas = async () => {
  loading.value = true
  error.value = null
  try {
    const response: any = await getAtracciones({
      page: 1,
      limit: 10,
      visible_publico: 'true',
      status: 'activa',
      sort_by: 'created_at',
      sort_order: 'desc'
    })

    if (response.success) {
      recomendadas.value = response.data?.atracciones || []
      return
    }

    error.value = response?.error?.message || 'No se pudieron cargar las atracciones'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar las atracciones'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRecomendadas()
})
</script>
