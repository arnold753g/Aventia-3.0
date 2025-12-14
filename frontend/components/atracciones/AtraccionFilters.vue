<template>
  <Card class="mb-6">
    <template #content>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <!-- Búsqueda -->
        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Buscar
          </label>
          <span class="p-input-icon-left w-full">
            <i class="pi pi-search" />
            <InputText
              v-model="localFilters.search"
              placeholder="Nombre o descripción..."
              class="w-full"
              @input="debouncedFilter"
            />
          </span>
        </div>

        <!-- Departamento -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Departamento
          </label>
          <Dropdown
            v-model="localFilters.departamento_id"
            :options="departamentos"
            optionLabel="nombre"
            optionValue="id"
            placeholder="Todos los departamentos"
            class="w-full"
            showClear
            @change="onDepartamentoChange"
          />
        </div>

        <!-- Provincia -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Provincia
          </label>
          <Dropdown
            v-model="localFilters.provincia_id"
            :options="provincias"
            optionLabel="nombre"
            optionValue="id"
            placeholder="Todas las provincias"
            class="w-full"
            showClear
            :disabled="!localFilters.departamento_id"
            @change="onFilterChange"
          />
        </div>

        <!-- Categoría -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Categoría
          </label>
          <Dropdown
            v-model="localFilters.categoria_id"
            :options="categorias"
            optionLabel="nombre"
            optionValue="id"
            placeholder="Todas las categorías"
            class="w-full"
            showClear
            @change="onCategoriaChange"
          />
        </div>

        <!-- Subcategoría -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Subcategoría
          </label>
          <Dropdown
            v-model="localFilters.subcategoria_id"
            :options="subcategorias"
            optionLabel="nombre"
            optionValue="id"
            placeholder="Todas las subcategorías"
            class="w-full"
            showClear
            :disabled="!localFilters.categoria_id"
            @change="onFilterChange"
          />
        </div>

        <!-- Nivel de dificultad -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Dificultad
          </label>
          <Dropdown
            v-model="localFilters.nivel_dificultad"
            :options="nivelesOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="Todos los niveles"
            class="w-full"
            showClear
            @change="onFilterChange"
          />
        </div>

        <!-- Estado -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Estado
          </label>
          <Dropdown
            v-model="localFilters.status"
            :options="statusOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="Todos los estados"
            class="w-full"
            showClear
            @change="onFilterChange"
          />
        </div>
      </div>

      <!-- Chips de filtros activos -->
      <div class="mt-4 flex items-center justify-between">
        <div class="flex flex-wrap gap-2">
          <Chip
            v-if="localFilters.search"
            :label="`Búsqueda: ${localFilters.search}`"
            removable
            @remove="clearSearch"
          />
          <Chip
            v-if="localFilters.departamento_id"
            :label="`Depto: ${getDepartamentoLabel(localFilters.departamento_id)}`"
            removable
            @remove="clearDepartamento"
          />
          <Chip
            v-if="localFilters.provincia_id"
            :label="`Provincia: ${getProvinciaLabel(localFilters.provincia_id)}`"
            removable
            @remove="clearProvincia"
          />
          <Chip
            v-if="localFilters.categoria_id"
            :label="`Categoría: ${getCategoriaLabel(localFilters.categoria_id)}`"
            removable
            @remove="clearCategoria"
          />
          <Chip
            v-if="localFilters.subcategoria_id"
            :label="`Subcategoría: ${getSubcategoriaLabel(localFilters.subcategoria_id)}`"
            removable
            @remove="clearSubcategoria"
          />
          <Chip
            v-if="localFilters.nivel_dificultad"
            :label="`Dificultad: ${getNivelDificultadLabel(localFilters.nivel_dificultad)}`"
            removable
            @remove="clearNivel"
          />
          <Chip
            v-if="localFilters.status"
            :label="`Estado: ${getStatusAtraccionLabel(localFilters.status)}`"
            removable
            @remove="clearStatus"
          />
        </div>
        <Button
          v-if="hasActiveFilters"
          label="Limpiar Filtros"
          icon="pi pi-filter-slash"
          severity="secondary"
          text
          @click="clearFilters"
        />
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getNivelDificultadLabel, getStatusAtraccionLabel } from '~/utils/formatters-atraccion'

const emit = defineEmits(['filter-change'])
const { getDepartamentos, getProvincias, getCategorias, getSubcategorias } = useAtracciones()

const localFilters = ref({
  search: '',
  departamento_id: '',
  provincia_id: '',
  categoria_id: '',
  subcategoria_id: '',
  nivel_dificultad: '',
  status: ''
})

