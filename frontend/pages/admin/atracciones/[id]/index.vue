<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center min-h-screen">
      <ProgressSpinner />
    </div>

    <!-- Content -->
    <div v-else-if="atraccion">
      <!-- Header -->
      <div class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 py-6">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
              <Button
                icon="pi pi-arrow-left"
                text
                rounded
                @click="navigateTo('/admin/atracciones')"
              />
              <div>
                <h1 class="text-3xl font-bold text-gray-900">{{ atraccion.nombre }}</h1>
                <div class="flex items-center gap-2 mt-1">
                  <i class="pi pi-map-marker text-gray-500"></i>
                  <span class="text-gray-600">
                    {{ atraccion.provincia?.nombre }}, {{ atraccion.provincia?.departamento?.nombre }}
                  </span>
                </div>
              </div>
            </div>
            <div class="flex gap-2">
              <Button
                label="Editar"
                icon="pi pi-pencil"
                @click="navigateTo(`/admin/atracciones/${atraccion.id}/editar`)"
              />
              <Button
                label="Desactivar"
                icon="pi pi-times"
                severity="danger"
                @click="confirmDelete"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="max-w-7xl mx-auto px-4 py-8">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Columna Principal -->
          <div class="lg:col-span-2 space-y-6">
            <!-- Galería de Fotos -->
            <Card>
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-images text-red-600"></i>
                  <span>Galería de Fotos</span>
                </div>
              </template>
              <template #content>
                <AtraccionesFotosGallery
                  :fotos="atraccion.fotos || []"
                  :editable="false"
                />
              </template>
            </Card>

            <!-- Descripción -->
            <Card>
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-align-left text-blue-600"></i>
                  <span>Descripción</span>
                </div>
              </template>
              <template #content>
                <p class="text-gray-700 whitespace-pre-line">{{ atraccion.descripcion }}</p>
              </template>
            </Card>

            <!-- Mapa -->
            <Card v-if="atraccion.latitud && atraccion.longitud">
              <template #title>
                <div class="flex items-center gap-2">
                  <i class="pi pi-map text-green-600"></i>
                  <span>Ubicación</span>
                </div>
              </template>
              <template #content>
                <AtraccionesAtraccionMap
                  :latitud="atraccion.latitud"
                  :longitud="atraccion.longitud"
                  :editable="false"
                  height="400px"
                />
              </template>
            </Card>
          </div>

          <!-- Columna Lateral -->
          <div class="space-y-6">
            <!-- Estado y Características -->
            <Card>
              <template #title>Información General</template>
              <template #content>
                <div class="space-y-4">
                  <div>
                    <label class="text-sm text-gray-500">Estado</label>
                    <div class="mt-1">
                      <span
                        :class="[
                          'px-3 py-1 rounded-full text-sm font-semibold',
                          `bg-${getStatusAtraccionColor(atraccion.status)}-100`,
                          `text-${getStatusAtraccionColor(atraccion.status)}-700`
                        ]"
                      >
                        {{ getStatusAtraccionLabel(atraccion.status) }}
                      </span>
                    </div>
                  </div>

                  <Divider />

                  <div>
                    <label class="text-sm text-gray-500">Precio de Entrada</label>
                    <p class="text-2xl font-bold text-blue-600 mt-1">
                      {{ formatPrecioBoliviano(atraccion.precio_entrada) }}
                    </p>
                  </div>

                  <div v-if="atraccion.nivel_dificultad">
                    <label class="text-sm text-gray-500">Nivel de Dificultad</label>
                    <div class="mt-1">
                      <span
                        :class="[
                          'px-3 py-1 rounded-full text-sm font-semibold',
                          `bg-${getNivelDificultadColor(atraccion.nivel_dificultad)}-100`,
                          `text-${getNivelDificultadColor(atraccion.nivel_dificultad)}-700`
                        ]"
                      >
                        {{ getNivelDificultadLabel(atraccion.nivel_dificultad) }}
                      </span>
                    </div>
                  </div>

                  <Divider />

                  <div class="space-y-2">
                    <div class="flex items-center gap-2">
                      <i :class="[
                        'pi',
                        atraccion.requiere_agencia ? 'pi-check-circle text-green-600' : 'pi-times-circle text-gray-400'
                      ]"></i>
                      <span class="text-sm">Requiere Agencia</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <i :class="[
                        'pi',
                        atraccion.acceso_particular ? 'pi-check-circle text-green-600' : 'pi-times-circle text-gray-400'
                      ]"></i>
                      <span class="text-sm">Acceso Particular</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <i :class="[
                        'pi',
                        atraccion.visible_publico ? 'pi-check-circle text-green-600' : 'pi-times-circle text-gray-400'
                      ]"></i>
                      <span class="text-sm">Visible al Público</span>
                    </div>
                  </div>
                </div>
              </template>
            </Card>

            <!-- Horarios -->
            <Card v-if="atraccion.horario_apertura || atraccion.horario_cierre">
              <template #title>Horarios</template>
              <template #content>
                <div class="space-y-3">
                  <div class="flex items-center gap-2">
                    <i class="pi pi-clock text-purple-600"></i>
                    <span>{{ formatHorario(atraccion.horario_apertura, atraccion.horario_cierre) }}</span>
                  </div>

                  <div v-if="atraccion.dias && atraccion.dias.length > 0">
                    <label class="text-sm text-gray-500 mb-2 block">Días de Apertura</label>
                    <div class="flex flex-wrap gap-2">
                      <Chip
                        v-for="dia in atraccion.dias"
                        :key="dia.id"
                        :label="dia.nombre"
                      />
                    </div>
                  </div>
                </div>
              </template>
            </Card>

            <!-- Mejor Época -->
            <Card v-if="atraccion.mes_inicio || atraccion.mes_fin">
              <template #title>Mejor Época de Visita</template>
              <template #content>
                <div class="flex items-center gap-2">
                  <i class="pi pi-calendar text-teal-600"></i>
                  <span>{{ formatMejorEpoca(atraccion.mes_inicio, atraccion.mes_fin) }}</span>
                </div>
              </template>
            </Card>

            <!-- Subcategorías -->
            <Card v-if="atraccion.subcategorias && atraccion.subcategorias.length > 0">
              <template #title>Categorías</template>
              <template #content>
                <div class="space-y-2">
                  <div
                    v-for="subcat in atraccion.subcategorias"
                    :key="subcat.id"
                    class="flex items-center justify-between p-2 bg-gray-50 rounded"
                  >
                    <span class="text-sm">{{ subcat.subcategoria?.nombre }}</span>
                    <span v-if="subcat.es_principal" class="px-2 py-1 bg-blue-500 text-white rounded text-xs">
                      Principal
                    </span>
                  </div>
                </div>
              </template>
            </Card>

            <!-- Contacto -->
            <Card>
              <template #title>Información de Contacto</template>
              <template #content>
                <div class="space-y-3">
                  <div v-if="atraccion.telefono" class="flex items-center gap-2">
                    <i class="pi pi-phone text-gray-600"></i>
                    <a :href="`tel:${atraccion.telefono}`" class="text-blue-600 hover:underline">
                      {{ atraccion.telefono }}
                    </a>
                  </div>

                  <div v-if="atraccion.email" class="flex items-center gap-2">
                    <i class="pi pi-envelope text-gray-600"></i>
                    <a :href="`mailto:${atraccion.email}`" class="text-blue-600 hover:underline">
                      {{ atraccion.email }}
                    </a>
                  </div>

                  <div v-if="atraccion.sitio_web" class="flex items-center gap-2">
                    <i class="pi pi-globe text-gray-600"></i>
                    <a :href="atraccion.sitio_web" target="_blank" class="text-blue-600 hover:underline">
                      Sitio Web
                    </a>
                  </div>

                  <div v-if="atraccion.facebook" class="flex items-center gap-2">
                    <i class="pi pi-facebook text-gray-600"></i>
                    <span>{{ atraccion.facebook }}</span>
                  </div>

                  <div v-if="atraccion.instagram" class="flex items-center gap-2">
                    <i class="pi pi-instagram text-gray-600"></i>
                    <span>{{ atraccion.instagram }}</span>
                  </div>

                  <div v-if="!atraccion.telefono && !atraccion.email && !atraccion.sitio_web && !atraccion.facebook && !atraccion.instagram">
                    <p class="text-sm text-gray-500">No hay información de contacto</p>
                  </div>
                </div>
              </template>
            </Card>

            <!-- Auditoría -->
            <Card>
              <template #title>Información del Sistema</template>
              <template #content>
                <div class="space-y-2 text-sm text-gray-600">
                  <div>
                    <span class="font-medium">ID:</span> {{ atraccion.id }}
                  </div>
                  <div>
                    <span class="font-medium">Creado:</span> {{ formatDateTime(atraccion.created_at) }}
                  </div>
                  <div>
                    <span class="font-medium">Actualizado:</span> {{ formatDateTime(atraccion.updated_at) }}
                  </div>
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>
    </div>

    <!-- Dialog confirmar eliminación -->
    <Dialog
      v-model:visible="showDeleteDialog"
      header="Confirmar Desactivación"
      :modal="true"
      :style="{ width: '450px' }"
    >
      <div class="flex items-start gap-4">
        <i class="pi pi-exclamation-triangle text-orange-500 text-3xl"></i>
        <div>
          <p class="mb-2">¿Está seguro de desactivar esta atracción?</p>
          <p class="text-sm text-gray-500">
            La atracción quedará inactiva y no será visible al público.
          </p>
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancelar"
          severity="secondary"
          @click="showDeleteDialog = false"
        />
        <Button
          label="Desactivar"
          severity="danger"
          @click="handleDelete"
          :loading="deleting"
        />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useRoute } from 'vue-router'
