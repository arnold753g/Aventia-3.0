<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="flex items-center justify-between mb-6">
      <div>
        <p class="text-sm text-gray-500">Detalle de agencia</p>
        <h1 class="text-3xl font-bold text-gray-900">{{ agencia?.nombre_comercial || 'Agencia' }}</h1>
      </div>
      <div class="flex gap-2">
        <Button label="Editar" icon="pi pi-pencil" severity="warning" @click="navigateTo(`/admin/agencias/${id}/editar`)" />
        <Button label="Volver" icon="pi pi-arrow-left" outlined @click="navigateTo('/admin/agencias')" />
      </div>
    </div>

    <Card class="mb-4">
      <template #content>
        <div class="flex flex-wrap gap-2">
          <Button
            v-for="tab in tabs"
            :key="tab.value"
            :label="tab.label"
            :icon="tab.icon"
            :severity="activeTab === tab.value ? 'primary' : 'secondary'"
            outlined
            @click="selectTab(tab.value)"
          />
        </div>
      </template>
    </Card>

    <Card v-if="activeTab === 'general'">
      <template #title>Informacion general</template>
      <template #content>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div class="space-y-2">
            <p><strong>Nombre:</strong> {{ agencia?.nombre_comercial }}</p>
            <p><strong>Descripcion:</strong> {{ agencia?.descripcion || 'N/D' }}</p>
            <p><strong>Direccion:</strong> {{ agencia?.direccion }}</p>
            <p><strong>Departamento:</strong> {{ agencia?.departamento?.nombre || 'N/D' }}</p>
            <p><strong>Dias:</strong> {{ agencia?.dias?.map((d: any) => d.nombre).join(', ') || 'N/D' }}</p>
          </div>
          <div class="space-y-2">
            <p><strong>Telefono:</strong> {{ agencia?.telefono }}</p>
            <p><strong>Email:</strong> {{ agencia?.email }}</p>
            <p><strong>Sitio web:</strong> {{ agencia?.sitio_web || 'N/D' }}</p>
            <p><strong>Licencia:</strong> {{ agencia?.licencia_turistica ? 'Si' : 'No' }}</p>
            <p><strong>Estado:</strong> {{ agencia?.status }}</p>
          </div>
        </div>
      </template>
    </Card>

    <div v-if="activeTab === 'general'" class="mt-6">
      <div class="flex items-center gap-2 mb-4">
        <i class="pi pi-eye text-blue-600"></i>
        <h2 class="text-xl font-semibold text-gray-900">Vista publica</h2>
      </div>
      <AgenciasAgenciaPublicView v-if="agencia" :agencia="agencia" :agenciaId="id" />
      <Skeleton v-else height="320px" />
    </div>

    <Card v-if="activeTab === 'fotos'" class="mb-4">
      <template #title>Fotos</template>
      <template #content>
        <div class="space-y-4">
          <div class="flex flex-wrap items-center gap-3">
            <input
              ref="fileInput"
              type="file"
              accept="image/jpeg,image/png,image/webp"
              multiple
              class="hidden"
              :disabled="maxFotosReached || uploading"
              @change="handleFileSelect"
            />
            <Button
              label="Seleccionar Fotos"
              icon="pi pi-upload"
              outlined
              :disabled="maxFotosReached || uploading"
              @click="fileInput?.click()"
            />
            <Button
              :label="uploadButtonLabel"
              icon="pi pi-cloud-upload"
              :loading="uploading"
              :disabled="selectedFotos.length === 0 || uploading"
              @click="uploadSelectedFotos"
            />
            <span class="text-sm text-gray-600">
              {{ totalFotos }} / 10 en total ({{ selectedFotos.length }} nuevas)
            </span>
            <small class="text-gray-500">Max 5MB c/u. Max 10 fotos.</small>
          </div>

          <div class="flex items-center gap-2">
            <Checkbox v-model="setPrimeraComoPrincipal" binary inputId="setPrimeraComoPrincipal" :disabled="uploading" />
            <label for="setPrimeraComoPrincipal" class="text-sm text-gray-700">
              Marcar la primera como principal
            </label>
          </div>

          <!-- Previews de fotos seleccionadas -->
          <div v-if="selectedFotos.length > 0">
            <p class="text-sm text-gray-600 mb-3">
              <i class="pi pi-info-circle"></i>
              Usa las flechas para reordenar.
              <span v-if="setPrimeraComoPrincipal">La primera foto será la principal.</span>
              <span v-else>La foto principal no se modificará.</span>
            </p>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div
                v-for="(foto, index) in selectedFotos"
                :key="index"
                class="relative group border-2 rounded-lg hover:border-blue-300 transition-all"
              >
                <img
                  :src="foto.preview"
                  class="w-full h-32 object-cover rounded-t"
                  :alt="foto.file.name"
                />
                <!-- Badge Principal y Número -->
                <div class="absolute top-2 left-2 flex flex-col gap-1">
                  <span class="px-2 py-1 bg-gray-900 bg-opacity-70 text-white rounded text-xs font-mono">
                    {{ index + 1 }}
                  </span>
                  <span v-if="setPrimeraComoPrincipal && index === 0" class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
                    PRINCIPAL
                  </span>
                </div>
                <!-- Botones de Control -->
                <div class="absolute top-2 right-2 flex flex-col gap-1">
                  <Button
                    icon="pi pi-chevron-up"
                    severity="secondary"
                    rounded
                    size="small"
                    :disabled="uploading || index === 0"
                    class="bg-white"
                    @click="moveSelectedFotoUp(index)"
                  />
                  <Button
                    icon="pi pi-chevron-down"
                    severity="secondary"
                    rounded
                    size="small"
                    :disabled="uploading || index === selectedFotos.length - 1"
                    class="bg-white"
                    @click="moveSelectedFotoDown(index)"
                  />
                  <Button
                    icon="pi pi-times"
                    severity="danger"
                    rounded
                    size="small"
                    :disabled="uploading"
                    class="bg-white"
                    @click="removeSelectedFoto(index)"
                  />
                </div>
                <!-- Nombre del archivo -->
                <div class="bg-gray-50 p-2 rounded-b">
                  <p class="text-xs text-gray-600 truncate" :title="foto.file.name">
                    {{ foto.file.name }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div
            v-if="fotosOrdenadas.length > 0"
            :class="selectedFotos.length > 0 ? 'border-t border-gray-200 pt-4' : ''"
          >
            <h3 class="text-lg font-semibold text-gray-900 mb-3">
              Fotos actuales ({{ fotosOrdenadas.length }}/10)
            </h3>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div
                v-for="(foto, index) in fotosOrdenadas"
                :key="foto.id"
                class="relative group border-2 rounded-lg hover:border-blue-300 transition-all"
              >
                <img
                  :src="resolveFotoUrl(foto.foto_url || foto.foto)"
                  class="w-full h-32 object-cover rounded-t"
                  :alt="foto.titulo || `Foto ${index + 1}`"
                />
                <!-- Badge Principal y Número -->
                <div class="absolute top-2 left-2 flex flex-col gap-1">
                  <span class="px-2 py-1 bg-gray-900 bg-opacity-70 text-white rounded text-xs font-mono">
                    {{ index + 1 }}
                  </span>
                  <span v-if="foto.es_principal" class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
                    PRINCIPAL
                  </span>
                </div>
                <!-- Botón eliminar -->
                <div class="absolute top-2 right-2 flex flex-col gap-1">
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    rounded
                    size="small"
                    :disabled="uploading"
                    class="bg-white"
                    @click="removeFoto(foto)"
                  />
                </div>
                <!-- Título -->
                <div class="bg-gray-50 p-2 rounded-b">
                  <p class="text-xs text-gray-600 truncate" :title="foto.titulo || foto.foto_url || foto.foto">
                    {{ foto.titulo || 'Foto' }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty state -->
          <div v-if="selectedFotos.length === 0 && fotosOrdenadas.length === 0" class="text-center py-8 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
            <i class="pi pi-images text-4xl text-gray-400 mb-2"></i>
            <p class="text-gray-600 text-sm">Haga clic en "Seleccionar Fotos" para agregar imágenes</p>
          </div>
        </div>
      </template>
    </Card>

    <Card v-if="activeTab === 'especialidades'" class="mb-4">
      <template #title>Especialidades</template>
      <template #content>
        <div class="flex gap-3 mb-4">
          <Select
            v-model="especialidadForm.categoria_id"
            :options="categoriasDisponibles"
            optionLabel="nombre"
            optionValue="id"
            placeholder="Seleccionar categoria"
            class="w-64"
            showClear
          />
          <Checkbox v-model="especialidadForm.es_principal" binary inputId="principal" />
          <label for="principal" class="text-sm text-gray-700">Principal</label>
          <Button label="Agregar" icon="pi pi-plus" :loading="addingEspecialidad" @click="addEspecialidadAction" />
        </div>

        <div class="flex flex-wrap gap-2">
          <Tag
            v-for="esp in agencia?.especialidades || []"
            :key="esp.id"
            :severity="esp.es_principal ? 'success' : 'info'"
          >
            <template #icon>
              <i class="pi pi-star" v-if="esp.es_principal"></i>
            </template>
            <span class="inline-flex items-center gap-2">
              <span class="p-tag-label text-sm">{{ getEspecialidadNombre(esp) }}</span>
              <Button icon="pi pi-times" text rounded severity="danger" @click="removeEspecialidadAction(esp)" />
            </span>
          </Tag>
        </div>
      </template>
    </Card>

    <Card v-if="activeTab === 'pagos'" class="mb-4">
      <template #title>Metodos de pago</template>
      <template #content>
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-2">
            <Checkbox v-model="pagoForm.acepta_qr" binary inputId="qr" />
            <label for="qr">QR</label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox v-model="pagoForm.acepta_transferencia" binary inputId="trans" />
            <label for="trans">Transferencia</label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox v-model="pagoForm.acepta_efectivo" binary inputId="cash" />
            <label for="cash">Efectivo</label>
          </div>
          <div class="flex justify-end">
            <Button label="Guardar cambios" icon="pi pi-save" :loading="savingPagos" @click="savePagos" />
          </div>
        </div>
      </template>
    </Card>

    <Card v-if="activeTab === 'capacidad'" class="mb-4">
      <template #title>Capacidad operativa</template>
      <template #content>
        <div v-if="capacidadLoading" class="space-y-3">
          <Skeleton height="200px" />
        </div>

        <div v-else class="space-y-5">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Maximo salidas por dia</label>
              <InputNumber v-model="capacidadForm.max_salidas_por_dia" :min="0" :max="1000" class="w-full" />
              <p class="text-xs text-gray-500 mt-1">Maximo total de salidas que la agencia puede manejar en un dia.</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Maximo salidas por horario</label>
              <InputNumber v-model="capacidadForm.max_salidas_por_horario" :min="0" :max="1000" class="w-full" />
              <p class="text-xs text-gray-500 mt-1">Limite simultaneo por horario (manana, tarde, todo_dia).</p>
            </div>
          </div>

          <Message
            v-if="(capacidadForm.max_salidas_por_horario ?? 0) > (capacidadForm.max_salidas_por_dia ?? 0)"
            severity="warn"
            :closable="false"
          >
            El maximo por horario no puede ser mayor al maximo por dia.
          </Message>

          <div class="flex justify-end">
            <Button label="Guardar cambios" icon="pi pi-save" :loading="savingCapacidad" @click="saveCapacidad" />
          </div>
        </div>
      </template>
    </Card>

    <Card v-if="activeTab === 'encargado'">
      <template #title>Encargado principal</template>
      <template #content>
        <div v-if="agencia?.encargado_principal" class="space-y-2 text-sm">
          <p><strong>Nombre:</strong> {{ agencia.encargado_principal.nombre }} {{ agencia.encargado_principal.apellido_paterno }}</p>
          <p><strong>Email:</strong> {{ agencia.encargado_principal.email }}</p>
          <p><strong>Telefono:</strong> {{ agencia.encargado_principal.phone || 'N/D' }}</p>
          <p><strong>CI:</strong> {{ agencia.encargado_principal.ci || 'N/D' }}</p>
        </div>
        <div v-else class="text-gray-500">Sin encargado asignado.</div>
      </template>
    </Card>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const authStore = useAuthStore()
if (!authStore.isAdmin && !authStore.isEncargado) {
  navigateTo('/dashboard')
}

const route = useRoute()
const id = Number(route.params.id)
const toast = useToast()
const {
  getAgencia,
  uploadFoto: uploadFotoApi,
  removeFoto: removeFotoApi,
  addEspecialidad,
  removeEspecialidad,
  updateAgencia,
  getCategorias,
  getAgenciaCapacidad,
  updateAgenciaCapacidad
} = useAgencias()

const agencia = ref<any>(null)
const categorias = ref<any[]>([])
const uploading = ref(false)
const uploadProgress = ref<{ current: number; total: number } | null>(null)
const addingEspecialidad = ref(false)
const savingPagos = ref(false)
const capacidadLoading = ref(false)
const savingCapacidad = ref(false)
const activeTab = ref<'general' | 'fotos' | 'especialidades' | 'pagos' | 'capacidad' | 'encargado'>('general')

const especialidadForm = ref({
  categoria_id: null as number | null,
  es_principal: false
})

const pagoForm = ref({
  acepta_qr: false,
  acepta_transferencia: false,
  acepta_efectivo: false
})

const capacidadForm = ref({
  max_salidas_por_dia: 5 as number | null,
  max_salidas_por_horario: 3 as number | null
})

const tabs = [
  { label: 'General', value: 'general', icon: 'pi pi-info-circle' },
  { label: 'Fotos', value: 'fotos', icon: 'pi pi-image' },
  { label: 'Especialidades', value: 'especialidades', icon: 'pi pi-tag' },
  { label: 'Pagos', value: 'pagos', icon: 'pi pi-wallet' },
  { label: 'Capacidad', value: 'capacidad', icon: 'pi pi-sliders-h' },
  { label: 'Encargado', value: 'encargado', icon: 'pi pi-user' }
]

const selectTab = (value: 'general' | 'fotos' | 'especialidades' | 'pagos' | 'capacidad' | 'encargado') => {
  activeTab.value = value
  if (value === 'capacidad') {
    void loadCapacidad()
  }
}

const categoriasById = computed(() => {
  const map = new Map<number, any>()
  for (const c of categorias.value || []) {
    const id = Number(c?.id)
    if (Number.isFinite(id)) map.set(id, c)
  }
  return map
})

const getEspecialidadNombre = (esp: any) => {
  const id = Number(esp?.categoria_id)
  return esp?.categoria?.nombre || categoriasById.value.get(id)?.nombre || (Number.isFinite(id) ? `ID ${id}` : 'Especialidad')
}

const categoriasDisponibles = computed(() => {
  const usadas = new Set<number>(
    (agencia.value?.especialidades || [])
      .map((e: any) => Number(e?.categoria_id))
      .filter((id: number) => Number.isFinite(id))
  )
  return (categorias.value || []).filter((c) => !usadas.has(Number(c?.id)))
})

const loadAgencia = async () => {
  try {
    const response: any = await getAgencia(id)
    if (response.success) {
      agencia.value = response.data
      pagoForm.value = {
        acepta_qr: response.data.acepta_qr,
        acepta_transferencia: response.data.acepta_transferencia,
        acepta_efectivo: response.data.acepta_efectivo
      }
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo cargar la agencia', life: 3000 })
  }
}

const loadCategorias = async () => {
  try {
    const response: any = await getCategorias()
    categorias.value = response?.data || []
  } catch (error) {
    categorias.value = []
  }
}

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFotos = ref<Array<{ file: File; preview: string }>>([])
const setPrimeraComoPrincipal = ref(true)
const maxFotos = 10

const fotosOrdenadas = computed(() => {
  const fotos = (agencia.value?.fotos || []) as any[]
  return [...fotos].sort((a, b) => {
    const aPrincipal = a?.es_principal ? 1 : 0
    const bPrincipal = b?.es_principal ? 1 : 0
    if (aPrincipal !== bPrincipal) return bPrincipal - aPrincipal

    const aOrden = Number(a?.orden ?? 0)
    const bOrden = Number(b?.orden ?? 0)
    if (aOrden !== bOrden) return aOrden - bOrden

    const aCreated = a?.created_at ? new Date(a.created_at).getTime() : 0
    const bCreated = b?.created_at ? new Date(b.created_at).getTime() : 0
    if (aCreated !== bCreated) return aCreated - bCreated

    return Number(a?.id ?? 0) - Number(b?.id ?? 0)
  })
})

const totalFotos = computed(() => fotosOrdenadas.value.length + selectedFotos.value.length)
const maxFotosReached = computed(() => totalFotos.value >= maxFotos)

const uploadButtonLabel = computed(() => {
  if (!uploading.value) return 'Subir seleccionadas'
  if (!uploadProgress.value) return 'Subiendo...'
  return `Subiendo ${uploadProgress.value.current}/${uploadProgress.value.total}`
})

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const files = input.files
  if (!files) return

  const remainingSlots = maxFotos - fotosOrdenadas.value.length - selectedFotos.value.length
  if (remainingSlots <= 0) {
    toast.add({ severity: 'warn', summary: 'Límite', detail: 'Máximo 10 fotos', life: 2500 })
    input.value = ''
    return
  }

  const filesToAdd = Array.from(files).slice(0, remainingSlots)

  filesToAdd.forEach((file) => {
    if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
      toast.add({
        severity: 'warn',
        summary: 'Formato no permitido',
        detail: `${file.name} no es un formato válido`,
        life: 3000
      })
      return
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.add({
        severity: 'warn',
        summary: 'Archivo muy grande',
        detail: `${file.name} supera 5MB`,
        life: 3000
      })
      return
    }

    const reader = new FileReader()
    reader.onload = (e) => {
      selectedFotos.value.push({
        file,
        preview: e.target?.result as string
      })
    }
    reader.readAsDataURL(file)
  })

  input.value = ''
}

const removeSelectedFoto = (index: number) => {
  selectedFotos.value.splice(index, 1)
}

const moveSelectedFotoUp = (index: number) => {
  if (index === 0) return
  const temp = selectedFotos.value[index]
  selectedFotos.value[index] = selectedFotos.value[index - 1]
  selectedFotos.value[index - 1] = temp
}

const moveSelectedFotoDown = (index: number) => {
  if (index === selectedFotos.value.length - 1) return
  const temp = selectedFotos.value[index]
  selectedFotos.value[index] = selectedFotos.value[index + 1]
  selectedFotos.value[index + 1] = temp
}

const getNextOrdenBase = () => {
  const fotos = (agencia.value?.fotos || []) as any[]
  const maxOrden = fotos.reduce((max, f) => {
    const orden = Number(f?.orden)
    if (!Number.isFinite(orden)) return max
    return Math.max(max, orden)
  }, -1)
  return maxOrden + 1
}

const uploadSelectedFotos = async () => {
  if (selectedFotos.value.length === 0) {
    toast.add({ severity: 'warn', summary: 'Fotos', detail: 'Seleccione al menos una imagen', life: 2500 })
    return
  }

  if (totalFotos.value > maxFotos) {
    toast.add({ severity: 'warn', summary: 'Límite', detail: 'Máximo 10 fotos', life: 2500 })
    return
  }

  uploading.value = true
  uploadProgress.value = { current: 0, total: selectedFotos.value.length }

  const baseOrden = getNextOrdenBase()
  const total = selectedFotos.value.length
  const failed: Array<{ file: File; preview: string }> = []
  let okCount = 0

  for (let i = 0; i < total; i++) {
    const item = selectedFotos.value[i]
    const formData = new FormData()
    formData.append('foto', item.file)
    formData.append('titulo', item.file.name)
    formData.append('orden', String(baseOrden + i))
    formData.append('es_principal', String(!!setPrimeraComoPrincipal.value && i === 0))

    try {
      const response: any = await uploadFotoApi(id, formData)
      if (response.success) {
        okCount++
      } else {
        failed.push(item)
      }
    } catch (error: any) {
      failed.push(item)
    } finally {
      uploadProgress.value = { current: i + 1, total }
    }
  }

  selectedFotos.value = failed
  if (fileInput.value) fileInput.value.value = ''

  try {
    await loadAgencia()
  } catch {
    // ignore
  }

  if (okCount > 0 && failed.length === 0) {
    toast.add({ severity: 'success', summary: 'Subida', detail: `Se subieron ${okCount} fotos`, life: 2500 })
  } else if (okCount > 0) {
    toast.add({ severity: 'warn', summary: 'Subida parcial', detail: `Subidas: ${okCount}. Fallidas: ${failed.length}`, life: 3000 })
  } else {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo subir ninguna foto', life: 3000 })
  }

  uploading.value = false
  uploadProgress.value = null
}

const removeFoto = async (foto: any) => {
  try {
    const response: any = await removeFotoApi(id, foto.id)
    if (response.success) {
      agencia.value.fotos = (agencia.value.fotos || []).filter((f: any) => f.id !== foto.id)
      toast.add({ severity: 'success', summary: 'Eliminada', detail: 'Foto eliminada', life: 2500 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo eliminar', life: 3000 })
  }
}

const addEspecialidadAction = async () => {
  if (!especialidadForm.value.categoria_id) {
    toast.add({ severity: 'warn', summary: 'Seleccione categoria', detail: '', life: 2500 })
    return
  }
  addingEspecialidad.value = true
  try {
    const response: any = await addEspecialidad(id, {
      categoria_id: especialidadForm.value.categoria_id,
      es_principal: especialidadForm.value.es_principal
    })
    if (response.success) {
      agencia.value.especialidades = [...(agencia.value.especialidades || []), response.data]
      especialidadForm.value = { categoria_id: null, es_principal: false }
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo agregar', life: 3000 })
  } finally {
    addingEspecialidad.value = false
  }
}

const removeEspecialidadAction = async (esp: any) => {
  try {
    const response: any = await removeEspecialidad(id, esp.id)
    if (response.success) {
      agencia.value.especialidades = (agencia.value.especialidades || []).filter((e: any) => e.id !== esp.id)
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo eliminar', life: 3000 })
  }
}

const savePagos = async () => {
  savingPagos.value = true
  try {
    const response: any = await updateAgencia(id, { ...pagoForm.value })
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Metodos de pago actualizados', life: 2500 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo guardar', life: 3000 })
  } finally {
    savingPagos.value = false
  }
}

const loadCapacidad = async () => {
  capacidadLoading.value = true
  try {
    const response: any = await getAgenciaCapacidad(id)
    if (response.success) {
      capacidadForm.value = {
        max_salidas_por_dia: response.data?.max_salidas_por_dia ?? 5,
        max_salidas_por_horario: response.data?.max_salidas_por_horario ?? 3
      }
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo cargar la capacidad',
      life: 3000
    })
  } finally {
    capacidadLoading.value = false
  }
}

const saveCapacidad = async () => {
  const maxDia = Number(capacidadForm.value.max_salidas_por_dia ?? 0)
  const maxHorario = Number(capacidadForm.value.max_salidas_por_horario ?? 0)

  if (maxDia < 0 || maxHorario < 0) {
    toast.add({ severity: 'warn', summary: 'Validacion', detail: 'Los valores no pueden ser negativos', life: 3000 })
    return
  }

  if (maxHorario > maxDia) {
    toast.add({ severity: 'warn', summary: 'Validacion', detail: 'El maximo por horario no puede ser mayor al maximo por dia', life: 3500 })
    return
  }

  savingCapacidad.value = true
  try {
    const response: any = await updateAgenciaCapacidad(id, {
      max_salidas_por_dia: maxDia,
      max_salidas_por_horario: maxHorario
    })
    if (response.success) {
      capacidadForm.value = {
        max_salidas_por_dia: response.data?.max_salidas_por_dia ?? maxDia,
        max_salidas_por_horario: response.data?.max_salidas_por_horario ?? maxHorario
      }
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Capacidad actualizada', life: 2500 })
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo guardar la capacidad',
      life: 3000
    })
  } finally {
    savingCapacidad.value = false
  }
}

const resolveFotoUrl = (path?: string) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  const base = new URL(useRuntimeConfig().public.apiBase).origin
  const clean = path.replace(/^\.?\//, '')
  return `${base}/${clean}`
}

onMounted(() => {
  loadAgencia()
  loadCategorias()
})
</script>
