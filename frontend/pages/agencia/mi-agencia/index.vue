<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="flex items-center justify-between mb-6">
      <div>
        <p class="text-sm text-gray-500">Mi agencia</p>
        <h1 class="text-3xl font-bold text-gray-900">{{ agencia?.nombre_comercial || 'Agencia' }}</h1>
      </div>
      <div class="flex gap-2">
        <Button label="Editar" icon="pi pi-pencil" severity="warning" @click="navigateTo('/agencia/mi-agencia/editar')" />
        <Button label="Volver" icon="pi pi-arrow-left" outlined @click="navigateTo('/agencia/dashboard')" />
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
            @click="activeTab = tab.value"
          />
        </div>
      </template>
    </Card>

    <div v-if="loading" class="space-y-4">
      <Skeleton height="400px" />
    </div>

    <div v-else-if="!agencia" class="space-y-3">
      <Message severity="warn" :closable="false">No se pudo cargar tu agencia.</Message>
      <Button label="Reintentar" icon="pi pi-refresh" outlined @click="loadAgencia" />
    </div>

    <template v-else>
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

            <div v-if="selectedFotos.length > 0">
              <p class="text-sm text-gray-600 mb-3">
                <i class="pi pi-info-circle"></i>
                Usa las flechas para reordenar.
              </p>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div
                  v-for="(foto, index) in selectedFotos"
                  :key="index"
                  class="relative group border-2 rounded-lg hover:border-blue-300 transition-all overflow-hidden"
                >
                  <img :src="foto.preview" class="w-full h-32 object-cover" :alt="foto.file.name" />

                  <div class="absolute top-2 left-2 flex flex-col gap-1">
                    <span class="px-2 py-1 bg-gray-900 bg-opacity-70 text-white rounded text-xs font-mono">
                      {{ index + 1 }}
                    </span>
                    <span v-if="setPrimeraComoPrincipal && index === 0" class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
                      PRINCIPAL
                    </span>
                  </div>

                  <div class="absolute top-2 right-2 flex flex-col gap-1">
                    <Button icon="pi pi-arrow-up" rounded size="small" text class="bg-white" :disabled="index === 0 || uploading" @click="moveSelectedFotoUp(index)" />
                    <Button icon="pi pi-arrow-down" rounded size="small" text class="bg-white" :disabled="index === selectedFotos.length - 1 || uploading" @click="moveSelectedFotoDown(index)" />
                    <Button icon="pi pi-trash" severity="danger" rounded size="small" :disabled="uploading" class="bg-white" @click="removeSelectedFoto(index)" />
                  </div>

                  <div class="bg-gray-50 p-2">
                    <p class="text-xs text-gray-600 truncate" :title="foto.file.name">{{ foto.file.name }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="fotosOrdenadas.length > 0">
              <h4 class="font-semibold text-gray-900">Fotos actuales</h4>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div
                  v-for="foto in fotosOrdenadas"
                  :key="foto.id"
                  class="relative group border-2 rounded-lg overflow-hidden"
                >
                  <img
                    :src="resolveFotoUrl(foto.foto_url || foto.foto)"
                    class="w-full h-32 object-cover"
                    :alt="foto.titulo || 'Foto'"
                  />
                  <div class="absolute top-2 left-2 flex flex-col gap-1">
                    <span v-if="foto.es_principal" class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
                      PRINCIPAL
                    </span>
                  </div>
                  <div class="absolute top-2 right-2">
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
                  <div class="bg-gray-50 p-2">
                    <p class="text-xs text-gray-600 truncate" :title="foto.titulo || foto.foto_url">
                      {{ foto.titulo || 'Foto' }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="selectedFotos.length === 0 && fotosOrdenadas.length === 0" class="text-center py-8 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
              <i class="pi pi-images text-4xl text-gray-400 mb-2"></i>
              <p class="text-gray-600 text-sm">Haga clic en "Seleccionar Fotos" para agregar imagenes</p>
            </div>
          </div>
        </template>
      </Card>

      <Card v-if="activeTab === 'especialidades'" class="mb-4">
        <template #title>Especialidades</template>
        <template #content>
          <div class="flex flex-wrap gap-3 mb-4 items-center">
            <Select
              v-model="especialidadForm.categoria_id"
              :options="categoriasDisponibles"
              optionLabel="nombre"
              optionValue="id"
              placeholder="Seleccionar categoria"
              class="w-64"
              showClear
            />
            <div class="flex items-center gap-2">
              <Checkbox v-model="especialidadForm.es_principal" binary inputId="principal" />
              <label for="principal" class="text-sm text-gray-700">Principal</label>
            </div>
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

      <Card v-if="activeTab === 'politicas'" class="mb-4">
        <template #title>Políticas de paquetes</template>
        <template #content>
          <div v-if="politicasLoading" class="space-y-3">
            <Skeleton height="240px" />
          </div>

          <div v-else class="space-y-5">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Edad mínima de pago (niños)
              </label>
              <InputNumber v-model="politicasForm.edad_minima_pago" :min="0" :max="18" class="w-full" />
              <p class="text-xs text-gray-500 mt-1">Menores de esta edad no pagan.</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Recargo por paquete privado (%)
              </label>
              <InputNumber
                v-model="politicasForm.recargo_privado_porcentaje"
                mode="decimal"
                suffix="%"
                :min="0"
                :max="100"
                :minFractionDigits="2"
                :maxFractionDigits="2"
                class="w-full"
              />
              <p class="text-xs text-gray-500 mt-1">Se aplica como porcentaje adicional cuando el paquete es privado.</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Política de cancelación
              </label>
              <Textarea
                v-model="politicasForm.politica_cancelacion"
                rows="5"
                autoResize
                class="w-full"
                placeholder="Ej: Cancelación sin costo hasta 24h antes. Luego se cobra 50%..."
              />
              <p class="text-xs text-gray-500 mt-1">Aplica a todos los paquetes de tu agencia.</p>
            </div>

            <div class="flex justify-end">
              <Button label="Guardar cambios" icon="pi pi-save" :loading="savingPoliticas" @click="savePoliticas" />
            </div>
          </div>
        </template>
      </Card>
    </template>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const {
  getMiAgencia,
  uploadFoto: uploadFotoApi,
  removeFoto: removeFotoApi,
  addEspecialidad,
  removeEspecialidad,
  updateAgencia,
  getCategorias,
  getPaquetePoliticas,
  updatePaquetePoliticas
} = useAgencias()

const agencia = ref<any>(null)
const categorias = ref<any[]>([])
const loading = ref(true)
const uploading = ref(false)
const uploadProgress = ref<{ current: number; total: number } | null>(null)
const addingEspecialidad = ref(false)
const savingPagos = ref(false)
const politicasLoading = ref(false)
const savingPoliticas = ref(false)
const activeTab = ref<'general' | 'fotos' | 'especialidades' | 'pagos' | 'politicas'>('general')

const especialidadForm = ref({
  categoria_id: null as number | null,
  es_principal: false
})

const pagoForm = ref({
  acepta_qr: false,
  acepta_transferencia: false,
  acepta_efectivo: false
})

const politicasForm = ref({
  edad_minima_pago: 6 as number | null,
  recargo_privado_porcentaje: 0 as number | null,
  politica_cancelacion: ''
})

const tabs = [
  { label: 'General', value: 'general', icon: 'pi pi-info-circle' },
  { label: 'Fotos', value: 'fotos', icon: 'pi pi-image' },
  { label: 'Especialidades', value: 'especialidades', icon: 'pi pi-tag' },
  { label: 'Pagos', value: 'pagos', icon: 'pi pi-wallet' },
  { label: 'Políticas', value: 'politicas', icon: 'pi pi-file' }
]

const agenciaId = computed(() => Number(agencia.value?.id || 0))

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
  loading.value = true
  try {
    const response: any = await getMiAgencia()
    if (response.success) {
      agencia.value = response.data
      pagoForm.value = {
        acepta_qr: response.data.acepta_qr,
        acepta_transferencia: response.data.acepta_transferencia,
        acepta_efectivo: response.data.acepta_efectivo
      }
      await loadPoliticas(Number(response.data.id))
    } else {
      agencia.value = null
    }
  } catch (error: any) {
    agencia.value = null
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo cargar la agencia', life: 3000 })
  } finally {
    loading.value = false
  }
}

const loadCategorias = async () => {
  try {
    const response: any = await getCategorias()
    categorias.value = response?.data || []
  } catch {
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
    toast.add({ severity: 'warn', summary: 'Limite', detail: 'Maximo 10 fotos', life: 2500 })
    input.value = ''
    return
  }

  const filesToAdd = Array.from(files).slice(0, remainingSlots)

  filesToAdd.forEach((file) => {
    if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
      toast.add({ severity: 'warn', summary: 'Formato no permitido', detail: `${file.name} no es valido`, life: 3000 })
      return
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.add({ severity: 'warn', summary: 'Archivo muy grande', detail: `${file.name} supera 5MB`, life: 3000 })
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
  const id = agenciaId.value
  if (!id) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Agencia no cargada', life: 2500 })
    return
  }

  if (selectedFotos.value.length === 0) {
    toast.add({ severity: 'warn', summary: 'Fotos', detail: 'Seleccione al menos una imagen', life: 2500 })
    return
  }

  if (totalFotos.value > maxFotos) {
    toast.add({ severity: 'warn', summary: 'Limite', detail: 'Maximo 10 fotos', life: 2500 })
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
      if (response.success) okCount++
      else failed.push(item)
    } catch {
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
  const id = agenciaId.value
  if (!id) return
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
  const id = agenciaId.value
  if (!id) return
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
  const id = agenciaId.value
  if (!id) return
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
  const id = agenciaId.value
  if (!id) return
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

const loadPoliticas = async (idOverride?: number) => {
  const id = Number(idOverride || agenciaId.value)
  if (!id) return

  politicasLoading.value = true
  try {
    const response: any = await getPaquetePoliticas(id)
    if (response.success) {
      politicasForm.value = {
        edad_minima_pago: response.data?.edad_minima_pago ?? 6,
        recargo_privado_porcentaje: response.data?.recargo_privado_porcentaje ?? 0,
        politica_cancelacion: response.data?.politica_cancelacion || ''
      }
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudieron cargar las politicas',
      life: 3000
    })
  } finally {
    politicasLoading.value = false
  }
}

const savePoliticas = async () => {
  const id = agenciaId.value
  if (!id) return

  savingPoliticas.value = true
  try {
    const payload = {
      edad_minima_pago: Number(politicasForm.value.edad_minima_pago ?? 6),
      recargo_privado_porcentaje: Number(politicasForm.value.recargo_privado_porcentaje ?? 0),
      politica_cancelacion: politicasForm.value.politica_cancelacion?.trim() ? politicasForm.value.politica_cancelacion.trim() : null
    }

    const response: any = await updatePaquetePoliticas(id, payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Políticas actualizadas', life: 2500 })
      politicasForm.value = {
        edad_minima_pago: response.data?.edad_minima_pago ?? payload.edad_minima_pago,
        recargo_privado_porcentaje: response.data?.recargo_privado_porcentaje ?? payload.recargo_privado_porcentaje,
        politica_cancelacion: response.data?.politica_cancelacion || payload.politica_cancelacion || ''
      }
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudieron guardar las politicas',
      life: 3000
    })
  } finally {
    savingPoliticas.value = false
  }
}

const resolveFotoUrl = (path?: string) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  const base = new URL(useRuntimeConfig().public.apiBase).origin
  const clean = path.replace(/^\.?\//, '')
  return `${base}/${clean}`
}

onMounted(async () => {
  await Promise.all([loadCategorias(), loadAgencia()])
})
</script>
