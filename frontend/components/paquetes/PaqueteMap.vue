<template>
  <div ref="mapContainer" class="w-full rounded-lg border border-gray-300" :style="{ height }" />
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

const props = withDefaults(defineProps<{
  atracciones?: any[]
  agencia?: any
  height?: string
}>(), {
  atracciones: () => [],
  height: '420px'
})

const mapContainer = ref<HTMLElement>()
let map: any = null
let markers: any[] = []
let polyline: any = null
let L: any = null

const points = computed(() => {
  const list: Array<{
    lat: number
    lng: number
    label: string
    kind: 'atraccion' | 'agencia'
    dia: number
    orden: number
  }> = []

  const items = (props.atracciones || []) as any[]
  items.forEach((item) => {
    const atr = item?.atraccion || item
    const lat = Number(atr?.latitud)
    const lng = Number(atr?.longitud)
    if (!Number.isFinite(lat) || !Number.isFinite(lng)) return
    list.push({
      lat,
      lng,
      label: atr?.nombre || `Atracción #${item?.atraccion_id || ''}`,
      kind: 'atraccion',
      dia: Number(item?.dia_numero || 0),
      orden: Number(item?.orden_visita || 0)
    })
  })

  const aLat = Number(props.agencia?.latitud)
  const aLng = Number(props.agencia?.longitud)
  if (Number.isFinite(aLat) && Number.isFinite(aLng)) {
    list.push({
      lat: aLat,
      lng: aLng,
      label: props.agencia?.nombre_comercial || 'Agencia',
      kind: 'agencia',
      dia: -1,
      orden: -1
    })
  }

  return list
    .slice()
    .sort((a, b) => (a.dia - b.dia) || (a.orden - b.orden))
})

onMounted(async () => {
  if (!process.client || !mapContainer.value) return

  try {
    const leafletModule = await import('leaflet')
    L = leafletModule.default
    await import('leaflet/dist/leaflet.css')

    // Fix de iconos en entornos bundler
    try {
      delete (L.Icon.Default.prototype as any)._getIconUrl
      L.Icon.Default.mergeOptions({
        iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
        iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
        shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png'
      })
    } catch {
      // ignore
    }

    map = L.map(mapContainer.value, { zoomControl: true }).setView([-16.5, -68.15], 6)
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map)

    refresh()
  } catch (error) {
    console.error('Error cargando mapa de paquete:', error)
  }
})

watch(points, () => {
  if (!map) return
  refresh()
}, { deep: true })

const clearLayers = () => {
  markers.forEach((m) => {
    try {
      m.remove()
    } catch {
      // ignore
    }
  })
  markers = []

  if (polyline) {
    try {
      polyline.remove()
    } catch {
      // ignore
    }
    polyline = null
  }
}

const refresh = () => {
  if (!map || !L) return
  clearLayers()

  const pts = points.value
  const routePts: Array<[number, number]> = []

  pts.forEach((p) => {
    const marker = L.marker([p.lat, p.lng], {
      title: p.label
    }).addTo(map)

    marker.bindPopup(`<strong>${p.label}</strong>`)
    markers.push(marker)

    if (p.kind === 'atraccion') routePts.push([p.lat, p.lng])
  })

  if (routePts.length >= 2) {
    polyline = L.polyline(routePts, {
      color: '#10b981',
      weight: 4,
      opacity: 0.75
    }).addTo(map)
  }

  if (markers.length > 0) {
    const group = L.featureGroup(markers)
    map.fitBounds(group.getBounds().pad(0.15))
  }
}
</script>

<style scoped>
:deep(.leaflet-container) {
  z-index: 0;
}
</style>
