<template>
  <Button
    :label="label"
    icon="pi pi-download"
    severity="secondary"
    outlined
    :loading="downloading"
    :disabled="disabled || downloading"
    @click="downloadPdf"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import type { PropType } from 'vue'

const props = defineProps({
  titulo: {
    type: String,
    default: 'Confirmacion de compra'
  },
  codigo: {
    type: String,
    required: true
  },
  rows: {
    type: Array as PropType<Array<[string, string]>>,
    default: () => []
  },
  fileName: {
    type: String,
    default: 'confirmacion_compra.pdf'
  },
  label: {
    type: String,
    default: 'Descargar PDF'
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const downloading = ref(false)
const toast = useToast()

const downloadPdf = async () => {
  if (downloading.value || typeof window === 'undefined') return
  downloading.value = true
  try {
    const [jsPDFModule, autoTableModule] = await Promise.all([import('jspdf'), import('jspdf-autotable')])
    const jsPDF = (jsPDFModule as any).default || (jsPDFModule as any).jsPDF
    const autoTable =
      (autoTableModule as any).default ||
      (autoTableModule as any).autoTable ||
      (autoTableModule as any)
    if (!jsPDF || !autoTable) {
      throw new Error('No se pudo cargar el modulo PDF')
    }

    const doc = new jsPDF()
    doc.setFontSize(16)
    doc.text(props.titulo, 14, 16)
    doc.setFontSize(10)
    doc.text(`Codigo: ${props.codigo}`, 14, 22)
    doc.text(`Generado: ${new Date().toLocaleString('es-BO')}`, 14, 27)

    const body = props.rows.map(([key, value]) => [String(key ?? ''), String(value ?? '')])

    autoTable(doc, {
      startY: 34,
      head: [['Detalle', 'Valor']],
      body,
      styles: { fontSize: 9 },
      headStyles: { fillColor: [16, 185, 129] }
    })

    doc.save(props.fileName)
  } catch (err: any) {
    toast.add({
      severity: 'warn',
      summary: 'No se pudo descargar',
      detail: err?.message || 'No se pudo generar el PDF',
      life: 3000
    })
  } finally {
    downloading.value = false
  }
}
</script>
