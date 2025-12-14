<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center gap-4">
          <Button
            icon="pi pi-arrow-left"
            text
            rounded
            @click="navigateTo('/admin/atracciones')"
          />
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Nueva Atracción Turística</h1>
            <p class="text-gray-600 mt-1">Complete el formulario para registrar una nueva atracción</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-5xl mx-auto px-4 py-8">
      <Form
        :resolver="zodResolver(createAtraccionSchema)"
        :initialValues="initialValues"
        @submit="onSubmit"
        class="space-y-6"
      >
        <!-- Información Básica -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-info-circle text-blue-600"></i>
              <span>Información Básica</span>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Nombre -->
              <FormField v-slot="$field" name="nombre" class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nombre de la Atracción *
                </label>
                <InputText
                  v-model="formData.nombre"
                  placeholder="Ej: Salar de Uyuni"
                  class="w-full"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Descripción -->
              <FormField v-slot="$field" name="descripcion" class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Descripción *
                </label>
                <Textarea
                  v-bind="$field"
                  rows="5"
                  placeholder="Describe la atracción turística..."
                  class="w-full"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Departamento -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Departamento *
                </label>
              <Dropdown
                v-model="selectedDepartamento"
                :options="departamentos"
                optionLabel="nombre"
                optionValue="id"
                placeholder="Seleccione departamento"
                class="w-full"
                :disabled="true"
                @change="loadProvincias"
              />
            </div>

              <!-- Provincia -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Provincia *
                </label>
                <Dropdown
                  v-model="localProvinciaId"
                  :options="provincias"
                  optionLabel="nombre"
                  optionValue="id"
                  placeholder="Seleccione provincia"
                  class="w-full"
                  :disabled="!selectedDepartamento"
                />
              </div>

              <!-- Dirección -->
              <FormField v-slot="$field" name="direccion" class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Dirección
                </label>
                <InputText
                  v-bind="$field"
                  placeholder="Dirección específica"
                  class="w-full"
                />
              </FormField>
            </div>
          </template>
        </Card>

        <!-- Ubicación GPS -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-map text-green-600"></i>
              <span>Ubicación GPS</span>
            </div>
          </template>
          <template #content>
            <FormField v-slot="$field" name="latitud">
              <input type="hidden" v-bind="$field" />
            </FormField>
            <FormField v-slot="$field" name="longitud">
              <input type="hidden" v-bind="$field" />
            </FormField>

            <ClientOnly>
              <AtraccionMap
                v-model:latitud="formData.latitud"
                v-model:longitud="formData.longitud"
                :editable="true"
                :show-coordinate-inputs="false"
                height="400px"
              />
              <template #fallback>
                <div class="flex items-center justify-center h-96 bg-gray-100 rounded">
                  <i class="pi pi-spin pi-spinner text-4xl text-gray-400"></i>
                </div>
              </template>
            </ClientOnly>
          </template>
        </Card>

        <!-- Horarios y Precios -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-clock text-purple-600"></i>
              <span>Horarios y Precios</span>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <!-- Horario Apertura -->
              <FormField v-slot="$field" name="horario_apertura">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Horario Apertura
                </label>
                <input type="hidden" v-bind="$field" v-model="$field.value" />
                <DatePicker
                  :modelValue="aperturaPickerValue"
                  timeOnly
                  hourFormat="24"
                  iconDisplay="input"
                  showIcon
                  icon="pi pi-clock"
                  fluid
                  @update:modelValue="(val) => { aperturaPickerValue = val as any; $field.value = dateToTimeString(val as Date | null) }"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Horario Cierre -->
              <FormField v-slot="$field" name="horario_cierre">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Horario Cierre
                </label>
                <input type="hidden" v-bind="$field" v-model="$field.value" />
                <DatePicker
                  :modelValue="cierrePickerValue"
                  timeOnly
                  hourFormat="24"
                  iconDisplay="input"
                  showIcon
                  icon="pi pi-clock"
                  fluid
                  @update:modelValue="(val) => { cierrePickerValue = val as any; $field.value = dateToTimeString(val as Date | null) }"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Precio -->
              <FormField v-slot="$field" name="precio_entrada">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Precio Entrada (Bs.)
                </label>
                <InputNumber
                  v-bind="$field"
                  :min="0"
                  :max="10000"
                  mode="decimal"
                  :minFractionDigits="2"
                  placeholder="0.00"
                  class="w-full"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>
            </div>
          </template>
        </Card>

        <!-- Características -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-star text-orange-600"></i>
              <span>Características</span>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Nivel Dificultad -->
              <FormField v-slot="$field" name="nivel_dificultad">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nivel de Dificultad
                </label>
                <Dropdown
                  v-bind="$field"
                  :options="nivelesOptions"
                  optionLabel="label"
                  optionValue="value"
                  placeholder="Seleccione nivel"
                  class="w-full"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Checkboxes -->
              <div class="space-y-4">
                <FormField v-slot="$field" name="requiere_agencia">
                  <div class="flex items-center gap-2">
                    <Checkbox v-bind="$field" inputId="requiere_agencia" binary />
                    <label for="requiere_agencia" class="text-sm font-medium">
                      Requiere Agencia Turística
                    </label>
                  </div>
                </FormField>

                <FormField v-slot="$field" name="acceso_particular">
                  <div class="flex items-center gap-2">
                    <Checkbox v-bind="$field" inputId="acceso_particular" binary />
                    <label for="acceso_particular" class="text-sm font-medium">
                      Permite Acceso Particular
                    </label>
                  </div>
                </FormField>

                <FormField v-slot="$field" name="visible_publico">
                  <div class="flex items-center gap-2">
                    <Checkbox v-bind="$field" inputId="visible_publico" binary />
                    <label for="visible_publico" class="text-sm font-medium">
                      Visible al Público
                    </label>
                  </div>
                </FormField>
              </div>

              <!-- Estado -->
              <FormField v-slot="$field" name="status">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Estado
                </label>
                <Dropdown
                  v-bind="$field"
                  :options="statusOptions"
                  optionLabel="label"
                  optionValue="value"
                  placeholder="Seleccione estado"
                  class="w-full"
                />
              </FormField>
            </div>
          </template>
        </Card>

        <!-- Mejor Época -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-calendar text-teal-600"></i>
              <span>Mejor Época de Visita</span>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <FormField v-slot="$field" name="mes_inicio_id">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Mes Inicio
                </label>
                <Dropdown
                  v-bind="$field"
                  :options="meses"
                  optionLabel="nombre"
                  optionValue="id"
                  placeholder="Seleccione mes"
                  class="w-full"
                  showClear
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <FormField v-slot="$field" name="mes_fin_id">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Mes Fin
                </label>
                <Dropdown
                  v-bind="$field"
                  :options="meses"
                  optionLabel="nombre"
                  optionValue="id"
                  placeholder="Seleccione mes"
                  class="w-full"
                  showClear
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>
            </div>
          </template>
        </Card>

        <!-- Días de Apertura -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-calendar-plus text-indigo-600"></i>
              <span>Días de Apertura</span>
            </div>
          </template>
          <template #content>
            <FormField v-slot="$field" name="dias_ids">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Seleccione los días *
              </label>
              <div class="flex flex-wrap gap-3">
                <div
                  v-for="dia in dias"
                  :key="dia.id"
                  class="flex items-center"
                >
                  <Checkbox
                    v-model="$field.value"
                    :inputId="`dia-${dia.id}`"
                    :value="dia.id"
                  />
                  <label :for="`dia-${dia.id}`" class="ml-2 text-sm">
                    {{ dia.nombre }}
                  </label>
                </div>
              </div>
              <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                {{ $field.error?.message }}
              </Message>
            </FormField>
          </template>
        </Card>

        <!-- Subcategorías -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-tags text-pink-600"></i>
              <span>Categorías</span>
            </div>
          </template>
          <template #content>
            <SubcategoriasSelector
              v-model="localSubcategoriasIds"
            />
            <Message v-if="subcategoriasError" severity="error" size="small" variant="simple">
              {{ subcategoriasError }}
            </Message>
          </template>
        </Card>

        <!-- Fotos -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-images text-red-600"></i>
              <span>Fotos</span>
            </div>
          </template>
          <template #content>
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Fotografías (Máximo 10)
                </label>
                <p class="text-xs text-gray-500 mb-3">
                  Formatos permitidos: JPG, PNG, WebP. La primera foto será la principal.
                </p>
                <div class="flex items-center gap-3">
                  <input
                    ref="fileInput"
                    type="file"
                    accept="image/jpeg,image/png,image/webp"
                    multiple
                    :disabled="selectedFotos.length >= 10"
                    class="hidden"
                    @change="handleFileSelect"
                  />
                  <Button
                    label="Seleccionar Fotos"
                    icon="pi pi-upload"
                    outlined
                    :disabled="selectedFotos.length >= 10"
                    @click="fileInput?.click()"
                  />
                  <span class="text-sm text-gray-600">
                    {{ selectedFotos.length }} / 10 fotos seleccionadas
                  </span>
                </div>
              </div>

              <!-- Previews de fotos seleccionadas -->
              <div v-if="selectedFotos.length > 0">
                <p class="text-sm text-gray-600 mb-3">
                  <i class="pi pi-info-circle"></i>
                  Usa las flechas para reordenar. La primera foto será la principal.
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
                      <span v-if="index === 0" class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
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
                        :disabled="index === 0"
                        class="bg-white"
                        @click="moveFotoUp(index)"
                      />
                      <Button
                        icon="pi pi-chevron-down"
                        severity="secondary"
                        rounded
                        size="small"
                        :disabled="index === selectedFotos.length - 1"
                        class="bg-white"
                        @click="moveFotoDown(index)"
                      />
                      <Button
                        icon="pi pi-times"
                        severity="danger"
                        rounded
                        size="small"
                        class="bg-white"
                        @click="removeFoto(index)"
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

              <!-- Empty state -->
              <div v-else class="text-center py-8 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
                <i class="pi pi-images text-4xl text-gray-400 mb-2"></i>
                <p class="text-gray-600 text-sm">Haga clic en "Seleccionar Fotos" para agregar imágenes</p>
              </div>
            </div>
          </template>
        </Card>

        <!-- Información de Contacto -->
        <Card>
          <template #title>
            <div class="flex items-center gap-2">
              <i class="pi pi-phone text-cyan-600"></i>
              <span>Información de Contacto</span>
            </div>
          </template>
          <template #content>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Teléfono -->
              <FormField v-slot="$field" name="telefono">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Teléfono
                </label>
                <InputText
                  v-bind="$field"
                  placeholder="70123456"
                  class="w-full"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Email -->
              <FormField v-slot="$field" name="email">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Email
                </label>
                <InputText
                  v-bind="$field"
                  type="email"
                  placeholder="info@atraccion.com"
                  class="w-full"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Sitio Web -->
              <FormField v-slot="$field" name="sitio_web">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Sitio Web
                </label>
                <InputText
                  v-bind="$field"
                  placeholder="https://www.ejemplo.com"
                  class="w-full"
                />
                <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
                  {{ $field.error?.message }}
                </Message>
              </FormField>

              <!-- Facebook -->
              <FormField v-slot="$field" name="facebook">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Facebook
                </label>
                <InputText
                  v-bind="$field"
                  placeholder="@atraccion"
                  class="w-full"
                />
              </FormField>

              <!-- Instagram -->
              <FormField v-slot="$field" name="instagram" class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Instagram
                </label>
                <InputText
                  v-bind="$field"
                  placeholder="@atraccion"
                  class="w-full"
                />
              </FormField>
            </div>
          </template>
        </Card>

        <!-- Botones -->
        <div class="flex justify-end gap-4">
          <Button
            label="Cancelar"
            severity="secondary"
            @click="navigateTo('/admin/atracciones')"
          />
          <Button
            type="submit"
            label="Crear Atracción"
            icon="pi pi-check"
            :loading="submitting"
          />
        </div>
      </Form>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, defineAsyncComponent, computed } from 'vue'
