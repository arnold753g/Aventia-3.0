<template>
  <div class="page-shell">
    <div v-if="loading" class="max-w-7xl mx-auto px-4 py-10">
      <Skeleton height="340px" class="mb-6" />
      <Skeleton height="2rem" width="55%" class="mb-3" />
      <Skeleton height="1rem" width="100%" class="mb-2" />
      <Skeleton height="1rem" width="80%" />
    </div>

    <div v-else-if="error" class="max-w-7xl mx-auto px-4 py-16 text-center space-y-4">
      <i class="pi pi-exclamation-triangle text-5xl text-orange-500"></i>
      <h2 class="text-2xl font-bold text-gray-900">No se pudo cargar la agencia</h2>
      <p class="muted">{{ error }}</p>
      <div class="flex justify-center gap-2">
        <Button label="Volver" icon="pi pi-arrow-left" severity="secondary" outlined @click="navigateTo('/turista/agencias')" />
        <Button label="Reintentar" icon="pi pi-refresh" @click="loadAgencia" />
      </div>
    </div>

    <div v-else-if="agencia">
      <section class="relative bg-white border-b border-gray-200">
        <div class="relative h-[52vh] min-h-[320px] overflow-hidden">
          <img
            v-if="heroFoto"
            :src="heroFoto"
            :alt="agencia.nombre_comercial"
            class="w-full h-full object-cover"
            loading="lazy"
          />
          <div
            v-else
            class="w-full h-full bg-gradient-to-br from-sky-100 to-emerald-100 flex items-center justify-center"
          >
            <i class="pi pi-image text-6xl text-gray-400"></i>
          </div>

          <div class="absolute inset-0 bg-gradient-to-t from-black/75 via-black/25 to-transparent pointer-events-none" />

          <Button
            icon="pi pi-arrow-left"
            class="!absolute top-4 left-4 z-20 !bg-white/80 backdrop-blur !border-gray-200"
            rounded
            @click="navigateTo('/turista/agencias')"
          />

          <div class="absolute bottom-0 left-0 right-0 z-10">
            <div class="max-w-7xl mx-auto px-4 pb-8">
              <div class="flex flex-wrap gap-2 mb-3">
                <Tag v-if="agencia.licencia_turistica" value="Licencia verificada" severity="success" icon="pi pi-check-circle" />
                <Tag v-if="agencia.departamento?.nombre" :value="agencia.departamento.nombre" severity="info" />
              </div>

              <h1 class="text-4xl md:text-5xl font-bold text-white leading-tight">
                {{ agencia.nombre_comercial }}
              </h1>

              <div class="flex flex-wrap items-center gap-4 text-white/90 mt-3">
                <div class="flex items-center gap-2">
                  <i class="pi pi-map-marker"></i>
                  <span class="font-medium">
                    {{ agencia.departamento?.nombre || 'Ubicacion por definir' }}
                  </span>
                </div>
                <div v-if="agencia.direccion" class="flex items-center gap-2">
                  <i class="pi pi-compass"></i>
                  <span>{{ agencia.direccion }}</span>
                </div>
                <Button
                  v-if="paquetes.length"
                  label="Ver paquetes"
                  icon="pi pi-briefcase"
                  text
                  class="!text-white/90"
                  @click="scrollToPaquetes"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      <div class="max-w-7xl mx-auto px-4 py-10">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div class="lg:col-span-2 space-y-6">
            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-info-circle text-emerald-600"></i>
                  <span>Descripcion</span>
                </div>
              </template>
              <template #content>
                <p class="text-gray-700 whitespace-pre-line">
                  {{ agencia.descripcion || 'No hay descripcion disponible.' }}
                </p>
              </template>
            </Card>

            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-star text-blue-600"></i>
                  <span>Especialidades</span>
                </div>
              </template>
              <template #content>
                <div v-if="especialidades.length" class="flex flex-wrap gap-2">
                  <Tag v-for="(item, idx) in especialidades" :key="idx" :value="item" severity="secondary" />
                </div>
                <p v-else class="text-sm text-gray-500">Sin especialidades registradas.</p>
              </template>
            </Card>

            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-clock text-emerald-600"></i>
                  <span>Horarios y dias</span>
                </div>
              </template>
              <template #content>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <p class="text-sm text-gray-500 mb-1">Horario</p>
                    <p class="text-gray-900 font-semibold">
                      {{ horarioLabel }}
                    </p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500 mb-1">Dias de atencion</p>
                    <div v-if="diasOperacion.length" class="flex flex-wrap gap-2">
                      <Tag v-for="(dia, idx) in diasOperacion" :key="idx" :value="dia" severity="info" />
                    </div>
                    <p v-else class="text-sm text-gray-500">Por definir</p>
                  </div>
                </div>
              </template>
            </Card>

            <Card v-if="fotosOrdenadas.length" class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-images text-blue-600"></i>
                  <span>Galeria de fotos</span>
                </div>
              </template>
              <template #content>
                <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
                  <a
                    v-for="(foto, idx) in fotosOrdenadas.slice(0, 9)"
                    :key="foto.id || idx"
                    :href="resolveFotoUrl(foto.foto_url || foto.foto)"
                    target="_blank"
                    rel="noopener"
                    class="relative group rounded-xl overflow-hidden border border-gray-200"
                  >
                    <img
                      :src="resolveFotoUrl(foto.foto_url || foto.foto)"
                      alt="Foto de la agencia"
                      class="w-full h-28 object-cover group-hover:scale-105 transition-transform duration-300"
                      loading="lazy"
                    />
                    <div class="absolute inset-0 bg-black/0 group-hover:bg-black/15 transition-colors" />
                  </a>
                </div>
              </template>
            </Card>

            <Card v-if="hasMap" class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-map text-emerald-600"></i>
                  <span>Ubicacion</span>
                </div>
              </template>
              <template #content>
                <ClientOnly>
                  <AgenciaMap :latitud="Number(agencia.latitud)" :longitud="Number(agencia.longitud)" :editable="false" :show-coordinate-inputs="false" height="360px" />
                </ClientOnly>
              </template>
            </Card>


          </div>

          <div class="space-y-6">
            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-phone text-emerald-600"></i>
                  <span>Contacto</span>
                </div>
              </template>
              <template #content>
                <div class="space-y-3 text-sm">
                  <div class="flex items-center gap-2">
                    <i class="pi pi-map-marker text-gray-500"></i>
                    <span>{{ agencia.direccion || 'Direccion por definir' }}</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <i class="pi pi-phone text-gray-500"></i>
                    <a v-if="telefonoLink" :href="telefonoLink" class="text-emerald-700 font-semibold hover:underline">
                      {{ agencia.telefono }}
                    </a>
                    <span v-else class="font-semibold text-gray-900">Por definir</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <i class="pi pi-envelope text-gray-500"></i>
                    <a v-if="agencia.email" :href="`mailto:${agencia.email}`" class="text-emerald-700 font-semibold hover:underline">
                      {{ agencia.email }}
                    </a>
                    <span v-else class="font-semibold text-gray-900">Por definir</span>
                  </div>
                </div>
              </template>
            </Card>

            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-file text-blue-600"></i>
                  <span>Politicas de la agencia</span>
                </div>
              </template>
              <template #content>
                <div class="space-y-3 text-sm text-gray-700">
                  <div class="flex items-start gap-2">
                    <i class="pi pi-users text-emerald-600 mt-0.5"></i>
                    <div>
                      <p class="text-xs text-gray-500">Edad minima para pago de ninos</p>
                      <p class="font-semibold text-gray-900">{{ edadMinimaPagoLabel }}</p>
                    </div>
                  </div>
                  <div class="flex items-start gap-2">
                    <i class="pi pi-info-circle text-blue-600 mt-0.5"></i>
                    <div>
                      <p class="text-xs text-gray-500">Politica de cancelacion</p>
                      <p class="text-gray-700 whitespace-pre-line">{{ politicaCancelacionLabel }}</p>
                    </div>
                  </div>
                </div>
              </template>
            </Card>

            <Card v-if="hasRedes" class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-globe text-blue-600"></i>
                  <span>Redes y sitio web</span>
                </div>
              </template>
              <template #content>
                <div class="space-y-3 text-sm">
                  <a v-if="agencia.sitio_web" :href="toExternalUrl(agencia.sitio_web)" target="_blank" rel="noopener" class="text-blue-700 font-semibold hover:underline">
                    Sitio web
                  </a>
                  <a v-if="agencia.facebook" :href="toExternalUrl(agencia.facebook)" target="_blank" rel="noopener" class="text-blue-700 font-semibold hover:underline">
                    Facebook
                  </a>
                  <a v-if="agencia.instagram" :href="toExternalUrl(agencia.instagram)" target="_blank" rel="noopener" class="text-blue-700 font-semibold hover:underline">
                    Instagram
                  </a>
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>
      <div ref="paquetesSection" class="scroll-mt-24">
        <Card class="surface-card">
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-briefcase text-blue-600"></i>
              <span>Paquetes de esta agencia</span>
            </div>
          </template>
          <template #content>
            <div v-if="paquetesLoading && paquetes.length === 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <Skeleton v-for="n in 4" :key="n" height="220px" />
            </div>

            <div v-else-if="paquetesError" class="space-y-3">
              <Message severity="warn" :closable="false">{{ paquetesError }}</Message>
              <Button label="Reintentar" icon="pi pi-refresh" outlined @click="loadPaquetes(true)" />
            </div>

            <div v-else>
              <div v-if="paquetes.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <PaquetesPaqueteCard v-for="paquete in paquetes" :key="paquete.id" :paquete="paquete" />
              </div>
              <div v-else class="text-sm text-gray-600">
                Esta agencia aun no tiene paquetes publicados.
              </div>

              <div v-if="paquetesPagination.total > 0" class="mt-4 text-sm text-gray-500">
                Mostrando {{ paquetes.length }} de {{ paquetesPagination.total }} paquetes
              </div>

              <div class="flex justify-center mt-6">
                <Button
                    v-if="canLoadMorePaquetes"
                    label="Cargar mas"
                    icon="pi pi-plus"
                    :loading="paquetesLoadingMore"
                    @click="loadMorePaquetes"
                />
              </div>
            </div>
          </template>
        </Card>
      </div>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, defineAsyncComponent } from 'vue'
