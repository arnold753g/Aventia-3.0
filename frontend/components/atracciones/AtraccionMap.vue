<template>
  <div>
    <div ref="mapContainer" class="w-full rounded-lg border border-gray-300" :style="{ height: height }"></div>

    <div v-if="editable && showInputs" class="mt-4 grid grid-cols-2 gap-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Latitud
        </label>
        <InputText
          :modelValue="localLatitud"
          @update:modelValue="updateLatitud"
          type="number"
          step="0.000001"
          placeholder="-22.90 a -9.67"
          class="w-full"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Longitud
        </label>
        <InputText
          :modelValue="localLongitud"
          @update:modelValue="updateLongitud"
          type="number"
          step="0.000001"
          placeholder="-69.65 a -57.45"
          class="w-full"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'

const props = defineProps<{
  latitud?: number
  longitud?: number
  editable?: boolean
  showCoordinateInputs?: boolean
  height?: string
}>()

const emit = defineEmits(['update:latitud', 'update:longitud'])

const mapContainer = ref<HTMLElement>()
let map: any = null
let marker: any = null
let L: any = null

const localLatitud = ref(props.latitud)
const localLongitud = ref(props.longitud)

const defaultCenter: [number, number] = [-16.5, -68.15] // Centro de Bolivia
const defaultZoom = 6
const showInputs = computed(() => props.showCoordinateInputs !== false)

onMounted(async () => {
  // Solo importar Leaflet en el cliente
  if (process.client) {
    try {
      const leafletModule = await import('leaflet')
      L = leafletModule.default
      await import('leaflet/dist/leaflet.css')
      console.log('✅ Leaflet cargado:', L)
      initMap()
    } catch (error) {
      console.error('❌ Error al cargar Leaflet:', error)
    }
  }
})

const initMap = () => {
  if (!mapContainer.value || !L) return

  // Crear mapa centrado en Bolivia
  map = L.map(mapContainer.value).setView(defaultCenter, defaultZoom)

  // Agregar capa de OpenStreetMap
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
  }).addTo(map)

  // Si hay coordenadas, agregar marcador
  if (props.latitud && props.longitud) {
    addMarker(props.latitud, props.longitud)
    map.setView([props.latitud, props.longitud], 12)
  }

  // Si es editable, permitir click en mapa
  if (props.editable) {
    map.on('click', (e: any) => {
      const { lat, lng } = e.latlng
      updateCoordinates(lat, lng)
    })
  }
}

const addMarker = (lat: number, lng: number) => {
  if (!map || !L) return

  // Remover marcador anterior si existe
  if (marker) {
    map.removeLayer(marker)
  }

  // Crear nuevo marcador
  marker = L.marker([lat, lng], {
    draggable: props.editable
  }).addTo(map)

  if (props.editable) {
    marker.on('dragend', (e: any) => {
      const { lat, lng } = e.target.getLatLng()
      updateCoordinates(lat, lng)
    })
  }
}

const updateCoordinates = (lat: number, lng: number) => {
  localLatitud.value = lat
  localLongitud.value = lng
  emit('update:latitud', lat)
  emit('update:longitud', lng)
  addMarker(lat, lng)
}

const updateLatitud = (value: any) => {
  const lat = parseFloat(value)
  if (!isNaN(lat) && localLongitud.value) {
    updateCoordinates(lat, localLongitud.value)
    map?.setView([lat, localLongitud.value], 12)
  }
}

const updateLongitud = (value: any) => {
  const lng = parseFloat(value)
  if (!isNaN(lng) && localLatitud.value) {
    updateCoordinates(localLatitud.value, lng)
    map?.setView([localLatitud.value, lng], 12)
  }
}

watch([() => props.latitud, () => props.longitud], ([newLat, newLng]) => {
  if (newLat && newLng && map) {
    addMarker(newLat, newLng)
    map.setView([newLat, newLng], 12)
  }
})
</script>

<style scoped>
:deep(.leaflet-container) {
  z-index: 0;
}
</style>