import { zodResolver } from '@primevue/forms/resolvers/zod'
import { useToast } from 'primevue/usetoast'
import { createAtraccionSchema } from '~/utils/validations/atraccion'
import SubcategoriasSelector from '~/components/atracciones/SubcategoriasSelector.vue'

// Componente lazy-loaded solo para el cliente
const AtraccionMap = defineAsyncComponent(() =>
  import('~/components/atracciones/AtraccionMap.vue')
)

definePageMeta({
  middleware: 'auth',
  layout: 'admin'
})

const toast = useToast()
const { createAtraccion, getDepartamentos, getProvincias, getDias, getMeses } = useAtracciones()

const submitting = ref(false)
const selectedDepartamento = ref<number | null>(DEFAULT_DEPARTAMENTO_ID)
const fileInput = ref<HTMLInputElement>()
const selectedFotos = ref<Array<{ file: File; preview: string }>>([])

const departamentos = ref<any[]>([])
const provincias = ref<any[]>([])
const dias = ref<any[]>([])
const meses = ref<any[]>([])
const localSubcategoriasIds = ref<number[]>([])
const subcategoriasError = ref<string>('')
const localProvinciaId = ref<number | null>(null)

const initialValues = reactive({
  nombre: '',
  descripcion: '',
  provincia_id: null,
  direccion: '',
  latitud: null,
  longitud: null,
  horario_apertura: '08:00',
  horario_cierre: '18:30',
  precio_entrada: 0,
  nivel_dificultad: 'facil',
  requiere_agencia: false,
  acceso_particular: true,
  mes_inicio_id: null,
  mes_fin_id: null,
  status: 'activa',
  visible_publico: true,
  telefono: '',
  email: '',
  sitio_web: '',
  facebook: '',
  instagram: '',
  subcategorias_ids: [],
  dias_ids: [],
  fotos: []
})

