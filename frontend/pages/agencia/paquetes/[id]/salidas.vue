<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-7xl mx-auto px-4">
      <!-- Header -->
      <div class="mb-6">
        <Button
          label="Volver a paquetes"
          icon="pi pi-arrow-left"
          text
          @click="navigateTo('/agencia/paquetes')"
        />
        <h1 class="text-3xl font-bold text-gray-900 mt-4">
          Gestión de Salidas - {{ paquete?.nombre }}
        </h1>
        <p class="text-gray-600 mt-2">
          Crea y gestiona las salidas habilitadas para este paquete
        </p>
      </div>

      <!-- Info Card -->
      <Card class="mb-6">
        <template #content>
          <div class="bg-blue-50 border-l-4 border-blue-500 p-4">
            <div class="flex items-start">
              <i class="pi pi-info-circle text-blue-500 text-xl mr-3 mt-1"></i>
              <div>
                <h3 class="font-semibold text-blue-900 mb-2">
                  Nueva Política: Salidas Pre-creadas
                </h3>
                <p class="text-blue-800 text-sm">
                  A partir de ahora, las salidas compartidas SOLO se crean manualmente desde este panel.
                  Los turistas solo podrán comprar en salidas que ya hayas habilitado.
                  Esto evita cancelaciones por falta de cupo mínimo.
                </p>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Formulario para crear salida -->
        <div class="lg:col-span-1">
          <Card>
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-plus-circle text-green-600"></i>
                <span>Crear Nueva Salida</span>
              </div>
            </template>
            <template #content>
              <form @submit.prevent="crearSalida" class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Fecha de Salida
                  </label>
                  <Calendar
                    v-model="form.fechaSalida"
                    dateFormat="yy-mm-dd"
                    :minDate="minDate"
                    showIcon
                    placeholder="Seleccionar fecha"
                    class="w-full"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Cupo Máximo
                  </label>
                  <InputNumber
                    v-model="form.cupoMaximo"
                    :min="1"
                    :max="paquete?.cupo_maximo || 50"
                    showButtons
                    class="w-full"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Cupo Mínimo (opcional)
                  </label>
                  <InputNumber
                    v-model="form.cupoMinimo"
                    :min="1"
                    :max="form.cupoMaximo"
                    showButtons
                    class="w-full"
                    placeholder="Usar cupo del paquete"
                  />
                  <small class="text-gray-500">
                    Por defecto: {{ paquete?.cupo_minimo }} personas
                  </small>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Descripción (opcional)
                  </label>
                  <Textarea
                    v-model="form.descripcion"
                    rows="3"
                    class="w-full"
                    placeholder="Ej: Salida especial Navidad 2025"
                  />
                </div>

                <Button
                  type="submit"
                  label="Crear Salida"
                  icon="pi pi-check"
                  :loading="loading"
                  class="w-full"
                  severity="success"
                />
              </form>
            </template>
          </Card>
        </div>

        <!-- Lista de salidas -->
        <div class="lg:col-span-2">
          <Card>
            <template #title>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <i class="pi pi-calendar text-blue-600"></i>
                  <span>Salidas Habilitadas</span>
                </div>
                <Button
                  icon="pi pi-refresh"
                  text
                  rounded
                  @click="cargarSalidas"
                  :loading="loadingSalidas"
                />
              </div>
            </template>
            <template #content>
              <div v-if="loadingSalidas" class="text-center py-8">
                <ProgressSpinner />
              </div>

              <div v-else-if="salidas.length === 0" class="text-center py-12">
                <i class="pi pi-inbox text-6xl text-gray-300 mb-4"></i>
                <p class="text-gray-500">
                  No hay salidas habilitadas. Crea la primera salida.
                </p>
              </div>

              <div v-else class="space-y-3">
                <Card
                  v-for="salida in salidas"
                  :key="salida.id"
                  class="border border-gray-200"
                >
                  <template #content>
                    <div class="flex items-start justify-between">
                      <div class="flex-1">
                        <div class="flex items-center gap-3 mb-2">
                          <h3 class="text-lg font-semibold text-gray-900">
                            {{ formatFecha(salida.fecha_salida) }}
                          </h3>
                          <Tag
                            :value="salida.estado"
                            :severity="getEstadoSeverity(salida.estado)"
                          />
                          <Tag
                            v-if="salida.creada_manualmente"
                            value="Manual"
                            severity="info"
                            icon="pi pi-user"
                          />
                        </div>

                        <div class="grid grid-cols-2 gap-4 text-sm">
                          <div>
                            <span class="text-gray-600">Cupos:</span>
                            <span class="font-semibold ml-2">
                              {{ salida.cupos_confirmados }} / {{ salida.cupo_maximo }}
                            </span>
                          </div>
                          <div>
                            <span class="text-gray-600">Reservados:</span>
                            <span class="font-semibold ml-2">
                              {{ salida.cupos_reservados }}
                            </span>
                          </div>
                          <div>
                            <span class="text-gray-600">Disponibles:</span>
                            <span class="font-semibold ml-2 text-green-600">
                              {{ salida.cupo_maximo - salida.cupos_confirmados - salida.cupos_reservados }}
                            </span>
                          </div>
                          <div>
                            <span class="text-gray-600">Mínimo:</span>
                            <span class="font-semibold ml-2">
                              {{ salida.cupo_minimo }}
                            </span>
                          </div>
                        </div>

                        <div v-if="salida.descripcion_salida" class="mt-2 text-sm text-gray-600">
                          <i class="pi pi-info-circle mr-1"></i>
                          {{ salida.descripcion_salida }}
                        </div>
                      </div>

                      <div class="flex gap-2 ml-4">
                        <Button
                          icon="pi pi-pencil"
                          text
                          rounded
                          severity="info"
                          @click="editarSalida(salida)"
                        />
                        <Button
                          icon="pi pi-times"
                          text
                          rounded
                          severity="danger"
                          @click="confirmarCancelar(salida)"
                          :disabled="salida.estado === 'cancelada'"
                        />
                      </div>
                    </div>
                  </template>
                </Card>
              </div>
            </template>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const config = useRuntimeConfig()

