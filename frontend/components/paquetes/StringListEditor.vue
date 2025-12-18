<template>
  <div class="space-y-2">
    <label v-if="label" class="block text-sm font-medium text-gray-700">{{ label }}</label>

    <div class="flex flex-col sm:flex-row gap-2">
      <InputText
        v-model="newItem"
        class="w-full"
        :placeholder="placeholder"
        :disabled="disabled"
        @keyup.enter="addItem"
      />
      <Button
        :label="addButtonLabel"
        icon="pi pi-plus"
        :disabled="disabled || !newItem.trim()"
        @click="addItem"
      />
    </div>

    <small v-if="help" class="text-gray-500">{{ help }}</small>

    <div v-if="items.length" class="flex flex-wrap gap-2">
      <Tag v-for="(item, idx) in items" :key="`${item}-${idx}`" severity="secondary">
        <span class="inline-flex items-center gap-2">
          <span class="p-tag-label text-sm">{{ item }}</span>
          <Button
            icon="pi pi-times"
            text
            rounded
            severity="danger"
            :disabled="disabled"
            @click="removeItem(idx)"
          />
        </span>
      </Tag>
    </div>

    <div v-else class="text-sm text-gray-500">
      {{ emptyText }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: string[]
    label?: string
    placeholder?: string
    help?: string
    addButtonLabel?: string
    emptyText?: string
    disabled?: boolean
  }>(),
  {
    placeholder: 'Agregar...',
    help: '',
    addButtonLabel: 'Agregar',
    emptyText: 'Sin elementos',
    disabled: false
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const newItem = ref('')

const items = computed(() => props.modelValue || [])

const addItem = () => {
  const value = newItem.value.trim()
  if (!value) return

  const next = [...items.value]
  if (next.some((x) => x.trim().toLowerCase() === value.toLowerCase())) {
    newItem.value = ''
    return
  }

  next.push(value)
  emit('update:modelValue', next)
  newItem.value = ''
}

const removeItem = (idx: number) => {
  const next = [...items.value]
  next.splice(idx, 1)
  emit('update:modelValue', next)
}
</script>