const formData = reactive({
  nombre: '',
  latitud: null as number | null,
  longitud: null as number | null
})

const nivelesOptions = [
  { label: 'Fácil', value: 'facil' },
  { label: 'Medio', value: 'medio' },
  { label: 'Difícil', value: 'dificil' },
  { label: 'Extremo', value: 'extremo' }
]

const statusOptions = [
  { label: 'Activa', value: 'activa' },
  { label: 'Inactiva', value: 'inactiva' },
  { label: 'Mantenimiento', value: 'mantenimiento' },
  { label: 'Fuera de Temporada', value: 'fuera_temporada' }
]

const DEFAULT_DEPARTAMENTO_ID = 6 // Tarija

const timeStringToDate = (timeStr: string | null | undefined) => {
  if (!timeStr) return null
  const [hours, minutes] = timeStr.split(':').map(Number)
  if (Number.isNaN(hours) || Number.isNaN(minutes)) return null
  const date = new Date()
  date.setHours(hours, minutes, 0, 0)
  return date
}

const dateToTimeString = (date: Date | null) => {
  if (!date) return ''
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  return `${hh}:${mm}`
}

const aperturaPickerValue = computed({
  get: () => timeStringToDate(initialValues.horario_apertura),
  set: (val: Date | null) => {
    initialValues.horario_apertura = dateToTimeString(val)
  }
})