const paqueteId = computed(() => route.params.id)
const paquete = ref<any>(null)
const salidas = ref<any[]>([])
const loading = ref(false)
const loadingSalidas = ref(false)

const minDate = ref(new Date(Date.now() + 86400000)) // Mañana

const form = ref({
  fechaSalida: null as Date | null,
  cupoMaximo: null as number | null,
  cupoMinimo: null as number | null,
  descripcion: ''
})

onMounted(() => {
  cargarPaquete()
  cargarSalidas()
})

const cargarPaquete = async () => {
  try {
    const response = await $fetch(`${config.public.apiBase}/agencias/${authStore.user?.agencia_id}/paquetes/${paqueteId.value}`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    if (response.success) {
      paquete.value = response.data
      form.value.cupoMaximo = paquete.value.cupo_maximo
    }
  } catch (error) {
    console.error('Error cargando paquete:', error)
  }
}

const cargarSalidas = async () => {
  loadingSalidas.value = true
  try {
    const response = await $fetch<any>(
      `${config.public.apiBase}/agencias/paquetes/${paqueteId.value}/salidas-manuales`,
      {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      }
    )
    if (response.success) {
      salidas.value = response.data
    }
  } catch (error) {
    console.error('Error cargando salidas:', error)
  } finally {
    loadingSalidas.value = false
  }
}

const crearSalida = async () => {
  if (!form.value.fechaSalida || !form.value.cupoMaximo) {
    alert('Por favor complete los campos requeridos')
    return
  }

  loading.value = true
  try {
    const payload = {
      fecha_salida: formatDateToSQL(form.value.fechaSalida),
      cupo_maximo: form.value.cupoMaximo,
      cupo_minimo: form.value.cupoMinimo || 0,
      descripcion: form.value.descripcion || null
    }

    const response = await $fetch(
      `${config.public.apiBase}/agencias/paquetes/${paqueteId.value}/salidas-manuales`,
      {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${authStore.token}`
        },
        body: payload
      }
    )

    if (response.success) {
      alert('Salida creada exitosamente')
      // Reset form
      form.value = {
        fechaSalida: null,
        cupoMaximo: paquete.value.cupo_maximo,
        cupoMinimo: null,
        descripcion: ''
      }
      cargarSalidas()
    }
  } catch (error: any) {
    alert(error.data?.message || 'Error al crear salida')
  } finally {
    loading.value = false
  }
}

const confirmarCancelar = (salida: any) => {
  const razon = prompt('¿Por qué deseas cancelar esta salida?')
  if (!razon) return

  cancelarSalida(salida.id, razon)
}

const cancelarSalida = async (salidaId: number, razon: string) => {
  try {
    const response = await $fetch(
      `${config.public.apiBase}/agencias/salidas/${salidaId}/cancelar`,
      {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${authStore.token}`
        },
        body: { razon }
      }
    )

    if (response.success) {
      alert('Salida cancelada exitosamente')
      cargarSalidas()
    }
  } catch (error: any) {
    alert(error.data?.message || 'Error al cancelar salida')
  }
}

const formatFecha = (fecha: string) => {
  const date = new Date(fecha)
  return date.toLocaleDateString('es-ES', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const formatDateToSQL = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const getEstadoSeverity = (estado: string) => {
  const map: Record<string, any> = {
    pendiente: 'warning',
    activa: 'success',
    completada: 'info',
    cancelada: 'danger'
  }
  return map[estado] || 'secondary'
}

const editarSalida = (salida: any) => {
  // TODO: Implementar edición
  console.log('Editar salida:', salida)
}
</script>
