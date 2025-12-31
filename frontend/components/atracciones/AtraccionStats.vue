<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
    <Card class="border-l-4 border-green-500">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500 mb-1">Total Atracciones</p>
            <p class="text-3xl font-bold text-gray-900">{{ stats.total_atracciones || 0 }}</p>
          </div>
          <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
            <i class="pi pi-map-marker text-green-600 text-2xl"></i>
          </div>
        </div>
      </template>
    </Card>

    <Card class="border-l-4 border-blue-500">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500 mb-1">Activas</p>
            <p class="text-3xl font-bold text-blue-600">{{ stats.atracciones_activas || 0 }}</p>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
            <i class="pi pi-check-circle text-blue-600 text-2xl"></i>
          </div>
        </div>
      </template>
    </Card>

    <Card class="border-l-4 border-orange-500">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500 mb-1">Requieren Agencia</p>
            <p class="text-3xl font-bold text-orange-600">{{ stats.requieren_agencia || 0 }}</p>
          </div>
          <div class="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
            <i class="pi pi-building text-orange-600 text-2xl"></i>
          </div>
        </div>
      </template>
    </Card>

    <Card class="border-l-4 border-purple-500">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500 mb-1">Acceso Particular</p>
            <p class="text-3xl font-bold text-purple-600">{{ stats.acceso_particular || 0 }}</p>
          </div>
          <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
            <i class="pi pi-car text-purple-600 text-2xl"></i>
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  stats: any
}>()
</script>
