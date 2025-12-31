<template>
  <Card class="surface-card">
    <template #title>
      <div class="flex items-center gap-2">
        <i class="pi pi-credit-card text-blue-600"></i>
        <span>Registrar pago</span>
      </div>
    </template>

    <template #content>
      <div class="space-y-4">
        <Message v-if="error" severity="error" :closable="false">
          {{ error }}
        </Message>

        <Message v-if="warning" severity="warn" :closable="false">
          {{ warning }}
        </Message>

        <Message v-if="successMessage" severity="success" :closable="false">
          {{ successMessage }}
        </Message>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Método de pago</label>
          <Dropdown
            v-model="metodoPago"
            :options="metodos"
            optionLabel="label"
            optionValue="value"
            :disabled="isLocked"
            class="w-full"
            placeholder="Seleccione"
          />
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Monto (Bs.)</label>
          <InputNumber v-model="montoLocal" class="w-full" :min="0" :maxFractionDigits="2" :disabled="isLocked" />
          <p class="text-xs text-gray-500">Debe coincidir con el total de la compra.</p>
        </div>

        <div v-if="needsComprobante" class="space-y-2">
          <label class="block text-sm font-medium text-gray-700">Comprobante</label>
          <input
            ref="fileInput"
            type="file"
            accept="image/png,image/jpeg,image/webp"
            class="block w-full text-sm text-gray-700 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-gray-100 file:text-gray-700 hover:file:bg-gray-200"
            :disabled="isLocked"
            @change="onFileChange"
          />
          <p class="text-xs text-gray-500">PNG/JPG/WebP. Máx. 5MB.</p>
        </div>

        <div v-if="needsComprobante && comprobantePreviewUrl" class="rounded-lg border border-gray-200 bg-gray-50 p-3 space-y-2">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-700">Vista previa</p>
            <Button v-if="!isLocked" label="Quitar" icon="pi pi-times" severity="secondary" text @click="clearComprobante" />
          </div>
          <img :src="comprobantePreviewUrl" alt="Comprobante" class="w-full max-h-80 object-contain rounded-md bg-white" />
          <p v-if="comprobante" class="text-xs text-gray-500 truncate">{{ comprobante.name }}</p>
        </div>

        <div class="flex justify-end gap-2">
          <Button
            label="Registrar pago"
            icon="pi pi-check"
            :loading="loading"
            :disabled="!canSubmit"
            @click="submit"
          />
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import type { CrearPagoRequest } from '~/types/pago'

const props = defineProps<{
  compraId: number
  monto: number
}>()

const emit = defineEmits<{
  (e: 'pago-registrado'): void
}>()

const toast = useToast()
const { crearPago } = usePago()

const metodoPago = ref<'efectivo' | 'qr' | 'transferencia' | null>(null)
const comprobante = ref<File | null>(null)
const comprobantePreviewUrl = ref<string | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const montoLocal = ref<number>(Number(props.monto || 0))

const loading = ref(false)
const error = ref<string | null>(null)
const warning = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const metodos = [
  // { label: 'Efectivo', value: 'efectivo' },
  { label: 'QR', value: 'qr' },
  { label: 'Transferencia', value: 'transferencia' }
]

const needsComprobante = computed(() => metodoPago.value === 'qr' || metodoPago.value === 'transferencia')

const isLocked = computed(() => loading.value || !!successMessage.value || !!warning.value)

const clearComprobante = () => {
  if (comprobantePreviewUrl.value) URL.revokeObjectURL(comprobantePreviewUrl.value)
  comprobantePreviewUrl.value = null
  comprobante.value = null
  if (fileInput.value) fileInput.value.value = ''
}

const canSubmit = computed(() => {
  if (isLocked.value) return false
  if (!metodoPago.value) return false
  if (!montoLocal.value || montoLocal.value <= 0) return false
  if (needsComprobante.value && !comprobante.value) return false
  return true
})

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0] || null
  if (!file) {
    clearComprobante()
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    clearComprobante()
    error.value = 'El archivo supera el máximo permitido (5MB)'
    return
  }

  const allowedTypes = ['image/png', 'image/jpeg', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    error.value = 'Formato no permitido. Use PNG/JPG/WebP.'
    clearComprobante()
    return
  }

  error.value = null
  comprobante.value = file
  if (comprobantePreviewUrl.value) URL.revokeObjectURL(comprobantePreviewUrl.value)
  comprobantePreviewUrl.value = URL.createObjectURL(file)
}

const submit = async () => {
  if (isLocked.value) return

  error.value = null

  if (!canSubmit.value || !metodoPago.value) return

  loading.value = true
  try {
    const payload: CrearPagoRequest = {
      compra_id: props.compraId,
      metodo_pago: metodoPago.value,
      monto: Number(montoLocal.value || 0),
      comprobante: needsComprobante.value ? comprobante.value : undefined
    }

    const response: any = await crearPago(payload)
    if (response?.success) {
      const msg = response.message || response.data?.mensaje || 'Pago registrado.'
      successMessage.value = msg
      toast.add({ severity: 'success', summary: 'Pago registrado', detail: msg, life: 4000 })
      emit('pago-registrado')
      return
    }
    error.value = response?.error?.message || 'No se pudo registrar el pago'
  } catch (err: any) {
    const msg = err?.data?.error?.message || err?.message || 'No se pudo registrar el pago'
    if (String(msg).toLowerCase().includes('ya existe un pago pendiente')) {
      warning.value = msg
      toast.add({ severity: 'warn', summary: 'Pago pendiente', detail: msg, life: 4500 })
      return
    }
    error.value = msg
  } finally {
    loading.value = false
  }
}

watch(needsComprobante, (needs) => {
  if (!needs) clearComprobante()
})

onBeforeUnmount(() => {
  clearComprobante()
})
</script>
