<template>
  <div>
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 mb-2">
        Subcategorías ({{ selected.length }}/4) *
      </label>
      <p class="text-xs text-gray-500 mb-3">
        Seleccione hasta 4 subcategorías. La primera será la principal.
      </p>
    </div>

    <!-- Categorías -->
    <div class="mb-4">
      <Dropdown
        v-model="selectedCategoria"
        :options="categorias"
        optionLabel="nombre"
        optionValue="id"
        placeholder="Seleccione una categoría"
        class="w-full"
        @change="loadSubcategorias"
      />
    </div>

    <!-- Subcategorías disponibles -->
    <div v-if="selectedCategoria" class="mb-4">
      <label class="block text-sm font-medium text-gray-700 mb-2">
        Subcategorías disponibles
      </label>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-2 max-h-48 overflow-y-auto p-2 border rounded-lg">
        <Button
          v-for="subcat in availableSubcategorias"
          :key="subcat.id"
          :label="subcat.nombre"
          :disabled="selected.length >= 4"
          outlined
          size="small"
          @click="addSubcategoria(subcat)"
        />
      </div>
    </div>

    <!-- Subcategorías seleccionadas -->
    <div v-if="selected.length > 0">
      <label class="block text-sm font-medium text-gray-700 mb-2">
        Seleccionadas (Orden: 1° es Principal)
      </label>
      <div class="space-y-2">
        <div
          v-for="(subcat, index) in selected"
          :key="subcat.id"
          class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border hover:border-blue-300 transition-colors"
        >
          <div class="flex items-center gap-3 flex-1">
            <div class="flex flex-col gap-1">
              <Button
                icon="pi pi-chevron-up"
                text
                rounded
                size="small"
                :disabled="index === 0"
                class="h-5"
                @click="moveUp(index)"
              />
              <Button
                icon="pi pi-chevron-down"
                text
                rounded
                size="small"
                :disabled="index === selected.length - 1"
                class="h-5"
                @click="moveDown(index)"
              />
            </div>
            <div class="flex items-center gap-2">
              <span class="text-gray-500 font-mono text-sm">{{ index + 1 }}.</span>
              <span v-if="index === 0" class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
                PRINCIPAL
              </span>
              <div class="flex flex-col">
                <span class="font-medium">{{ subcat.nombre }}</span>
                <span class="text-xs text-gray-500">{{ subcat.categoria?.nombre }}</span>
              </div>
            </div>
          </div>
          <Button
            icon="pi pi-times"
            severity="danger"
            text
            rounded
            size="small"
            @click="removeSubcategoria(index)"
          />
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else class="text-center py-8 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
      <i class="pi pi-tags text-4xl text-gray-400 mb-2"></i>
      <p class="text-gray-600 text-sm">Seleccione subcategorías de la lista</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: number[]
}>(), {
  modelValue: () => []
})

const emit = defineEmits(['update:modelValue'])

const { getCategorias, getSubcategorias } = useAtracciones()

const categorias = ref<any[]>([])
const subcategorias = ref<any[]>([])
const selected = ref<any[]>([])
const selectedCategoria = ref<number | null>(null)

const availableSubcategorias = computed(() => {
  return subcategorias.value.filter(s => !selected.value.find(sel => sel.id === s.id))
})

const loadCategorias = async () => {
  console.log('🔄 Intentando cargar categorías...')
  try {
    const response: any = await getCategorias()
    console.log('📡 Respuesta del servidor:', response)
    const data = response?.data ?? response
    console.log('📦 Datos extraídos:', data)
    categorias.value = Array.isArray(data) ? data : []
    console.log('✅ Categorías cargadas:', categorias.value)
  } catch (error) {
    console.error('❌ Error al cargar categorías', error)
    categorias.value = []
  }
}

const loadSubcategorias = async () => {
  if (!selectedCategoria.value) {
    subcategorias.value = []
    return
  }

  try {
    const response: any = await getSubcategorias(selectedCategoria.value)
    const data = response?.data ?? response
    subcategorias.value = Array.isArray(data) ? data : []
  } catch (error) {
    console.error('Error al cargar subcategorías', error)
    subcategorias.value = []
  }
}

const addSubcategoria = (subcat: any) => {
  if (selected.value.length >= 4) return
  selected.value.push(subcat)
  emitValue()
}

const removeSubcategoria = (index: number) => {
  selected.value.splice(index, 1)
  emitValue()
}

const moveUp = (index: number) => {
  console.log('⬆️ moveUp llamado con index:', index)
  if (index === 0) {
    console.log('⚠️ No se puede mover arriba, ya está en la primera posición')
    return
  }
  console.log('📊 Array antes:', selected.value.map(s => s.nombre))
  // Crear un nuevo array para forzar reactividad
  const newArray = [...selected.value]
  const temp = newArray[index]
  newArray[index] = newArray[index - 1]
  newArray[index - 1] = temp
  selected.value = newArray
  console.log('📊 Array después:', selected.value.map(s => s.nombre))
  emitValue()
}

const moveDown = (index: number) => {
  console.log('⬇️ moveDown llamado con index:', index)
  if (index === selected.value.length - 1) {
    console.log('⚠️ No se puede mover abajo, ya está en la última posición')
    return
  }
  console.log('📊 Array antes:', selected.value.map(s => s.nombre))
  // Crear un nuevo array para forzar reactividad
  const newArray = [...selected.value]
  const temp = newArray[index]
  newArray[index] = newArray[index + 1]
  newArray[index + 1] = temp
  selected.value = newArray
  console.log('📊 Array después:', selected.value.map(s => s.nombre))
  emitValue()
}

const emitValue = () => {
  const ids = selected.value.map(s => s.id)
  console.log('📤 SubcategoriasSelector emitiendo IDs:', ids, 'tipo:', Array.isArray(ids) ? 'array' : typeof ids)
  emit('update:modelValue', ids)
}

// Cargar subcategorías iniciales si hay valores
watch(() => props.modelValue, async (newVal) => {
  if (newVal && newVal.length > 0) {
    const response: any = await getSubcategorias()
    const data = response?.data ?? response
    if (Array.isArray(data)) {
      // Respetar el orden de los IDs en modelValue
      selected.value = newVal.map(id => data.find((s: any) => s.id === id)).filter(Boolean)
      console.log('🔄 Watcher actualizó selected con orden:', selected.value.map(s => s.nombre))
    }
  }
}, { immediate: true })

onMounted(() => {
  console.log('🎨 SubcategoriasSelector montado')
  loadCategorias()
})
</script>
