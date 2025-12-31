<template>
  <div class="flex items-center gap-3">
    <div v-if="slug" class="flex-1 flex items-center gap-2 px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg">
      <i class="pi pi-link text-gray-500"></i>
      <span class="text-sm text-gray-700 font-mono truncate">{{ displayUrl }}</span>
    </div>

    <Button
      :label="copyLabel"
      :icon="copyIcon"
      :severity="copiedRecently ? 'success' : 'secondary'"
      :disabled="!slug"
      @click="handleCopy"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useToast } from 'primevue/usetoast'

const props = defineProps<{
  slug?: string
}>()

const toast = useToast()
const { copyToClipboard, getShareUrl } = useAgenciaSlug()

const copiedRecently = ref(false)
const copyIcon = computed(() => copiedRecently.value ? 'pi pi-check' : 'pi pi-copy')
const copyLabel = computed(() => copiedRecently.value ? 'Copiado!' : 'Copiar enlace')

const displayUrl = computed(() => {
  if (!props.slug) return ''
  if (process.client) {
    const origin = window.location.origin
    return `${origin}/agencias/${props.slug}`
  }
  return `/agencias/${props.slug}`
})

const handleCopy = async () => {
  if (!props.slug) return

  const success = await copyToClipboard(props.slug)

  if (success) {
    copiedRecently.value = true
    toast.add({
      severity: 'success',
      summary: 'Enlace copiado',
      detail: 'El enlace de tu agencia se copió al portapapeles',
      life: 3000
    })

    setTimeout(() => {
      copiedRecently.value = false
    }, 3000)
  } else {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'No se pudo copiar el enlace',
      life: 3000
    })
  }
}
</script>
