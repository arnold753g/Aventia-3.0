<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <h3 class="text-lg font-semibold text-gray-900">
        Fotos ({{ fotos.length }}/10)
      </h3>
      <Button
        v-if="editable && fotos.length < 10"
        label="Agregar Foto"
        icon="pi pi-plus"
        size="small"
        @click="showAddDialog = true"
      />
    </div>

    <!-- Gallery -->
    <div v-if="fotos.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div
        v-for="(foto, index) in fotos"
        :key="foto.id || index"
        class="relative group"
      >
        <img
          :src="resolveFotoUrl(foto.foto)"
          :alt="`Foto ${index + 1}`"
          class="w-full h-48 object-cover rounded-lg border border-gray-200"
        />

        <!-- Badge principal -->
        <div v-if="foto.es_principal" class="absolute top-2 left-2">
          <span class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
            PRINCIPAL
          </span>
        </div>

        <!-- Overlay con acciones -->
        <div
          v-if="editable"
          class="absolute inset-0 bg-black bg-opacity-50 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg flex items-center justify-center gap-2"
        >
          <Button
            v-if="!foto.es_principal"
            icon="pi pi-star"
            rounded
            severity="warning"
            v-tooltip.top="'Hacer principal'"
            @click="$emit('set-principal', foto.id)"
          />
          <Button
            icon="pi pi-trash"
            rounded
            severity="danger"
            v-tooltip.top="'Eliminar'"
            @click="confirmDelete(foto)"
          />
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else class="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
      <i class="pi pi-image text-6xl text-gray-400 mb-4"></i>
      <p class="text-gray-600 mb-4">No hay fotos agregadas</p>
      <Button
        v-if="editable"
        label="Agregar Primera Foto"
        icon="pi pi-plus"
        @click="showAddDialog = true"
      />
    </div>

    <!-- Dialog agregar foto -->
    <Dialog
      v-model:visible="showAddDialog"
      header="Agregar Foto"
      :modal="true"
      :style="{ width: '500px' }"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            URL de la foto
          </label>
          <InputText
            v-model="newFotoUrl"
            placeholder="https://ejemplo.com/foto.jpg"
            class="w-full"
          />
        </div>

        <div v-if="newFotoUrl" class="border rounded-lg p-2">
          <img :src="newFotoUrl" alt="Preview" class="w-full h-48 object-cover rounded" />
        </div>

        <div class="flex items-center gap-2">
          <Checkbox v-model="newFotoEsPrincipal" inputId="principal" binary />
          <label for="principal" class="text-sm">Marcar como foto principal</label>
        </div>
      </div>

      <template #footer>
        <Button
          label="Cancelar"
          severity="secondary"
          @click="showAddDialog = false"
        />
        <Button
          label="Agregar"
          icon="pi pi-check"
          :disabled="!newFotoUrl"
          @click="handleAddFoto"
        />
      </template>
    </Dialog>

    <!-- Dialog confirmar eliminación -->
    <Dialog
      v-model:visible="showDeleteDialog"
      header="Confirmar Eliminación"
      :modal="true"
      :style="{ width: '450px' }"
    >
      <div class="flex items-start gap-4">
        <i class="pi pi-exclamation-triangle text-orange-500 text-3xl"></i>
        <div>
          <p class="mb-2">¿Está seguro de eliminar esta foto?</p>
          <p class="text-sm text-gray-500">Esta acción no se puede deshacer.</p>
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancelar"
          severity="secondary"
          @click="showDeleteDialog = false"
        />
        <Button
          label="Eliminar"
          severity="danger"
          @click="handleDeleteFoto"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  fotos: any[]
  editable?: boolean
}>()

const emit = defineEmits(['add-foto', 'delete-foto', 'set-principal'])

// Log para depuración
watch(() => props.fotos, (newValue) => {
  console.log('🎨 FotosGallery recibió fotos:', newValue)
  console.log('🔢 Cantidad de fotos en FotosGallery:', newValue?.length)
  if (newValue && newValue.length > 0) {
    console.log('📷 Primera foto:', newValue[0])
  }
}, { immediate: true })

const showAddDialog = ref(false)
const showDeleteDialog = ref(false)
const newFotoUrl = ref('')
const newFotoEsPrincipal = ref(false)
const selectedFoto = ref<any>(null)
const config = useRuntimeConfig()
const assetsBase = config.public.apiBase.replace(/\/api\/v1\/?$/, '')

const resolveFotoUrl = (path?: string) => {
  if (!path) return ''
  let normalized = path.replace(/\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const handleAddFoto = () => {
  emit('add-foto', {
    foto: newFotoUrl.value,
    es_principal: newFotoEsPrincipal.value,
    orden: props.fotos.length
  })

  newFotoUrl.value = ''
  newFotoEsPrincipal.value = false
  showAddDialog.value = false
}

const confirmDelete = (foto: any) => {
  selectedFoto.value = foto
  showDeleteDialog.value = true
}

const handleDeleteFoto = () => {
  emit('delete-foto', selectedFoto.value.id)
  showDeleteDialog.value = false
  selectedFoto.value = null
}
</script>