import { useToast } from 'primevue/usetoast'

const AgenciaMap = defineAsyncComponent(() => import('~/components/atracciones/AtraccionMap.vue'))

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const route = useRoute()
const config = useRuntimeConfig()
const assetsBase = String(config.public.apiBase || '').replace(/\/api\/v1\/?$/, '')

const { getAgencia } = useAgencias()
const { getPaquetes } = usePaquetesTuristicos()

const agencia = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const paquetes = ref<any[]>([])
const paquetesError = ref<string | null>(null)
const paquetesLoading = ref(false)
const paquetesLoadingMore = ref(false)
const paquetesPagination = ref({
  page: 1,
  limit: 6,
  total: 0,
  total_pages: 0
})

const paquetesSection = ref<HTMLElement | null>(null)

const agenciaId = computed(() => Number(route.params.id))

const resolveFotoUrl = (path?: string) => {
  if (!path) return ''
  let normalized = String(path).replace(/\\\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const fotosOrdenadas = computed(() => {
  const fotos = (agencia.value?.fotos || []).slice()
  fotos.sort((a: any, b: any) => {
    if (!!a.es_principal !== !!b.es_principal) return a.es_principal ? -1 : 1
    return (a.orden || 0) - (b.orden || 0)
  })
  return fotos
})

const heroFoto = computed(() => {
  const principal = fotosOrdenadas.value[0]
  const path = principal?.foto_url || principal?.foto
  return path ? resolveFotoUrl(path) : ''
})

const especialidades = computed(() => {
  return (agencia.value?.especialidades || [])
    .map((item: any) => item?.categoria?.nombre || item?.categoria?.nombre_categoria || '')
    .filter(Boolean)
})

const diasOperacion = computed(() => {
  return (agencia.value?.dias || []).map((item: any) => item?.nombre).filter(Boolean)
})

const formatHora = (raw?: any) => {
  if (!raw) return ''
  const value = String(raw).trim()
  const match = value.match(/(\d{1,2}):(\d{2})/)
  if (match) return `${match[1].padStart(2, '0')}:${match[2]}`
  return value
}

const horarioLabel = computed(() => {
  const apertura = formatHora(agencia.value?.horario_apertura)
  const cierre = formatHora(agencia.value?.horario_cierre)
  if (apertura && cierre) return `${apertura} - ${cierre}`
  if (apertura) return `Desde ${apertura}`
  if (cierre) return `Hasta ${cierre}`
  return 'Por definir'
})

const telefonoLink = computed(() => {
  const digits = String(agencia.value?.telefono || '').replace(/\D/g, '')
  return digits ? `tel:${digits}` : ''
})

const edadMinimaPagoLabel = computed(() => {
  const value = Number(agencia.value?.politicas?.edad_minima_pago || 0)
  if (Number.isFinite(value) && value > 0) return `${value} anos`
  return 'Por definir'
})

const politicaCancelacionLabel = computed(() => {
  const raw = agencia.value?.politicas?.politica_cancelacion
  const value = raw ? String(raw).trim() : ''
  return value || 'Por definir'
})

const hasRedes = computed(() => {
  return !!(agencia.value?.sitio_web || agencia.value?.facebook || agencia.value?.instagram)
})

const hasMap = computed(() => {
  const lat = Number(agencia.value?.latitud)
  const lng = Number(agencia.value?.longitud)
  return Number.isFinite(lat) && Number.isFinite(lng)
})

const canLoadMorePaquetes = computed(() => {
  if (!paquetesPagination.value.total_pages) return false
  return paquetesPagination.value.page < paquetesPagination.value.total_pages
})

const toExternalUrl = (value?: string) => {
  if (!value) return ''
  const trimmed = String(value).trim()
  if (!trimmed) return ''
  if (/^https?:\/\//i.test(trimmed)) return trimmed
  return `https://${trimmed}`
}

const loadAgencia = async () => {
  loading.value = true
  error.value = null
  agencia.value = null
  try {
    const response: any = await getAgencia(agenciaId.value)
    if (response.success) {
      const data = response.data
      if (data?.visible_publico === false || (data?.status && data.status !== 'activa')) {
        error.value = 'La agencia no esta disponible para el publico'
        return
      }
      agencia.value = data
      await loadPaquetes(true)
      return
    }
    error.value = response?.error?.message || 'No se pudo cargar la agencia'
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudo cargar la agencia'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

const loadPaquetes = async (reset = false) => {
  if (!agenciaId.value) return
  paquetesError.value = null
  paquetesLoading.value = true
  if (reset) {
    paquetes.value = []
    paquetesPagination.value.page = 1
  }
  try {
    const response: any = await getPaquetes({
      page: paquetesPagination.value.page,
      limit: paquetesPagination.value.limit,
      agencia_id: agenciaId.value
    } as any)
    if (response.success) {
      const data = response.data?.paquetes || []
      paquetes.value = reset ? data : [...paquetes.value, ...data]
      paquetesPagination.value = { ...paquetesPagination.value, ...(response.data?.pagination || {}) }
      return
    }
    paquetesError.value = response?.error?.message || 'No se pudieron cargar los paquetes'
  } catch (err: any) {
    paquetesError.value = err?.data?.error?.message || err?.message || 'No se pudieron cargar los paquetes'
  } finally {
    paquetesLoading.value = false
  }
}

const loadMorePaquetes = async () => {
  if (!canLoadMorePaquetes.value) return
  paquetesLoadingMore.value = true
  paquetesPagination.value.page += 1
  try {
    await loadPaquetes(false)
  } finally {
    paquetesLoadingMore.value = false
  }
}

const scrollToPaquetes = () => {
  paquetesSection.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

onMounted(() => {
  loadAgencia()
})
</script>