const departamentos = ref<any[]>([])
const provincias = ref<any[]>([])
const categorias = ref<any[]>([])
const subcategorias = ref<any[]>([])

const nivelesOptions = [
  { label: 'Todos los niveles', value: '' },
  { label: 'Fácil', value: 'facil' },
  { label: 'Medio', value: 'medio' },
  { label: 'Difícil', value: 'dificil' },
  { label: 'Extremo', value: 'extremo' }
]

const statusOptions = [
  { label: 'Todos los estados', value: '' },
  { label: 'Activa', value: 'activa' },
  { label: 'Inactiva', value: 'inactiva' },
  { label: 'Mantenimiento', value: 'mantenimiento' },
  { label: 'Fuera de Temporada', value: 'fuera_temporada' }
]

const hasActiveFilters = computed(() => {
  return Object.values(localFilters.value).some(v => v !== '')
})

let debounceTimeout: NodeJS.Timeout

const debouncedFilter = () => {
  clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    onFilterChange()
  }, 500)
}

const onFilterChange = () => {
  emit('filter-change', { ...localFilters.value })
}

const onDepartamentoChange = async () => {
  localFilters.value.provincia_id = ''
  if (localFilters.value.departamento_id) {
    await loadProvincias(parseInt(localFilters.value.departamento_id))
  } else {
    provincias.value = []
  }
  onFilterChange()
}

const onCategoriaChange = async () => {
  localFilters.value.subcategoria_id = ''
  if (localFilters.value.categoria_id) {
    await loadSubcategorias(parseInt(localFilters.value.categoria_id))
  } else {
    subcategorias.value = []
  }
  onFilterChange()
}

const loadDepartamentos = async () => {
  try {
    const response = await getDepartamentos()
    if ((response as any).success) {
      departamentos.value = (response as any).data
    }
  } catch (error) {
    console.error('Error al cargar departamentos')
  }
}

const loadProvincias = async (departamentoId: number) => {
  try {
    const response = await getProvincias(departamentoId)
    if ((response as any).success) {
      provincias.value = (response as any).data
    }
  } catch (error) {
    console.error('Error al cargar provincias')
  }
}

const loadCategorias = async () => {
  try {
    const response = await getCategorias()
    if ((response as any).success) {
      categorias.value = (response as any).data
    }
  } catch (error) {
    console.error('Error al cargar categorías')
  }
}

const loadSubcategorias = async (categoriaId: number) => {
  try {
    const response = await getSubcategorias(categoriaId)
    if ((response as any).success) {
      subcategorias.value = (response as any).data
    }
  } catch (error) {
    console.error('Error al cargar subcategorías')
  }
}

// Labels helpers
const getDepartamentoLabel = (id: string) => {
  return departamentos.value.find(d => d.id === parseInt(id))?.nombre || id
}

const getProvinciaLabel = (id: string) => {
  return provincias.value.find(p => p.id === parseInt(id))?.nombre || id
}

const getCategoriaLabel = (id: string) => {
  return categorias.value.find(c => c.id === parseInt(id))?.nombre || id
}

const getSubcategoriaLabel = (id: string) => {
  return subcategorias.value.find(s => s.id === parseInt(id))?.nombre || id
}

// Clear functions
const clearSearch = () => {
  localFilters.value.search = ''
  onFilterChange()
}

const clearDepartamento = () => {
  localFilters.value.departamento_id = ''
  localFilters.value.provincia_id = ''
  provincias.value = []
  onFilterChange()
}

const clearProvincia = () => {
  localFilters.value.provincia_id = ''
  onFilterChange()
}

const clearCategoria = () => {
  localFilters.value.categoria_id = ''
  localFilters.value.subcategoria_id = ''
  subcategorias.value = []
  onFilterChange()
}

const clearSubcategoria = () => {
  localFilters.value.subcategoria_id = ''
  onFilterChange()
}

const clearNivel = () => {
  localFilters.value.nivel_dificultad = ''
  onFilterChange()
}

const clearStatus = () => {
  localFilters.value.status = ''
  onFilterChange()
}

const clearFilters = () => {
  localFilters.value = {
    search: '',
    departamento_id: '',
    provincia_id: '',
    categoria_id: '',
    subcategoria_id: '',
    nivel_dificultad: '',
    status: ''
  }
  provincias.value = []
  subcategorias.value = []
  onFilterChange()
}

onMounted(async () => {
  await loadDepartamentos()
  await loadCategorias()
})
</script>
