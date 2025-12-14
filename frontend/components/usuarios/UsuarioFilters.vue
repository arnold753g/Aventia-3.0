<template>
  <Card class="surface-card">
    <template #content>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <!-- Búsqueda con debounce -->
        <div class="md:col-span-2">
          <IconField>
            <InputIcon class="pi pi-search" />
            <InputText
              v-model="localFilters.search"
              placeholder="Buscar por nombre, email o CI..."
              class="w-full"
              @input="debouncedFilter"
            />
          </IconField>
        </div>

        <!-- Filtro por rol -->
        <Select
          v-model="localFilters.rol"
          :options="roleOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Filtrar por rol"
          class="w-full"
          showClear
          @change="onFilterChange"
        />

        <!-- Filtro por estado -->
        <Select
          v-model="localFilters.status"
          :options="statusOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Filtrar por estado"
          class="w-full"
          showClear
          @change="onFilterChange"
        />
      </div>

      <!-- Chips de filtros activos -->
      <div v-if="hasActiveFilters" class="flex gap-2 mt-4 items-center justify-between">
        <div class="flex gap-2">
          <Chip
            v-if="localFilters.search"
            :label="`Búsqueda: ${localFilters.search}`"
            removable
            @remove="clearSearch"
          />
          <Chip
            v-if="localFilters.rol"
            :label="`Rol: ${getRolLabel(localFilters.rol)}`"
            removable
            @remove="clearRol"
          />
          <Chip
            v-if="localFilters.status"
            :label="`Estado: ${getStatusLabel(localFilters.status)}`"
            removable
            @remove="clearStatus"
          />
        </div>
        <Button
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
import { ref, computed } from 'vue'
import { getRolLabel, getStatusLabel } from '~/utils/formatters'

const emit = defineEmits(['filter-change'])

const localFilters = ref({
  search: '',
  rol: '',
  status: ''
})

const roleOptions = [
  { label: 'Administrador', value: 'admin' },
  { label: 'Turista', value: 'turista' },
  { label: 'Encargado de Agencia', value: 'encargado_agencia' }
]

const statusOptions = [
  { label: 'Activo', value: 'active' },
  { label: 'Inactivo', value: 'inactive' },
  { label: 'Suspendido', value: 'suspended' }
]

const hasActiveFilters = computed(() => {
  return !!(localFilters.value.search || localFilters.value.rol || localFilters.value.status)
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

const clearSearch = () => {
  localFilters.value.search = ''
  onFilterChange()
}

const clearRol = () => {
  localFilters.value.rol = ''
  onFilterChange()
}

const clearStatus = () => {
  localFilters.value.status = ''
  onFilterChange()
}

const clearFilters = () => {
  localFilters.value = {
    search: '',
    rol: '',
    status: ''
  }
  onFilterChange()
}
</script>
