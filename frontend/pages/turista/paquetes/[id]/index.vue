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
      <h2 class="text-2xl font-bold text-gray-900">No se pudo cargar el paquete</h2>
      <p class="muted">{{ error }}</p>
      <div class="flex justify-center gap-2">
        <Button label="Volver" icon="pi pi-arrow-left" severity="secondary" outlined @click="navigateTo('/turista/paquetes')" />
        <Button label="Reintentar" icon="pi pi-refresh" @click="loadPaquete" />
      </div>
    </div>

    <div v-else-if="paquete">
      <!-- Hero -->
      <section class="relative bg-white border-b border-gray-200">
        <div class="relative h-[62vh] min-h-[380px] overflow-hidden">
          <ClientOnly>
            <Carousel v-if="fotosOrdenadas.length" :items-to-show="1" :wrap-around="true" :transition="650" class="h-full">
              <Slide v-for="(foto, idx) in fotosOrdenadas" :key="foto.id || idx">
                <div class="relative h-[62vh] min-h-[380px]">
                  <img
                    :src="resolveFotoUrl(foto.foto)"
                    :alt="paquete.nombre"
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

          <div
            v-if="!fotosOrdenadas.length"
            class="w-full h-full bg-gradient-to-br from-emerald-100 to-blue-100 flex items-center justify-center"
          >
            <i class="pi pi-image text-6xl text-gray-400"></i>
          </div>

          <div class="absolute inset-0 bg-gradient-to-t from-black/75 via-black/20 to-transparent pointer-events-none" />

          <Button
            icon="pi pi-arrow-left"
            class="!absolute top-4 left-4 z-20 !bg-white/80 backdrop-blur !border-gray-200"
            rounded
            @click="navigateTo('/turista/paquetes')"
          />

          <Button
            label="Comprar"
            icon="pi pi-shopping-cart"
            class="!absolute top-4 right-4 z-20 !bg-white/80 backdrop-blur !border-gray-200"
            @click="navigateTo(`/turista/paquetes/${paquete.id}/comprar`)"
          />

          <div class="absolute bottom-0 left-0 right-0 z-10">
            <div class="max-w-7xl mx-auto px-4 pb-8">
              <div class="flex flex-wrap gap-2 mb-3">
                <Tag :value="frecuenciaLabel" severity="info" />
                <Tag v-if="duracionChip" :value="duracionChip" severity="secondary" />
                <Tag v-if="dificultadLabel" :value="dificultadLabel" :severity="dificultadSeverity" />
                <Tag v-if="paquete.permite_privado" value="Privado disponible" severity="warning" icon="pi pi-lock" />
                <Tag
                  v-if="paquete.frecuencia === 'salida_unica' && paquete.fecha_salida_fija"
                  :value="`Salida: ${formatFecha(paquete.fecha_salida_fija)}`"
                  severity="secondary"
                  icon="pi pi-calendar"
                />
              </div>

              <h1 ref="heroTitle" class="text-4xl md:text-5xl font-bold text-white leading-tight">
                {{ paquete.nombre }}
              </h1>

              <div class="flex flex-wrap items-center gap-4 text-white/90 mt-3">
                <div v-if="paquete.agencia?.nombre_comercial" class="flex items-center gap-2">
                  <i class="pi pi-building"></i>
                  <span class="font-medium">{{ paquete.agencia.nombre_comercial }}</span>
                  <span v-if="paquete.agencia?.departamento?.nombre" class="text-white/70">· {{ paquete.agencia.departamento.nombre }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <i class="pi pi-tag"></i>
                  <span class="font-semibold">Bs. {{ formatMoney(paquete.precio_base_nacionales) }}</span>
                  <span v-if="Number(paquete.precio_adicional_extranjeros || 0) > 0" class="text-white/70">
                    (+ Bs. {{ formatMoney(paquete.precio_adicional_extranjeros) }} extranjeros)
                  </span>
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
                  <i class="pi pi-info-circle text-emerald-600"></i>
                  <span>Descripción</span>
                </div>
              </template>
              <template #content>
                <p class="text-gray-700 whitespace-pre-line">
                  {{ paquete.descripcion || 'No hay descripción disponible.' }}
                </p>
              </template>
            </Card>

            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-sparkles text-blue-600"></i>
                  <span>Qué incluye / No incluye / Qué llevar</span>
                </div>
              </template>
              <template #content>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div>
                    <p class="text-sm font-semibold text-gray-900 mb-2">Incluye</p>
                    <ul v-if="(paquete.incluye || []).length" class="space-y-2">
                      <li v-for="(item, idx) in paquete.incluye" :key="idx" class="flex items-start gap-2 text-sm text-gray-700">
                        <i class="pi pi-check-circle text-emerald-600 mt-0.5"></i>
                        <span>{{ item }}</span>
                      </li>
                    </ul>
                    <p v-else class="text-sm text-gray-500">Sin especificar.</p>
                  </div>

                  <div>
                    <p class="text-sm font-semibold text-gray-900 mb-2">No incluye</p>
                    <ul v-if="(paquete.no_incluye || []).length" class="space-y-2">
                      <li v-for="(item, idx) in paquete.no_incluye" :key="idx" class="flex items-start gap-2 text-sm text-gray-700">
                        <i class="pi pi-times-circle text-rose-500 mt-0.5"></i>
                        <span>{{ item }}</span>
                      </li>
                    </ul>
                    <p v-else class="text-sm text-gray-500">Sin especificar.</p>
                  </div>

                  <div>
                    <p class="text-sm font-semibold text-gray-900 mb-2">Qué llevar</p>
                    <ul v-if="(paquete.que_llevar || []).length" class="space-y-2">
                      <li v-for="(item, idx) in paquete.que_llevar" :key="idx" class="flex items-start gap-2 text-sm text-gray-700">
                        <i class="pi pi-briefcase text-slate-600 mt-0.5"></i>
                        <span>{{ item }}</span>
                      </li>
                    </ul>
                    <p v-else class="text-sm text-gray-500">Sin especificar.</p>
                  </div>
                </div>
              </template>
            </Card>

            <Card class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-route text-emerald-600"></i>
                  <span>Itinerario</span>
                </div>
              </template>
              <template #content>
                <div v-if="!isMultiDay" class="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div>
                    <p class="text-sm text-gray-500 mb-1">Horario</p>
                    <p class="text-gray-900 font-semibold">{{ horarioLabel }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500 mb-1">Hora de salida</p>
                    <p class="text-gray-900 font-semibold">{{ horaSalidaLabel }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500 mb-1">Duración</p>
                    <p class="text-gray-900 font-semibold">{{ duracionHorasLabel }}</p>
                  </div>
                </div>

                <div v-else class="space-y-4">
                  <div v-if="(paquete.itinerario || []).length" class="space-y-3">
                    <div
                      v-for="item in itinerarioOrdenado"
                      :key="item.id"
                      class="rounded-xl border border-gray-200 bg-white p-4"
                    >
                      <div class="flex items-start justify-between gap-3">
                        <div>
                          <p class="text-xs uppercase tracking-wider text-emerald-700/80 font-semibold">Día {{ item.dia_numero }}</p>
                          <h3 class="text-lg font-semibold text-gray-900 mt-1">{{ item.titulo }}</h3>
                        </div>
                        <Tag :value="`Día ${item.dia_numero}`" severity="secondary" />
                      </div>

                      <p v-if="item.descripcion" class="text-sm text-gray-700 mt-2 whitespace-pre-line">
                        {{ item.descripcion }}
                      </p>

                      <div v-if="(item.actividades || []).length" class="mt-3">
                        <p class="text-sm font-semibold text-gray-900 mb-2">Actividades</p>
                        <ul class="grid grid-cols-1 md:grid-cols-2 gap-2">
                          <li v-for="(act, idx) in item.actividades" :key="idx" class="text-sm text-gray-700 flex items-start gap-2">
                            <i class="pi pi-check text-emerald-600 mt-0.5"></i>
                            <span>{{ act }}</span>
                          </li>
                        </ul>
                      </div>

                      <div v-if="item.hospedaje_info" class="mt-3 p-3 rounded-lg bg-slate-50 border border-slate-200">
                        <p class="text-xs text-slate-500 mb-1">Hospedaje</p>
                        <p class="text-sm text-slate-700 whitespace-pre-line">{{ item.hospedaje_info }}</p>
                      </div>
                    </div>
                  </div>

                  <div v-else class="text-sm text-gray-600">
                    Este paquete no tiene itinerario detallado aún. Revisa las atracciones incluidas.
                  </div>
                </div>
              </template>
            </Card>

            <Card v-if="atraccionesOrdenadas.length" class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-map-marker text-blue-600"></i>
                  <span>Atracciones incluidas</span>
                </div>
              </template>
              <template #content>
                <div class="space-y-6">

                  <div v-if="isMultiDay" class="space-y-6">
                  <div v-for="dia in diasAtracciones" :key="dia" class="space-y-3">
                    <div class="flex items-center justify-between">
                      <h3 class="text-base font-semibold text-gray-900">Día {{ dia }}</h3>
                      <span class="text-xs text-gray-500">{{ (atraccionesPorDia[dia] || []).length }} atracciones</span>
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div
                        v-for="item in atraccionesPorDia[dia]"
                        :key="item.id"
                        class="rounded-xl border border-gray-200 bg-white overflow-hidden hover:shadow-md transition-shadow"
                      >
                        <div class="flex gap-4 p-4">
                          <div class="w-20 h-20 rounded-lg overflow-hidden bg-slate-100 flex-shrink-0">
                            <img
                              v-if="getAtraccionCover(item.atraccion)"
                              :src="getAtraccionCover(item.atraccion)"
                              :alt="item.atraccion?.nombre || 'Atracción'"
                              class="w-full h-full object-cover"
                              loading="lazy"
                            />
                            <div v-else class="w-full h-full flex items-center justify-center">
                              <i class="pi pi-image text-gray-400 text-2xl"></i>
                            </div>
                          </div>

                          <div class="min-w-0 flex-1">
                            <div class="flex items-start justify-between gap-2">
                              <p class="font-semibold text-gray-900 truncate">{{ item.atraccion?.nombre || `Atracción #${item.atraccion_id}` }}</p>
                              <span class="text-xs text-gray-500">#{{ item.orden_visita }}</span>
                            </div>
                            <p class="text-xs text-gray-500 mt-1">
                              {{ item.atraccion?.provincia?.nombre || '' }}
                              <span v-if="item.atraccion?.provincia?.departamento?.nombre">
                                , {{ item.atraccion.provincia.departamento.nombre }}
                              </span>
                            </p>

                            <div class="flex flex-wrap items-center gap-3 mt-3 text-xs text-gray-700">
                              <span class="inline-flex items-center gap-1">
                                <i class="pi pi-clock"></i>{{ item.duracion_estimada_horas || 1 }}h
                              </span>
                              <Button
                                icon="pi pi-external-link"
                                label="Ver"
                                size="small"
                                text
                                @click="goToAtraccion(item.atraccion_id)"
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  </div>

                  <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div
                    v-for="item in atraccionesOrdenadas"
                    :key="item.id"
                    class="rounded-xl border border-gray-200 bg-white overflow-hidden hover:shadow-md transition-shadow"
                  >
                    <div class="flex gap-4 p-4">
                      <div class="w-20 h-20 rounded-lg overflow-hidden bg-slate-100 flex-shrink-0">
                        <img
                          v-if="getAtraccionCover(item.atraccion)"
                          :src="getAtraccionCover(item.atraccion)"
                          :alt="item.atraccion?.nombre || 'Atracción'"
                          class="w-full h-full object-cover"
                          loading="lazy"
                        />
                        <div v-else class="w-full h-full flex items-center justify-center">
                          <i class="pi pi-image text-gray-400 text-2xl"></i>
                        </div>
                      </div>

                      <div class="min-w-0 flex-1">
                        <div class="flex items-start justify-between gap-2">
                          <p class="font-semibold text-gray-900 truncate">{{ item.atraccion?.nombre || `Atracción #${item.atraccion_id}` }}</p>
                          <span class="text-xs text-gray-500">#{{ item.orden_visita }}</span>
                        </div>
                        <p class="text-xs text-gray-500 mt-1">
                          {{ item.atraccion?.provincia?.nombre || '' }}
                          <span v-if="item.atraccion?.provincia?.departamento?.nombre">
                            , {{ item.atraccion.provincia.departamento.nombre }}
                          </span>
                        </p>

                        <div class="flex flex-wrap items-center gap-3 mt-3 text-xs text-gray-700">
                          <span class="inline-flex items-center gap-1">
                            <i class="pi pi-clock"></i>{{ item.duracion_estimada_horas || 1 }}h
                          </span>
                          <Button icon="pi pi-external-link" label="Ver" size="small" text @click="goToAtraccion(item.atraccion_id)" />
                        </div>
                      </div>
                    </div>
                  </div>
                  </div>
                </div>
              </template>
            </Card>

            <Card v-if="hasMap" class="surface-card">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-map text-emerald-600"></i>
                  <span>Mapa del recorrido</span>
                </div>
              </template>
              <template #content>
                <ClientOnly>
                  <PaquetesPaqueteMap :atracciones="paquete.atracciones || []" :agencia="paquete.agencia" height="460px" />
                </ClientOnly>
              </template>
            </Card>
          </div>

          <!-- Sidebar -->
          <div class="sticky top-24 z-10">
            <div class="space-y-6">
              <Card class="surface-card">
              <template #title>Información rápida</template>
              <template #content>
                <div class="space-y-4">
                  <div>
                    <p class="text-xs text-gray-500 mb-1">Precio base (nacionales)</p>
                    <p class="text-2xl font-bold text-emerald-600">Bs. {{ formatMoney(paquete.precio_base_nacionales) }}</p>
                    <p v-if="Number(paquete.precio_adicional_extranjeros || 0) > 0" class="text-sm text-gray-600 mt-1">
                      Extranjeros: +Bs. {{ formatMoney(paquete.precio_adicional_extranjeros) }}
                    </p>
                  </div>

                  <Divider />

                  <div class="grid grid-cols-2 gap-4">
                    <div>
                      <p class="text-xs text-gray-500 mb-1">Duración</p>
                      <p class="text-sm font-semibold text-gray-900">{{ duracionTexto }}</p>
                    </div>
                    <div>
                      <p class="text-xs text-gray-500 mb-1">Cupos</p>
                      <p class="text-sm font-semibold text-gray-900">{{ paquete.cupo_minimo }}-{{ paquete.cupo_maximo }}</p>
                    </div>
                    <div>
                      <p class="text-xs text-gray-500 mb-1">Compra con</p>
                      <p class="text-sm font-semibold text-gray-900">{{ paquete.dias_previos_compra || 1 }} días ant.</p>
                    </div>
                    <div v-if="paquete.nivel_dificultad">
                      <p class="text-xs text-gray-500 mb-1">Dificultad</p>
                      <Tag :value="dificultadLabel" :severity="dificultadSeverity" class="w-full justify-center" />
                    </div>
                  </div>

                  <Divider />

                  <div v-if="paquete.politicas" class="space-y-2">
                    <p class="text-xs uppercase tracking-wider text-gray-500 font-semibold">Políticas</p>
                    <div class="text-sm text-gray-700 space-y-2">
                      <div class="flex items-start gap-2">
                        <i class="pi pi-users text-emerald-600 mt-0.5"></i>
                        <span>Niños pagan desde: <strong>{{ paquete.politicas.edad_minima_pago }} años</strong></span>
                      </div>
                      <div v-if="paquete.permite_privado" class="flex items-start gap-2">
                        <i class="pi pi-lock text-amber-600 mt-0.5"></i>
                        <span>Recargo privado: <strong>{{ formatMoney(paquete.politicas.recargo_privado_porcentaje) }}%</strong></span>
                      </div>
                      <div v-if="paquete.politicas.politica_cancelacion" class="flex items-start gap-2">
                        <i class="pi pi-info-circle text-blue-600 mt-0.5"></i>
                        <span class="whitespace-pre-line">{{ paquete.politicas.politica_cancelacion }}</span>
                      </div>
                    </div>
                  </div>

                  <Divider v-if="paquete.agencia" />

                  <div v-if="paquete.agencia" class="space-y-3">
                    <p class="text-xs uppercase tracking-wider text-gray-500 font-semibold">Adquiere tu paquete </p>

                    <Button
                      label="Comprar paquete"
                      icon="pi pi-shopping-cart"
                      class="w-full"
                      @click="navigateTo(`/turista/paquetes/${paquete.id}/comprar`)"
                    />
                  </div>
                </div>
              </template>
              </Card>

              <Card v-if="fotosOrdenadas.length" class="surface-card">
                <template #title>
                  <div class="flex items-center gap-2">
                    <i class="pi pi-images text-blue-600"></i>
                    <span>Galería</span>
                  </div>
                </template>
                <template #content>
                  <div class="grid grid-cols-3 gap-3">
                    <button
                      v-for="(foto, idx) in fotosOrdenadas.slice(0, 9)"
                      :key="foto.id || idx"
                      type="button"
                      class="relative group rounded-xl overflow-hidden border border-gray-200"
                      @click="openFoto(resolveFotoUrl(foto.foto))"
                    >
                      <img
                        :src="resolveFotoUrl(foto.foto)"
                        :alt="paquete.nombre"
                        class="w-full h-24 object-cover group-hover:scale-105 transition-transform duration-300"
                        loading="lazy"
                      />
                      <div class="absolute inset-0 bg-black/0 group-hover:bg-black/15 transition-colors" />
                    </button>
                  </div>
                </template>
              </Card>
            </div>
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
import Timeline from 'primevue/timeline'
import { Carousel, Slide, Pagination, Navigation } from 'vue3-carousel'
import 'vue3-carousel/dist/carousel.css'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const route = useRoute()
const toast = useToast()
const { getPaquete } = usePaquetesTuristicos()
const config = useRuntimeConfig()
const assetsBase = String(config.public.apiBase || '').replace(/\/api\/v1\/?$/, '')

const loading = ref(true)
const error = ref<string | null>(null)
const paquete = ref<any>(null)

const fotoDialog = ref(false)
const fotoSeleccionada = ref<string | null>(null)
const heroTitle = ref<HTMLElement | null>(null)

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
  const fotos = (paquete.value?.fotos || []).slice()
  fotos.sort((a: any, b: any) => {
    if (!!a.es_principal !== !!b.es_principal) return a.es_principal ? -1 : 1
    return (a.orden || 0) - (b.orden || 0)
  })
  return fotos
})

const itinerarioOrdenado = computed(() => {
  const items = (paquete.value?.itinerario || []).slice()
  items.sort((a: any, b: any) => (a.dia_numero || 0) - (b.dia_numero || 0))
  return items
})

const atraccionesOrdenadas = computed(() => {
  const items = (paquete.value?.atracciones || []).slice()
  items.sort((a: any, b: any) => {
    const diaA = Number(a?.dia_numero || 0)
    const diaB = Number(b?.dia_numero || 0)
    if (diaA !== diaB) return diaA - diaB
    return Number(a?.orden_visita || 0) - Number(b?.orden_visita || 0)
  })
  return items
})

const atraccionesTimeline = computed(() => {
  return atraccionesOrdenadas.value.map((item: any, index: number) => {
    const diaNumero = Math.max(1, Number(item?.dia_numero || 1))
    const ordenVisita = Number(item?.orden_visita || 0)
    const atraccion = item?.atraccion || {}
    const provincia = atraccion?.provincia?.nombre || ''
    const departamento = atraccion?.provincia?.departamento?.nombre || ''
    const timelineLocation = [provincia, departamento].filter(Boolean).join(', ')
    return {
      ...item,
      diaNumero,
      timelineOrden: ordenVisita > 0 ? ordenVisita : index + 1,
      timelineNombre: atraccion?.nombre || `Atraccion #${item?.atraccion_id}`,
      timelineLocation
    }
  })
})

const isMultiDay = computed(() => Number(paquete.value?.duracion_dias || 1) > 1)

const duracionChip = computed(() => {
  const dias = Number(paquete.value?.duracion_dias || 1)
  if (!Number.isFinite(dias) || dias <= 1) return null
  const noches = Number(paquete.value?.duracion_noches || (dias - 1))
  return `${dias}D/${Math.max(0, noches)}N`
})

const duracionTexto = computed(() => {
  const dias = Number(paquete.value?.duracion_dias || 1)
  if (!Number.isFinite(dias) || dias <= 1) return '1 día'
  const noches = Number(paquete.value?.duracion_noches || (dias - 1))
  return `${dias} días / ${Math.max(0, noches)} noches`
})

const frecuenciaLabel = computed(() => {
  const map: Record<string, string> = { salida_diaria: 'Salida diaria', salida_unica: 'Salida única' }
  return map[paquete.value?.frecuencia || ''] || (paquete.value?.frecuencia || 'N/D')
})

const dificultadLabel = computed(() => {
  const nivel = paquete.value?.nivel_dificultad
  if (!nivel) return null
  const map: Record<string, string> = { facil: 'Fácil', medio: 'Medio', dificil: 'Difícil', extremo: 'Extremo' }
  return map[nivel] || String(nivel)
})

const dificultadSeverity = computed(() => {
  const nivel = paquete.value?.nivel_dificultad
  const map: Record<string, string> = { facil: 'success', medio: 'info', dificil: 'warning', extremo: 'danger' }
  return map[nivel] || 'secondary'
})

const horarioLabel = computed(() => {
  const h = paquete.value?.horario
  if (!h) return 'N/D'
  const map: Record<string, string> = { mañana: 'Mañana', tarde: 'Tarde', todo_dia: 'Todo el día' }
  return map[h] || String(h)
})

const horaSalidaLabel = computed(() => {
  const raw = (paquete.value as any)?.hora_salida ?? (paquete.value as any)?.horaSalida ?? (paquete.value as any)?.HoraSalida
  if (raw === null || raw === undefined) return 'N/D'
  const value = String(raw).trim()
  if (!value) return 'N/D'

  const match = value.match(/(\d{1,2}):(\d{2})/)
  if (match) {
    const hh = match[1].padStart(2, '0')
    return `${hh}:${match[2]}`
  }

  const parsed = new Date(value)
  if (!Number.isNaN(parsed.getTime())) {
    const hh = String(parsed.getHours()).padStart(2, '0')
    const mm = String(parsed.getMinutes()).padStart(2, '0')
    return `${hh}:${mm}`
  }

  return value
})

const duracionHorasLabel = computed(() => {
  const raw = paquete.value?.duracion_horas
  if (!raw) return 'N/D'
  const m = String(raw).match(/(\d+)/)
  if (!m) return String(raw)
  const hours = Number(m[1])
  return hours === 1 ? '1 hora' : `${hours} horas`
})

const atraccionesPorDia = computed(() => {
  const map: Record<number, any[]> = {}
  atraccionesOrdenadas.value.forEach((item: any) => {
    const dia = Math.max(1, Number(item?.dia_numero || 1))
    map[dia] ||= []
    map[dia].push(item)
  })
  return map
})

const diasAtracciones = computed(() => Object.keys(atraccionesPorDia.value).map((k) => Number(k)).sort((a, b) => a - b))

const hasMap = computed(() => {
  const items = atraccionesOrdenadas.value
  const hasAtr = items.some((it: any) => Number.isFinite(Number(it?.atraccion?.latitud)) && Number.isFinite(Number(it?.atraccion?.longitud)))
  const hasAg = Number.isFinite(Number(paquete.value?.agencia?.latitud)) && Number.isFinite(Number(paquete.value?.agencia?.longitud))
  return hasAtr || hasAg
})

const formatMoney = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatFecha = (value?: any) => {
  if (!value) return ''
  const raw = String(value)
  const datePart = raw.split('T')[0].split(' ')[0]
  const match = datePart.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`
  return datePart || raw
}

const getAtraccionCover = (atraccion: any) => {
  const fotos = atraccion?.fotos || []
  const principal = fotos.find((f: any) => f.es_principal) || fotos[0]
  return principal?.foto ? resolveFotoUrl(principal.foto) : ''
}

const openFoto = (url: string) => {
  fotoSeleccionada.value = url
  fotoDialog.value = true
}

const goToAtraccion = (atraccionId: any) => {
  if (!atraccionId) return
  navigateTo(`/turista/atracciones/${atraccionId}`)
}

const animateHero = async () => {
  if (!heroTitle.value) return
  try {
    const { gsap } = await import('gsap')
    gsap.from(heroTitle.value, { opacity: 0, y: 14, duration: 0.55, ease: 'power2.out' })
  } catch {
    // no-op
  }
}

const loadPaquete = async () => {
  loading.value = true
  error.value = null
  paquete.value = null
  try {
    const id = Number(route.params.id)
    const response: any = await getPaquete(id)
    if (response.success) {
      paquete.value = response.data
      await animateHero()
    } else {
      error.value = response.error?.message || 'No se pudo cargar el paquete'
    }
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error inesperado'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3500 })
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadPaquete()
})
</script>