const cierrePickerValue = computed({
  get: () => timeStringToDate(initialValues.horario_cierre),
  set: (val: Date | null) => {
    initialValues.horario_cierre = dateToTimeString(val)
  }
})

// Función para capitalizar palabras
const capitalize = (text: string) => {
  if (!text) return text
  return text
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

// Watcher para capitalizar nombre en tiempo real
watch(() => formData.nombre, (newVal) => {
  const capitalized = capitalize(newVal)
  if (capitalized !== newVal) {
    formData.nombre = capitalized
  }
  // SIEMPRE sincronizar con initialValues
  initialValues.nombre = capitalized
})

// Sincronizar subcategorías con initialValues
watch(localSubcategoriasIds, (newVal) => {
  initialValues.subcategorias_ids = newVal
  // Limpiar error si hay subcategorías seleccionadas
  if (newVal.length > 0) {
    subcategoriasError.value = ''
  }
})

// Sincronizar provincia con initialValues
watch(localProvinciaId, (newVal) => {
  initialValues.provincia_id = newVal
  console.log('📍 Provincia seleccionada:', newVal)
})

const loadDepartamentos = async () => {
  try {
    const response: any = await getDepartamentos()
    if (response.success) {
      departamentos.value = response.data
      const tarija = response.data.find((d: any) => d.nombre?.toLowerCase() === 'tarija')
      selectedDepartamento.value = tarija?.id || DEFAULT_DEPARTAMENTO_ID
      await loadProvincias()
    }
  } catch (error) {
    console.error('Error al cargar departamentos')
  }
}

const loadProvincias = async () => {
  if (!selectedDepartamento.value) {
    provincias.value = []
    return
  }

  try {
    const response: any = await getProvincias(selectedDepartamento.value)
    if (response.success) {
      provincias.value = response.data
    }
  } catch (error) {
    console.error('Error al cargar provincias')
  }
}

const loadDias = async () => {
  try {
    const response: any = await getDias()
    if (response.success) {
      dias.value = response.data
    }
  } catch (error) {
    console.error('Error al cargar días')
  }
}

const loadMeses = async () => {
  try {
    const response: any = await getMeses()
    if (response.success) {
      meses.value = response.data
    }
  } catch (error) {
    console.error('Error al cargar meses')
  }
}

// Manejar selección de archivos
const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const files = input.files
  if (!files) return

  const remainingSlots = 10 - selectedFotos.value.length
  const filesToAdd = Array.from(files).slice(0, remainingSlots)

  filesToAdd.forEach(file => {
    // Validar tipo de archivo
    if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
      toast.add({
        severity: 'warn',
        summary: 'Formato no permitido',
        detail: `${file.name} no es un formato válido`,
        life: 3000
      })
      return
    }

    // Crear preview
    const reader = new FileReader()
    reader.onload = (e) => {
      selectedFotos.value.push({
        file,
        preview: e.target?.result as string
      })
    }
    reader.readAsDataURL(file)
  })

  // Limpiar input
  input.value = ''
}

