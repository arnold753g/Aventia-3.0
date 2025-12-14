<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
    <Card class="bg-gradient-to-br from-purple-50 to-purple-100">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-purple-600 font-medium">Total Usuarios</p>
            <p class="text-3xl font-bold text-purple-700">{{ stats.total || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-purple-200 flex items-center justify-center">
            <i class="pi pi-users text-2xl text-purple-700"></i>
          </div>
        </div>
      </template>
    </Card>

    <Card class="bg-gradient-to-br from-green-50 to-green-100">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-green-600 font-medium">Activos</p>
            <p class="text-3xl font-bold text-green-700">{{ stats.active || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-green-200 flex items-center justify-center">
            <i class="pi pi-check-circle text-2xl text-green-700"></i>
          </div>
        </div>
      </template>
    </Card>

    <Card class="bg-gradient-to-br from-blue-50 to-blue-100">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-blue-600 font-medium">Turistas</p>
            <p class="text-3xl font-bold text-blue-700">{{ stats.turistas || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-blue-200 flex items-center justify-center">
            <i class="pi pi-user text-2xl text-blue-700"></i>
          </div>
        </div>
      </template>
    </Card>

    <Card class="bg-gradient-to-br from-orange-50 to-orange-100">
      <template #content>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-orange-600 font-medium">Administradores</p>
            <p class="text-3xl font-bold text-orange-700">{{ stats.admins || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-orange-200 flex items-center justify-center">
            <i class="pi pi-shield text-2xl text-orange-700"></i>
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  stats: Record<string, any>
}>()
</script>