import {
  getNivelDificultadLabel,
  getNivelDificultadColor,
  getStatusAtraccionLabel,
  getStatusAtraccionColor,
  formatPrecioBoliviano,
  formatHorario,
  formatMejorEpoca
} from '~/utils/formatters-atraccion'
import { formatDateTime } from '~/utils/formatters'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const route = useRoute()
const toast = useToast()
const { getAtraccion, deleteAtraccion } = useAtracciones()

const loading = ref(true)
const deleting = ref(false)
const atraccion = ref<any>(null)
const showDeleteDialog = ref(false)

const loadAtraccion = async () => {
  loading.value = true
  try {
    const id = parseInt(route.params.id as string)
    const response: any = await getAtraccion(id)

    console.log('📡 Respuesta completa:', response)
    console.log('📦 Datos de atracción:', response.data)
    console.log('🖼️ Fotos recibidas:', response.data?.fotos)
    console.log('🔢 Cantidad de fotos:', response.data?.fotos?.length)

    if (response.success) {
      atraccion.value = response.data
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Error al cargar atracción',
      life: 3000
    })
    navigateTo('/admin/atracciones')
  } finally {
    loading.value = false
  }
}

const confirmDelete = () => {
  showDeleteDialog.value = true
}

const handleDelete = async () => {
  if (!atraccion.value) return

  deleting.value = true
  try {
    const response: any = await deleteAtraccion(atraccion.value.id)

    if (response.success) {
      toast.add({
        severity: 'success',
        summary: 'Atracción Desactivada',
        detail: 'La atracción ha sido desactivada exitosamente',
        life: 3000
      })

      setTimeout(() => {
        navigateTo('/admin/atracciones')
      }, 1500)
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'Error al desactivar atracción',
      life: 3000
    })
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadAtraccion()
})
</script>