// Remover foto seleccionada
const removeFoto = (index: number) => {
  selectedFotos.value.splice(index, 1)
}

// Mover foto hacia arriba
const moveFotoUp = (index: number) => {
  if (index === 0) return
  const temp = selectedFotos.value[index]
  selectedFotos.value[index] = selectedFotos.value[index - 1]
  selectedFotos.value[index - 1] = temp
}

// Mover foto hacia abajo
const moveFotoDown = (index: number) => {
  if (index === selectedFotos.value.length - 1) return
  const temp = selectedFotos.value[index]
  selectedFotos.value[index] = selectedFotos.value[index + 1]
  selectedFotos.value[index + 1] = temp
}


// Normalizar valores del submit (PrimeVue a veces no llena event.values)
const buildSubmitValues = (event: any) => {
  const stateValues = event?.states
    ? Object.fromEntries(
        Object.entries(event.states).map(([key, state]: [string, any]) => [key, (state as any)?.value])
      )
    : {}

  return {
    ...stateValues,
    ...(event?.values || {}),
    subcategorias_ids: localSubcategoriasIds.value || [],
    provincia_id: localProvinciaId.value,
    latitud: formData.latitud,
    longitud: formData.longitud
  }
}

const onSubmit = async (event: any) => {
  console.log('Raw event:', event)
  const values = buildSubmitValues(event)
  console.log('Valores normalizados:', values)

  submitting.value = true

  try {
    if (!values.provincia_id) {
      toast.add({
        severity: 'warn',
        summary: 'Validacion',
        detail: 'Debe seleccionar una provincia',
        life: 3000
      })
      submitting.value = false
      return
    }

    if (!values.subcategorias_ids || values.subcategorias_ids.length === 0) {
      subcategoriasError.value = 'Debe seleccionar al menos una subcategoria'
      toast.add({
        severity: 'warn',
        summary: 'Validacion',
        detail: subcategoriasError.value,
        life: 3000
      })
      submitting.value = false
      return
    }

    if (!values.nivel_dificultad) {
      toast.add({
        severity: 'warn',
        summary: 'Validacion',
        detail: 'Debe seleccionar un nivel de dificultad',
        life: 3000
      })
      submitting.value = false
      return
    }

    // Si hay fotos, enviar como FormData
    if (selectedFotos.value.length > 0) {
      const formData = new FormData()
      const precioEntrada = values.precio_entrada ?? 0

      // Agregar campos de texto
      formData.append('nombre', values.nombre || '')
      formData.append('descripcion', values.descripcion || '')
      formData.append('provincia_id', String(values.provincia_id))
      formData.append('direccion', values.direccion || '')

      if (values.latitud !== null && values.latitud !== undefined) formData.append('latitud', String(values.latitud))
      if (values.longitud !== null && values.longitud !== undefined) formData.append('longitud', String(values.longitud))

      formData.append('horario_apertura', values.horario_apertura || '')
      formData.append('horario_cierre', values.horario_cierre || '')
      formData.append('precio_entrada', String(precioEntrada))
      formData.append('nivel_dificultad', values.nivel_dificultad)
      formData.append('requiere_agencia', String(!!values.requiere_agencia))
      formData.append('acceso_particular', String(!!values.acceso_particular))
      formData.append('visible_publico', String(!!values.visible_publico))
      formData.append('status', values.status || 'activa')

      if (values.mes_inicio_id != null) formData.append('mes_inicio_id', String(values.mes_inicio_id))
      if (values.mes_fin_id != null) formData.append('mes_fin_id', String(values.mes_fin_id))

      formData.append('telefono', values.telefono || '')
      formData.append('email', values.email || '')
      formData.append('sitio_web', values.sitio_web || '')
      formData.append('facebook', values.facebook || '')
      formData.append('instagram', values.instagram || '')

      // Agregar arrays
      formData.append('subcategorias_ids', values.subcategorias_ids.join(','))
      if (values.dias_ids && values.dias_ids.length > 0) {
        formData.append('dias_ids', values.dias_ids.join(','))
      }

      // Agregar fotos
      selectedFotos.value.forEach(foto => {
        formData.append('fotos', foto.file)
      })

      console.log('Enviando FormData con', selectedFotos.value.length, 'fotos')

      const response: any = await createAtraccion(formData)

      if (response.success) {
        toast.add({
          severity: 'success',
          summary: 'Atraccion Creada',
          detail: 'La atraccion ha sido registrada exitosamente',
          life: 3000
        })

        setTimeout(() => {
          navigateTo('/admin/atracciones')
        }, 1500)
      }
    } else {
      // Sin fotos, enviar JSON
      const data = {
        ...values,
        precio_entrada: values.precio_entrada ?? 0
      }

      console.log('Data a enviar:', JSON.stringify(data, null, 2))

      const response: any = await createAtraccion(data)

      if (response.success) {
        toast.add({
          severity: 'success',
          summary: 'Atraccion Creada',
          detail: 'La atraccion ha sido registrada exitosamente',
          life: 3000
        })

        setTimeout(() => {
          navigateTo('/admin/atracciones')
        }, 1500)
      }
    }
  } catch (error: any) {
    console.error('Error completo:', error)
    console.error('Error data:', error.data)
    console.error('Error response:', error.response)

    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || error.data?.message || 'Error al crear atraccion',
      life: 3000
    })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadDepartamentos()
  loadDias()
  loadMeses()
})
</script>
